package main

import (
	"fmt"
	"log"
	"net/http"
	"social-network/pkg/database"
	"social-network/pkg/middleware"
	"user-service/internal/handlers"
)

func main() {
	dbPath := "../../database/forum.db"
	database.InitDB(dbPath)
	fmt.Println("User Service connected to database at", dbPath)

	mux := http.NewServeMux()

	// Routes
	// Protected by AuthMiddleware
	mux.HandleFunc("/profile/", middleware.AuthMiddleware(handlers.Profile))
	mux.HandleFunc("/toggle-follow", middleware.AuthMiddleware(handlers.SetFollowers))
	mux.HandleFunc("/request_follow", middleware.AuthMiddleware(handlers.RequestFollow))
	mux.HandleFunc("/Friends", middleware.AuthMiddleware(handlers.GetFriendList))
	mux.HandleFunc("/users-sug", middleware.AuthMiddleware(handlers.GetSuggestedUsers))
	mux.HandleFunc("/update-privacy", middleware.AuthMiddleware(handlers.UpdatePrivacy))
	mux.HandleFunc("/update-profile", middleware.AuthMiddleware(handlers.UpdateProfile))
	mux.HandleFunc("/followers", middleware.AuthMiddleware(handlers.FollowersList))
	mux.HandleFunc("/following", middleware.AuthMiddleware(handlers.FollowingList))
	// Add other handlers...

	handler := middleware.CORSMiddleware(mux)

	port := ":8082" // User Service Port
	fmt.Printf("User Service running on %s\n", port)
	if err := http.ListenAndServe(port, handler); err != nil {
		log.Fatal(err)
	}
}
