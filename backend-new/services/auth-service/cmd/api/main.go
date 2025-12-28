package main

import (
	"auth-service/internal/handlers"
	"fmt"
	"log"
	"net/http"
	"os"
	"social-network/pkg/database"
	"social-network/pkg/middleware"
)

func main() {
	// Initialize Database (Shared DB)
	// Assuming running from root of backend-new or services/auth-service
	// We'll use absolute path or relative to root if running via `go run`
	// Ideally pass via ENV var.
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "../../database/social.db"
	} 

	database.InitDB("/home/soubaaiss/Desktop/rox-network/backend-new/database/social.db")
	fmt.Println("Auth Service connected to database at", dbPath)

	// Run Migrations
	migrationsPath := os.Getenv("MIGRATIONS_PATH")
	if migrationsPath == "" {
		migrationsPath = "home/soubaaiss/Desktop/rox-network/backend-new/pkg/db/migrations/sqlite"
	}
	if err := database.RunMigrations(database.DB, "/home/soubaaiss/Desktop/rox-network/backend-new/pkg/db/migrations/sqlite"); err != nil {
		log.Printf("Warning: Failed to run migrations: %v", err)
	} else {
		fmt.Println("Database migrations applied successfully.")
	}

	mux := http.NewServeMux()

	// Routes
	mux.HandleFunc("/login", handlers.Login)
	mux.HandleFunc("/register", handlers.Register)
	mux.HandleFunc("/logout", handlers.Logout)
	mux.HandleFunc("/checkstate", handlers.CheckState)

	// Wrap with CORS
	handler := middleware.CORSMiddleware(mux)

	port := ":8081" // Auth Service Port
	fmt.Printf("Auth Service running on %s\n", port)
	if err := http.ListenAndServe(port, handler); err != nil {
		log.Fatal(err)
	}
}
