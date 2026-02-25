#!/bin/bash
# Usage: ./create_migration.sh migration_name
# Example: ./create_migration.sh add_users_table

if [ -z "$1" ]; then
  echo "Usage: $0 migration_name"
  exit 1
fi

MIGRATIONS_DIR="$(dirname "$0")/migrations"
DATE=$(date +"%Y%m%d%H%M%S")
FILENAME="${DATE}_$1.sql"

mkdir -p "$MIGRATIONS_DIR"
touch "$MIGRATIONS_DIR/$FILENAME"
echo "Created migration: $MIGRATIONS_DIR/$FILENAME"