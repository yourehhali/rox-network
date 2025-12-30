package repository

import (
	"log"
	"social-network/pkg/database"
)


func GetUserIDBySessionToken(token string) (int, error) {
	var userID int
	err := database.DB.QueryRow("SELECT user_id FROM sessions WHERE session_token = ?", token).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}


func UpdateUserSession(id int, token string) error {
	tx, err := database.DB.Begin()
	if err != nil {
		log.Println("Failed to start transaction:", err)
		return err
	}

	
	res, err := tx.Exec(`UPDATE sessions SET expires_at = DATETIME('now', '+1 hour'), session_token = ? WHERE user_id = ?`, token, id)
	if err != nil {
		log.Println("Update error:", err)
		tx.Rollback()
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		log.Println("Error getting affected rows:", err)
		tx.Rollback()
		return err
	}

	if count > 0 {
		return tx.Commit()
	}

	
	_, err = tx.Exec(`INSERT INTO sessions (user_id, session_token, expires_at) VALUES (?, ?, DATETIME('now', '+1 hour'))`, id, token)
	if err != nil {
		log.Println("Insert error:", err)
		tx.Rollback()
		return err
	}
	return tx.Commit()
}


func ResetUserSession(session_id string) error {
	_, err := database.DB.Exec(`UPDATE sessions SET session_token = NULL WHERE session_token = ?`, session_id)
	return err
}
