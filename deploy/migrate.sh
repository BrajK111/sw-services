#!/usr/bin/env bash
# Applies migrations/ddl/*.sql in order against the configured Postgres DB.
# Usage: DATABASE_URL=postgres://user:pass@host:port/dbname ./deploy/migrate.sh
set -euo pipefail

DATABASE_URL="${DATABASE_URL:-postgres://postgres:postgres@localhost:5432/rainmaker_new?sslmode=disable}"
DDL_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../migrations/ddl" && pwd)"

for f in "$DDL_DIR"/*.sql; do
  echo "Applying $(basename "$f")"
  psql "$DATABASE_URL" -v ON_ERROR_STOP=1 -f "$f"
done

echo "Migrations complete."
