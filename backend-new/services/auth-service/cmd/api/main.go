package main

import (
	"social-network/services/auth-service/internal/handlers"
	"fmt"
	"log"
	"net/http"
	"os"
	"social-network/pkg/database"
	"social-network/pkg/middleware"
)

func main() {

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "../../database/social.db"
	}
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		log.Printf("Warning: Database file %s does not exist. Creating it now.", dbPath)
		if err := os.MkdirAll("../../database", 0755); err != nil {
			log.Fatalf("Failed to create database directory: %v", err)
		}
		if _,err := os.Create(dbPath); err != nil {
			log.Fatalf("Failed to create database file: %v", err)
		}
	}

	database.InitDB(dbPath)
	fmt.Println("Auth Service connected to database at", dbPath)

	migrationsPath := os.Getenv("MIGRATIONS_PATH")
	if migrationsPath == "" {
		migrationsPath = "../../pkg/db/migrations/sqlite"
	}
	log.Println("Migrations Path:", migrationsPath)
	if err := database.RunMigrations(database.DB, migrationsPath); err != nil {
		log.Printf("Warning: Failed to run migrations: %v", err)
	} else {
		fmt.Println("Database migrations applied successfully.")
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/login", handlers.Login)
	mux.HandleFunc("/register", handlers.Register)
	mux.HandleFunc("/logout", handlers.Logout)
	mux.HandleFunc("/checkstate", handlers.CheckState)

	handler := middleware.CORSMiddleware(mux)

	port := ":8081"
	fmt.Printf("Auth Service running on %s\n", port)
	if err := http.ListenAndServe(port, handler); err != nil {
		log.Fatal(err)
	}
}
