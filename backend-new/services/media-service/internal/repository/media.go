package repository

import (
	"database/sql"
	"social-network/pkg/database"
)

type PostMeta struct {
	ID      int
	UserID  int
	Privacy int
}

type GroupPostMeta struct {
	GroupID int
}

type UserMeta struct {
	ID        int
	IsPrivate int
}

// GetPostByImage finds a post by its image path
func GetPostByImage(path string) (*PostMeta, error) {
	var p PostMeta
	err := database.DB.QueryRow("SELECT id, user_id, privacy FROM posts WHERE image_path = ?", path).Scan(&p.ID, &p.UserID, &p.Privacy)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// GetGroupPostByImage finds a group post by its image path
func GetGroupPostByImage(path string) (*GroupPostMeta, error) {
	var gp GroupPostMeta
	// Note: group_posts table has 'image' column, not 'image_path' based on migration 000009
	err := database.DB.QueryRow("SELECT group_id FROM group_posts WHERE image = ?", path).Scan(&gp.GroupID)
	if err != nil {
		return nil, err
	}
	return &gp, nil
}

// GetUserByAvatar finds a user by their avatar path
func GetUserByAvatar(path string) (*UserMeta, error) {
	var u UserMeta
	err := database.DB.QueryRow("SELECT id, is_private FROM users WHERE avatar = ?", path).Scan(&u.ID, &u.IsPrivate)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// IsPostViewer checks if a user is in the post_viewers table for a specific post
func IsPostViewer(postID, userID int) (bool, error) {
	var exists int
	err := database.DB.QueryRow("SELECT 1 FROM post_viewers WHERE post_id = ? AND user_id = ?", postID, userID).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// IsGroupMember checks if a user is a member of a group
func IsGroupMember(groupID, userID int) (bool, error) {
	var exists int
	// Check if member and status is 'member'
	err := database.DB.QueryRow("SELECT 1 FROM group_members WHERE group_id = ? AND user_id = ? AND status = 'member'", groupID, userID).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
