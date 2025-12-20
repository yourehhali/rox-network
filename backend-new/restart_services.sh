#!/bin/bash
pkill -9 -f "go run cmd/api/main.go"
pkill -9 -f "api-gateway"
# Kill any lingering processes on ports
lsof -ti:8080 | xargs kill -9
lsof -ti:8081 | xargs kill -9
lsof -ti:8082 | xargs kill -9
lsof -ti:8083 | xargs kill -9
lsof -ti:8084 | xargs kill -9
lsof -ti:8085 | xargs kill -9
lsof -ti:8086 | xargs kill -9

echo "Stopped services."
sleep 5

# Start services using subshells with nohup
(cd services/auth-service && nohup go run cmd/api/main.go > auth.log 2>&1 &)
echo "Started Auth Service"

(cd services/user-service && nohup go run cmd/api/main.go > user.log 2>&1 &)
echo "Started User Service"

(cd services/post-service && nohup go run cmd/api/main.go > post.log 2>&1 &)
echo "Started Post Service"

(cd services/media-service && nohup go run cmd/api/main.go > media.log 2>&1 &)
echo "Started Media Service"

(cd services/group-service && nohup go run cmd/api/main.go > group.log 2>&1 &)
echo "Started Group Service"

(cd services/chat-service && nohup go run cmd/api/main.go > chat.log 2>&1 &)
echo "Started Chat Service"

(cd api-gateway && nohup go run cmd/api/main.go > gateway.log 2>&1 &)
echo "Started API Gateway"

echo "Waiting for services to initialize..."
sleep 5

echo "Running Database Seeder..."
(cd seeder && go run main.go)

echo "All services started in background."

