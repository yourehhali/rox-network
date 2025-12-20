package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB(datasource string) {
	var err error
	DB, err = sql.Open("sqlite3", datasource)
	if err != nil {
		log.Fatalf("Failed open database: %v", err)
	}
	// Verify connection
	if err = DB.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
}

// ExecuteSchema reads a schema file and executes it.
// Useful for initial setup, though migrations are better.
func ExecuteSchema(db *sql.DB, schemaPath string) error {
	schema, err := os.ReadFile(schemaPath)
	if err != nil {
		return err
	}
	_, err = db.Exec(string(schema))
	return err
}
