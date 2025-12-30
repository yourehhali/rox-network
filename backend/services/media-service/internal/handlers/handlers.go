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
	
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	
	path := strings.TrimPrefix(r.URL.Path, "/")
	filename := path
	if filename == "" {
		http.Error(w, "File not specified", http.StatusBadRequest)
		return
	}

	
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") {
		http.Error(w, "Invalid filename", http.StatusBadRequest)
		return
	}

	dbPath := "uploads/" + filename

	

	
	post, err := repository.GetPostByImage(dbPath)
	if err == nil {
		allowed := false
		switch post.Privacy {
		case 0: 
			allowed = true
		case 2: 
			if post.UserID == userID {
				allowed = true
			}
		case 1: 
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
		
	} else if err != sql.ErrNoRows {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	
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

	
	{
		u, err := repository.GetUserByAvatar(dbPath)
		if err == nil {
			_ = u 
			serveFile(w, r, filename)
			return
		} else if err != sql.ErrNoRows {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
	http.Error(w, "Not Found", http.StatusNotFound)
}
