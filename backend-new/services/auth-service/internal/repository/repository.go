package repository

import (
	"database/sql"
	"errors"
	"log"
	"social-network/pkg/database"
)

// Queries
const (
	SELECT_USER_BY_SESSION_TOKEN             = `SELECT user_id FROM sessions WHERE session_token = ?`
	UPDATE_SESSION_EXPIRING_TIME             = `UPDATE sessions SET expires_at = DATETIME('now', '+1 hour'), session_token = ? WHERE user_id = ?`
	INSERT_NEW_SESSION                       = `INSERT INTO sessions (user_id, session_token, expires_at) VALUES (?, ?, DATETIME('now', '+1 hour'))`
	RESET_USER_SESSION_TOKEN                 = `UPDATE sessions SET session_token = NULL WHERE session_token = ?`
	SELECT_USER_COUNT_BY_USERNAME_EMAIL      = `SELECT COUNT(*) FROM users WHERE username = ? OR email = ?`
	SELECT_USERID_PASSHASH_BY_USERNAME_EMAIL = `SELECT id,password_hash FROM users WHERE username = ? OR email = ?`
	INSERT_USERNAME_EMAIL_PASSHASH           = `
    INSERT INTO users (username, email, password_hash, first_name, last_name, age, gender, nickname, about, avatar, is_private)
    	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0)
	`
)

func SelectUserSession(session_id string) (int, bool, error) {
	var userId int
	err := database.DB.QueryRow(SELECT_USER_BY_SESSION_TOKEN, session_id).Scan(&userId)
	if err != nil {
		return userId, false, err
	}
	return userId, true, nil
}

func UpdateUserSession(id int, token string) error {
	tx, err := database.DB.Begin()
	if err != nil {
		log.Println("Failed to start transaction:", err)
		return err
	}
	res, err := tx.Exec(UPDATE_SESSION_EXPIRING_TIME, token, id)
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
	_, err = tx.Exec(INSERT_NEW_SESSION, id, token)
	if err != nil {
		log.Println("Insert error:", err)
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func ResetUserSession(session_id string) (bool, error) {
	_, err := database.DB.Exec(RESET_USER_SESSION_TOKEN, session_id)
	if err != nil {
		return false, err
	}
	return true, nil
}

func AlreadyExists(username, email string) (bool, error) {
	var count int
	err := database.DB.QueryRow(SELECT_USER_COUNT_BY_USERNAME_EMAIL, username, email).Scan(&count)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return count > 0, nil
}

func GetUserHashByUsername(username string) (int, string, error) {
	var hash string
	var id int
	err := database.DB.QueryRow(SELECT_USERID_PASSHASH_BY_USERNAME_EMAIL, username, username).Scan(&id, &hash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return id, hash, nil
		}
		return id, hash, err
	}
	return id, hash, nil
}

func AddNewUser(username, email, hashedPass, firstName, lastName, gender string, age int, nickname, about, avatar string) error {
	_, err := database.DB.Exec(
		INSERT_USERNAME_EMAIL_PASSHASH,
		username, email, hashedPass, firstName, lastName, age, gender, nickname, about, avatar,
	)
	return err
}
