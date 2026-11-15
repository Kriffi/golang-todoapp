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

case "$DATABASE_URL" in
    *\?*) separator='&' ;;
    *) separator='?' ;;
esac

migrate -path /app/migrations -database "${DATABASE_URL}${separator}search_path=public" up
export HTTP_ADDR=:${PORT:-10000}
exec /app/exe
