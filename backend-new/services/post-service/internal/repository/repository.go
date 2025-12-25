package repository

import (
	"database/sql"
	"log"
	"social-network/pkg/database"
	"social-network/pkg/models"
)

var (
	INSERT_NEW_POST    = `INSERT INTO posts (user_id, title, content, image_path, privacy) VALUES (?, ?, ?, ?, ?)`
	INSERT_POST_VIEWER = `INSERT INTO post_viewers (post_id, user_id) VALUES (?, ?)`
	INSERT_NEW_COMMENT = `INSERT INTO comments (user_id, post_id, content, image_url) VALUES (?, ?, ?, ?)`
	SELECT_POST_BY_ID  = `
        SELECT p.id, p.user_id, u.username, u.first_name || ' ' || u.last_name AS full_name,
               p.title, p.content, image_path, p.created_at, p.updated_at, p.privacy
        FROM posts p
        JOIN users u ON p.user_id = u.id
        WHERE p.id = ?
	`
	SELECT_ALL_POSTS = `
        SELECT DISTINCT p.id, p.user_id, u.username, u.first_name || ' ' || u.last_name AS full_name,
               p.title, p.content, p.image_path, p.created_at, p.updated_at,
               (SELECT COUNT(*) FROM reactions r WHERE r.post_id = p.id AND r.reaction_type='like') AS likes_count,
               (SELECT COUNT(*) FROM reactions r WHERE r.post_id = p.id AND r.reaction_type='dislike') AS dislikes_count,
               (SELECT reaction_type FROM reactions r WHERE r.post_id = p.id AND r.user_id = ?) AS user_reaction,
               p.privacy
        FROM posts p
        JOIN users u ON p.user_id = u.id
        LEFT JOIN followers f ON f.followed_id = p.user_id AND f.follower_id = ? AND f.status = 'accepted'
        LEFT JOIN post_viewers pv ON pv.post_id = p.id AND pv.user_id = ?
        WHERE 
           p.user_id = ? 
           OR (p.privacy = 0 AND (u.is_private = 0 OR f.id IS NOT NULL))
           OR (p.privacy = 1 AND f.id IS NOT NULL) 
           OR (p.privacy = 2 AND pv.id IS NOT NULL)
        ORDER BY p.created_at DESC
    `
	SELECT_USER_POSTS = `
        SELECT DISTINCT p.id, p.user_id, u.username, u.first_name || ' ' || u.last_name AS full_name,
               p.title, p.content, p.image_path, p.created_at, p.updated_at,
               (SELECT COUNT(*) FROM reactions r WHERE r.post_id = p.id AND r.reaction_type='like') AS likes_count,
               (SELECT COUNT(*) FROM reactions r WHERE r.post_id = p.id AND r.reaction_type='dislike') AS dislikes_count,
               (SELECT reaction_type FROM reactions r WHERE r.post_id = p.id AND r.user_id = ?) AS user_reaction,
               p.privacy
        FROM posts p
        JOIN users u ON p.user_id = u.id
        LEFT JOIN followers f ON f.followed_id = p.user_id AND f.follower_id = ? AND f.status = 'accepted'
        LEFT JOIN post_viewers pv ON pv.post_id = p.id AND pv.user_id = ?
        WHERE p.user_id = ?
          AND (
             p.user_id = ? 
             OR (p.privacy = 0 AND (u.is_private = 0 OR f.id IS NOT NULL))
             OR (p.privacy = 1 AND f.id IS NOT NULL) 
             OR (p.privacy = 2 AND pv.id IS NOT NULL)
          )
        ORDER BY p.created_at DESC
    `
	SELECT_COMMENT_BY_ID = `SELECT c.id, c.user_id, c.post_id, c.content, c.created_at, c.image_url, u.username, u.avatar 
		 FROM comments c 
		 JOIN users u ON c.user_id = u.id 
		 WHERE c.id = ?`
	SELECT_COMMENTS_BY_POST_ID = `SELECT c.id, c.user_id, c.post_id, c.content, c.created_at, c.image_url, u.username, u.avatar 
		 FROM comments c 
		 JOIN users u ON c.user_id = u.id 
		 WHERE c.post_id = ? 
		 ORDER BY c.created_at DESC`
	SELECT_REACTION = "SELECT reaction_type FROM reactions WHERE user_id = ? AND post_id = ?"
	INSERT_REACTION = "INSERT INTO reactions (user_id, post_id, reaction_type) VALUES (?, ?, ?)"
	DELETE_REACTION = "DELETE FROM reactions WHERE user_id = ? AND post_id = ?"
	UPDATE_REACTION = "UPDATE reactions SET reaction_type = ? WHERE user_id = ? AND post_id = ?"
	COUNT_LIKES     = "SELECT COUNT(*) FROM reactions WHERE post_id = ? AND reaction_type = 'like'"
	COUNT_DISLIKES  = "SELECT COUNT(*) FROM reactions WHERE post_id = ? AND reaction_type = 'dislike'"
)

func AddNewPost(userId int, title string, content string, image string, privacy int, allowedUsers []string) (int, error) {
	tx, err := database.DB.Begin()
	if err != nil {
		return -1, err
	}

	res, err := tx.Exec(INSERT_NEW_POST, userId, title, content, image, privacy)
	if err != nil {
		tx.Rollback()
		return -1, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		log.Println("Warning: could not get LastInsertId:", err)
		return 0, nil
	}

	if privacy == 2 && len(allowedUsers) > 0 {
		stmt, err := tx.Prepare(INSERT_POST_VIEWER)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		defer stmt.Close()
		for _, uidStr := range allowedUsers {
			if uidStr == "" {
				continue
			}
			_, err = stmt.Exec(id, uidStr)
			if err != nil {
				tx.Rollback()
				return 0, err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return -1, err
	}

	return int(id), nil
}

func GetAddedPost(id int) (*models.Posts, error) {
	row := database.DB.QueryRow(SELECT_POST_BY_ID, id)
	var post models.Posts
	err := row.Scan(&post.ID, &post.UserID, &post.UserName, &post.FullName,
		&post.Title, &post.Content, &post.ImagePath, &post.CreatedAt, &post.UpdatedAt, &post.Privacy)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func GetAllPosts(viewerID int) ([]models.Posts, error) {
	rows, err := database.DB.Query(SELECT_ALL_POSTS, viewerID, viewerID, viewerID, viewerID)
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
			&p.LikesCount, &p.DislikesCount, &userReaction, &p.Privacy,
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

func GetUserPosts(targetUserID, viewerID int) ([]models.Posts, error) {
	rows, err := database.DB.Query(SELECT_USER_POSTS, viewerID, viewerID, viewerID, targetUserID, viewerID)
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
			&p.LikesCount, &p.DislikesCount, &userReaction, &p.Privacy,
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
	var commentText, createdAt, username string
	var imagePath, avatar sql.NullString
	err := database.DB.QueryRow(SELECT_COMMENT_BY_ID, id).Scan(&commentID, &uid, &pid, &commentText, &createdAt, &imagePath, &username, &avatar)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id":         commentID,
		"user_id":    username,
		"username":   username,
		"avatar":     avatar.String,
		"post_id":    pid,
		"text":       commentText,
		"image_path": imagePath.String,
		"created_at": createdAt,
	}, nil
}

func GetCommentsByPostID(postID int) ([]map[string]interface{}, error) {
	rows, err := database.DB.Query(SELECT_COMMENTS_BY_POST_ID, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []map[string]interface{}
	for rows.Next() {
		var id, uid, pid int
		var text, createdAt, username string
		var imagePath, avatar sql.NullString
		if err := rows.Scan(&id, &uid, &pid, &text, &createdAt, &imagePath, &username, &avatar); err != nil {
			return nil, err
		}
		comments = append(comments, map[string]interface{}{
			"id":         id,
			"user_id":    username,
			"username":   username,
			"avatar":     avatar.String,
			"post_id":    pid,
			"text":       text,
			"image_path": imagePath.String,
			"created_at": createdAt,
		})
	}
	return comments, nil
}
