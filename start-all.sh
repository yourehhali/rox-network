#!/bin/bash

# Function to kill processes on exit
cleanup() {
    echo "Stopping all services..."
    lsof -ti:3000,8080,8081,8082,8083,8084,8085,8086 | xargs kill -9 2>/dev/null
    echo "All services stopped."
    exit
}

# Trap SIGINT (Ctrl+C)
trap cleanup SIGINT

echo "Starting Backend Services..."

cd backend-new

# Set DB_PATH for all services
# Ensure the path is absolute
export DB_PATH=$(pwd)/database/forum.db
export MIGRATIONS_PATH=$(pwd)/pkg/db/migrations/sqlite
export UPLOAD_DIR=$(pwd)/uploads
echo "Database Path: $DB_PATH"
echo "Uploads Path: $UPLOAD_DIR"
echo "Migrations Path: $MIGRATIONS_PATH"

# Start services
echo "Starting Auth Service (8081)..."
go run services/auth-service/cmd/api/main.go &

echo "Starting User Service (8082)..."
go run services/user-service/cmd/api/main.go &

echo "Starting Post Service (8083)..."
go run services/post-service/cmd/api/main.go &

echo "Starting Media Service (8084)..."
go run services/media-service/cmd/api/main.go &

echo "Starting Chat Service (8085)..."
go run services/chat-service/cmd/api/main.go &

echo "Starting Group Service (8086)..."
go run services/group-service/cmd/api/main.go &

# Wait a bit for services to initialize
sleep 3

echo "Starting API Gateway (8080)..."
go run api-gateway/cmd/api/main.go &

cd ..

echo "Starting Frontend (3000)..."
cd ui-frontend
# Using 'npm run dev' for development mode (hot reloading)
# If you prefer the production build, use 'npm run build && npm start'
npm run dev &

echo "All services started. Press Ctrl+C to stop."
wait
