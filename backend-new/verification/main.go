package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const (
	GatewayURL = "http://127.0.0.1:8080"
	ChatWsURL  = "ws://127.0.0.1:8085/ws" // Direct to chat service
)

type UserClient struct {
	Username      string
	Email         string
	Password      string
	Client        *http.Client
	UserID        int
	WsConn        *websocket.Conn
	Notifications chan map[string]interface{}
}

func NewUserClient(username, email, password string) *UserClient {
	jar, _ := cookiejar.New(nil)
	return &UserClient{
		Username:      username,
		Email:         email,
		Password:      password,
		Client:        &http.Client{Jar: jar},
		Notifications: make(chan map[string]interface{}, 100),
	}
}

func (u *UserClient) Register() error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("username", u.Username)
	writer.WriteField("first_name", "TestFirst")
	writer.WriteField("last_name", "TestLast")
	writer.WriteField("email", u.Email)
	writer.WriteField("password", u.Password)
	writer.WriteField("confirmpassword", u.Password)
	writer.WriteField("age", "30")
	writer.WriteField("gender", "male")
	writer.WriteField("nickname", u.Username)
	writer.WriteField("about", "tester")
	writer.Close()

	resp, err := u.Client.Post(GatewayURL+"/register", writer.FormDataContentType(), body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		if strings.Contains(string(b), "already used") {
			// Try to login to get ID if already registered
			return nil
		}
		return fmt.Errorf("register failed: code=%d body=%s", resp.StatusCode, string(b))
	}
	return nil
}

func (u *UserClient) Login() error {
	data := map[string]string{
		"username": u.Username,
		"password": u.Password,
	}
	jsonData, _ := json.Marshal(data)
	resp, err := u.Client.Post(GatewayURL+"/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var resBody map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&resBody); err != nil {
		return err
	}

	if status, ok := resBody["status"].(string); !ok || status != "ok" {
		return fmt.Errorf("login failed: %v", resBody)
	}

	if uid, ok := resBody["user_id"].(float64); ok {
		u.UserID = int(uid)
	} else {
		return fmt.Errorf("user_id not found in response")
	}
	return nil
}

func (u *UserClient) ConnectWS() error {
	uUrl, _ := url.Parse(GatewayURL)
	cookies := u.Client.Jar.Cookies(uUrl)
	var sessionCookie *http.Cookie
	for _, c := range cookies {
		if c.Name == "session" {
			sessionCookie = c
			break
		}
	}
	if sessionCookie == nil {
		return fmt.Errorf("no session cookie")
	}

	header := http.Header{}
	header.Add("Cookie", sessionCookie.String())

	conn, _, err := websocket.DefaultDialer.Dial(ChatWsURL, header)
	if err != nil {
		return err
	}
	u.WsConn = conn

	go func() {
		defer close(u.Notifications)
		for {
			var msg map[string]interface{}
			if err := conn.ReadJSON(&msg); err != nil {
				return
			}
			u.Notifications <- msg
		}
	}()
	return nil
}

func (u *UserClient) WaitForNotification(key, value, msgType string, timeout time.Duration) (map[string]interface{}, error) {
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case msg := <-u.Notifications:
			fmt.Printf("DEBUG: Received notification: %+v\n", msg)
			if msgType != "" {
				if t, ok := msg["type"].(string); !ok || t != msgType {
					continue
				}
			}
			if key != "" {
				if v, ok := msg[key]; ok {
					if strings.Contains(fmt.Sprintf("%v", v), value) {
						return msg, nil
					}
				}
			} else {
				return msg, nil
			}
		case <-timer.C:
			return nil, fmt.Errorf("timeout waiting for notification type %s", msgType)
		}
	}
}

func main() {
	// Setup Users
	// Use unique usernames to avoid conflicts on repeated runs if DB persists
	ts := time.Now().Unix()
	u1 := NewUserClient(fmt.Sprintf("userA_%d", ts), fmt.Sprintf("usera_%d@test.com", ts), "Password123!")
	u2 := NewUserClient(fmt.Sprintf("userB_%d", ts), fmt.Sprintf("userb_%d@test.com", ts), "Password123!")

	fmt.Println("Registering users...")
	if err := u1.Register(); err != nil {
		log.Fatal(err)
	}
	if err := u2.Register(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Logging in users...")
	if err := u1.Login(); err != nil {
		log.Fatal(err)
	}
	if err := u2.Login(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Logged in: U1=%d, U2=%d\n", u1.UserID, u2.UserID)

	// Connect WS for U2 (receiver)
	fmt.Println("Connecting WS for U2...")
	if err := u2.ConnectWS(); err != nil {
		log.Fatal(err)
	}

	// --- Test 1: Follow / Cancel Follow ---
	fmt.Println("\n--- Test 1: Follow Request & Cancel ---")
	followReq := map[string]interface{}{"followed_id": u2.UserID}
	jsonFollow, _ := json.Marshal(followReq)

	fmt.Println("Sending Follow Request...")
	resp, err := u1.Client.Post(GatewayURL+"/toggle-follow", "application/json", bytes.NewBuffer(jsonFollow))
	if err != nil {
		fmt.Printf("FAIL: Failed to send follow request: %v\n", err)
	} else {
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Printf("Follow Request Status: %s Body: %s\n", resp.Status, string(bodyBytes))
		resp.Body.Close()
	}

	msg, err := u2.WaitForNotification("type", "new_follower", "new_follower", 5*time.Second)
	if err != nil {
		log.Printf("FAIL: Follow notification: %v", err)
	} else {
		fmt.Printf("PASS: Received Follow Notification: %+v\n", msg)
	}

	fmt.Println("Cancelling Follow Request (toggling)...")
	u1.Client.Post(GatewayURL+"/toggle-follow", "application/json", bytes.NewBuffer(jsonFollow))

	msg, err = u2.WaitForNotification("", "", "notification_removed", 5*time.Second)
	if err != nil {
		log.Printf("FAIL: Removal notification: %v", err)
	} else {
		fmt.Printf("PASS: Received Removal Notification: %+v\n", msg)
	}

	// --- Test 2: Group Invite & Accept ---
	fmt.Println("\n--- Test 2: Group Invite & Accept ---")
	groupReq := map[string]string{"Title": fmt.Sprintf("TestGroup_%d", ts), "Description": "Test Desc"}
	jsonGroup, _ := json.Marshal(groupReq)
	resp, _ = u1.Client.Post(GatewayURL+"/Create_Group", "application/json", bytes.NewBuffer(jsonGroup))
	var groupRes map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&groupRes)
	// Response might vary, let's print it if needed. Assuming standard ID return.
	// Actually Create_Group handler returns "Group created successfully" usually.
	// Let's check if it returns ID. If not, we might need to fetch groups.
	// Handlers often just say "ok".
	// Let's assume we can fetch U1's groups to find the ID.

	resp, _ = u1.Client.Get(GatewayURL + "/Get_Groups")
	var groups []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&groups)
	var groupID int
	if len(groups) > 0 {
		groupID = int(groups[len(groups)-1]["id"].(float64)) // Last one

		// Verify Creator Membership
		lastGroup := groups[len(groups)-1]
		if status, ok := lastGroup["user_status"].(string); ok && status == "member" {
			fmt.Printf("PASS: Creator is member of Group %d\n", groupID)
		} else {
			fmt.Printf("FAIL: Creator is NOT member of Group %d. Status: %v\n", groupID, lastGroup["user_status"])
		}
	}
	fmt.Printf("Group Created: %d\n", groupID)

	inviteReq := map[string]interface{}{"user_id": u2.UserID, "group_id": groupID}
	jsonInvite, _ := json.Marshal(inviteReq)
	fmt.Println("Inviting U2...")
	u1.Client.Post(GatewayURL+"/groups/invite", "application/json", bytes.NewBuffer(jsonInvite))

	msg, err = u2.WaitForNotification("type", "group_invite", "group_invite", 5*time.Second)
	if err != nil {
		log.Printf("FAIL: Group invite notification: %v", err)
	} else {
		fmt.Printf("PASS: Received Group Invite Notification: %+v\n", msg)
	}

	// Respond to Invite (Accept)
	fmt.Println("U2 Accepting Invite...")
	respondReq := map[string]interface{}{
		"group_id":        groupID,
		"status":          "accept",
		"notification_id": 0, // Not needed by logic?
	}
	jsonRespond, _ := json.Marshal(respondReq)
	u2.Client.Post(GatewayURL+"/group/invite/respond", "application/json", bytes.NewBuffer(jsonRespond))

	msg, err = u2.WaitForNotification("", "", "notification_removed", 5*time.Second)
	if err != nil {
		log.Printf("FAIL: Invite Removal notification: %v", err)
	} else {
		fmt.Printf("PASS: Received Invite Removal Notification: %+v\n", msg)
	}

	// --- Test 3: Event Creation & Respond ---
	fmt.Println("\n--- Test 3: Event Creation & Respond ---")
	eventReq := map[string]interface{}{
		"group_id":    groupID,
		"title":       "Party",
		"description": "Fun",
		"date":        "2025-12-31T20:00:00Z",
	}
	jsonEvent, _ := json.Marshal(eventReq)
	fmt.Println("U1 Creating Event...")
	resp, err = u1.Client.Post(GatewayURL+"/create_event", "application/json", bytes.NewBuffer(jsonEvent))
	if err != nil {
		fmt.Printf("FAIL: Failed to create event: %v\n", err)
	} else {
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Printf("Create Event Status: %s Body: %s\n", resp.Status, string(bodyBytes))
		resp.Body.Close()
	}

	msg, err = u2.WaitForNotification("type", "event_created", "event_created", 5*time.Second)
	if err != nil {
		log.Printf("FAIL: Event notification: %v", err)
		return
	} else {
		fmt.Printf("PASS: Received Event Notification: %+v\n", msg)
	}

	// Get Event ID from notification
	eventID := int(msg["event_id"].(float64))

	fmt.Println("U2 Responding 'going'...")
	eventRespondReq := map[string]interface{}{
		"event_id": eventID,
		"status":   "going",
	}
	jsonEventRespond, _ := json.Marshal(eventRespondReq)
	u2.Client.Post(GatewayURL+"/event-invite-response", "application/json", bytes.NewBuffer(jsonEventRespond))

	msg, err = u2.WaitForNotification("", "", "notification_removed", 5*time.Second)
	if err != nil {
		log.Printf("FAIL: Event Removal notification: %v", err)
	} else {
		fmt.Printf("PASS: Received Event Removal Notification: %+v\n", msg)
	}

	fmt.Println("\nDONE")
}
