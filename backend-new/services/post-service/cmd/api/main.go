package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"social-network/services/post-service/internal/handlers"
	"social-network/pkg/database"
	"social-network/pkg/middleware"
)

func main() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "../../database/social.db"
	}
	database.InitDB(dbPath)
	fmt.Println("Post Service connected to database at", dbPath)

	mux := http.NewServeMux()

	// Routes
	mux.HandleFunc("/Createpost", middleware.AuthMiddleware(handlers.CreatePost))
	mux.HandleFunc("/getposts", middleware.AuthMiddleware(handlers.GetPosts))
	mux.HandleFunc("/GetCUser/", middleware.AuthMiddleware(handlers.GetPostsByUser))
	mux.HandleFunc("/reactions", middleware.AuthMiddleware(handlers.Reaction))
	// Pattern for comments is /posts/{id}/comments
	mux.HandleFunc("/posts/", middleware.AuthMiddleware(handlers.SubmitComment))

	handler := middleware.CORSMiddleware(mux)

	port := ":8083" // Post Service Port
	fmt.Printf("Post Service running on %s\n", port)
	if err := http.ListenAndServe(port, handler); err != nil {
		log.Fatal(err)
	}
}
