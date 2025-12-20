package repository

import (
	"database/sql"
	"fmt"
	"social-network/pkg/database"
	"social-network/pkg/models"
	sharedRepo "social-network/pkg/repository"
)

// Queries
const (
	INSERT_NEW_POST    = `INSERT INTO posts (user_id, title, content, image_path) VALUES (?, ?, ?, ?)`
	INSERT_NEW_COMMENT = `INSERT INTO comments (user_id, post_id, content, image_url) VALUES (?, ?, ?, ?)`
	SELECT_POST_BY_ID  = `
        SELECT p.id, p.user_id, u.username, u.first_name || ' ' || u.last_name AS full_name,
               p.title, p.content, image_path, p.created_at, p.updated_at
        FROM posts p
        JOIN users u ON p.user_id = u.id
        WHERE p.id = ?
	`
	SELECT_ALL_POSTS = `
        SELECT p.id, p.user_id, u.username, u.first_name || ' ' || u.last_name AS full_name,
               p.title, p.content, p.image_path, p.created_at, p.updated_at,
               (SELECT COUNT(*) FROM reactions r WHERE r.post_id = p.id AND r.reaction_type='like') AS likes_count,
               (SELECT COUNT(*) FROM reactions r WHERE r.post_id = p.id AND r.reaction_type='dislike') AS dislikes_count,
               (SELECT reaction_type FROM reactions r WHERE r.post_id = p.id AND r.user_id = ?) AS user_reaction
        FROM posts p
        JOIN users u ON p.user_id = u.id
        ORDER BY p.created_at DESC
    `
	SELECT_REACTION = "SELECT reaction_type FROM reactions WHERE user_id = ? AND post_id = ?"
	INSERT_REACTION = "INSERT INTO reactions (user_id, post_id, reaction_type) VALUES (?, ?, ?)"
	DELETE_REACTION = "DELETE FROM reactions WHERE user_id = ? AND post_id = ?"
	UPDATE_REACTION = "UPDATE reactions SET reaction_type = ? WHERE user_id = ? AND post_id = ?"
	COUNT_LIKES     = "SELECT COUNT(*) FROM reactions WHERE post_id = ? AND reaction_type = 'like'"
	COUNT_DISLIKES  = "SELECT COUNT(*) FROM reactions WHERE post_id = ? AND reaction_type = 'dislike'"
)

func AddNewPost(userId int, title string, content string, image string) (int, error) {
	res, err := database.DB.Exec(INSERT_NEW_POST, userId, title, content, image)
	if err != nil {
		return -1, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		fmt.Println("Warning: could not get LastInsertId:", err)
		return 0, nil
	}
	return int(id), nil
}

func GetAddedPost(id int) (*models.Posts, error) {
	row := database.DB.QueryRow(SELECT_POST_BY_ID, id)
	var post models.Posts
	err := row.Scan(&post.ID, &post.UserID, &post.UserName, &post.FullName,
		&post.Title, &post.Content, &post.ImagePath, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func GetAllPosts(userID int) ([]models.Posts, error) {
	rows, err := database.DB.Query(SELECT_ALL_POSTS, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []models.Posts
	for rows.Next() {
		var p models.Posts
		var userReaction sql.NullString

		if err := rows.Scan(
			&p.ID, &p.UserID, &p.UserName, &p.FullName,
			&p.Title, &p.Content, &p.ImagePath, &p.CreatedAt, &p.UpdatedAt,
			&p.LikesCount, &p.DislikesCount, &userReaction,
		); err != nil {
			return nil, err
		}
		if userReaction.Valid {
			p.UserReaction = &userReaction.String
		} else {
			p.UserReaction = nil
		}

		posts = append(posts, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

func GetUserPosts(userID int) ([]models.Posts, error) {
	query := `
        SELECT p.id, p.user_id, u.username, u.first_name || ' ' || u.last_name AS full_name,
               p.title, p.content, p.image_path, p.created_at, p.updated_at,
               (SELECT COUNT(*) FROM reactions r WHERE r.post_id = p.id AND r.reaction_type='like') AS likes_count,
               (SELECT COUNT(*) FROM reactions r WHERE r.post_id = p.id AND r.reaction_type='dislike') AS dislikes_count,
               (SELECT reaction_type FROM reactions r WHERE r.post_id = p.id AND r.user_id = ?) AS user_reaction
        FROM posts p
        JOIN users u ON p.user_id = u.id
        WHERE p.user_id = ?
        ORDER BY p.created_at DESC
    `

	rows, err := database.DB.Query(query, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Posts
	for rows.Next() {
		var p models.Posts
		var userReaction sql.NullString
		var imagePath sql.NullString

		if err := rows.Scan(
			&p.ID, &p.UserID, &p.UserName, &p.FullName,
			&p.Title, &p.Content, &imagePath, &p.CreatedAt, &p.UpdatedAt,
			&p.LikesCount, &p.DislikesCount, &userReaction,
		); err != nil {
			return nil, err
		}

		if userReaction.Valid {
			p.UserReaction = &userReaction.String
		}
		if imagePath.Valid {
			p.ImagePath = &imagePath.String
		}

		posts = append(posts, p)
	}

	return posts, nil
}

func AddComment(userID, postID int, content, imagePath string) (int, error) {
	res, err := database.DB.Exec(INSERT_NEW_COMMENT, userID, postID, content, imagePath)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

func GetCommentByID(id int) (map[string]interface{}, error) {
	var commentID, uid, pid int
	var commentText, createdAt string
	var imagePath sql.NullString
	err := database.DB.QueryRow(
		`SELECT id, user_id, post_id, content, created_at, image_url FROM comments WHERE id = ?`, id,
	).Scan(&commentID, &uid, &pid, &commentText, &createdAt, &imagePath)
	if err != nil {
		return nil, err
	}

	username, _ := sharedRepo.GetUsernameByID(uid)

	return map[string]interface{}{
		"id":         commentID,
		"user_id":    username, // Legacy behavior: return username in user_id field
		"post_id":    pid,
		"text":       commentText,
		"image_path": imagePath.String,
		"created_at": createdAt,
	}, nil
}

func GetCommentsByPostID(postID int) ([]map[string]interface{}, error) {
	rows, err := database.DB.Query(
		`SELECT id, user_id, post_id, content, created_at, image_url FROM comments WHERE post_id = ? ORDER BY created_at ASC`, postID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []map[string]interface{}
	for rows.Next() {
		var id, uid, pid int
		var text, createdAt string
		var imagePath sql.NullString
		if err := rows.Scan(&id, &uid, &pid, &text, &createdAt, &imagePath); err != nil {
			return nil, err
		}
		username, _ := sharedRepo.GetUsernameByID(uid)
		comments = append(comments, map[string]interface{}{
			"id":         id,
			"user_id":    username, // Legacy behavior
			"post_id":    pid,
			"text":       text,
			"image_path": imagePath.String,
			"created_at": createdAt,
		})
	}
	return comments, nil
}
