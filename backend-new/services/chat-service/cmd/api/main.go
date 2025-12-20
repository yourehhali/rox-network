package main

import (
	"chat-service/internal/handlers"
	"fmt"
	"log"
	"net/http"
	"social-network/pkg/database"
)

func main() {
	// Initialize Database
	database.InitDB("../../database/forum.db")

	// Routes
	http.HandleFunc("/chat/ws", handlers.ChatWebSocket)
	http.HandleFunc("/ws", handlers.NotificationWebSocket)
	http.HandleFunc("/notifications", handlers.GetNotifications)
	http.HandleFunc("/chat/private/history", handlers.GetPrivateChatHistory)
	http.HandleFunc("/chat/group/history", handlers.GetGroupChatHistory)

	// Internal API (Not exposed by Gateway)
	http.HandleFunc("/internal/notification", handlers.HandleInternalNotification)

	port := ":8085"
	fmt.Printf("Chat Service running on %s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
