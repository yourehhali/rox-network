# API Documentation

- Base URL via API Gateway: http://localhost:8080
- Authentication: cookie-based session; credentials: include on requests
- Content types: JSON unless specified as multipart/form-data

## Auth Service (8081)

- POST /login
  - Body: { "username": string, "password": string }
  - Responses: 200 { status, user_id } | 401 error
  - Handler: [handlers.go](file:///Users/macadmin/Desktop/my-area/repos-to-rust/rox-network/backend/services/auth-service/internal/handlers/handlers.go#L26-L82)
- POST /register (multipart/form-data)
  - Fields: username, first_name, last_name, email, password, confirmpassword, gender, nickname?, about?, age, avatar?
  - Responses: 200 { status, message } | 4xx validation errors
  - Handler: [handlers.go](file:///Users/macadmin/Desktop/my-area/repos-to-rust/rox-network/backend/services/auth-service/internal/handlers/handlers.go#L84-L206)
- POST /logout
  - Clears session cookie
  - Handler: [handlers.go](file:///Users/macadmin/Desktop/my-area/repos-to-rust/rox-network/backend/services/auth-service/internal/handlers/handlers.go#L208-L219)
- GET /checkstate
  - Response: { authenticated: boolean, userID? }
  - Handler: [handlers.go](file:///Users/macadmin/Desktop/my-area/repos-to-rust/rox-network/backend/services/auth-service/internal/handlers/handlers.go#L221-L245)

## User Service (8082)

- GET /profile?id=:userId
  - Returns profile, followers/following counts, friend/pending flags
  - Handler: [handlers.go](file:///Users/macadmin/Desktop/my-area/repos-to-rust/rox-network/backend/services/user-service/internal/handlers/handlers.go#L19-L85)
- POST /toggle-follow
  - Body: { followed_id: number }
  - Response: { following, status, followers_Count, following_Count }
  - Handler: [handlers.go](file:///Users/macadmin/Desktop/my-area/repos-to-rust/rox-network/backend/services/user-service/internal/handlers/handlers.go#L98-L192)
- POST /request_follow
  - Body: { sender_id, notif_id, status: "accept"|"reject" }
  - Handler: [handlers.go](file:///Users/macadmin/Desktop/my-area/repos-to-rust/rox-network/backend/services/user-service/internal/handlers/handlers.go#L195-L252)
- GET /Friends
  - Returns accepted friends list
  - Handler: [handlers.go](file:///Users/macadmin/Desktop/my-area/repos-to-rust/rox-network/backend/services/user-service/internal/handlers/handlers.go#L267-L282)
- GET /users-sug
  - Suggested users to follow
  - Handler: [handlers.go](file:///Users/macadmin/Desktop/my-area/repos-to-rust/rox-network/backend/services/user-service/internal/handlers/handlers.go#L284-L299)
- POST /update-privacy
  - Body: { is_private: 0|1 }
  - Handler: [handlers.go](file:///Users/macadmin/Desktop/my-area/repos-to-rust/rox-network/backend/services/user-service/internal/handlers/handlers.go#L310-L341)
- POST /update-profile (multipart/form-data)
  - Fields: first_name, last_name, nickname?, about?, is_private?, avatar?
  - Handler: [handlers.go](file:///Users/macadmin/Desktop/my-area/repos-to-rust/rox-network/backend/services/user-service/internal/handlers/handlers.go#L343-L436)
- GET /followers?id=:userId
  - Handler: [handlers.go](file:///Users/macadmin/Desktop/my-area/repos-to-rust/rox-network/backend/services/user-service/internal/handlers/handlers.go#L438-L458)
- GET /following?id=:userId
  - Handler: [handlers.go](file:///Users/macadmin/Desktop/my-area/repos-to-rust/rox-network/backend/services/user-service/internal/handlers/handlers.go#L460-L480)

## Post Service (8083)

- POST /Createpost (multipart/form-data)
  - Fields: title?, content, privacy (0 public, 1 almost, 2 private), allowed_users? (comma IDs), image?
  - Handler: [handlers.go](file:///Users/macadmin/Desktop/my-area/repos-to-rust/rox-network/backend/services/post-service/internal/handlers/handlers.go#L20-L111)
- GET /GetPost?id=:postId
  - Returns single post with visibility checks
  - Handler: [handlers.go](file:///Users/macadmin/Desktop/my-area/repos-to-rust/rox-network/backend/services/post-service/internal/handlers/handlers.go#L113-L144)
- GET /getposts?cursor=:offset&limit=:n
  - Returns feed for current user
  - Handler: [handlers.go](file:///Users/macadmin/Desktop/my-area/repos-to-rust/rox-network/backend/services/post-service/internal/handlers/handlers.go#L146-L181)
- GET /GetCUser/:id?cursor=&limit=
  - Returns posts by user with viewer-based visibility
  - Handler: [handlers.go](file:///Users/macadmin/Desktop/my-area/repos-to-rust/rox-network/backend/services/post-service/internal/handlers/handlers.go#L183-L230)
- POST /reactions
  - Body: { post_id, reaction: "like"|"dislike" }
  - Handler: [handlers.go](file:///Users/macadmin/Desktop/my-area/repos-to-rust/rox-network/backend/services/post-service/internal/handlers/handlers.go#L232-L290)
- /posts/:id/comments
  - POST: multipart or JSON; adds a comment with optional image
  - GET: list comments with pagination
  - Handler: [handlers.go](file:///Users/macadmin/Desktop/my-area/repos-to-rust/rox-network/backend/services/post-service/internal/handlers/handlers.go#L292-L410)

## Media Service (8084)

- GET /uploads/:filename
  - Auth required; enforces post privacy, group membership, or user avatar access
  - Service setup: [main.go](file:///Users/macadmin/Desktop/my-area/repos-to-rust/rox-network/backend/services/media-service/cmd/api/main.go#L30-L37)
  - Handler: [handlers.go](file:///Users/macadmin/Desktop/my-area/repos-to-rust/rox-network/backend/services/media-service/internal/handlers/handlers.go#L22-L113)

## Chat Service (8085)

- WS /chat/ws
  - Auth via session cookie or initial {"type":"auth","session_token": "..."}
  - Message types:
    - private_message: { receiver_id, content, is_emoji }
    - group_message: { group_id, content, is_emoji }
  - Handler: [handlers.go](file:///Users/macadmin/Desktop/my-area/repos-to-rust/rox-network/backend/services/chat-service/internal/handlers/handlers.go#L24-L77)
- GET /chat/private/history?user_id=:id&limit=:n&offset=:m
  - Handler: [handlers.go](file:///Users/macadmin/Desktop/my-area/repos-to-rust/rox-network/backend/services/chat-service/internal/handlers/handlers.go#L185-L239)
- GET /chat/group/history?group_id=:id&limit=:n&offset=:m
  - Handler: [handlers.go](file:///Users/macadmin/Desktop/my-area/repos-to-rust/rox-network/backend/services/chat-service/internal/handlers/handlers.go#L242-L296)

## Group Service (8086)

- POST /Create_Group
  - Body: { Title, Description }
  - Handler: [handlers.go](file:///Users/macadmin/Desktop/my-area/repos-to-rust/rox-network/backend/services/group-service/internal/handlers/handlers.go#L31-L67)
- GET /Get_Groups
  - Lists groups with current user status
  - Handler: [handlers.go](file:///Users/macadmin/Desktop/my-area/repos-to-rust/rox-network/backend/services/group-service/internal/handlers/handlers.go#L70-L85)
- POST /join
  - Body: { group_id }
  - Handler: [handlers.go](file:///Users/macadmin/Desktop/my-area/repos-to-rust/rox-network/backend/services/group-service/internal/handlers/handlers.go#L88-L133)
- POST /accept-reject-join
  - Body: { notif_id, group_id, user_id, action: "accept"|"reject" }
  - Handler: [handlers.go](file:///Users/macadmin/Desktop/my-area/repos-to-rust/rox-network/backend/services/group-service/internal/handlers/handlers.go#L136-L187)
- POST /group/invite/respond
  - Body: { group_id, notif_id, status: "accept"|"reject" }
  - Handler: [handlers.go](file:///Users/macadmin/Desktop/my-area/repos-to-rust/rox-network/backend/services/group-service/internal/handlers/handlers.go#L190-L223)
- GET /Get_Group_By_ID?id=:groupId or /Get_Group_By_ID/:id
  - Returns group details, members, and current user status
  - Handler: [handlers.go](file:///Users/macadmin/Desktop/my-area/repos-to-rust/rox-network/backend/services/group-service/internal/handlers/handlers.go#L226-L274)
- POST /CreateGroupPost (multipart/form-data)
  - Fields: title, content, group_id, image?
  - Handler: [handlers.go](file:///Users/macadmin/Desktop/my-area/repos-to-rust/rox-network/backend/services/group-service/internal/handlers/handlers.go#L277-L333)
- GET /GetGroupPosts?group_id=:id&cursor=:offset&limit=:n
  - Requires group membership
  - Handler: [handlers.go](file:///Users/macadmin/Desktop/my-area/repos-to-rust/rox-network/backend/services/group-service/internal/handlers/handlers.go#L336-L377)
- POST /groups/not-member
  - Body: { user_id } to list groups where user is not a member (for invites)
  - Handler: [handlers.go](file:///Users/macadmin/Desktop/my-area/repos-to-rust/rox-network/backend/services/group-service/internal/handlers/handlers.go#L380-L401)
- POST /groups/invite
  - Body: { group_id, user_id }
  - Handler: [handlers.go](file:///Users/macadmin/Desktop/my-area/repos-to-rust/rox-network/backend/services/group-service/internal/handlers/handlers.go#L404-L451)
- POST /create_event
  - Body: { group_id, title, description, event_time }
  - Handler: [handlers.go](file:///Users/macadmin/Desktop/my-area/repos-to-rust/rox-network/backend/services/group-service/internal/handlers/handlers.go#L454-L522)
- GET /group-events?group_id=:id&cursor=:offset&limit=:n
  - Handler: [handlers.go](file:///Users/macadmin/Desktop/my-area/repos-to-rust/rox-network/backend/services/group-service/internal/handlers/handlers.go#L525-L553)
- POST /event-invite-response
  - Body: { event_id, notif_id, status }
  - Handler: [handlers.go](file:///Users/macadmin/Desktop/my-area/repos-to-rust/rox-network/backend/services/group-service/internal/handlers/handlers.go#L556-L589)

## Notification Service (8087)

- WS /ws
  - Requires session cookie; subscribes user to notifications
  - Handler: [handlers.go](file:///Users/macadmin/Desktop/my-area/repos-to-rust/rox-network/backend/services/notification-service/internal/handlers/handlers.go#L20-L49)
- GET /notifications
  - Returns unread notifications enriched with sender info
  - Handler: [handlers.go](file:///Users/macadmin/Desktop/my-area/repos-to-rust/rox-network/backend/services/notification-service/internal/handlers/handlers.go#L107-L137)
- POST /notifications/read
  - Body: { notification_id }
  - Handler: [handlers.go](file:///Users/macadmin/Desktop/my-area/repos-to-rust/rox-network/backend/services/notification-service/internal/handlers/handlers.go#L139-L168)
- POST /internal/notification
  - Internal API for services to send/broadcast notifications
  - Handler: [handlers.go](file:///Users/macadmin/Desktop/my-area/repos-to-rust/rox-network/backend/services/notification-service/internal/handlers/handlers.go#L51-L105)

## Gateway Routing

- Rules in [api-gateway](file:///Users/macadmin/Desktop/my-area/repos-to-rust/rox-network/backend/api-gateway/cmd/api/main.go#L48-L120)
- CORS headers are set per request; cookies forwarded

