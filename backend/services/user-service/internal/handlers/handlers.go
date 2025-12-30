package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"social-network/pkg/middleware"
	"social-network/pkg/models"
	"social-network/services/user-service/internal/repository"
	"strconv"
)

func Profile(w http.ResponseWriter, r *http.Request) {
	
	viewerID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Internal server error"})
		return
	}

	var targetUserID int
	userIDStr := r.URL.Query().Get("id")
	if userIDStr == "" {
		targetUserID = viewerID
	} else {
		var err error
		targetUserID, err = strconv.Atoi(userIDStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"message": "Invalid user ID"})
			return
		}
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
	isMyFollower, _ := repository.IsFollowing(targetUserID, viewerID)
	postsCount, _ := repository.GetUserPostsCount(targetUserID)

	response := struct {
		User           models.User `json:"user"`
		FollowersCount int         `json:"followers_count"`
		FollowingCount int         `json:"following_count"`
		IsFriend       bool        `json:"isfriend"`
		IsPending      bool        `json:"p"` 
		HasPosts       bool        `json:"has_posts"`
		IsMyFollower   bool        `json:"ismyfollower"`
	}{
		User:           userData,
		FollowersCount: followersCount,
		FollowingCount: followingCount,
		IsFriend:       isFollowing,
		IsPending:      isPending,
		HasPosts:       postsCount > 0,
		IsMyFollower:   isMyFollower,
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
	FollowingCount int    `json:"following_Count"`
}

func SetFollowers(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req FollowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var status string
	isPending, _ := repository.IsPending(userID, req.FollowedID)

	if isPending {
		
		repository.RemoveFollow(userID, req.FollowedID)
		status = "cancelled"
		
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
		
		_ = sendInternalNotification(models.Notification{
			UserID:   req.FollowedID,
			Type:     "notification_removed",
			Message:  "started following you",
			SenderID: userID,
		})
	} else if !isPending {
		
		isPrivate, err := repository.IsUserPrivate(req.FollowedID)
		if err != nil {
			log.Printf("IsUserPrivate error for user %d: %v\n", req.FollowedID, err)
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


func RequestFollow(w http.ResponseWriter, r *http.Request) {
	receiverID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var input struct {
		SenderID int    `json:"sender_id"`
		NotifID  int    `json:"notif_id"`
		Status   string `json:"status"` 
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if input.SenderID == 0 || (input.Status != "accept" && input.Status != "reject") {
		http.Error(w, "Invalid parameters", http.StatusBadRequest)
		return
	}

	
	isPending, _ := repository.IsPending(input.SenderID, receiverID)
	if !isPending {
		http.Error(w, "No pending request", http.StatusBadRequest)
		return
	}

	var err error
	if input.Status == "accept" {
		err = repository.UpdateFollowStatus(input.SenderID, receiverID, "accepted")

	} else { 
		err = repository.RemoveFollow(input.SenderID, receiverID)
	}
	if err != nil {
		http.Error(w, "Operation failed", http.StatusInternalServerError)
		return
	}

	
	_ = sendInternalNotification(models.Notification{
		ID:       input.NotifID,
		UserID:   receiverID,
		Type:     "notification_removed",
		Message:  "Invitation_friendships",
		SenderID: input.SenderID,
	})

	
	followersCount, _ := repository.GetFollowersCount(receiverID)
	followingCount, _ := repository.GetFollowingCount(receiverID)
	json.NewEncoder(w).Encode(map[string]any{
		"status":          input.Status,
		"followers_Count": followersCount,
		"following_Count": followingCount,
	})
}


func sendInternalNotification(n models.Notification) error {
	b, _ := json.Marshal(n)
	
	resp, err := http.Post("http://localhost:8087/internal/notification", "application/json", bytes.NewReader(b))
	if err != nil {
		log.Printf("Error sending notification: %v\n", err)
		return err
	}
	defer resp.Body.Close()
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
		log.Println("DB update error:", err)
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
