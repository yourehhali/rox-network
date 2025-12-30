#!/bin/bash

# Kill existing processes on ports 8080-8087
echo "Killing existing processes..."
lsof -ti:8080-8087 | xargs kill -9 2>/dev/null

# Set environment variables
export DB_PATH=$(pwd)/social-network.db
export UPLOAD_DIR=$(pwd)/uploads

# Create uploads directory if it doesn't exist
mkdir -p $UPLOAD_DIR

# Start services
echo "Starting services..."

# Auth Service (8081)
cd services/auth-service
go run cmd/api/main.go &
echo "Auth Service started on 8081"
cd ../..

# User Service (8082)
cd services/user-service
go run cmd/api/main.go &
echo "User Service started on 8082"
cd ../..

# Post Service (8083)
cd services/post-service
go run cmd/api/main.go &
echo "Post Service started on 8083"
cd ../..

# Media Service (8084)
cd services/media-service
go run cmd/api/main.go &
echo "Media Service started on 8084"
cd ../..

# Chat Service (8085)
cd services/chat-service
go run cmd/api/main.go &
echo "Chat Service started on 8085"
cd ../..

# Group Service (8086)
cd services/group-service
go run cmd/api/main.go &
echo "Group Service started on 8086"
cd ../..

# Notification Service (8087)
cd services/notification-service
go run cmd/api/main.go &
echo "Notification Service started on 8087"
cd ../..

# Wait for services to initialize
sleep 2

# API Gateway (8080)
cd api-gateway
go run cmd/api/main.go &
echo "API Gateway started on 8080"
cd ..

echo "All services started."
