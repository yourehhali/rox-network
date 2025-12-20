package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"post-service/internal/repository"
	"social-network/pkg/database"
	"social-network/pkg/middleware"
	"strconv"
	"strings"
	"time"
)

func CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err := r.ParseMultipartForm(10 << 20) // 10MB max size
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")

	var imagePath string
	file, handler, err := r.FormFile("image")
	if err == nil {
		defer file.Close()
		// Using relative path for migration simplicity.
		// Ensure uploads directory exists.
		uploadDir := "../../uploads"
		if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
			os.MkdirAll(uploadDir, 0755)
		}

		filename := fmt.Sprintf("%d_%s", time.Now().Unix(), handler.Filename)
		imagePath = "uploads/" + filename // Path stored in DB
		fullPath := uploadDir + "/" + filename

		f, err := os.Create(fullPath)
		if err != nil {
			fmt.Println("Error creating file:", err)
			http.Error(w, "Cannot save image", http.StatusInternalServerError)
			return
		}
		defer f.Close()
		io.Copy(f, file)
	}

	ID, err := repository.AddNewPost(userID, title, content, imagePath)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		fmt.Println("There is a problem in the database:", err)
		return
	}

	post, _ := repository.GetAddedPost(ID)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(post); err != nil {
		fmt.Println("Error encoding post:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func GetPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	posts, err := repository.GetAllPosts(userID)

	if err != nil {
		fmt.Println("There is a problem in get posts:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(posts); err != nil {
		fmt.Println("Error encoding posts:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func GetPostsByUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userIDStr := r.URL.Query().Get("id")
	if userIDStr == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	posts, err := repository.GetUserPosts(userID)
	if err != nil {
		fmt.Println("Error getting user posts:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(posts); err != nil {
		fmt.Println("Error encoding posts:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func Reaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var input struct {
		PostID   int    `json:"post_id"`
		Reaction string `json:"reaction"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	var currentReaction string
	err := database.DB.QueryRow(
		repository.SELECT_REACTION,
		userID, input.PostID,
	).Scan(&currentReaction)

	if err != nil {
		if err == sql.ErrNoRows {
			currentReaction = ""
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
	}

	if currentReaction == "" {
		_, err = database.DB.Exec(repository.INSERT_REACTION, userID, input.PostID, input.Reaction)
	} else if currentReaction == input.Reaction {
		_, err = database.DB.Exec(repository.DELETE_REACTION, userID, input.PostID)
	} else {
		_, err = database.DB.Exec(repository.UPDATE_REACTION, input.Reaction, userID, input.PostID)
	}
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	var likes, dislikes int
	_ = database.DB.QueryRow(repository.COUNT_LIKES, input.PostID).Scan(&likes)
	_ = database.DB.QueryRow(repository.COUNT_DISLIKES, input.PostID).Scan(&dislikes)

	var userReaction sql.NullString
	_ = database.DB.QueryRow(repository.SELECT_REACTION, userID, input.PostID).Scan(&userReaction)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"likes_count":    likes,
		"dislikes_count": dislikes,
		"userReaction":   userReaction.String,
	})
}

func SubmitComment(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	// Expected /posts/{id}/comments
	if len(pathParts) < 3 || pathParts[2] != "comments" {
		http.Error(w, "Invalid route", http.StatusNotFound)
		return
	}

	postID, err := strconv.Atoi(pathParts[1])
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodPost:
		userID, ok := r.Context().Value(middleware.UserIDKey).(int)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Handle Multipart Form (consistent with old subcomment.go)
		err := r.ParseMultipartForm(10 << 20) // 10MB
		var content string
		var imagePath string

		if err == nil {
			content = r.FormValue("content")
			file, handler, err := r.FormFile("image")
			if err == nil {
				defer file.Close()
				uploadDir := "../../uploads"
				if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
					os.MkdirAll(uploadDir, 0755)
				}
				filename := fmt.Sprintf("comment_%d_%s", time.Now().Unix(), handler.Filename)
				imagePath = "uploads/" + filename
				fullPath := uploadDir + "/" + filename
				
				f, err := os.Create(fullPath)
				if err == nil {
					defer f.Close()
					io.Copy(f, file)
				}
			}
		} else {
			// Fallback to JSON if multipart fails or is not used (for backward compatibility if needed)
			// But since we want to support images, we prefer multipart.
			// If ParseMultipartForm failed, it might be because it's JSON.
			// Let's check Content-Type.
			if strings.Contains(r.Header.Get("Content-Type"), "application/json") {
				var input struct {
					Content string `json:"content"`
				}
				if err := json.NewDecoder(r.Body).Decode(&input); err == nil {
					content = input.Content
				}
			}
		}

		if content == "" && imagePath == "" {
			http.Error(w, "Empty comment", http.StatusBadRequest)
			return
		}

		id, err := repository.AddComment(userID, postID, content, imagePath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		comment, err := repository.GetCommentByID(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(comment)

	case http.MethodGet:
		comments, err := repository.GetCommentsByPostID(postID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if comments == nil {
			comments = []map[string]interface{}{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(comments)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
