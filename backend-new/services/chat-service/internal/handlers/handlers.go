package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	sharedRepo "social-network/pkg/repository"
	"social-network/pkg/utils"
	repository "social-network/services/chat-service/internal/repository"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// ChatWebSocket handles chat connections
func ChatWebSocket(w http.ResponseWriter, r *http.Request) {
	// 1. Try to authenticate via Cookie first
	var userID int
	var sessionToken string
	if cookie, err := r.Cookie("session"); err == nil {
		sessionToken = cookie.Value
		userID, _ = sharedRepo.GetUserIDBySessionToken(sessionToken)
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
		token, _ := authMsg["session_token"].(string)

		if msgType != "auth" || token == "" {
			conn.Close()
			return
		}
		sessionToken = token
		userID, err = sharedRepo.GetUserIDBySessionToken(sessionToken)
		if err != nil {
			conn.Close()
			return
		}
	}

	repository.AddChatClient(userID, sessionToken, conn)
	log.Printf("User %d connected to Chat (Session: %s)", userID, sessionToken)

	// Cleanup on exit
	defer repository.RemoveChatClient(userID, sessionToken, conn)

	// Message Loop
	for {
		var msg map[string]interface{}
		if err := conn.ReadJSON(&msg); err != nil {
			break
		}
		handleChatMessage(userID, msg, w)
	}
}

func handleChatMessage(senderID int, msg map[string]interface{}, w http.ResponseWriter) {
	msgType, _ := msg["type"].(string)
	switch msgType {
	case "private_message":
		handlePrivateMessage(senderID, msg, w)
	case "group_message":
		handleGroupMessage(senderID, msg, w)
	}
}

func handlePrivateMessage(senderID int, msg map[string]interface{}, w http.ResponseWriter) {
	receiverIDFloat, _ := msg["receiver_id"].(float64)
	receiverID := int(receiverIDFloat)
	content, _ := msg["content"].(string)
	isEmoji, _ := msg["is_emoji"].(bool)

	if !utils.ValidChatMessage(content) { 
		http.Error(w, "message too long", http.StatusBadRequest) 
		return
	}

	canChat, _ := repository.CanUsersChat(senderID, receiverID)
	if !canChat {
		return
	}

	msgID, _ := repository.SaveMessage(senderID, &receiverID, nil, content, isEmoji)
	sender, _ := sharedRepo.GetBasicUserInfo(senderID)

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

	// Send to receiver (all sessions)
	for _, conn := range repository.GetChatClients(receiverID) {
		conn.WriteJSON(resp)
	}
	// Send back to sender (all sessions)
	for _, conn := range repository.GetChatClients(senderID) {
		conn.WriteJSON(resp)
	}
}

func handleGroupMessage(senderID int, msg map[string]interface{}, w http.ResponseWriter) {
	groupIDFloat, _ := msg["group_id"].(float64)
	groupID := int(groupIDFloat)
	content, _ := msg["content"].(string)
	isEmoji, _ := msg["is_emoji"].(bool)
	 
	if !utils.ValidChatMessage(content) { 
		http.Error(w, "message too long", http.StatusBadRequest) 
		return
	}

	// Verify membership
	isMember, err := repository.IsUserGroupMember(senderID, groupID)
	if err != nil || !isMember {
		log.Printf("User %d is not a member of group %d or error: %v", senderID, groupID, err)
		return
	}

	msgID, _ := repository.SaveMessage(senderID, nil, &groupID, content, isEmoji)
	sender, _ := sharedRepo.GetBasicUserInfo(senderID)

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
		for _, conn := range repository.GetChatClients(memberID) {
			conn.WriteJSON(resp)
		}
	}
}

// NotificationWebSocket handles notification connections

// GetPrivateChatHistory returns chat history between two users
func GetPrivateChatHistory(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	var userID int
	if err == nil {
		userID, _ = sharedRepo.GetUserIDBySessionToken(cookie.Value)
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
		userID, _ = sharedRepo.GetUserIDBySessionToken(cookie.Value)
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
