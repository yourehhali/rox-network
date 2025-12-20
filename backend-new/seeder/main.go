package main

import (
	"fmt"
	"log"
	"social-network/pkg/database"
	"social-network/pkg/utils"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	dbPath := "../database/forum.db"
	database.InitDB(dbPath)
	fmt.Println("Seeder connected to database at", dbPath)

	// Wait for tables to be created by Auth Service if running concurrently,
	// but we assume this runs after Auth Service has started and InitDB'd.
	// We can also ensure tables exist here if we want to be safe, but schema.sql is authoritative.

	seedUsers()
	seedPosts()
	seedGroups()
}

func seedUsers() {
	users := []struct {
		Username  string
		Email     string
		Password  string
		FirstName string
		LastName  string
		Age       int
		Gender    string
		Nickname  string
		About     string
		Avatar    string
	}{
		{"alice", "alice@example.com", "password123", "Alice", "Wonderland", 25, "female", "Ali", "I love coding", ""},
		{"bob", "bob@example.com", "password123", "Bob", "Builder", 30, "male", "Bobby", "Can we fix it?", ""},
		{"charlie", "charlie@example.com", "password123", "Charlie", "Chaplin", 35, "male", "Chuck", "Silent comedian", ""},
	}

	stmt, err := database.DB.Prepare(`
		INSERT INTO users (username, email, password_hash, first_name, last_name, age, gender, nickname, about, avatar, is_private)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0)
	`)
	if err != nil {
		// Table might not exist if Auth service hasn't run schema yet.
		log.Printf("Error preparing user insert (tables missing?): %v", err)
		return
	}
	defer stmt.Close()

	for _, u := range users {
		// Check if exists
		var exists int
		database.DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", u.Username).Scan(&exists)
		if exists > 0 {
			fmt.Printf("User %s already exists, skipping.\n", u.Username)
			continue
		}

		hash, _ := utils.HashPassword(u.Password)
		_, err := stmt.Exec(u.Username, u.Email, hash, u.FirstName, u.LastName, u.Age, u.Gender, u.Nickname, u.About, u.Avatar)
		if err != nil {
			log.Printf("Failed to insert user %s: %v", u.Username, err)
		} else {
			fmt.Printf("Inserted user: %s\n", u.Username)
		}
	}

	// Create Follows
	// Alice follows Bob
	aliceID := getUserID("alice")
	bobID := getUserID("bob")
	if aliceID > 0 && bobID > 0 {
		_, err := database.DB.Exec("INSERT OR IGNORE INTO followers (follower_id, followed_id, status) VALUES (?, ?, 'accepted')", aliceID, bobID)
		if err == nil {
			fmt.Println("Alice follows Bob")
		}
	}
}

func seedPosts() {
	posts := []struct {
		Username string
		Title    string
		Content  string
		Image    string
	}{
		{"alice", "My First Post", "Hello everyone! This is my first post on Rox Social.", ""},
		{"bob", "Building things", "Just finished a new project. It's amazing!", ""},
		{"charlie", "Silent Movie", "...", "uploads/silent_movie.jpg"},
	}

	stmt, err := database.DB.Prepare(`INSERT INTO posts (user_id, title, content, image_path) VALUES (?, ?, ?, ?)`)
	if err != nil {
		log.Printf("Error preparing post insert: %v", err)
		return
	}
	defer stmt.Close()

	for _, p := range posts {
		uid := getUserID(p.Username)
		if uid == 0 {
			continue
		}
		_, err := stmt.Exec(uid, p.Title, p.Content, p.Image)
		if err != nil {
			log.Printf("Failed to insert post for %s: %v", p.Username, err)
		} else {
			fmt.Printf("Inserted post by %s\n", p.Username)
		}
	}
}

func seedGroups() {
	// Create Group "Tech Talk" by Alice
	aliceID := getUserID("alice")
	bobID := getUserID("bob")
	charlieID := getUserID("charlie")

	if aliceID == 0 {
		return
	}

	res, err := database.DB.Exec(`INSERT INTO groups (title, description, creator_id) VALUES (?, ?, ?)`, "Tech Talk", "Discussions about technology", aliceID)
	if err != nil {
		log.Printf("Failed to insert group: %v", err)
		return
	}
	groupID, _ := res.LastInsertId()
	fmt.Printf("Inserted Group: Tech Talk (ID: %d)\n", groupID)

	// Add Members
	// Alice (Creator) -> Member
	database.DB.Exec(`INSERT INTO group_members (group_id, user_id, status) VALUES (?, ?, 'member')`, groupID, aliceID)
	// Bob -> Member
	if bobID > 0 {
		database.DB.Exec(`INSERT INTO group_members (group_id, user_id, status) VALUES (?, ?, 'member')`, groupID, bobID)
		fmt.Println("Bob joined Tech Talk")
	}
	// Charlie -> Invited
	if charlieID > 0 {
		database.DB.Exec(`INSERT INTO group_members (group_id, user_id, status) VALUES (?, ?, 'invited')`, groupID, charlieID)
		fmt.Println("Charlie invited to Tech Talk")
	}

	// Create Event
	database.DB.Exec(`INSERT INTO events (group_id, creator_id, title, description, event_time) VALUES (?, ?, ?, ?, ?)`,
		groupID, aliceID, "Go Workshop", "Learning Go Language", time.Now().Add(24*time.Hour).Format(time.RFC3339))
	fmt.Println("Created Event: Go Workshop")
}

func getUserID(username string) int {
	var id int
	err := database.DB.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&id)
	if err != nil {
		return 0
	}
	return id
}
