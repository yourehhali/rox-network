package repository

import (
	"social-network/pkg/database"
	"social-network/pkg/models"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	ChatClients = make(map[int]map[string]*websocket.Conn)
	ChatMutex   sync.Mutex

	SELECT_CAN_USERS_CHAT = `
		SELECT u1.is_private, u2.is_private,
			   (SELECT COUNT(*) FROM followers WHERE follower_id = ? AND followed_id = ? AND status = 'accepted') as u1_follows_u2,
			   (SELECT COUNT(*) FROM followers WHERE follower_id = ? AND followed_id = ? AND status = 'accepted') as u2_follows_u1
		FROM users u1, users u2
		WHERE u1.id = ? AND u2.id = ?
	`
	INSERT_MESSAGE = `
		INSERT INTO messages (sender_id, receiver_id, group_id, content, is_emoji, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	SELECT_GROUP_MEMBERS    = "SELECT user_id FROM group_members WHERE group_id = ? AND status = 'member'"
	IS_USER_GROUP_MEMBER    = "SELECT COUNT(*) FROM group_members WHERE user_id = ? AND group_id = ? AND status = 'member'"
	SELECT_PRIVATE_MESSAGES = `
		SELECT m.id, m.sender_id, m.content, m.is_emoji, m.created_at,
		       u.id, u.username, u.first_name, u.last_name, u.avatar
		FROM messages m
		JOIN users u ON m.sender_id = u.id
		WHERE (m.sender_id = ? AND m.receiver_id = ?) OR (m.sender_id = ? AND m.receiver_id = ?)
		ORDER BY m.created_at DESC
		LIMIT ? OFFSET ?
	`
	SELECT_GROUP_MESSAGES = `
		SELECT m.id, m.sender_id, m.content, m.is_emoji, m.created_at,
		       u.id, u.username, u.first_name, u.last_name, u.avatar
		FROM messages m
		JOIN users u ON m.sender_id = u.id
		WHERE m.group_id = ?
		ORDER BY m.created_at DESC
		LIMIT ? OFFSET ?
	`

	SELECT_UNREAD_COUNT_PRIVATE = `
		SELECT COUNT(*) FROM messages 
		WHERE receiver_id = ? AND is_read = 0
	`
	UPDATE_MESSAGES_READ_PRIVATE = `
		UPDATE messages SET is_read = 1 
		WHERE receiver_id = ? AND sender_id = ? AND is_read = 0
	`
)

func GetUnreadMessageCount(userID int) (int, error) {
	var count int
	err := database.DB.QueryRow(SELECT_UNREAD_COUNT_PRIVATE, userID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func MarkPrivateMessagesRead(userID, senderID int) error {
	_, err := database.DB.Exec(UPDATE_MESSAGES_READ_PRIVATE, userID, senderID)
	return err
}

func AddChatClient(userID int, sessionID string, conn *websocket.Conn) {
	ChatMutex.Lock()
	defer ChatMutex.Unlock()
	if _, ok := ChatClients[userID]; !ok {
		ChatClients[userID] = make(map[string]*websocket.Conn)
	}
	ChatClients[userID][sessionID] = conn
}

func RemoveChatClient(userID int, sessionID string, conn *websocket.Conn) {
	ChatMutex.Lock()
	defer ChatMutex.Unlock()
	if sessions, ok := ChatClients[userID]; ok {
		if currentConn, ok := sessions[sessionID]; ok {
			if currentConn == conn {
				currentConn.Close()
				delete(sessions, sessionID)
			}
		}
		if len(sessions) == 0 {
			delete(ChatClients, userID)
		}
	}
}

func GetChatClients(userID int) []*websocket.Conn {
	ChatMutex.Lock()
	defer ChatMutex.Unlock()
	var conns []*websocket.Conn
	if sessions, ok := ChatClients[userID]; ok {
		for _, conn := range sessions {
			conns = append(conns, conn)
		}
	}
	return conns
}

func CanUsersChat(user1, user2 int) (bool, error) {
	var u1Private, u2Private bool
	var u1FollowsU2, u2FollowsU1 int
	err := database.DB.QueryRow(SELECT_CAN_USERS_CHAT, user1, user2, user2, user1, user1, user2).Scan(
		&u1Private, &u2Private, &u1FollowsU2, &u2FollowsU1)
	if err != nil {
		return false, err
	}
	return u1FollowsU2 > 0 || u2FollowsU1 > 0 || !u2Private, nil
}

func SaveMessage(senderID int, receiverID, groupID *int, content string, isEmoji bool) (int, error) {
	result, err := database.DB.Exec(INSERT_MESSAGE, senderID, receiverID, groupID, content, isEmoji, time.Now())
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return int(id), err
}

func DeleteNotificationByID(id int) error {
	_, err := database.DB.Exec("DELETE FROM notifications WHERE id = ?", id)
	return err
}

func GetGroupMembers(groupID int) ([]int, error) {
	rows, err := database.DB.Query(SELECT_GROUP_MEMBERS, groupID)
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

func IsUserGroupMember(userID, groupID int) (bool, error) {
	var count int
	err := database.DB.QueryRow(IS_USER_GROUP_MEMBER, userID, groupID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func GetPrivateMessages(userID, otherUserID, limit, offset int) ([]models.Message, error) {
	rows, err := database.DB.Query(SELECT_PRIVATE_MESSAGES, userID, otherUserID, otherUserID, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := make([]models.Message, 0)
	for rows.Next() {
		var msg models.Message
		var createdAt time.Time
		err := rows.Scan(
			&msg.ID, &msg.SenderID, &msg.Content, &msg.IsEmoji, &createdAt,
			&msg.Sender.ID, &msg.Sender.Username, &msg.Sender.FirstName, &msg.Sender.LastName, &msg.Sender.Avatar,
		)
		if err != nil {
			continue
		}
		msg.CreatedAt = createdAt.Format(time.RFC3339)
		msg.ReceiverID = &otherUserID
		messages = append(messages, msg)
	}

	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

func GetGroupMessages(groupID, limit, offset int) ([]models.Message, error) {
	rows, err := database.DB.Query(SELECT_GROUP_MESSAGES, groupID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := make([]models.Message, 0)
	for rows.Next() {
		var msg models.Message
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

	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}
