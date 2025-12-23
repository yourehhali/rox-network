package database

import (
	"database/sql"
	"fmt"
	"log"

	migrate "github.com/rubenv/sql-migrate"
)

// RunMigrations applies all up migrations to the database
func RunMigrations(db *sql.DB, migrationsPath string) error {
	migrations := &migrate.FileMigrationSource{
		Dir: migrationsPath,
	}

	// Apply migrations
	n, err := migrate.Exec(db, "sqlite3", migrations, migrate.Up)
	if err != nil {
		return fmt.Errorf("could not apply migrations: %w", err)
	}

	if n > 0 {
		log.Printf("Applied %d migrations!\n", n)
	} else {
		log.Println("No new migrations to apply.")
	}

	return nil
}
