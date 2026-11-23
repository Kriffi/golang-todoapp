#!/bin/sh
set -eu

: "${DATABASE_URL:?Set DATABASE_URL to the PostgreSQL connection string}"

# The Go drivers in this image do not support Neon's channel_binding option.
# Verify the PostgreSQL server certificate and hostname instead.
case "$DATABASE_URL" in
    *channel_binding=require*)
        DATABASE_URL=$(printf '%s' "$DATABASE_URL" | sed -E \
            -e 's/([?&])channel_binding=require&/\1/' \
            -e 's/[?&]channel_binding=require$//' \
            -e 's/([?&])sslmode=require(&|$)/\1sslmode=verify-full\2/')
        export DATABASE_URL
        ;;
esac

# Migration locks must stay on one PostgreSQL session. Neon pooler hostnames
# can be converted to their direct endpoint without changing the app connection.
migration_url=$(printf '%s' "$DATABASE_URL" | sed -E \
    's~^([a-zA-Z]+://[^/@]+@ep-[a-z0-9-]+)-pooler(\.([a-z0-9-]+\.)*neon\.tech)([:/?]|$)~\1\2\4~')

case "$migration_url" in
    *\?*) separator='&' ;;
    *) separator='?' ;;
esac

if /app/exe check-migrations /app/migrations; then
    :
else
    status=$?
    if [ "$status" -ne 1 ]; then
        exit "$status"
    fi
    migrate -path /app/migrations -database "${migration_url}${separator}search_path=public" up
fi
export HTTP_ADDR=:${PORT:-10000}
exec /app/exe
