# Rox Network

- Microservices social network with API Gateway and a Next.js frontend
- Services: auth, user, post, media, chat, group, notification
- Database: SQLite via migrate; uploads served from media-service

## Architecture

- API Gateway on port 8080 proxies requests to services
- Service ports: 8081 auth, 8082 user, 8083 post, 8084 media, 8085 chat, 8086 group, 8087 notification
- Frontend on port 3000 talks to API Gateway at NEXT_PUBLIC_API_URL

## Quick Start

- Start all services and frontend
  - Run [start-all.sh](file:///Users/macadmin/Desktop/my-area/repos-to-rust/rox-network/start-all.sh)
- Environment variables
  - DB_PATH, MIGRATIONS_PATH, UPLOAD_DIR for backend services
  - NEXT_PUBLIC_API_URL for frontend (defaults to http://localhost:8080)

## API Overview

- Full endpoint reference in [API.md](file:///Users/macadmin/Desktop/my-area/repos-to-rust/rox-network/API.md)
- Gateway routing rules in [main.go](file:///Users/macadmin/Desktop/my-area/repos-to-rust/rox-network/backend/api-gateway/cmd/api/main.go#L48-L120)

## Frontend

- Next.js 16 with app router
- Setup and usage documented in [README.md](file:///Users/macadmin/Desktop/my-area/repos-to-rust/rox-network/frontend/README.md)

## Development Notes

- Migrations live in [sqlite migrations](file:///Users/macadmin/Desktop/my-area/repos-to-rust/rox-network/backend/pkg/db/migrations/sqlite)
- Media files stored under backend/uploads and served via media-service
- WebSockets:
  - Chat: /chat/ws via chat-service
  - Notifications: /ws via notification-service

