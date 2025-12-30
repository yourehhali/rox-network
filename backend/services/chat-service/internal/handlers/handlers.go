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


func ChatWebSocket(w http.ResponseWriter, r *http.Request) {
	
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

	
	defer repository.RemoveChatClient(userID, sessionToken, conn)

	
	for {
		var msg map[string]interface{}
		if err := conn.ReadJSON(&msg); err != nil {
			break
		}
		handleChatMessage(userID, msg, conn)
	}
}

func handleChatMessage(senderID int, msg map[string]interface{}, conn *websocket.Conn) {
	msgType, _ := msg["type"].(string)
	switch msgType {
	case "private_message":
		handlePrivateMessage(senderID, msg, conn)
	case "group_message":
		handleGroupMessage(senderID, msg, conn)
	}
}

func handlePrivateMessage(senderID int, msg map[string]interface{}, conn *websocket.Conn) {
	receiverIDFloat, _ := msg["receiver_id"].(float64)
	receiverID := int(receiverIDFloat)
	content, _ := msg["content"].(string)
	isEmoji, _ := msg["is_emoji"].(bool)

	if !utils.ValidChatMessage(content) {
		conn.WriteJSON(map[string]interface{}{
			"type":    "error",
			"message": "Message too long (max 250 characters)",
		})
		return
	}

	canChat, _ := repository.CanUsersChat(senderID, receiverID)
	if !canChat {
		conn.WriteJSON(map[string]interface{}{
			"type":    "error",
			"message": "You cannot chat with this user",
		})
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

	
	for _, conn := range repository.GetChatClients(receiverID) {
		conn.WriteJSON(resp)
	}
	
	for _, conn := range repository.GetChatClients(senderID) {
		conn.WriteJSON(resp)
	}
}

func handleGroupMessage(senderID int, msg map[string]interface{}, conn *websocket.Conn) {
	groupIDFloat, _ := msg["group_id"].(float64)
	groupID := int(groupIDFloat)
	content, _ := msg["content"].(string)
	isEmoji, _ := msg["is_emoji"].(bool)

	if !utils.ValidChatMessage(content) {
		conn.WriteJSON(map[string]interface{}{
			"type":    "error",
			"message": "Message too long (max 250 characters)",
		})
		return
	}

	
	isMember, err := repository.IsUserGroupMember(senderID, groupID)
	if err != nil || !isMember {
		conn.WriteJSON(map[string]interface{}{
			"type":    "error",
			"message": "You are not a member of this group",
		})
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

	limitStr := r.URL.Query().Get("limit")
	limit := 10
	if limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err == nil && l > 0 {
			limit = l
		}
	}

	offsetStr := r.URL.Query().Get("offset")
	offset := 0
	if offsetStr != "" {
		o, err := strconv.Atoi(offsetStr)
		if err == nil && o >= 0 {
			offset = o
		}
	}

	messages, err := repository.GetPrivateMessages(userID, otherUserID, limit, offset)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}


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

	limitStr := r.URL.Query().Get("limit")
	limit := 10
	if limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err == nil && l > 0 {
			limit = l
		}
	}

	offsetStr := r.URL.Query().Get("offset")
	offset := 0
	if offsetStr != "" {
		o, err := strconv.Atoi(offsetStr)
		if err == nil && o >= 0 {
			offset = o
		}
	}

	messages, err := repository.GetGroupMessages(groupID, limit, offset)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
