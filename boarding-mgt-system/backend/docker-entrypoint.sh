#!/bin/bash
set -e

echo "Waiting for database to be ready..."
until nc -z ${DB_HOST:-postgres} ${DB_PORT:-5432}; do
  echo "Database is unavailable - sleeping"
  sleep 2
done

echo "Database is ready!"

# Run migrations
echo "Running database migrations..."
./migrate up

# Start the server
echo "Starting server..."
exec ./server