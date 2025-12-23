# Rox Network

A full-featured social network application built with a Microservices architecture.

## 🚀 Tech Stack

### Backend
- **Language:** Go (Golang) 1.24+
- **Architecture:** Microservices
- **Database:** SQLite (Shared via Docker Volume)
- **Communication:** REST API

### Frontend
- **Framework:** Next.js 14 (App Router)
- **Styling:** CSS Modules

### Infrastructure
- **Containerization:** Docker & Docker Compose
- **Gateway:** Custom Go API Gateway

## 🏗 Project Structure

```
rox-network/
├── backend-new/            # Backend Microservices
│   ├── api-gateway/        # Entry point for all requests
│   ├── pkg/                # Shared Go packages (DB, Middleware, Models)
│   ├── services/           # Individual microservices
│   │   ├── auth-service/   # Authentication & Session management
│   │   ├── chat-service/   # Real-time chat functionality
│   │   ├── group-service/  # Group management
│   │   ├── media-service/  # File uploads & serving
│   │   ├── post-service/   # Posts, comments & feed
│   │   └── user-service/   # User profiles & relationships
│   ├── seeder/             # Database seeding tool
├── ui-frontend/            # Next.js Frontend Application
├── docker-compose.yml      # Docker orchestration config
└── README.md               # Project documentation
```

## 🛠 Getting Started

### Prerequisites
- [Docker](https://www.docker.com/get-started) installed on your machine.

### Running the Application

1. **Clone the repository** (if you haven't already):
   ```bash
   git clone <repository-url>
   cd rox-network
   ```

2. **Start the services** using Docker Compose:
   ```bash
   docker-compose up --build
   ```
   This command will:
   - Build all backend microservices images.
   - Build the frontend image.
   - Create shared volumes for the database and uploads.
   - Start all services.

3. **Access the Application**:
   - **Frontend:** [http://localhost:3000](http://localhost:3000)
   - **API Gateway:** [http://localhost:8080](http://localhost:8080)

### 🛑 Stopping the Application

To stop the containers:
```bash
docker-compose down
```

## 🔌 Services & Ports

| Service | Internal Port | Description |
|---------|---------------|-------------|
| **Frontend** | 3000 | Next.js Web UI |
| **API Gateway** | 8080 | Main Entry Point |
| **Auth Service** | 8081 | Authentication |
| **User Service** | 8082 | User Management |
| **Post Service** | 8083 | Posts & Feed |
| **Media Service** | 8084 | File Uploads |
| **Chat Service** | 8085 | Messaging |
| **Group Service** | 8086 | Groups |

## 📦 Database

The project uses **SQLite**.
- In Docker: The database file is persisted in a volume mapped to `/app/database/forum.db` inside containers.
- Shared Access: All microservices share the same database file (monolithic database pattern for simplicity in this microservices setup).

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request
