# Social Network Backend (Microservices)

This repository contains the backend for the Social Network application, refactored into a microservices architecture.

## Structure

The backend is organized into the following services:

- **api-gateway** (Port 8080): The entry point for all client requests. It routes requests to the appropriate microservice.
- **services/auth-service** (Port 8081): Handles user authentication (login, register, logout, session check).
- **services/user-service** (Port 8082): Manages user profiles, privacy, followers, and friends.
- **services/post-service** (Port 8083): Handles posts, comments, and reactions.
- **services/media-service** (Port 8084): Serves uploaded media files.
- **pkg**: Shared code (models, database connections, middleware, utilities).

## Prerequisites

- Go 1.24+
- SQLite3

## Getting Started

### 1. Database Setup

The services currently share a single SQLite database for migration simplicity.
Ensure the database file exists at `backend/database/forum.db` (relative to the old structure) or update the paths in `cmd/api/main.go` of each service.

### 2. Running the Services

You can run each service in a separate terminal.

**Option A: Using Go Workspace (Recommended)**

The project is set up as a Go workspace. You can run services from the `backend-new` directory.

Terminal 1 (Auth Service):
```bash
cd services/auth-service
go run cmd/api/main.go
```

Terminal 2 (User Service):
```bash
cd services/user-service
go run cmd/api/main.go
```

Terminal 3 (Post Service):
```bash
cd services/post-service
go run cmd/api/main.go
```

Terminal 4 (Media Service):
```bash
cd services/media-service
go run cmd/api/main.go
```

Terminal 5 (API Gateway):
```bash
cd api-gateway
go run cmd/api/main.go
```

**Option B: Build and Run**

Build each service:
```bash
cd services/auth-service && go build -o auth-service cmd/api/main.go
cd ../user-service && go build -o user-service cmd/api/main.go
cd ../post-service && go build -o post-service cmd/api/main.go
cd ../media-service && go build -o media-service cmd/api/main.go
cd ../../api-gateway && go build -o api-gateway cmd/api/main.go
```

Then run the binaries.

### 3. API Usage

All requests should be sent to the **API Gateway** at `http://localhost:8080`.

**Auth Endpoints:**
- POST `/login`
- POST `/register`
- GET `/logout`
- GET `/checkstate`

**User Endpoints:**
- GET `/profile/?id={id}`
- POST `/toggle-follow`
- GET `/Friends`
- GET `/users-sug`
- POST `/update-privacy`

**Post Endpoints:**
- POST `/Createpost`
- GET `/getposts`
- GET `/GetCUser/?id={id}`
- POST `/posts/{id}/comments`
- POST `/reactions`

**Media Endpoints:**
- GET `/uploads/{filename}`

## Development

- Shared code is in `pkg`. If you modify it, all services using it will be affected.
- `go.work` file in the root connects all modules.
