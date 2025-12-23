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

	// Create tables (Simplified schema for testing)
	schema := `
	CREATE TABLE users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL,
		email TEXT NOT NULL,
		password_hash TEXT NOT NULL,
		first_name TEXT,
		last_name TEXT,
		avatar TEXT
	);
	CREATE TABLE posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		image_path TEXT,
        privacy INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
    CREATE TABLE post_viewers (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        post_id INTEGER NOT NULL,
        user_id INTEGER NOT NULL
    );
    CREATE TABLE followers (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        follower_id INTEGER NOT NULL,
        followed_id INTEGER NOT NULL,
        status TEXT
    );
	CREATE TABLE comments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		post_id INTEGER NOT NULL,
		content TEXT NOT NULL,
		image_url TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE TABLE reactions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		post_id INTEGER NOT NULL,
		reaction_type TEXT
	);
	`
	_, err = database.DB.Exec(schema)
	if err != nil {
		t.Fatalf("Failed to create tables: %v", err)
	}
}

func TestAddNewPost(t *testing.T) {
	setupTestDB(t)
	defer database.DB.Close()

	// Insert User
	_, err := database.DB.Exec(`INSERT INTO users (username, email, password_hash) VALUES ('testuser', 'test@test.com', 'hash')`)
	if err != nil {
		t.Fatalf("Failed to insert user: %v", err)
	}

	// Test AddNewPost (Public)
	postID, err := AddNewPost(1, "New Title", "New Content", "image.png", 0, nil)
	if err != nil {
		t.Fatalf("AddNewPost failed: %v", err)
	}

	if postID <= 0 {
		t.Errorf("Expected valid post ID, got %d", postID)
	}

	// Verify insertion
	var title, content, imagePath string
	var privacy int
	err = database.DB.QueryRow("SELECT title, content, image_path, privacy FROM posts WHERE id = ?", postID).Scan(&title, &content, &imagePath, &privacy)
	if err != nil {
		t.Fatalf("Failed to retrieve added post: %v", err)
	}

	if title != "New Title" {
		t.Errorf("Expected title 'New Title', got '%s'", title)
	}
	if content != "New Content" {
		t.Errorf("Expected content 'New Content', got '%s'", content)
	}
	if imagePath != "image.png" {
		t.Errorf("Expected imagePath 'image.png', got '%s'", imagePath)
	}
	if privacy != 0 {
		t.Errorf("Expected privacy 0, got %d", privacy)
	}

	// Test AddNewPost (Private with viewers)
	postID2, err := AddNewPost(1, "Private Post", "Private Content", "", 2, []string{"2", "3"})
	if err != nil {
		t.Fatal(err)
	}

	// Check viewers
	var count int
	err = database.DB.QueryRow("SELECT COUNT(*) FROM post_viewers WHERE post_id = ?", postID2).Scan(&count)
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Errorf("Expected 2 post viewers, got %d", count)
	}
}

func TestGetAllPosts(t *testing.T) {
	setupTestDB(t)
	defer database.DB.Close()

	// Insert Users
	_, err := database.DB.Exec(`INSERT INTO users (username, email, password_hash, first_name, last_name) VALUES ('u1', 'u1@e.com', 'h', 'User', 'One')`)
	if err != nil {
		t.Fatal(err)
	}
	_, err = database.DB.Exec(`INSERT INTO users (username, email, password_hash, first_name, last_name) VALUES ('u2', 'u2@e.com', 'h', 'User', 'Two')`)
	if err != nil {
		t.Fatal(err)
	}
	_, err = database.DB.Exec(`INSERT INTO users (username, email, password_hash, first_name, last_name) VALUES ('u3', 'u3@e.com', 'h', 'User', 'Three')`)
	if err != nil {
		t.Fatal(err)
	}

	// Insert Posts
	// Post 1: Public by User 1
	_, err = database.DB.Exec(`INSERT INTO posts (user_id, title, content, privacy) VALUES (1, 'Post 1', 'Content 1', 0)`)
	if err != nil {
		t.Fatal(err)
	}
	// Post 2: Almost Private by User 1 (Followers only)
	_, err = database.DB.Exec(`INSERT INTO posts (user_id, title, content, privacy) VALUES (1, 'Post 2', 'Content 2', 1)`)
	if err != nil {
		t.Fatal(err)
	}
	// Post 3: Private by User 1 (Only User 2 allowed)
	_, err = database.DB.Exec(`INSERT INTO posts (user_id, title, content, privacy) VALUES (1, 'Post 3', 'Content 3', 2)`)
	if err != nil {
		t.Fatal(err)
	}
	_, err = database.DB.Exec(`INSERT INTO post_viewers (post_id, user_id) VALUES (3, 2)`)
	if err != nil {
		t.Fatal(err)
	}

	// Insert Reactions
	// User 2 likes Post 1
	_, err = database.DB.Exec(`INSERT INTO reactions (user_id, post_id, reaction_type) VALUES (2, 1, 'like')`)
	if err != nil {
		t.Fatal(err)
	}

	// Setup Followers: User 2 follows User 1
	_, err = database.DB.Exec(`INSERT INTO followers (follower_id, followed_id, status) VALUES (2, 1, 'accepted')`)
	if err != nil {
		t.Fatal(err)
	}

	// Test GetAllPosts for User 2 (Should see Post 1 (Public), Post 2 (Follower), Post 3 (Allowed)) -> Total 3
	posts, err := GetAllPosts(2)
	if err != nil {
		t.Fatalf("GetAllPosts failed: %v", err)
	}
	if len(posts) != 3 {
		t.Errorf("Expected 3 posts for User 2, got %d", len(posts))
	}

	// Test GetAllPosts for User 3 (Should see Post 1 (Public) ONLY) -> Total 1
	// User 3 is NOT follower, NOT allowed in Post 3
	posts3, err := GetAllPosts(3)
	if err != nil {
		t.Fatal(err)
	}
	if len(posts3) != 1 {
		t.Errorf("Expected 1 post for User 3, got %d", len(posts3))
	}
	if len(posts3) > 0 && posts3[0].ID != 1 {
		t.Errorf("Expected User 3 to see Post 1, got Post %d", posts3[0].ID)
	}
}

func TestGetUserPosts(t *testing.T) {
	setupTestDB(t)
	defer database.DB.Close()

	// Insert Users
	_, err := database.DB.Exec(`INSERT INTO users (username, email, password_hash, first_name, last_name) VALUES ('u1', 'u1@e.com', 'h', 'User', 'One')`)
	if err != nil {
		t.Fatal(err)
	}
	_, err = database.DB.Exec(`INSERT INTO users (username, email, password_hash, first_name, last_name) VALUES ('u2', 'u2@e.com', 'h', 'User', 'Two')`)
	if err != nil {
		t.Fatal(err)
	}

	// Insert Posts for User 1
	_, err = database.DB.Exec(`INSERT INTO posts (user_id, title, content, privacy) VALUES (1, 'Post 1', 'Content 1', 0)`)
	if err != nil {
		t.Fatal(err)
	}
	_, err = database.DB.Exec(`INSERT INTO posts (user_id, title, content, privacy) VALUES (1, 'Post 2', 'Content 2', 2)`) // Private
	if err != nil {
		t.Fatal(err)
	}

	// Test GetUserPosts for User 1 viewing their own profile (should see both)
	posts, err := GetUserPosts(1, 1)
	if err != nil {
		t.Fatalf("GetUserPosts failed: %v", err)
	}
	if len(posts) != 2 {
		t.Errorf("Expected 2 posts for user 1 (own profile), got %d", len(posts))
	}

	// Test GetUserPosts for User 2 viewing User 1 (should see 1 public post)
	posts2, err := GetUserPosts(1, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(posts2) != 1 {
		t.Errorf("Expected 1 post for user 2 viewing user 1, got %d", len(posts2))
	}
}

func TestAddComment(t *testing.T) {
	setupTestDB(t)
	defer database.DB.Close()

	// Insert User & Post
	_, err := database.DB.Exec(`INSERT INTO users (username, email, password_hash) VALUES ('u', 'e', 'h')`)
	if err != nil {
		t.Fatal(err)
	}
	_, err = database.DB.Exec(`INSERT INTO posts (user_id, title, content) VALUES (1, 'T', 'C')`)
	if err != nil {
		t.Fatal(err)
	}

	// Test AddComment
	commentID, err := AddComment(1, 1, "Nice post", "img.png")
	if err != nil {
		t.Fatalf("AddComment failed: %v", err)
	}
	if commentID <= 0 {
		t.Errorf("Expected valid comment ID, got %d", commentID)
	}

	// Verify insertion
	var content string
	err = database.DB.QueryRow("SELECT content FROM comments WHERE id = ?", commentID).Scan(&content)
	if err != nil {
		t.Fatalf("Failed to verify comment: %v", err)
	}
	if content != "Nice post" {
		t.Errorf("Expected content 'Nice post', got '%s'", content)
	}
}

func TestGetCommentByID(t *testing.T) {
	setupTestDB(t)
	defer database.DB.Close()

	// Insert User
	_, err := database.DB.Exec(`INSERT INTO users (username, email, password_hash, first_name, last_name, avatar) VALUES ('alice', 'alice@example.com', 'hash', 'Alice', 'Wonderland', 'avatar.png')`)
	if err != nil {
		t.Fatalf("Failed to insert user: %v", err)
	}

	// Insert Post
	_, err = database.DB.Exec(`INSERT INTO posts (user_id, title, content) VALUES (1, 'Test Post', 'Content')`)
	if err != nil {
		t.Fatalf("Failed to insert post: %v", err)
	}

	// Insert Comment
	_, err = database.DB.Exec(`INSERT INTO comments (user_id, post_id, content, image_url) VALUES (1, 1, 'Test Comment', 'comment.jpg')`)
	if err != nil {
		t.Fatalf("Failed to insert comment: %v", err)
	}

	// Test GetCommentByID
	comment, err := GetCommentByID(1)
	if err != nil {
		t.Fatalf("GetCommentByID failed: %v", err)
	}

	if comment["username"] != "alice" {
		t.Errorf("Expected username alice, got %v", comment["username"])
	}
	if comment["avatar"] != "avatar.png" {
		t.Errorf("Expected avatar avatar.png, got %v", comment["avatar"])
	}
	if comment["text"] != "Test Comment" {
		t.Errorf("Expected text Test Comment, got %v", comment["text"])
	}
	if comment["image_path"] != "comment.jpg" {
		t.Errorf("Expected image_path comment.jpg, got %v", comment["image_path"])
	}
}

func TestGetCommentsByPostID(t *testing.T) {
	setupTestDB(t)
	defer database.DB.Close()

	// Insert User
	_, err := database.DB.Exec(`INSERT INTO users (username, email, password_hash, first_name, last_name, avatar) VALUES ('bob', 'bob@example.com', 'hash', 'Bob', 'Builder', 'bob.png')`)
	if err != nil {
		t.Fatalf("Failed to insert user: %v", err)
	}

	// Insert Post
	_, err = database.DB.Exec(`INSERT INTO posts (user_id, title, content) VALUES (1, 'Bob Post', 'Building things')`)
	if err != nil {
		t.Fatalf("Failed to insert post: %v", err)
	}

	// Insert Comments
	_, err = database.DB.Exec(`INSERT INTO comments (user_id, post_id, content) VALUES (1, 1, 'Comment 1')`)
	if err != nil {
		t.Fatalf("Failed to insert comment 1: %v", err)
	}
	_, err = database.DB.Exec(`INSERT INTO comments (user_id, post_id, content) VALUES (1, 1, 'Comment 2')`)
	if err != nil {
		t.Fatalf("Failed to insert comment 2: %v", err)
	}

	// Test GetCommentsByPostID
	comments, err := GetCommentsByPostID(1)
	if err != nil {
		t.Fatalf("GetCommentsByPostID failed: %v", err)
	}

	if len(comments) != 2 {
		t.Errorf("Expected 2 comments, got %d", len(comments))
	}
	if len(comments) > 0 {
		if comments[0]["username"] != "bob" {
			t.Errorf("Expected username bob, got %v", comments[0]["username"])
		}
		if comments[0]["avatar"] != "bob.png" {
			t.Errorf("Expected avatar bob.png, got %v", comments[0]["avatar"])
		}
	}
}
