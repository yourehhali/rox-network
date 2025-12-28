package handlers

import (
	"auth-service/internal/repository"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	sharedRepo "social-network/pkg/repository"
	"social-network/pkg/utils"
	"strconv"
	"strings"
	"time"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	exist, err := repository.AlreadyExists(req.Username, req.Username)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	if (!utils.ValidUsername(req.Username) && !utils.ValidEmail(req.Username)) || !utils.ValidPassword(req.Password) || !exist {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]any{
			"status":  "error",
			"message": "invalid credentials, try again",
		})
		return
	}
	userId, hash, err := repository.GetUserHashByUsername(req.Username)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	if !utils.CheckPassword(req.Password, hash) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]any{
			"status":  "error",
			"message": "wrong password, try again",
		})
		return
	}
	session := GenerateToken(32)
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    session,
		Expires:  time.Now().Add(time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/", // Important for cookie to be sent to all paths
	})
	err = sharedRepo.UpdateUserSession(userId, session)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]any{
		"status":  "ok",
		"user_id": userId,
	})
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil { // limit 10MB
		http.Error(w, `{"status":"error","message":"Cannot parse form"}`, http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	email := r.FormValue("email")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirmpassword")
	gender := r.FormValue("gender")
	nickname := r.FormValue("nickname")
	about := r.FormValue("about")
	age, _ := strconv.Atoi(r.FormValue("age"))

	// === Handle Avatar Upload ===
	file, handler, err := r.FormFile("avatar") 
	kb := 1024;
	mb := 1024 * kb;

	if (handler.Size > int64(15*mb)) {
		http.Error(w, `{"status":"error","message":"Invalid avatar SIZE MAX 15mb "}`, http.StatusBadRequest) 
		return
	} 
	var avatarPath string
	if err == nil {
		defer file.Close()

		if !utils.IsAllowedImageExtension(handler.Filename) {
			http.Error(w, `{"status":"error","message":"Invalid avatar file type. Allowed: .jpg, .jpeg, .png, .gif, .webp"}`, http.StatusBadRequest)
			return
		}

		// Note: In microservices, this should ideally be handled by Media Service.
		// For now, we save it to a shared location or local uploads.
		uploadDir := os.Getenv("UPLOAD_DIR") // Relative to where binary runs, or absolute path
		if uploadDir == "" {
			uploadDir = "../../uploads"
		}
		if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
			os.MkdirAll(uploadDir, 0755)
		}

		filename := fmt.Sprintf("%s_%s", username, filepath.Base(handler.Filename))
		avatarPath = filepath.Join("uploads", filename) // Path stored in DB
		fullPath := filepath.Join(uploadDir, filename)

		dst, err := os.Create(fullPath)
		if err != nil {
			log.Println("Error creating file:", err)
			http.Error(w, `{"status":"error","message":"Cannot save avatar"}`, http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		io.Copy(dst, file)
	} else {
		avatarPath = ""
	}

	// === Validation ===
	if !utils.ValidUsername(username) || !utils.ValidName(firstName) || !utils.ValidName(lastName) {
		http.Error(w, `{"status":"error","message":"Invalid username or name"}`, http.StatusBadRequest)
		return
	}
	if !utils.ValidEmail(email) {
		http.Error(w, `{"status":"error","message":"Invalid email"}`, http.StatusBadRequest)
		return
	}
	if !utils.ValidPassword(password) {
		http.Error(w, `{"status":"error","message":"Invalid password"}`, http.StatusBadRequest)
		return
	}
	if password != confirmPassword {
		http.Error(w, `{"status":"error","message":"Passwords do not match"}`, http.StatusBadRequest)
		return
	}

	hash, err := utils.HashPassword(password)
	if err != nil {
		http.Error(w, `{"status":"error","message":"Server error"}`, http.StatusInternalServerError)
		return
	}

	err = repository.AddNewUser(username, email, hash, firstName, lastName, gender, age, nickname, about, avatarPath)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			http.Error(w, `{"status":"error","message":"Username or email already used"}`, http.StatusConflict)
			return
		}
		log.Println("DB Error:", err)
		http.Error(w, `{"status":"error","message":"Server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status":  "ok",
		"message": "Registration successful!",
	})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Logged out successfully"}`))
}

func CheckState(w http.ResponseWriter, r *http.Request) {
	// CheckState ...
	session, err := r.Cookie("session")
	if err != nil {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"authenticated": false,
		})
		return
	}

	id, err := sharedRepo.GetUserIDBySessionToken(session.Value)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"authenticated": false,
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"authenticated": true,
		"userID":        id,
	})
}

func GenerateToken(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatalf("failed to generate token %v", err)
	}
	return base64.URLEncoding.EncodeToString(bytes)
}
