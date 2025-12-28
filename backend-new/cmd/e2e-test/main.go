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
	"time"

	"github.com/gorilla/websocket"
)

// Constants
const (
	AuthURL    = "http://localhost:8081"
	UserURL    = "http://localhost:8082"
	NotifURL   = "http://localhost:8087"
	NotifWSURL = "ws://localhost:8087/ws"
)

// Client wraps http.Client to handle cookies
type Client struct {
	HTTPClient *http.Client
}

func NewClient() *Client {
	jar, _ := cookiejar.New(nil)
	return &Client{
		HTTPClient: &http.Client{
			Jar:     jar,
			Timeout: 10 * time.Second,
		},
	}
}

// Login performs login and stores cookies
func (c *Client) Login(email, password string) error {
	reqBody, _ := json.Marshal(map[string]string{
		"username": email,
		"password": password,
	})

	resp, err := c.HTTPClient.Post(AuthURL+"/login", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("login failed: %s - %s", resp.Status, string(body))
	}
	return nil
}

// Register registers a new user using Multipart Form
func (c *Client) Register(email, password, first, last, nick, dob string) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	_ = writer.WriteField("email", email)
	_ = writer.WriteField("password", password)
	_ = writer.WriteField("confirmpassword", password)
	_ = writer.WriteField("first_name", first)
	_ = writer.WriteField("last_name", last)
	_ = writer.WriteField("nickname", nick)
	_ = writer.WriteField("username", nick)
	_ = writer.WriteField("age", "25")
	_ = writer.WriteField("gender", "male")
	_ = writer.WriteField("about", "test user")

	writer.Close()

	req, _ := http.NewRequest("POST", AuthURL+"/register", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("register failed: %s - %s", resp.Status, string(body))
	}
	return nil
}

// SetPrivacy updates user privacy
func (c *Client) SetPrivacy(isPrivate int) error {
	reqBody, _ := json.Marshal(map[string]int{
		"is_private": isPrivate,
	})
	req, _ := http.NewRequest("POST", UserURL+"/update-privacy", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("set privacy failed: %s", resp.Status)
	}
	return nil
}

// FollowUser sends a follow request
func (c *Client) FollowUser(targetID int) error {
	reqBody, _ := json.Marshal(map[string]int{
		"followed_id": targetID,
	})
	// User service handles /toggle-follow directly
	req, _ := http.NewRequest("POST", UserURL+"/toggle-follow", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("follow user failed: %s", resp.Status)
	}
	return nil
}

// GetNotifications fetches notifications via HTTP
func (c *Client) GetNotifications() ([]map[string]interface{}, error) {
	req, _ := http.NewRequest("GET", NotifURL+"/notifications", nil)
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get notifications failed: %s", resp.Status)
	}
	var notifs []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&notifs)
	return notifs, nil
}

func main() {
	log.Println("Starting End-to-End Notification Test (Direct Ports)...")

	// 1. Setup Users
	ts := time.Now().Unix()
	user1Nick := fmt.Sprintf("u1_%d", ts)
	user2Nick := fmt.Sprintf("u2_%d", ts)
	user1Email := fmt.Sprintf("%s@test.com", user1Nick)
	user2Email := fmt.Sprintf("%s@test.com", user2Nick)
	password := "password123"

	client1 := NewClient()
	client2 := NewClient()

	log.Println("Registering User 1...")
	if err := client1.Register(user1Email, password, "User", "One", user1Nick, "2000-01-01"); err != nil {
		log.Fatalf("Failed to register user 1: %v", err)
	}
	log.Println("Logging in User 1...")
	if err := client1.Login(user1Nick, password); err != nil {
		log.Fatalf("Failed to login user 1: %v", err)
	}

	log.Println("Registering User 2...")
	if err := client2.Register(user2Email, password, "User", "Two", user2Nick, "2000-01-01"); err != nil {
		log.Fatalf("Failed to register user 2: %v", err)
	}
	log.Println("Logging in User 2...")
	if err := client2.Login(user2Nick, password); err != nil {
		log.Fatalf("Failed to login user 2: %v", err)
	}

	// 2. Set User 2 as Private
	log.Println("Setting User 2 as Private...")
	if err := client2.SetPrivacy(1); err != nil {
		log.Fatalf("Failed to set privacy: %v", err)
	}

	// Get User 2 ID via Profile (User Service)
	resp, err := client2.HTTPClient.Get(UserURL + "/profile")
	if err != nil {
		log.Fatalf("Failed to get profile 2: %v", err)
	}
	var profile2 struct {
		User struct {
			ID int `json:"id"`
		} `json:"user"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&profile2); err != nil {
		log.Fatalf("Failed to decode profile: %v", err)
	}
	user2ID := profile2.User.ID
	log.Printf("User 2 ID: %d", user2ID)

	// 3. Connect User 2 to WebSocket
	log.Printf("Connecting User 2 to WebSocket: %s", NotifWSURL)

	// Ensure cookie is passed.
	// Note: Cookie domain is 'localhost'. It should be available.
	// We need to find the session cookie.
	var sessionCookie *http.Cookie
	// Check against AuthURL since that's where we logged in
	authReqURL, _ := http.NewRequest("GET", AuthURL, nil)
	for _, c := range client2.HTTPClient.Jar.Cookies(authReqURL.URL) {
		if c.Name == "session" {
			sessionCookie = c
			break
		}
	}

	header := http.Header{}
	if sessionCookie != nil {
		header.Add("Cookie", sessionCookie.String())
	} else {
		log.Println("Warning: No session cookie found for Auth URL")
	}

	ws, _, err := websocket.DefaultDialer.Dial(NotifWSURL, header)
	if err != nil {
		log.Fatalf("WebSocket connection failed: %v", err)
	}
	defer ws.Close()

	// Channel to receive notifications
	notifChan := make(chan map[string]interface{})
	go func() {
		for {
			var msg map[string]interface{}
			if err := ws.ReadJSON(&msg); err != nil {
				close(notifChan)
				return
			}
			notifChan <- msg
		}
	}()

	// 4. User 1 Follows User 2
	log.Println("User 1 sending Follow Request to User 2...")
	if err := client1.FollowUser(user2ID); err != nil {
		log.Fatalf("Failed to follow user: %v", err)
	}

	// 5. Verify Notification via WebSocket
	log.Println("Waiting for notification via WebSocket...")
	select {
	case msg := <-notifChan:
		log.Printf("Received WebSocket Message: %+v", msg)
		if msg["type"] == "Invitation_friendships" {
			log.Println("✅ Success: Received 'Invitation_friendships' notification via WebSocket")
		} else {
			log.Printf("⚠️ Received unexpected type: %v", msg["type"])
		}
	case <-time.After(5 * time.Second):
		log.Fatal("❌ Timeout: Did not receive notification via WebSocket")
	}

	// 6. Verify Notification via HTTP
	log.Println("Verifying notification via HTTP API...")
	notifs, err := client2.GetNotifications()
	if err != nil {
		log.Fatalf("Failed to get notifications: %v", err)
	}
	found := false
	for _, n := range notifs {
		if n["type"] == "Invitation_friendships" {
			found = true
			break
		}
	}
	if found {
		log.Println("✅ Success: Notification found in API response")
	} else {
		log.Fatal("❌ Failure: Notification NOT found in API response")
	}

	log.Println("🎉 ALL TESTS PASSED!")
}
