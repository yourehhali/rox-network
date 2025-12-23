package main

import (
	"fmt"
	"group-service/internal/handlers"
	"log"
	"net/http"
	"os"
	"social-network/pkg/database"
	"social-network/pkg/middleware"
)

func main() {
	// Initialize Database
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "../../database/forum.db"
	}
	database.InitDB(dbPath)

	mux := http.NewServeMux()

	// Group Routes
	mux.HandleFunc("/Create_Group", middleware.AuthMiddleware(handlers.CreateGroup))
	mux.HandleFunc("/Get_Groups", middleware.AuthMiddleware(handlers.GetGroups))
	mux.HandleFunc("/join", middleware.AuthMiddleware(handlers.RequestJoin))
	mux.HandleFunc("/accept-reject-join", middleware.AuthMiddleware(handlers.AcceptRejectJoin))
	mux.HandleFunc("/Get_Group_By_ID", middleware.AuthMiddleware(handlers.GetGroupByID))
	mux.HandleFunc("/CreateGroupPost", middleware.AuthMiddleware(handlers.CreateGroupPost))
	mux.HandleFunc("/GetGroupPosts", middleware.AuthMiddleware(handlers.GetGroupPosts))
	mux.HandleFunc("/groups/not-member", middleware.AuthMiddleware(handlers.FetchNotMembers))
	mux.HandleFunc("/groups/invite", middleware.AuthMiddleware(handlers.InviteToGroup))
	mux.HandleFunc("/group/invite/respond", middleware.AuthMiddleware(handlers.RespondToGroupInvite))

	// Event Routes
	mux.HandleFunc("/create_event", middleware.AuthMiddleware(handlers.CreateEvent))
	mux.HandleFunc("/group-events", middleware.AuthMiddleware(handlers.GetGroupEvents))
	mux.HandleFunc("/event-invite-response", middleware.AuthMiddleware(handlers.RespondToEvent))

	port := ":8086" // Group Service Port
	fmt.Printf("Group Service running on %s\n", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatal(err)
	}
}
