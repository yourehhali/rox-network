package main

import (
	"auth-service/internal/handlers"
	"fmt"
	"log"
	"net/http"
	"social-network/pkg/database"
	"social-network/pkg/middleware"
)

func main() {
	// Initialize Database (Shared DB)
	// Assuming running from root of backend-new or services/auth-service
	// We'll use absolute path or relative to root if running via `go run`
	// Ideally pass via ENV var.
	dbPath := "../../database/forum.db" 
	database.InitDB(dbPath)
	fmt.Println("Auth Service connected to database at", dbPath)

	// Execute Schema to ensure tables exist
	schemaPath := "../../database/schema.sql"
	if err := database.ExecuteSchema(database.DB, schemaPath); err != nil {
		log.Printf("Warning: Failed to execute schema: %v", err)
	} else {
		fmt.Println("Database schema executed successfully.")
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
