package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"group-service/internal/repository"
	"io"
	"net/http"
	"os"
	"social-network/pkg/middleware"
	"social-network/pkg/models"
	sharedRepo "social-network/pkg/repository"
	"strconv"
	"strings"
	"time"
)

// Helper to send notification to chat-service
func sendNotification(n models.Notification) {
	fmt.Printf("DEBUG: Sending notification: %+v\n", n)
	data, _ := json.Marshal(n)
	resp, err := http.Post("http://localhost:8085/internal/notification", "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Printf("DEBUG: Error sending notification: %v\n", err)
		return
	}
	defer resp.Body.Close()
	fmt.Println("DEBUG: Notification sent successfully")
}

// CreateGroup handles group creation
func CreateGroup(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Title       string `json:"Title"`
		Description string `json:"Description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	groupID, err := repository.InsertGroup(req.Title, req.Description, userID)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	if err := repository.InsertGroupMember(int(groupID), userID, "member"); err != nil {
		fmt.Printf("Error adding creator as member: %v\n", err)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"id": groupID, "status": "created"})
}

// GetGroups returns groups for the user
func GetGroups(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	groups, err := repository.GetGroups(userID)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(groups)
}

// RequestJoin handles join requests
func RequestJoin(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		GroupID int `json:"group_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Check if already member or pending
	groups, _ := repository.GetGroups(userID)
	for _, g := range groups {
		if g.ID == req.GroupID && (g.UserStatus == "member" || g.UserStatus == "requested" || g.UserStatus == "pending") {
			http.Error(w, "Already requested or member", http.StatusBadRequest)
			return
		}
	}

	err := repository.InsertGroupMember(req.GroupID, userID, "requested")
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Notify Creator
	group, _ := repository.GetGroupByID(req.GroupID)
	if group != nil {
		username, _ := sharedRepo.GetUsernameByID(userID)
		msg := fmt.Sprintf("The user %s sent a request to join group %s", username, group.Title)
		sendNotification(models.Notification{
			UserID:   group.CreatorID,
			SenderID: userID,
			Type:     "request_join_group",
			Message:  msg,
			GroupID:  req.GroupID,
		})
	}

	w.WriteHeader(http.StatusOK)
}

// AcceptRejectJoin handles join request decisions
func AcceptRejectJoin(w http.ResponseWriter, r *http.Request) {
	// Only Creator can accept/reject? Old code implies only creator gets notification.
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		NotifID int    `json:"notif_id"` // Used to identify request in old code, but here we need group_id and user_id
		GroupID int    `json:"group_id"`
		UserID  int    `json:"user_id"` // The user requesting join
		Action  string `json:"action"`  // 'accept' or 'reject'
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Verify current user is creator (or admin)
	group, err := repository.GetGroupByID(req.GroupID)
	if err != nil || group.CreatorID != userID {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	if req.Action == "accept" {
		repository.UpdateGroupMemberStatus(req.GroupID, req.UserID, "member")
		// Notify User
		sendNotification(models.Notification{
			UserID:   req.UserID,
			SenderID: userID,
			Type:     "group_join_accepted",
			Message:  fmt.Sprintf("Your request to join group %s was accepted", group.Title),
			GroupID:  req.GroupID,
		})
	} else {
		repository.RemoveGroupMember(req.GroupID, req.UserID)
	}

	// Remove the request notification for the Creator (myself)
	sendNotification(models.Notification{
		UserID:   userID,
		Type:     "notification_removed",
		Message:  "request_join_group",
		SenderID: req.UserID, // The user who made the request
		GroupID:  req.GroupID,
	})

	w.WriteHeader(http.StatusOK)
}

// RespondToGroupInvite handles invite response
func RespondToGroupInvite(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		GroupID int    `json:"group_id"`
		Status  string `json:"status"` // 'accept', 'reject'
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.Status == "accept" {
		repository.UpdateGroupMemberStatus(req.GroupID, userID, "member")
	} else {
		repository.RemoveGroupMember(req.GroupID, userID)
	}

	// Remove notification for this user
	sendNotification(models.Notification{
		UserID:  userID,
		Type:    "notification_removed",
		Message: "group_invite",
		GroupID: req.GroupID,
	})

	w.WriteHeader(http.StatusOK)
}

// GetGroupByID returns single group details
func GetGroupByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		// Try path param
		idStr = strings.TrimPrefix(r.URL.Path, "/Get_Group_By_ID/")
	}

	if idStr == "" || idStr == "/Get_Group_By_ID/" {
		http.Error(w, "Missing id", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	group, err := repository.GetGroupByID(id)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	// Fetch members
	members, _ := repository.GetGroupMembersDetails(id)

	// Check if current user is member
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	var userStatus string
	if ok {
		// We can reuse GetGroups or just query simply
		// But GetGroups gets all groups.
		// Let's just check quickly if they are in members list?
		// Or query DB.
		// For efficiency, since we have members list, we can check it.
		for _, m := range members {
			if m.ID == userID {
				userStatus = "member"
				break
			}
		}
		// If not found in members, they might be pending/requested?
		// GetGroupMembersDetails only returns 'member' status.
		// If we want accurate status (pending etc), we should query DB.
		// But let's stick to what we have or improve repository.
	}

	// Construct response
	response := map[string]interface{}{
		"group": map[string]interface{}{
			"id":          group.ID,
			"title":       group.Title,
			"description": group.Description,
			"creator_id":  group.CreatorID,
			"is_member":   userStatus == "member", // specific field for frontend
			"user_status": userStatus,
		},
		"members": members,
	}

	json.NewEncoder(w).Encode(response)
}

// CreateGroupPost handles posting to a group
func CreateGroupPost(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Parse error", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")
	groupID, _ := strconv.Atoi(r.FormValue("group_id"))

	// Handle Image
	var imagePath string
	file, handler, err := r.FormFile("image")
	if err == nil {
		defer file.Close()
		// Save to shared uploads dir? Or media service?
		// For now, save to local uploads and serve via Media Service?
		// Better: Forward to Media Service?
		// Simplest for migration: Save to "../media-service/uploads" relative path?
		// Or just "uploads/" and ensure Media Service serves it.
		// For now, save to shared uploads dir
		uploadDir := "../../uploads"
		os.MkdirAll(uploadDir, 0755)
		filename := fmt.Sprintf("%d_%s", time.Now().Unix(), handler.Filename)
		dstPath := fmt.Sprintf("%s/%s", uploadDir, filename)
		dst, _ := os.Create(dstPath)
		defer dst.Close()
		io.Copy(dst, file)
		imagePath = "uploads/" + filename // Path stored in DB
	}

	id, err := repository.InsertGroupPost(groupID, userID, title, content, imagePath)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	post, _ := repository.GetGroupPostByID(int(id))
	json.NewEncoder(w).Encode(post)
}

// GetGroupPosts retrieves posts for a group
func GetGroupPosts(w http.ResponseWriter, r *http.Request) {
	groupIDStr := r.URL.Query().Get("group_id")
	groupID, _ := strconv.Atoi(groupIDStr)

	posts, err := repository.GetGroupPosts(groupID)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

// FetchNotMembers returns groups where user can be invited
func FetchNotMembers(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		UserID int `json:"user_id"` // Target user to invite
	}
	json.NewDecoder(r.Body).Decode(&req)

	groups, err := repository.GetGroupsForInvite(userID, req.UserID)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(groups)
}

// InviteToGroup handles invitations
func InviteToGroup(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		GroupID int `json:"group_id"`
		UserID  int `json:"user_id"` // User to invite
	}
	json.NewDecoder(r.Body).Decode(&req)

	// Check if already member? (GetGroupsForInvite handles this logic partly, but safe to check)
	repository.InsertGroupMember(req.GroupID, req.UserID, "invited")

	// Notify Invitee
	group, _ := repository.GetGroupByID(req.GroupID)
	inviterName, _ := sharedRepo.GetUsernameByID(userID)
	msg := fmt.Sprintf("%s invited you to join group %s", inviterName, group.Title)

	sendNotification(models.Notification{
		UserID:   req.UserID,
		SenderID: userID,
		Type:     "group_invite",
		Message:  msg,
		GroupID:  req.GroupID,
	})

	w.WriteHeader(http.StatusOK)
}

// CreateEvent handles event creation
func CreateEvent(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		GroupID     int    `json:"group_id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		EventTime   string `json:"event_time"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	eventID, err := repository.InsertEvent(req.GroupID, userID, req.Title, req.Description, req.EventTime)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	repository.InsertEventMember(int(eventID), userID, "going")

	// Notify all group members
	members, _ := repository.GetGroupMembers(req.GroupID)
	for _, memberID := range members {
		if memberID != userID {
			sendNotification(models.Notification{
				UserID:   memberID,
				SenderID: userID,
				Type:     "event_created",
				Message:  fmt.Sprintf("New event in group: %s", req.Title),
				GroupID:  req.GroupID,
				EventID:  int(eventID),
			})
		}
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"status": "created"})
}

// GetGroupEvents retrieves events
func GetGroupEvents(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	groupIDStr := r.URL.Query().Get("group_id")
	groupID, _ := strconv.Atoi(groupIDStr)

	events, err := repository.GetGroupEvents(groupID, userID)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(events)
}

// RespondToEvent handles going/not_going
func RespondToEvent(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		EventID int    `json:"event_id"`
		Status  string `json:"status"` // 'going', 'not_going'
	}
	json.NewDecoder(r.Body).Decode(&req)

	err := repository.UpdateEventMemberStatus(req.EventID, userID, req.Status)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Remove "event_created" notification for this event
	sendNotification(models.Notification{
		UserID:  userID,
		Type:    "notification_removed",
		Message: "event_created",
		EventID: req.EventID,
	})

	w.WriteHeader(http.StatusOK)
}
