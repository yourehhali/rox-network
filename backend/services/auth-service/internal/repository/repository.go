package repository

import (
	"database/sql"
	"errors"
	"social-network/pkg/database"
)

var (
	SELECT_USER_COUNT_BY_USERNAME_EMAIL      = `SELECT COUNT(*) FROM users WHERE username = ? OR email = ?`
	SELECT_USERID_PASSHASH_BY_USERNAME_EMAIL = `SELECT id,password_hash FROM users WHERE username = ? OR email = ?`
	INSERT_USERNAME_EMAIL_PASSHASH           = `
    INSERT INTO users (username, email, password_hash, first_name, last_name, age, gender, nickname, about, avatar, is_private)
    	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0)
	`
)

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
