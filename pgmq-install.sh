#!/bin/sh
set -e

echo "Checking if PGMQ schema exists in database '$DB_NAME'..."
if PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -tAc "SELECT 1 FROM pg_namespace WHERE nspname = 'pgmq';" | grep -q 1; then
  echo "PGMQ is already installed. Skipping."
  exit 0
fi

echo "PGMQ not found. Installing PGMQ v1.9.0 from GitHub..."
pgmq-cli install -d "postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:5432/$DB_NAME" install-from-github -v 1.9.0
