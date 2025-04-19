#!/bin/sh

set -e

echo "Run db mirgration"
/app/migrate -path /app/migrations -database "$DB_SOURCE" -verbose up

echo "Start app"
exec "$@"