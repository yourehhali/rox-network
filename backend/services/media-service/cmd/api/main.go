package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"social-network/pkg/database"
	"social-network/pkg/middleware"
	"social-network/services/media-service/internal/handlers"
)

func main() {
	
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "../../database/social.db"
	}
	database.InitDB(dbPath)

	
	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "../../uploads"
	}
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Printf("Warning: could not create uploads dir: %v", err)
	}

	mux := http.NewServeMux()

	
	
	
	protectedMedia := middleware.AuthMiddleware(handlers.ServeMedia)
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", protectedMedia))

	port := ":8084"
	fmt.Printf("Media Service running on %s, serving %s\n", port, uploadDir)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatal(err)
	}
}
