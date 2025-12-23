package repository

import (
	"database/sql"
	"social-network/pkg/database"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) {
	var err error
	database.DB, err = sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	schema := `
	CREATE TABLE users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL,
		email TEXT NOT NULL,
		password_hash TEXT NOT NULL,
		first_name TEXT,
		last_name TEXT,
		avatar TEXT,
        is_private INTEGER DEFAULT 0,
        age INTEGER,
        gender TEXT,
        nickname TEXT,
        about TEXT,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE TABLE followers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		follower_id INTEGER NOT NULL,
		followed_id INTEGER NOT NULL,
		status TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err = database.DB.Exec(schema)
	if err != nil {
		t.Fatalf("Failed to create tables: %v", err)
	}
}

func TestGetFollowers(t *testing.T) {
	setupTestDB(t)
	defer database.DB.Close()

	// Insert Users
	_, err := database.DB.Exec(`INSERT INTO users (id, username, email, password_hash, first_name, last_name) VALUES (1, 'u1', 'e1', 'h', 'User', 'One')`)
    if err != nil { t.Fatal(err) }
	_, err = database.DB.Exec(`INSERT INTO users (id, username, email, password_hash, first_name, last_name) VALUES (2, 'u2', 'e2', 'h', 'User', 'Two')`)
    if err != nil { t.Fatal(err) }

	// User 2 follows User 1
	_, err = database.DB.Exec(`INSERT INTO followers (follower_id, followed_id, status) VALUES (2, 1, 'accepted')`)
    if err != nil { t.Fatal(err) }

	// Test GetFollowers for User 1
	followers, err := GetFollowers(1)
	if err != nil {
		t.Fatalf("GetFollowers failed: %v", err)
	}

	if len(followers) != 1 {
		t.Errorf("Expected 1 follower, got %d", len(followers))
	}
	if followers[0].ID != 2 {
		t.Errorf("Expected follower ID 2, got %d", followers[0].ID)
	}
    
    // Test GetFollowers for User 2 (should be 0)
    followers2, err := GetFollowers(2)
    if err != nil { t.Fatal(err) }
    if len(followers2) != 0 {
        t.Errorf("Expected 0 followers, got %d", len(followers2))
    }
}

func TestGetFollowing(t *testing.T) {
	setupTestDB(t)
	defer database.DB.Close()

	// Insert Users
	_, err := database.DB.Exec(`INSERT INTO users (id, username, email, password_hash, first_name, last_name) VALUES (1, 'u1', 'e1', 'h', 'User', 'One')`)
    if err != nil { t.Fatal(err) }
	_, err = database.DB.Exec(`INSERT INTO users (id, username, email, password_hash, first_name, last_name) VALUES (2, 'u2', 'e2', 'h', 'User', 'Two')`)
    if err != nil { t.Fatal(err) }

	// User 1 follows User 2
	_, err = database.DB.Exec(`INSERT INTO followers (follower_id, followed_id, status) VALUES (1, 2, 'accepted')`)
    if err != nil { t.Fatal(err) }

	// Test GetFollowing for User 1
	following, err := GetFollowing(1)
	if err != nil {
		t.Fatalf("GetFollowing failed: %v", err)
	}

	if len(following) != 1 {
		t.Errorf("Expected 1 following, got %d", len(following))
	}
	if following[0].ID != 2 {
		t.Errorf("Expected following ID 2, got %d", following[0].ID)
	}
}
