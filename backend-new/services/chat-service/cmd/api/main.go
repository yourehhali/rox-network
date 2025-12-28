package main

import (
	"social-network/services/chat-service/internal/handlers"
	"fmt"
	"log"
	"net/http"
	"os"
	"social-network/pkg/database"
)

func main() {
	// Initialize Database
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "../../database/social.db"
	}
	database.InitDB(dbPath)

	// Routes
	http.HandleFunc("/chat/ws", handlers.ChatWebSocket)
	http.HandleFunc("/chat/private/history", handlers.GetPrivateChatHistory)
	http.HandleFunc("/chat/group/history", handlers.GetGroupChatHistory)

	port := ":8085"
	fmt.Printf("Chat Service running on %s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
