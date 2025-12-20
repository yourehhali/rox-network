package handlers

import (
	"chat-service/internal/repository"
	"encoding/json"
	"log"
	"net/http"
	"social-network/pkg/models"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// ChatWebSocket handles chat connections
func ChatWebSocket(w http.ResponseWriter, r *http.Request) {
	// 1. Try to authenticate via Cookie first
	var userID int
	if cookie, err := r.Cookie("session"); err == nil {
		userID, _ = repository.VerifySession(cookie.Value)
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Chat Upgrade error: %v", err)
		return
	}

	// 2. If cookie auth failed, try JSON handshake (legacy/test support)
	if userID == 0 {
		var authMsg map[string]interface{}
		if err := conn.ReadJSON(&authMsg); err != nil {
			conn.Close()
			return
		}

		msgType, _ := authMsg["type"].(string)
		sessionToken, _ := authMsg["session_token"].(string)

		if msgType != "auth" || sessionToken == "" {
			conn.Close()
			return
		}

		userID, err = repository.VerifySession(sessionToken)
		if err != nil {
			conn.Close()
			return
		}
	}

	repository.AddChatClient(userID, conn)
	log.Printf("User %d connected to Chat", userID)

	// Cleanup on exit
	defer repository.RemoveChatClient(userID)

	// Message Loop
	for {
		var msg map[string]interface{}
		if err := conn.ReadJSON(&msg); err != nil {
			break
		}
		handleChatMessage(userID, msg)
	}
}

func handleChatMessage(senderID int, msg map[string]interface{}) {
	msgType, _ := msg["type"].(string)
	switch msgType {
	case "private_message":
		handlePrivateMessage(senderID, msg)
	case "group_message":
		handleGroupMessage(senderID, msg)
	}
}

func handlePrivateMessage(senderID int, msg map[string]interface{}) {
	receiverIDFloat, _ := msg["receiver_id"].(float64)
	receiverID := int(receiverIDFloat)
	content, _ := msg["content"].(string)
	isEmoji, _ := msg["is_emoji"].(bool)

	canChat, _ := repository.CanUsersChat(senderID, receiverID)
	if !canChat {
		return
	}

	msgID, _ := repository.SaveMessage(senderID, &receiverID, nil, content, isEmoji)
	sender, _ := repository.GetUserInfo(senderID)

	resp := map[string]interface{}{
		"type":        "private_message",
		"id":          msgID,
		"sender_id":   senderID,
		"receiver_id": receiverID,
		"content":     content,
		"is_emoji":    isEmoji,
		"created_at":  time.Now().Format(time.RFC3339),
		"sender":      sender,
	}

	// Send to receiver
	if conn := repository.GetChatClient(receiverID); conn != nil {
		conn.WriteJSON(resp)
	}
	// Send back to sender
	if conn := repository.GetChatClient(senderID); conn != nil {
		conn.WriteJSON(resp)
	}
}

func handleGroupMessage(senderID int, msg map[string]interface{}) {
	groupIDFloat, _ := msg["group_id"].(float64)
	groupID := int(groupIDFloat)
	content, _ := msg["content"].(string)
	isEmoji, _ := msg["is_emoji"].(bool)

	// Verify membership
	isMember, err := repository.IsUserGroupMember(senderID, groupID)
	if err != nil || !isMember {
		log.Printf("User %d is not a member of group %d or error: %v", senderID, groupID, err)
		return
	}

	msgID, _ := repository.SaveMessage(senderID, nil, &groupID, content, isEmoji)
	sender, _ := repository.GetUserInfo(senderID)

	resp := map[string]interface{}{
		"type":       "group_message",
		"id":         msgID,
		"group_id":   groupID,
		"sender_id":  senderID,
		"content":    content,
		"is_emoji":   isEmoji,
		"created_at": time.Now().Format(time.RFC3339),
		"sender":     sender,
	}

	members, _ := repository.GetGroupMembers(groupID)
	for _, memberID := range members {
		if conn := repository.GetChatClient(memberID); conn != nil {
			conn.WriteJSON(resp)
		}
	}
}

// NotificationWebSocket handles notification connections
func NotificationWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Notif Upgrade error: %v", err)
		return
	}

	// Auth Handshake (similar to Chat)
	// If frontend doesn't send auth message for notif, we might need cookie check.
	// But let's assume consistent auth flow or cookie availability.
	// Since standard WS client in JS doesn't send custom headers easily, query param or initial message is best.
	// Old code: socke_tunnel.go used r.Context().Value(UserIDKey).
	// This implies the upgrade request had cookies and passed middleware.
	// So we need to parse cookies manually or hope middleware ran.
	// But `chat-service` main.go might not have full AuthMiddleware if it's behind gateway.
	// Actually, Gateway proxies requests. Headers (Cookie) are preserved.
	// So we can check cookie manually or use middleware.
	// I'll implement a simple cookie check here if context is missing, or rely on auth message if frontend is updated.
	// Given I can't easily change frontend, I'll try to support both or stick to what `socke_tunnel.go` implied.
	// `socke_tunnel.go` had `ID_USER := r.Context().Value(key.UserIDKey)`.
	// This means I MUST wrap this handler in AuthMiddleware in main.go.

	// But wait, WS upgrade might happen before middleware can set context if not careful?
	// Standard http.Handler flow: Middleware -> Handler. So context should be set.
	// I'll assume AuthMiddleware is used in main.go.

	// Wait, if I rely on context, I need to import middleware.
	// But `handlers` package shouldn't import `middleware` if `middleware` imports `handlers` (cycle).
	// `middleware` is in `pkg/middleware`, so it's fine.

	// BUT, `socke_tunnel.go` used `key.UserIDKey`.
	// I need to know what key `pkg/middleware` uses.
	// `pkg/middleware` uses `UserIDKey`.

	// Issue: `r.Context().Value("userID")` (string) or type?
	// Let's check `pkg/middleware`.

	// I'll check `pkg/middleware` in a moment. For now, I'll assume I can get UserID.

	// If UserID is not in context (e.g. middleware not used), I'll try to read cookie.
	// But let's assume middleware is used.

	// Wait, `socke_tunnel.go` does NOT read an initial auth message. It just upgrades and registers.
	// So I should just register the user from Context.

	// Retrieve UserID (assuming middleware put it there)
	// I'll use a helper to get it safely.

	cookie, err := r.Cookie("session")
	var userID int
	if err == nil {
		userID, _ = repository.VerifySession(cookie.Value)
	}

	// If no cookie or invalid, try query param?
	if userID == 0 {
		token := r.URL.Query().Get("session")
		if token != "" {
			userID, _ = repository.VerifySession(token)
		}
	}

	if userID == 0 {
		// Try context if middleware was used
		// But I don't want to depend on middleware package here if I can avoid it to keep loose coupling?
		// Actually, standard is to use middleware.
		conn.Close()
		return
	}

	repository.AddNotificationClient(userID, conn)
	log.Printf("User %d connected to Notifications", userID)
	defer repository.RemoveNotificationClient(userID)

	// Send unread notifications immediately
	unreads, _ := repository.GetUnreadNotifications(userID)
	for _, n := range unreads {
		conn.WriteJSON(n)
	}

	// Keep alive / Read loop
	for {
		if _, _, err := conn.NextReader(); err != nil {
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

	// Save to DB
	// If it's an acceptance/rejection of a follow request, update the original request notification
	if n.Type == "Invitation_friendships" && (n.Message == "accepted" || n.Message == "rejected") {
		// The sender of this "accepted" notification is the one who accepted (Receiver of original request)
		// The user_id of this notification is the one receiving the "accepted" msg (Sender of original request)
		// Wait, we want to update the notification that was shown to the person who *received* the request.
		// Let's look at user-service RequestFollow:
		// input.SenderID is the one who asked to follow. receiverID is the one accepting.
		// Notification sent: UserID=input.SenderID (requester), SenderID=receiverID (accepter).
		
		// But we ALSO need to update the notification for the *accepter* so the "Accept/Reject" buttons go away?
		// Actually, the "Accept/Reject" UI usually depends on the list of *pending* requests.
		// If the user-service `RequestFollow` updates the DB status to 'accepted', the pending list logic (fetching from DB) should handle it.
		// However, the `notifications` table in chat-service might still show the "request" notification.
		// We should mark the "request" notification for the *accepter* as handled.
		
		// The `n` here is the notification being sent to the *Requester* saying "You were accepted".
		// We want to update the notification that was sent to the *Accepter* (n.SenderID) from the *Requester* (n.UserID).
		// So we update: UserID = n.SenderID, SenderID = n.UserID.
		_ = repository.UpdateInvitationStatus(n.SenderID, n.UserID, n.Message)
	}

	id, err := repository.InsertNotification(n)
	if err != nil {
		log.Printf("Error saving notification: %v", err)
		http.Error(w, "DB Error", http.StatusInternalServerError)
		return
	}
	n.ID = int(id)

	// Enrich with Sender info if missing (though sender usually sends ID)
	if n.Sender.Username == "" {
		sender, _ := repository.GetUserInfo(n.SenderID)
		if sender != nil {
			n.Sender = *sender
		}
	}
	// Enrich with timestamps
	n.CreatedAt = time.Now().Format(time.RFC3339)

	// Send to WebSocket
	if conn := repository.GetNotificationClient(n.UserID); conn != nil {
		conn.WriteJSON(n)
	}

	w.WriteHeader(http.StatusOK)
}

// GetNotifications returns unread notifications (HTTP API)
func GetNotifications(w http.ResponseWriter, r *http.Request) {
	// Need userID from context or cookie
	cookie, err := r.Cookie("session")
	var userID int
	if err == nil {
		userID, _ = repository.VerifySession(cookie.Value)
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifications)
}

// GetPrivateChatHistory returns chat history between two users
func GetPrivateChatHistory(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	var userID int
	if err == nil {
		userID, _ = repository.VerifySession(cookie.Value)
	}
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	otherUserStr := r.URL.Query().Get("user_id")
	otherUserID, err := strconv.Atoi(otherUserStr)
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	canChat, err := repository.CanUsersChat(userID, otherUserID)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	if !canChat {
		http.Error(w, "Cannot chat with this user", http.StatusForbidden)
		return
	}

	messages, err := repository.GetPrivateMessages(userID, otherUserID)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

// GetGroupChatHistory returns chat history for a group
func GetGroupChatHistory(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	var userID int
	if err == nil {
		userID, _ = repository.VerifySession(cookie.Value)
	}
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	groupIDStr := r.URL.Query().Get("group_id")
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		http.Error(w, "Invalid group_id", http.StatusBadRequest)
		return
	}

	isMember, err := repository.IsUserGroupMember(userID, groupID)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	if !isMember {
		http.Error(w, "Not a member", http.StatusForbidden)
		return
	}

	messages, err := repository.GetGroupMessages(groupID)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
