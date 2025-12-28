package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"social-network/pkg/models"
	sharedRepo "social-network/pkg/repository"
	"social-network/services/notification-service/internal/repository"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WSHandler handles WebSocket connections
func WSHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID, err := sharedRepo.GetUserIDBySessionToken(cookie.Value)
	if err != nil || userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	sessionID := cookie.Value
	repository.AddNotificationClient(userID, sessionID, conn)

	// Keep connection alive
	for {
		if _, _, err := conn.NextReader(); err != nil {
			repository.RemoveNotificationClient(userID, sessionID, conn)
			break
		}
	}
}

// HandleInternalNotification handles POST /internal/notification
func HandleInternalNotification(w http.ResponseWriter, r *http.Request) {
	var n models.Notification
	if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	switch n.Type {
	case "Invitation_friendships", "group_invite", "request_join_group", "event_created":
		// New Notification
		// Enrich sender info if missing
		if n.Sender.Username == "" && n.SenderID != 0 {
			sender, _ := sharedRepo.GetBasicUserInfo(n.SenderID)
			if sender != nil {
				n.Sender = *sender
			}
		}

		id, err := repository.InsertNotification(n)
		if err != nil {
			log.Printf("Error saving notification: %v", err)
			http.Error(w, "DB Error", http.StatusInternalServerError)
			return
		}
		n.ID = id

		// Broadcast to user
		for _, conn := range repository.GetNotificationClients(n.UserID) {
			conn.WriteJSON(n)
		}

	case "notification_removed":
		// Action taken, mark as seen
		// The message field usually contains the type of the original notification (e.g. "Invitation_friendships")
		err := repository.MarkAsSeenByContext(n.UserID, n.SenderID, n.GroupID, n.EventID, n.Message)
		if err != nil {
			log.Printf("Error marking seen: %v", err)
		}

		// Forward removal to UI
		// We send the "notification_removed" message as is, so frontend can remove it from UI list
		for _, conn := range repository.GetNotificationClients(n.UserID) {
			conn.WriteJSON(n)
		}

	default:
		// Ignore unsupported types
		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetNotifications returns unread notifications
func GetNotifications(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	var userID int
	if err == nil {
		userID, _ = sharedRepo.GetUserIDBySessionToken(cookie.Value)
	}
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	notifications, err := repository.GetUnreadNotifications(userID)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Enrich with sender info
	for i := range notifications {
		if notifications[i].SenderID != 0 {
			sender, _ := sharedRepo.GetBasicUserInfo(notifications[i].SenderID)
			if sender != nil {
				notifications[i].Sender = *sender
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifications)
}

// MarkNotificationAsRead handles POST /notifications/read
func MarkNotificationAsRead(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	var userID int
	if err == nil {
		userID, _ = sharedRepo.GetUserIDBySessionToken(cookie.Value)
	}
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		NotificationID int `json:"notification_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err = repository.UpdateNotificationStatus(req.NotificationID, "seen")
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
