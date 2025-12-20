package repository

import (
	"database/sql"
	"social-network/pkg/database"
	"social-network/pkg/models"
)

// Queries
const (
	SELECT_USER_BY_ID = `
    SELECT id, username, email, first_name, last_name, age, gender, 
           nickname, about, avatar, created_at, updated_at ,is_private
    FROM users 
    WHERE id = ?
	`
	SELECT_FOLLOWERS_COUNT = `SELECT COUNT(*) FROM followers WHERE followed_id = ? AND status = 'accepted'`
	SELECT_FOLLOWING_COUNT = `SELECT COUNT(*) FROM followers WHERE follower_id = ? AND status = 'accepted'`
	IS_FOLLOWING           = `SELECT COUNT(*) FROM followers WHERE follower_id = ? AND followed_id = ? AND status = 'accepted'`
	CHECK_PENDING          = `SELECT COUNT(*) FROM followers WHERE follower_id = ? AND followed_id = ? AND status = 'pending'` // assuming 'howa' or 'pending'
	INSERT_FOLLOW          = `INSERT INTO followers (follower_id, followed_id, status) VALUES (?, ?, ?)`
	DELETE_FOLLOW          = `DELETE FROM followers WHERE follower_id = ? AND followed_id = ?`
	UPDATE_FOLLOW_STATUS   = `UPDATE followers SET status = ? WHERE follower_id = ? AND followed_id = ?`
	CHECK_PRIVACY          = `SELECT is_private FROM users WHERE id = ?`
	// Suggestion queries, etc. can be added here
)

func GetUserInfo(userId int) (models.User, error) {
	var user models.User
	var age sql.NullInt64
	var gender sql.NullString
	var nickname sql.NullString
	var about sql.NullString
	var avatar sql.NullString
	var isPrivate sql.NullInt64

	err := database.DB.QueryRow(SELECT_USER_BY_ID, userId).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&age,
		&gender,
		&nickname,
		&about,
		&avatar,
		&user.CreatedAt,
		&user.UpdatedAt,
		&isPrivate,
	)

	if err != nil {
		return user, err
	}

	if age.Valid {
		a := int(age.Int64)
		user.Age = &a
	}
	if gender.Valid {
		user.Gender = &gender.String
	}
	if nickname.Valid {
		user.Nickname = &nickname.String
	}
	if about.Valid {
		user.About = &about.String
	}
	if avatar.Valid {
		user.Avatar = &avatar.String
	}
	if isPrivate.Valid {
		p := int(isPrivate.Int64)
		user.IsPrivate = &p
	}
	return user, nil
}

func GetFollowersCount(userID int) (int, error) {
	var count int
	err := database.DB.QueryRow(SELECT_FOLLOWERS_COUNT, userID).Scan(&count)
	return count, err
}

func GetFollowingCount(userID int) (int, error) {
	var count int
	err := database.DB.QueryRow(SELECT_FOLLOWING_COUNT, userID).Scan(&count)
	return count, err
}

func IsFollowing(followerID, followedID int) (bool, error) {
	var count int
	err := database.DB.QueryRow(IS_FOLLOWING, followerID, followedID).Scan(&count)
	return count > 0, err
}

func IsPending(followerID, followedID int) (bool, error) {
	var count int
	err := database.DB.QueryRow(CHECK_PENDING, followerID, followedID).Scan(&count)
	return count > 0, err
}

func ToggleFollow(followerID, followedID int) (string, error) {
	// Logic to check if following, if so unfollow, else follow (check privacy)
	// This is simplified, actual logic in handler calls repository primitive ops
	return "", nil
}

// AddFollow inserts a follow request
func AddFollow(followerID, followedID int, status string) error {
	_, err := database.DB.Exec(INSERT_FOLLOW, followerID, followedID, status)
	return err
}

// RemoveFollow deletes a follow relationship
func RemoveFollow(followerID, followedID int) error {
	_, err := database.DB.Exec(DELETE_FOLLOW, followerID, followedID)
	return err
}

func IsUserPrivate(userID int) (bool, error) {
	var isPrivate sql.NullInt64
	err := database.DB.QueryRow(CHECK_PRIVACY, userID).Scan(&isPrivate)
	if err != nil {
		return false, err
	}
	if !isPrivate.Valid {
		return false, nil
	}
	return isPrivate.Int64 == 1, nil
}

// UpdateFollowStatus sets status for an existing follow row
func UpdateFollowStatus(followerID, followedID int, status string) error {
	_, err := database.DB.Exec(UPDATE_FOLLOW_STATUS, status, followerID, followedID)
	return err
}

type Friend struct {
	ID        int     `json:"id"`
	Username  string  `json:"username"`
	FullName  string  `json:"full_name"`
	ImagePath *string `json:"image_path"`
}

func GetFriendList(userID int) ([]Friend, error) {
	query := `
        SELECT DISTINCT u.id, u.username, u.first_name || ' ' || ifnull(u.last_name, '') AS full_name, u.avatar
        FROM users u
        JOIN followers f ON (
            (f.follower_id = ? AND f.followed_id = u.id)
            OR
            (f.followed_id = ? AND f.follower_id = u.id)
        )
        WHERE f.status = 'accepted';
    `
	rows, err := database.DB.Query(query, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var friends []Friend
	for rows.Next() {
		var f Friend
		if err := rows.Scan(&f.ID, &f.Username, &f.FullName, &f.ImagePath); err != nil {
			return nil, err
		}
		friends = append(friends, f)
	}
	return friends, nil
}

func GetAllUsers(currentUserID int) ([]models.Sugg, error) {
	query := `
		SELECT id, username, first_name, last_name, avatar
		FROM users
	`
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.Sugg
	for rows.Next() {
		var u models.Sugg
		var firstName, lastName sql.NullString
		var imagePath sql.NullString

		err := rows.Scan(&u.UserID, &u.UserName, &firstName, &lastName, &imagePath)
		if err != nil {
			return nil, err
		}

		// Check if following or self
		isFollowing, _ := IsFollowing(currentUserID, u.UserID)
		if currentUserID == u.UserID || isFollowing {
			continue
		}

		fullName := ""
		if firstName.Valid {
			fullName += firstName.String
		}
		if lastName.Valid {
			if fullName != "" {
				fullName += " "
			}
			fullName += lastName.String
		}
		u.FullName = fullName

		if imagePath.Valid {
			u.ImagePath = &imagePath.String
		} else {
			u.ImagePath = nil
		}
		users = append(users, u)
	}
	return users, nil
}

func UpdatePrivacy(userID int, isPrivate int) error {
	_, err := database.DB.Exec(`UPDATE users SET is_private = ? WHERE id = ?`, isPrivate, userID)
	return err
}
