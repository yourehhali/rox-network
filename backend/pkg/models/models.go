package models

type User struct {
	ID        int     `json:"id"`
	Username  string  `json:"username"`
	Email     string  `json:"email"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Age       *int    `json:"age,omitempty"`
	Gender    *string `json:"gender,omitempty"`
	Nickname  *string `json:"nickname,omitempty"`
	About     *string `json:"about,omitempty"`
	Avatar    *string `json:"avatar,omitempty"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	IsPrivate *int    `json:"is_private"`
}

type Post struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Posts struct {
	ID            int     `json:"id"`
	UserID        int     `json:"user_id"`
	UserName      string  `json:"username"`
	FullName      string  `json:"full_name"`
	Title         string  `json:"title"`
	Content       string  `json:"content"`
	ImagePath     *string `json:"image_path"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
	LikesCount    int     `json:"likes_count"`
	DislikesCount int     `json:"dislikes_count"`
	UserReaction  *string `json:"userReaction"`
	Privacy       int     `json:"privacy"`
}

type Sugg struct {
	UserID    int     `json:"user_id"`
	UserName  string  `json:"username"`
	FullName  string  `json:"full_name"`
	ImagePath *string `json:"image_path"`
}

type Notification struct {
	ID                int       `json:"id"`
	UserID            int       `json:"user_id"`
	SenderID          int       `json:"sender_id"`
	Type              string    `json:"type"`
	Message           string    `json:"message"`
	State             string    `json:"state"`
	CreatedAt         string    `json:"created_at"` 
	EventID           int       `json:"event_id"`
	Sender            User      `json:"sender"`
	ReceiverIsPrivate bool      `json:"receiver_is_private"`
	GroupTitle        string    `json:"group_title"`
	GroupID           int       `json:"group_id"`
}

type Message struct {
	ID         int    `json:"id"`
	SenderID   int    `json:"sender_id"`
	ReceiverID *int   `json:"receiver_id,omitempty"` 
	GroupID    *int   `json:"group_id,omitempty"`    
	Content    string `json:"content"`
	IsEmoji    bool   `json:"is_emoji"`
	CreatedAt  string `json:"created_at"`
	Sender     User   `json:"sender"`
}
