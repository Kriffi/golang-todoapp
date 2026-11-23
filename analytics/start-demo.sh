#!/bin/bash
set -euo pipefail

runtime=$(mktemp -d)
export PATH="/usr/lib/postgresql/15/bin:$PATH"
export PYTHONPATH="/app:/analytics"
export REDASH_HOST="${REDASH_HOST:-https://${RENDER_EXTERNAL_HOSTNAME:-localhost}}"
export REDASH_PUBLIC_DEMO=1
export REDASH_DEMO_REDIRECT_FILE="$runtime/dashboard-path"
if [ -n "${RENDER_EXTERNAL_HOSTNAME:-}" ]; then
    export REDASH_SOURCE_URL="${REDASH_SOURCE_URL:-https://kriffi-todoapp-be3n.onrender.com/api/v1/statistics/public}"
fi
if [ -n "${REDASH_SOURCE_URL:-}" ]; then
    export REDASH_ENABLED_QUERY_RUNNERS=redash.query_runner.json_ds
else
    export REDASH_ENABLED_QUERY_RUNNERS=redash.query_runner.pg
fi
# A hostname preserves the URL scheme when Redash adds its query options.
export REDASH_REDIS_URL="unix://localhost$runtime/redis.sock?db=0"
export REDASH_LIMITER_STORAGE="redis+unix://localhost$runtime/redis.sock?db=0"
export REPORTING_HOST="$runtime"

random_secret() {
    python -c 'import secrets; print(secrets.token_hex(32))'
}
pg_password=$(random_secret)
export REPORTING_PASSWORD=$(random_secret)
export REDASH_SECRET_KEY=$(random_secret)
export REDASH_COOKIE_SECRET=$(random_secret)
export REDASH_ADMIN_PASSWORD=$(random_secret)
export REDASH_DATABASE_URL="postgresql://redash:$pg_password@/redash?host=$runtime"

cleanup() {
    for pid in "${server_pid:-}" "${worker_pid:-}" "${redis_pid:-}"; do
        if [ -n "$pid" ]; then
            kill -TERM "$pid" 2>/dev/null || true
        fi
    done
    pg_ctl -D "$runtime/pg" -m fast -w stop >/dev/null 2>&1 || true
    rm -r -- "$runtime"
}
trap 'exit 143' TERM
trap 'exit 130' INT
trap cleanup EXIT

initdb -D "$runtime/pg" --username=redash --auth-local=scram-sha-256 --auth-host=scram-sha-256 \
    --pwfile=<(printf '%s\n' "$pg_password") --no-locale --encoding=UTF8 > /dev/null
pg_ctl -D "$runtime/pg" -l "$runtime/postgres.log" \
    -o "-h '' -k $runtime -c shared_buffers=16MB -c max_connections=24" -w start > /dev/null
PGPASSWORD="$pg_password" createdb -h "$runtime" -U redash redash
if [ -z "${REDASH_SOURCE_URL:-}" ]; then
    PGPASSWORD="$pg_password" psql -h "$runtime" -U redash -d redash -v ON_ERROR_STOP=1 -f /analytics/init-demo.sql > /dev/null
fi

redis-server --port 0 --unixsocket "$runtime/redis.sock" --unixsocketperm 700 --save '' --appendonly no \
    --maxmemory 16mb --maxmemory-policy noeviction &
redis_pid=$!
/app/bin/docker-entrypoint create_db

/app/manage.py rq worker queries schemas default &
worker_pid=$!
python /analytics/bootstrap.py

gunicorn -b "0.0.0.0:${PORT:-5000}" --workers 1 --timeout 60 demo_wsgi:app &
server_pid=$!
wait -n "$redis_pid" "$worker_pid" "$server_pid"
