package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"social-network/pkg/middleware"
	"social-network/pkg/models"
	"strconv"
	"time"
	"user-service/internal/repository"
)

func Profile(w http.ResponseWriter, r *http.Request) {
	// Current user ID from context (viewer)
	viewerID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		// handle missing context if needed, but middleware should ensure it
	}

	userIDStr := r.URL.Query().Get("id")
	if userIDStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Missing user ID"})
		return
	}

	targetUserID, err := strconv.Atoi(userIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid user ID"})
		return
	}

	userData, err := repository.GetUserInfo(targetUserID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Internal server error"})
		return
	}

	followersCount, err := repository.GetFollowersCount(targetUserID)
	if err != nil {
		http.Error(w, "Failed to get followers count", http.StatusInternalServerError)
		return
	}
	followingCount, err := repository.GetFollowingCount(targetUserID)
	if err != nil {
		http.Error(w, "Failed to get following count", http.StatusInternalServerError)
		return
	}

	isPending, _ := repository.IsPending(viewerID, targetUserID)
	isFollowing, _ := repository.IsFollowing(viewerID, targetUserID)

	response := struct {
		User           models.User `json:"user"`
		FollowersCount int         `json:"followers_count"`
		FollowingCount int         `json:"following_count"`
		IsFriend       bool        `json:"isfriend"`
		IsPending      bool        `json:"p"` // matching frontend expectation "p"
	}{
		User:           userData,
		FollowersCount: followersCount,
		FollowingCount: followingCount,
		IsFriend:       isFollowing,
		IsPending:      isPending,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

type FollowRequest struct {
	FollowedID int `json:"followed_id"`
}

type FollowResponse struct {
	Following      bool   `json:"following"`
	Status         string `json:"status"`
	FollowersCount int    `json:"followers_Count"`
	FollowingCount int    `json:"following_Count"` // Note: Frontend might expect this casing
}

func SetFollowers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("DEBUG: SetFollowers handler called")
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		fmt.Println("DEBUG: Unauthorized access in SetFollowers")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req FollowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Printf("DEBUG: Invalid request body in SetFollowers: %v\n", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	fmt.Printf("DEBUG: User %d toggling follow for %d\n", userID, req.FollowedID)

	var status string
	isPending, _ := repository.IsPending(userID, req.FollowedID)

	if isPending {
		// Cancel request
		repository.RemoveFollow(userID, req.FollowedID)
		status = "cancelled"
		// Remove the pending request notification from the user
		_ = sendInternalNotification(models.Notification{
			UserID:   req.FollowedID,
			Type:     "notification_removed",
			Message:  "Invitation_friendships",
			SenderID: userID,
		})
	}

	isFollowing, err := repository.IsFollowing(userID, req.FollowedID)
	if err != nil {
		http.Error(w, "error checking follow status", http.StatusInternalServerError)
		return
	}

	if isFollowing {
		repository.RemoveFollow(userID, req.FollowedID)
		status = "unfollowed"
		// Remove the "new_follower" notification
		_ = sendInternalNotification(models.Notification{
			UserID:   req.FollowedID,
			Type:     "notification_removed",
			Message:  "started following you",
			SenderID: userID,
		})
	} else if !isPending {
		// Start following
		isPrivate, err := repository.IsUserPrivate(req.FollowedID)
		if err != nil {
			fmt.Printf("IsUserPrivate error for user %d: %v\n", req.FollowedID, err)
			http.Error(w, "error checking privacy", http.StatusInternalServerError)
			return
		}
		status = "accepted"
		if isPrivate {
			status = "pending"
		}
		err = repository.AddFollow(userID, req.FollowedID, status)
		if err != nil {
			http.Error(w, "error adding follow", http.StatusInternalServerError)
			return
		}
		// Fire notification to followed user if request is pending
		if status == "pending" {
			_ = sendInternalNotification(models.Notification{
				UserID:   req.FollowedID,
				SenderID: userID,
				Type:     "Invitation_friendships",
				Message:  "request",
			})
		} else if status == "accepted" {
			_ = sendInternalNotification(models.Notification{
				UserID:   req.FollowedID,
				SenderID: userID,
				Type:     "new_follower",
				Message:  "started following you",
			})
		}
	}

	followersCount, _ := repository.GetFollowersCount(req.FollowedID)
	followingCount, _ := repository.GetFollowingCount(req.FollowedID)

	// Re-check following status for response
	newIsFollowing, _ := repository.IsFollowing(userID, req.FollowedID)

	res := FollowResponse{
		Following:      newIsFollowing,
		Status:         status,
		FollowersCount: followersCount,
		FollowingCount: followingCount,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// RequestFollow handles accept/reject of follow requests by the receiver
func RequestFollow(w http.ResponseWriter, r *http.Request) {
	receiverID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var input struct {
		SenderID int    `json:"sender_id"`
		Status   string `json:"status"` // "accept" | "reject"
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if input.SenderID == 0 || (input.Status != "accept" && input.Status != "reject") {
		http.Error(w, "Invalid parameters", http.StatusBadRequest)
		return
	}

	// Validate there is a pending request
	isPending, _ := repository.IsPending(input.SenderID, receiverID)
	if !isPending {
		http.Error(w, "No pending request", http.StatusBadRequest)
		return
	}

	var err error
	if input.Status == "accept" {
		err = repository.UpdateFollowStatus(input.SenderID, receiverID, "accepted")
		if err == nil {
			_ = sendInternalNotification(models.Notification{
				UserID:   input.SenderID,
				SenderID: receiverID,
				Type:     "Invitation_friendships",
				Message:  "accepted",
			})
		}
	} else { // reject
		err = repository.RemoveFollow(input.SenderID, receiverID)
		if err == nil {
			_ = sendInternalNotification(models.Notification{
				UserID:   input.SenderID,
				SenderID: receiverID,
				Type:     "Invitation_friendships",
				Message:  "rejected",
			})
		}
	}
	if err != nil {
		http.Error(w, "Operation failed", http.StatusInternalServerError)
		return
	}

	// Remove the original request notification for the Receiver (me)
	_ = sendInternalNotification(models.Notification{
		UserID:   receiverID,
		Type:     "notification_removed",
		Message:  "Invitation_friendships",
		SenderID: input.SenderID,
	})

	// Return updated counts for receiver (optional)
	followersCount, _ := repository.GetFollowersCount(receiverID)
	followingCount, _ := repository.GetFollowingCount(receiverID)
	json.NewEncoder(w).Encode(map[string]any{
		"status":          input.Status,
		"followers_Count": followersCount,
		"following_Count": followingCount,
	})
}

// sendInternalNotification posts to chat-service internal notification endpoint
func sendInternalNotification(n models.Notification) error {
	fmt.Printf("DEBUG: Sending notification: %+v\n", n)
	b, _ := json.Marshal(n)
	// Direct call to chat-service internal endpoint
	resp, err := http.Post("http://localhost:8085/internal/notification", "application/json", bytes.NewReader(b))
	if err != nil {
		fmt.Printf("DEBUG: Error sending notification: %v\n", err)
		return err
	}
	defer resp.Body.Close()
	fmt.Println("DEBUG: Notification sent successfully")
	return nil
}

func GetFriendList(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	friends, err := repository.GetFriendList(userID)
	if err != nil {
		http.Error(w, "Query error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(friends)
}

func GetSuggestedUsers(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	users, err := repository.GetAllUsers(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

type PrivacyRequest struct {
	IsPrivate int `json:"is_private"`
}

type PrivacyResponse struct {
	IsPrivate int    `json:"is_private"`
	Message   string `json:"message"`
}

func UpdatePrivacy(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req PrivacyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.IsPrivate != 0 && req.IsPrivate != 1 {
		http.Error(w, "Invalid value for is_private", http.StatusBadRequest)
		return
	}

	err := repository.UpdatePrivacy(userID, req.IsPrivate)
	if err != nil {
		fmt.Println("DB update error:", err)
		http.Error(w, "Failed to update privacy", http.StatusInternalServerError)
		return
	}

	res := PrivacyResponse{
		IsPrivate: req.IsPrivate,
		Message:   "Successfully Updated Privacy",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	var u models.User
	u.FirstName = r.FormValue("first_name")
	u.LastName = r.FormValue("last_name")

	nick := r.FormValue("nickname")
	u.Nickname = &nick

	about := r.FormValue("about")
	u.About = &about

	isPrivateStr := r.FormValue("is_private")
	isPrivate, _ := strconv.Atoi(isPrivateStr)
	u.IsPrivate = &isPrivate

	// Handle avatar upload
	file, handler, err := r.FormFile("avatar")
	if err == nil {
		defer file.Close()
		uploadDir := "../../uploads"
		if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
			os.MkdirAll(uploadDir, 0755)
		}
		filename := fmt.Sprintf("%d_%s", time.Now().Unix(), handler.Filename)
		avatarPath := "uploads/" + filename
		fullPath := uploadDir + "/" + filename

		f, err := os.Create(fullPath)
		if err != nil {
			http.Error(w, "Cannot save image", http.StatusInternalServerError)
			return
		}
		defer f.Close()
		io.Copy(f, file)
		u.Avatar = &avatarPath
	} else {
		// Preserve existing avatar if not updated
		currentUser, err := repository.GetUserInfo(userID)
		if err == nil {
			u.Avatar = currentUser.Avatar
		}
	}

	err = repository.UpdateUserInfo(userID, u)
	if err != nil {
		fmt.Printf("UpdateUserInfo error: %v\n", err)
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Profile updated successfully"})
}

func FollowersList(w http.ResponseWriter, r *http.Request) {
	targetIDStr := r.URL.Query().Get("id")
	if targetIDStr == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}
	targetID, err := strconv.Atoi(targetIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	followers, err := repository.GetFollowers(targetID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(followers)
}

func FollowingList(w http.ResponseWriter, r *http.Request) {
	targetIDStr := r.URL.Query().Get("id")
	if targetIDStr == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}
	targetID, err := strconv.Atoi(targetIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	following, err := repository.GetFollowing(targetID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(following)
}
