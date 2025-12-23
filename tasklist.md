# Project Task List

Based on the codebase analysis and standard Social Network requirements.

## Frontend (Next.js)

### Authentication
- [x] Login Page
- [x] Register Page
- [x] Auth Context (Session Management)

### User Profile
- [x] Profile Page
- [x] Profile Card (User Info)
- [x] User's Posts Feed
- [x] Edit Profile (Privacy settings, etc.)

### Posts & Feed
- [x] Home Feed
- [x] Create Post (Modal)
- [x] Post Privacy (Public/Private/Almost Private)
- [x] Comments
- [x] Likes/Reactions (Basic)

### Social Features
- [x] Follow/Unfollow Users
- [x] Friend Suggestions
- [x] Followers/Following List

### Groups
- [x] List Groups
- [x] Create Group
- [x] Join/Request to Join Group
- [x] Group Page
- [x] Group Posts
- [x] Group Events (Create/View)
- [x] Group Chat (WebSocket integration)

### Chat
- [x] Real-time Private Chat (WebSocket)
- [x] Chat History
- [x] Group Chat Integration
- [x] Chat UI (Sidebar, Message Area)

### Notifications
- [x] Notification System (WebSocket)
- [x] Toast Notifications (UI)
- [x] Event/Group Invites

## Backend (Go Microservices)

### Services
- [x] **Auth Service**: Login, Register, Session
- [x] **User Service**: Profile, Followers, Privacy
- [x] **Post Service**: CRUD Posts, Comments, Likes
- [x] **Group Service**: Groups, Members, Events
- [x] **Chat Service**: WebSocket, History, Notifications
- [x] **Media Service**: File Uploads
- [x] **API Gateway**: Route proxying

### Infrastructure
- [x] SQLite Database
- [x] Database Migrations (`schema.sql`)
- [ ] **Dockerization** (Missing Dockerfiles & docker-compose.yml)
- [ ] **CI/CD** (No workflow files)

## Documentation & Testing
- [ ] **README.md** (Current file is corrupted/HTML)
- [ ] **USAGE.md** (Missing)
- [ ] **Unit Tests** (Partial coverage in `repository_test.go`)
- [ ] **Integration Tests**

## Known Issues / Refinements needed
- [ ] Fix `CreatePost` in Group Page (ensure FormData works correctly with `apiFetch`)
- [ ] Verify `toMediaUrl` handling for WebSocket URLs (Fixed in `api.js`?)
- [ ] Clean up root directory (remove `cookies.txt`, `plain ui` folder?)

