#!/bin/bash
pkill -f "go run cmd/api/main.go"
sleep 2

(cd backend-new/services/auth-service && go run cmd/api/main.go > auth.log 2>&1) &
echo "Auth Service started"
(cd backend-new/services/user-service && go run cmd/api/main.go > user.log 2>&1) &
echo "User Service started"
(cd backend-new/services/post-service && go run cmd/api/main.go > post.log 2>&1) &
echo "Post Service started"
(cd backend-new/services/media-service && go run cmd/api/main.go > media.log 2>&1) &
echo "Media Service started"
(cd backend-new/services/chat-service && go run cmd/api/main.go > chat.log 2>&1) &
echo "Chat Service started"
(cd backend-new/services/group-service && go run cmd/api/main.go > group.log 2>&1) &
echo "Group Service started"
sleep 5
(cd backend-new/api-gateway && go run cmd/api/main.go > gateway.log 2>&1) &
echo "API Gateway started"
sleep 5
