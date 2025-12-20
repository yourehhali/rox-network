# Social Network Backend (Microservices)

This is the refactored backend for the Social Network application, using a microservices architecture.

## Architecture

The application is split into the following microservices:

1.  **API Gateway** (`:8080`): The single entry point for the frontend. Routes requests to appropriate services.
2.  **Auth Service** (`:8081`): Handles user registration, login, logout, and session verification.
3.  **User Service** (`:8082`): Manages user profiles, followers, and privacy settings.
4.  **Post Service** (`:8083`): Manages posts, comments, and likes.
5.  **Media Service** (`:8084`): Serves uploaded static files (images).
6.  **Chat Service** (`:8085`): Handles real-time chat (private & group) and notifications via WebSockets.
7.  **Group Service** (`:8086`): Manages groups and events.

## Prerequisites

- Go 1.20+
- SQLite3
- Bash (for the start script)

## Directory Structure

```
backend-new/
├── api-gateway/       # Reverse proxy
├── pkg/               # Shared code (database, models, middleware)
├── services/          # Microservices
│   ├── auth-service/
│   ├── chat-service/
│   ├── group-service/
│   ├── media-service/
│   ├── post-service/
│   └── user-service/
├── uploads/           # Shared storage for uploaded files
├── go.work            # Go workspace configuration
└── start_services.sh  # Script to run all services
```

## Running the Application

### 1. Initialize the Database

The services share a SQLite database located at `pkg/database/forum.db`. The `pkg` initialization automatically handles schema migrations (if configured).

### 2. Start All Services

Run the provided script to start all services in the background:

```bash
chmod +x start_services.sh
./start_services.sh
```

You should see output indicating all services have started.

### 3. API Endpoints

All requests should be sent to the **API Gateway** at `http://localhost:8080`.

**Auth:**
- `POST /register`
- `POST /login`
- `POST /logout`
- `GET /checkstate`

**User:**
- `GET /profile`
- `POST /update-privacy`
- `POST /follow`

**Posts:**
- `POST /create_post`
- `GET /posts`
- `POST /submit_comment`

**Groups:**
- `POST /Create_Group`
- `GET /Get_Groups`

**Chat & Notifications:**
- `WS /chat/ws` (WebSocket)
- `WS /ws` (Notification WebSocket)

## Development Notes

- **Shared Code**: Common logic (DB connection, Models, Middleware) is in the `pkg` directory.
- **Go Workspace**: The project uses `go.work` to manage local dependencies. If you add a new module, update `go.work`.
- **Database**: Currently using a shared SQLite database for simplicity during migration. Ideally, each service should have its own DB.
