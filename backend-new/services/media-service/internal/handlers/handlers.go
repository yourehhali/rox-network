package handlers

import (
	"database/sql"
	"net/http"
	"os"
	"path/filepath"
	"social-network/pkg/middleware"
	"social-network/services/media-service/internal/repository"
	"strings"
)

func serveFile(w http.ResponseWriter, r *http.Request, filename string) {
	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "../../uploads"
	}
	fullPath := filepath.Join(uploadDir, filename)
	http.ServeFile(w, r, fullPath)
}

func ServeMedia(w http.ResponseWriter, r *http.Request) {
	// 1. Authenticate (handled by middleware, but we need ID)
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 2. Extract Filename
	path := strings.TrimPrefix(r.URL.Path, "/")
	filename := path
	if filename == "" {
		http.Error(w, "File not specified", http.StatusBadRequest)
		return
	}

	// Security: prevent directory traversal
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") {
		http.Error(w, "Invalid filename", http.StatusBadRequest)
		return
	}

	dbPath := "uploads/" + filename

	// 3. Check Contexts

	// A. Check Posts
	post, err := repository.GetPostByImage(dbPath)
	if err == nil {
		allowed := false
		switch post.Privacy {
		case 0: // Public
			allowed = true
		case 2: // Private
			if post.UserID == userID {
				allowed = true
			}
		case 1: // Almost Private
			if post.UserID == userID {
				allowed = true
			} else {
				isViewer, _ := repository.IsPostViewer(post.ID, userID)
				if isViewer {
					allowed = true
				}
			}
		}

		if allowed {
			serveFile(w, r, filename)
			return
		}
		// If not allowed, continue to check other contexts (e.g. might also be an avatar)
	} else if err != sql.ErrNoRows {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// B. Check Group Posts
	{
		gp, err := repository.GetGroupPostByImage(dbPath)
		if err == nil {
			allowed := false
			isMember, _ := repository.IsGroupMember(gp.GroupID, userID)
			if isMember {
				allowed = true
			}

			if allowed {
				serveFile(w, r, filename)
				return
			}
		} else if err != sql.ErrNoRows {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	// C. Check Users (Avatars)
	{
		u, err := repository.GetUserByAvatar(dbPath)
		if err == nil {
			_ = u // Unused
			serveFile(w, r, filename)
			return
		} else if err != sql.ErrNoRows {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
	http.Error(w, "Not Found", http.StatusNotFound)
}
