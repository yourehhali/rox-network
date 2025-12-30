package repository

import (
	"database/sql"
	"log"
	"social-network/pkg/database"
	"social-network/pkg/models"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	NotificationClients = make(map[int]map[string][]*websocket.Conn)
	NotifMutex          sync.Mutex
)



func AddNotificationClient(userID int, sessionID string, conn *websocket.Conn) {
	NotifMutex.Lock()
	defer NotifMutex.Unlock()
	if _, ok := NotificationClients[userID]; !ok {
		NotificationClients[userID] = make(map[string][]*websocket.Conn)
	}
	NotificationClients[userID][sessionID] = append(NotificationClients[userID][sessionID], conn)
}

func RemoveNotificationClient(userID int, sessionID string, conn *websocket.Conn) {
	NotifMutex.Lock()
	defer NotifMutex.Unlock()
	if sessions, ok := NotificationClients[userID]; ok {
		if currentConns, ok := sessions[sessionID]; ok {
			for i, currentConn := range currentConns {
				if currentConn == conn {
					currentConn.Close()
					sessions[sessionID] = append(currentConns[:i], currentConns[i+1:]...)
					break
				}
			}
		}
		if len(sessions) == 0 {
			delete(NotificationClients, userID)
		}
	}
}

func GetNotificationClients(userID int) []*websocket.Conn {
	NotifMutex.Lock()
	defer NotifMutex.Unlock()
	var conns []*websocket.Conn
	if sessions, ok := NotificationClients[userID]; ok {
		for _, conn := range sessions {
			conns = append(conns, conn...)
		}
	}
	return conns
}



func InsertNotification(n models.Notification) (int, error) {
	
	query := `SELECT id FROM notifications 
		WHERE user_id = ? AND sender_id = ? AND type = ? AND state = 'unread'`

	args := []interface{}{n.UserID, n.SenderID, n.Type}

	if n.GroupID != 0 {
		query += ` AND group_id = ?`
		args = append(args, n.GroupID)
	}
	if n.EventID != 0 {
		query += ` AND event_id = ?`
		args = append(args, n.EventID)
	}

	var existingID int
	err := database.DB.QueryRow(query, args...).Scan(&existingID)
	if err == nil {
		return existingID, nil 
	} else if err != sql.ErrNoRows {
		log.Printf("Error checking duplicate: %v", err)
		return 0, err
	}

	
	insertQuery := `INSERT INTO notifications (user_id, sender_id, type, message, state, group_id, event_id, created_at) 
		VALUES (?, ?, ?, ?, 'unread', ?, ?, CURRENT_TIMESTAMP)`

	res, err := database.DB.Exec(insertQuery, n.UserID, n.SenderID, n.Type, n.Message, n.GroupID, n.EventID)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func GetUnreadNotifications(userID int) ([]models.Notification, error) {
	query := `SELECT id, user_id, sender_id, type, message, state, group_id, event_id, created_at 
		FROM notifications 
		WHERE user_id = ? AND state = 'unread' 
		ORDER BY created_at DESC`

	rows, err := database.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []models.Notification
	for rows.Next() {
		var n models.Notification
		var groupID, eventID sql.NullInt64

		err := rows.Scan(&n.ID, &n.UserID, &n.SenderID, &n.Type, &n.Message, &n.State, &groupID, &eventID, &n.CreatedAt)
		if err != nil {
			log.Printf("Error scanning notification: %v", err)
			continue
		}

		if groupID.Valid {
			n.GroupID = int(groupID.Int64)
		}
		if eventID.Valid {
			n.EventID = int(eventID.Int64)
		}

		notifications = append(notifications, n)
	}
	return notifications, nil
}


func MarkAsSeenByContext(userID int, senderID int, groupID int, eventID int, msg string) error {
	query := `UPDATE notifications SET state = 'seen' WHERE user_id = ? AND state = 'unread'`
	args := []interface{}{userID}

	
	

	matched := false

	if msg == "Invitation_friendships" {
		query += ` AND type = 'Invitation_friendships' AND sender_id = ?`
		args = append(args, senderID)
		matched = true
	} else if msg == "group_invite" {
		query += ` AND type = 'group_invite' AND group_id = ?`
		args = append(args, groupID)
		matched = true
	} else if msg == "request_join_group" {
		query += ` AND type = 'request_join_group' AND group_id = ? AND sender_id = ?`
		args = append(args, groupID, senderID)
		matched = true
	} else if msg == "event_created" {
		query += ` AND type = 'event_created' AND event_id = ?`
		args = append(args, eventID)
		matched = true
	}

	if !matched {
		return nil 
	}

	_, err := database.DB.Exec(query, args...)
	return err
}

func UpdateNotificationStatus(id int, state string) error {
	query := `UPDATE notifications SET state = ? WHERE id = ?`
	_, err := database.DB.Exec(query, state, id)
	return err
}
