package repository

import (
	"social-network/pkg/database"
	"social-network/pkg/models"
)


func GetBasicUserInfo(userID int) (*models.User, error) {
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


func GetUsernameByID(userID int) (string, error) {
	var username string
	err := database.DB.QueryRow("SELECT username FROM users WHERE id = ?", userID).Scan(&username)
	if err != nil {
		return "", err
	}
	return username, nil
}
