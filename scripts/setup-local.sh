#!/bin/sh
# Helper script to bootstrap the local development environment.
# It starts the database and Ethereum services, waits for them to
# initialize, runs database migrations, and finally brings up the API.
set -e

echo "Starting dependencies (db, eth)..."
docker-compose up -d db eth

echo "Waiting for database to be ready..."
# Wait for Postgres to accept connections
until docker exec $(docker-compose ps -q db) pg_isready -U ${DB_USER:-postgres} > /dev/null 2>&1; do
  sleep 1
done

echo "Running database migrations..."
./scripts/migrate.sh

echo "Starting API service..."
docker-compose up -d api
echo "Local environment setup complete. API available on port ${APP_PORT:-8080}."