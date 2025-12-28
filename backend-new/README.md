# Social Network Backend (Microservices)

This is the backend for the Social Network application, architected as a set of microservices written in Go.

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

- Go 1.24+
- SQLite3
- Bash (for the start script)

## Directory Structure

```
backend-new/
├── api-gateway/       # Reverse proxy module
├── pkg/               # Shared code (database, models, middleware)
├── services/          # Microservices
│   ├── auth-service/
│   ├── chat-service/
│   ├── group-service/
│   ├── media-service/
│   ├── post-service/
│   └── user-service/
├── migrations/        # SQL migrations
├── uploads/           # Shared storage for uploaded files (created at runtime)
└── go.mod             # Root module definition
```

## Running the Application

### Option 1: Run Everything (Recommended)

Use the `start-all.sh` script in the project root to run both backend and frontend:

```bash
# From project root
./start-all.sh
```

### Option 2: Run Backend Only

If you want to run the backend services manually, you need to set up the environment variables and run each service.

1.  **Navigate to backend directory:**
    ```bash
    cd backend-new
    ```

2.  **Set Environment Variables:**
    ```bash
    export DB_PATH=$(pwd)/database/social.db
    export MIGRATIONS_PATH=$(pwd)/pkg/db/migrations/sqlite
    export UPLOAD_DIR=$(pwd)/uploads
    ```
    *Note: Ensure the `database` and `uploads` directories exist or are writable.*

3.  **Start Services:**
    You can run each service in a separate terminal or background them:

    ```bash
    go run services/auth-service/cmd/api/main.go &
    go run services/user-service/cmd/api/main.go &
    go run services/post-service/cmd/api/main.go &
    go run services/media-service/cmd/api/main.go &
    go run services/chat-service/cmd/api/main.go &
    go run services/group-service/cmd/api/main.go &
    ```

4.  **Start API Gateway:**
    ```bash
    go run api-gateway/cmd/api/main.go
    ```

## API Endpoints

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
- **Modules**: The project uses a root `go.mod` for shared code and services, and a separate `go.mod` for the `api-gateway`.
- **Database**: Uses a shared SQLite database for simplicity. The schema is managed via migrations in `pkg/db/migrations`.
