#!/bin/sh
# A simple migration script to create the initial database schema.
# It expects the environment variable DATABASE_URL or uses default
# Postgres connection parameters from the .env file.
set -e

DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-postgres}
DB_PASSWORD=${DB_PASSWORD:-postgres}
DB_NAME=${DB_NAME:-postgres}
DB_SSLMODE=${DB_SSLMODE:-disable}

export PGPASSWORD="$DB_PASSWORD"

psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" <<'SQL'
CREATE TABLE IF NOT EXISTS bridge_messages (
    id VARCHAR PRIMARY KEY,
    direction VARCHAR NOT NULL,
    payload JSONB NOT NULL,
    status VARCHAR NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
SQL

echo "Database migration completed."