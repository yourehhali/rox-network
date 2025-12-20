package repository

import (
	"database/sql"
	"log"
	"social-network/pkg/database"
	"social-network/pkg/models"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	// ChatClients stores WebSocket connections for chat: UserID -> Connection
	ChatClients = make(map[int]*websocket.Conn)
	ChatMutex   sync.Mutex

	// NotificationClients stores WebSocket connections for notifications: UserID -> Connection
	NotificationClients = make(map[int]*websocket.Conn)
	NotifMutex          sync.Mutex
)

// AddChatClient adds a client to the chat pool
func AddChatClient(userID int, conn *websocket.Conn) {
	ChatMutex.Lock()
	defer ChatMutex.Unlock()
	ChatClients[userID] = conn
}

// RemoveChatClient removes a client from the chat pool
func RemoveChatClient(userID int) {
	ChatMutex.Lock()
	defer ChatMutex.Unlock()
	if conn, ok := ChatClients[userID]; ok {
		conn.Close()
		delete(ChatClients, userID)
	}
}

// GetChatClient retrieves a chat client
func GetChatClient(userID int) *websocket.Conn {
	ChatMutex.Lock()
	defer ChatMutex.Unlock()
	return ChatClients[userID]
}

// AddNotificationClient adds a client to the notification pool
func AddNotificationClient(userID int, conn *websocket.Conn) {
	NotifMutex.Lock()
	defer NotifMutex.Unlock()
	NotificationClients[userID] = conn
}

// RemoveNotificationClient removes a client from the notification pool
func RemoveNotificationClient(userID int) {
	NotifMutex.Lock()
	defer NotifMutex.Unlock()
	if conn, ok := NotificationClients[userID]; ok {
		conn.Close()
		delete(NotificationClients, userID)
	}
}

// GetNotificationClient retrieves a notification client
func GetNotificationClient(userID int) *websocket.Conn {
	NotifMutex.Lock()
	defer NotifMutex.Unlock()
	return NotificationClients[userID]
}

// VerifySession verifies the session token and returns the user ID
func VerifySession(token string) (int, error) {
	var userID int
	err := database.DB.QueryRow("SELECT user_id FROM sessions WHERE session_token = ?", token).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

// GetUserInfo retrieves basic user info
func GetUserInfo(userID int) (*models.User, error) {
	var u models.User
	err := database.DB.QueryRow(`
		SELECT id, username, first_name, last_name, avatar, is_private 
		FROM users WHERE id = ?`, userID).
		Scan(&u.ID, &u.Username, &u.FirstName, &u.LastName, &u.Avatar, &u.IsPrivate)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// CanUsersChat checks if two users can chat
func CanUsersChat(user1, user2 int) (bool, error) {
	// Logic: Follow each other OR user2 is public
	// Note: Old logic was u1 follows u2 OR u2 follows u1 OR !u2Private
	// Let's stick to old logic:
	query := `
		SELECT u1.is_private, u2.is_private,
			   (SELECT COUNT(*) FROM followers WHERE follower_id = ? AND followed_id = ? AND status = 'accepted') as u1_follows_u2,
			   (SELECT COUNT(*) FROM followers WHERE follower_id = ? AND followed_id = ? AND status = 'accepted') as u2_follows_u1
		FROM users u1, users u2
		WHERE u1.id = ? AND u2.id = ?
	`
	var u1Private, u2Private bool
	var u1FollowsU2, u2FollowsU1 int
	err := database.DB.QueryRow(query, user1, user2, user2, user1, user1, user2).Scan(
		&u1Private, &u2Private, &u1FollowsU2, &u2FollowsU1)
	if err != nil {
		return false, err
	}
	return u1FollowsU2 > 0 || u2FollowsU1 > 0 || !u2Private, nil
}

// SaveMessage saves a message to the database
func SaveMessage(senderID int, receiverID, groupID *int, content string, isEmoji bool) (int, error) {
	query := `
		INSERT INTO messages (sender_id, receiver_id, group_id, content, is_emoji, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	result, err := database.DB.Exec(query, senderID, receiverID, groupID, content, isEmoji, time.Now())
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return int(id), err
}

// InsertNotification saves a notification
func InsertNotification(n models.Notification) (int, error) {
	// Check duplicates if needed (optional, based on old code)
	// For now, just insert
	query := `
		INSERT INTO notifications (user_id, sender_id, type, message, group_id, event_id, state, created_at)
		VALUES (?, ?, ?, ?, ?, ?, 'unread', ?)
	`
	result, err := database.DB.Exec(query, n.UserID, n.SenderID, n.Type, n.Message, n.GroupID, n.EventID, time.Now())
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return int(id), err
}

// GetUnreadNotifications retrieves unread notifications for a user
func GetUnreadNotifications(userID int) ([]models.Notification, error) {
	rows, err := database.DB.Query(`
		SELECT 
			n.id, n.user_id, n.type, n.message, n.state, n.event_id, n.created_at, 
			u.id, u.username, u.first_name, u.last_name, u.avatar, u.is_private,
			g.id, g.title
		FROM notifications n
		JOIN users u ON n.sender_id = u.id
		LEFT JOIN groups g ON n.group_id = g.id
		WHERE n.user_id = ? AND n.state = 'unread'
		ORDER BY n.created_at DESC
	`, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []models.Notification
	for rows.Next() {
		var n models.Notification
		var groupID sql.NullInt64
		var groupTitle sql.NullString
		var eventID sql.NullInt64 // handle null event_id
		var createdAt time.Time

		if err := rows.Scan(
			&n.ID,
			&n.UserID,
			&n.Type,
			&n.Message,
			&n.State,
			&eventID,
			&createdAt,
			&n.Sender.ID,
			&n.Sender.Username,
			&n.Sender.FirstName,
			&n.Sender.LastName,
			&n.Sender.Avatar,
			&n.Sender.IsPrivate,
			&groupID,
			&groupTitle,
		); err != nil {
			log.Printf("Error scanning notification: %v", err)
			continue
		}

		if groupID.Valid {
			n.GroupID = int(groupID.Int64)
		}
		if groupTitle.Valid {
			n.GroupTitle = groupTitle.String
		}
		if eventID.Valid {
			n.EventID = int(eventID.Int64)
		}
		n.CreatedAt = createdAt.Format(time.RFC3339)

		notifications = append(notifications, n)
	}

	return notifications, nil
}

// UpdateNotificationStatus updates the state of a notification
func UpdateNotificationStatus(id int, state string) error {
	_, err := database.DB.Exec("UPDATE notifications SET state = ? WHERE id = ?", state, id)
	return err
}

// UpdateInvitationStatus updates the state of an invitation (follow request)
func UpdateInvitationStatus(receiverID, senderID int, state string) error {
    _, err := database.DB.Exec("UPDATE notifications SET state = ? WHERE user_id = ? AND sender_id = ? AND type = 'Invitation_friendships' AND message = 'request'", state, receiverID, senderID)
    return err
}

// GetGroupMembers retrieves members of a group
func GetGroupMembers(groupID int) ([]int, error) {
	rows, err := database.DB.Query("SELECT user_id FROM group_members WHERE group_id = ? AND status = 'member'", groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var members []int
	for rows.Next() {
		var id int
		rows.Scan(&id)
		members = append(members, id)
	}
	return members, nil
}

// IsUserGroupMember checks if a user is a member of a group
func IsUserGroupMember(userID, groupID int) (bool, error) {
	var count int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM group_members WHERE user_id = ? AND group_id = ? AND status = 'member'", userID, groupID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Message struct for history (local definition to match repository return)
type Message struct {
	ID         int         `json:"id"`
	SenderID   int         `json:"sender_id"`
	ReceiverID *int        `json:"receiver_id,omitempty"`
	GroupID    *int        `json:"group_id,omitempty"`
	Content    string      `json:"content"`
	IsEmoji    bool        `json:"is_emoji"`
	CreatedAt  string      `json:"created_at"`
	Sender     models.User `json:"sender"`
}

// GetPrivateMessages retrieves private chat history
func GetPrivateMessages(userID, otherUserID int) ([]Message, error) {
	query := `
		SELECT m.id, m.sender_id, m.content, m.is_emoji, m.created_at,
		       u.id, u.username, u.first_name, u.last_name, u.avatar
		FROM messages m
		JOIN users u ON m.sender_id = u.id
		WHERE (m.sender_id = ? AND m.receiver_id = ?) OR (m.sender_id = ? AND m.receiver_id = ?)
		ORDER BY m.created_at ASC
	`
	rows, err := database.DB.Query(query, userID, otherUserID, otherUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := make([]Message, 0)
	for rows.Next() {
		var msg Message
		var createdAt time.Time
		err := rows.Scan(
			&msg.ID, &msg.SenderID, &msg.Content, &msg.IsEmoji, &createdAt,
			&msg.Sender.ID, &msg.Sender.Username, &msg.Sender.FirstName, &msg.Sender.LastName, &msg.Sender.Avatar,
		)
		if err != nil {
			continue
		}
		msg.CreatedAt = createdAt.Format(time.RFC3339)
		msg.ReceiverID = &otherUserID // Approximation for consistent struct
		messages = append(messages, msg)
	}
	return messages, nil
}

// GetGroupMessages retrieves group chat history
func GetGroupMessages(groupID int) ([]Message, error) {
	query := `
		SELECT m.id, m.sender_id, m.content, m.is_emoji, m.created_at,
		       u.id, u.username, u.first_name, u.last_name, u.avatar
		FROM messages m
		JOIN users u ON m.sender_id = u.id
		WHERE m.group_id = ?
		ORDER BY m.created_at ASC
	`
	rows, err := database.DB.Query(query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := make([]Message, 0)
	for rows.Next() {
		var msg Message
		var createdAt time.Time
		err := rows.Scan(
			&msg.ID, &msg.SenderID, &msg.Content, &msg.IsEmoji, &createdAt,
			&msg.Sender.ID, &msg.Sender.Username, &msg.Sender.FirstName, &msg.Sender.LastName, &msg.Sender.Avatar,
		)
		if err != nil {
			continue
		}
		msg.CreatedAt = createdAt.Format(time.RFC3339)
		msg.GroupID = &groupID
		messages = append(messages, msg)
	}
	return messages, nil
}
