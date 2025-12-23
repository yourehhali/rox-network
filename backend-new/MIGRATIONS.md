# Database Migrations Guide

## Overview

We have migrated from a single `schema.sql` file to a robust, version-controlled migration system using `rubenv/sql-migrate`. This allows us to track database changes over time, roll back changes if needed, and ensure all environments are consistent.

## Directory Structure

Migrations are located in `backend-new/migrations/sqlite`. Each migration is a single `.sql` file containing both Up and Down steps:

Example: `migrations/sqlite/000001_create_users_table.sql`

```sql
-- +migrate Up
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    ...
);

-- +migrate Down
DROP TABLE users;
```

## How It Works

### 1. Automatic Application
The **Auth Service** (`services/auth-service`) is currently configured to automatically run pending migrations on startup.
- When `auth-service` starts, it checks the `migrations/sqlite` directory.
- It compares the applied migrations (stored in `gorp_migrations` table by default) with the files on disk.
- It applies any new migrations in sequence.

### 2. Manual Control (CLI)
You can also use the `sql-migrate` CLI tool if installed.

## Transition from `schema.sql`

- **Old Way**: `schema.sql` was a single file executed blindly.
- **New Way**: Migrations are versioned and annotated.
- **Action Required**:
  - We have updated `auth-service` to use `RunMigrations` with `sql-migrate`.
  - **Do NOT** use `schema.sql` anymore.
  - If you need to change the DB schema, create a new `.sql` file with the next sequence number (e.g., `000014_add_column_xyz.sql`) and include `-- +migrate Up` and `-- +migrate Down` sections.

## Adding a New Migration

1. Create a new file with the next sequence number (e.g., `000014_...sql`).
2. Add the `-- +migrate Up` section with your changes.
3. Add the `-- +migrate Down` section with the revert logic.
4. Restart the `auth-service`, and the changes will be applied.
