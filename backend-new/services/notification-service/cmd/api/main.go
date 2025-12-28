package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"social-network/pkg/database"
	"social-network/services/notification-service/internal/handlers"
)

func main() {
	// Initialize DB
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "../../database/social.db"
	}
	database.InitDB(dbPath)
	fmt.Println("Notification Service connected to database at", dbPath)

	// Routes
	http.HandleFunc("/ws", handlers.WSHandler)
	http.HandleFunc("/notifications", handlers.GetNotifications)
	http.HandleFunc("/notifications/read", handlers.MarkNotificationAsRead)

	// Internal API
	http.HandleFunc("/internal/notification", handlers.HandleInternalNotification)

	port := ":8087"
	fmt.Printf("Notification Service running on %s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
