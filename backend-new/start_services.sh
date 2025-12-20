#!/bin/bash

# Function to kill all background processes on exit
cleanup() {
    echo "Stopping all services..."
    kill $(jobs -p)
    exit
}

trap cleanup SIGINT SIGTERM

echo "Starting Social Network Microservices..."

# Ensure uploads directory exists
mkdir -p uploads

# Start API Gateway
echo "Starting API Gateway..."
cd api-gateway
go run cmd/api/main.go &
cd ..

# Start Auth Service
echo "Starting Auth Service..."
cd services/auth-service
go run cmd/api/main.go &
cd ../..

# Start User Service
echo "Starting User Service..."
cd services/user-service
go run cmd/api/main.go &
cd ../..

# Start Post Service
echo "Starting Post Service..."
cd services/post-service
go run cmd/api/main.go &
cd ../..

# Start Media Service
echo "Starting Media Service..."
cd services/media-service
go run cmd/api/main.go &
cd ../..

# Start Chat Service
echo "Starting Chat Service..."
cd services/chat-service
go run cmd/api/main.go &
cd ../..

# Start Group Service
echo "Starting Group Service..."
cd services/group-service
go run cmd/api/main.go &
cd ../..

echo "All services started!"
echo "API Gateway running on http://localhost:8080"
echo "Press Ctrl+C to stop."

wait
