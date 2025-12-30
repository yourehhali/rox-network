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
	
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "../../database/social.db"
	}
	database.InitDB(dbPath)

	migrationsPath := os.Getenv("MIGRATIONS_PATH")
	if migrationsPath == "" {
		migrationsPath = "../../pkg/db/migrations/sqlite"
	}
	if err := database.RunMigrations(database.DB, migrationsPath); err != nil {
		log.Printf("Warning: Failed to run migrations: %v", err)
	} else {
		fmt.Println("Database migrations applied successfully.")
	}

	http.HandleFunc("/chat/ws", handlers.ChatWebSocket)
	http.HandleFunc("/chat/private/history", handlers.GetPrivateChatHistory)
	http.HandleFunc("/chat/group/history", handlers.GetGroupChatHistory)
	http.HandleFunc("/chat/unread/count", handlers.GetUnreadCount)
	http.HandleFunc("/chat/private/read", handlers.MarkPrivateRead)

	port := ":8085"
	fmt.Printf("Chat Service running on %s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
