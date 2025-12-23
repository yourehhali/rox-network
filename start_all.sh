#!/bin/bash

# Kill any existing processes to ensure a clean start
./stop_all.sh

echo "Starting all services in DEV mode..."

# Define common environment variables for services if needed
# Note: Services usually default to ../../database/forum.db if DB_PATH is not set,
# which is correct when running from the service directory.
# But we will set them explicitly to be safe.

ABS_ROOT=$(pwd)
DB_PATH="$ABS_ROOT/backend-new/database/forum.db"
UPLOAD_DIR="$ABS_ROOT/backend-new/uploads"

# 1. Start Auth Service (Port 8081)
echo "Starting Auth Service on :8081..."
(cd backend-new/services/auth-service && DB_PATH="$DB_PATH" go run cmd/api/main.go > "$ABS_ROOT/auth.log" 2>&1) &

# 2. Start User Service (Port 8082)
echo "Starting User Service on :8082..."
(cd backend-new/services/user-service && DB_PATH="$DB_PATH" go run cmd/api/main.go > "$ABS_ROOT/user.log" 2>&1) &

# 3. Start Post Service (Port 8083)
echo "Starting Post Service on :8083..."
(cd backend-new/services/post-service && DB_PATH="$DB_PATH" go run cmd/api/main.go > "$ABS_ROOT/post.log" 2>&1) &

# 4. Start Media Service (Port 8084)
echo "Starting Media Service on :8084..."
(cd backend-new/services/media-service && UPLOAD_DIR="$UPLOAD_DIR" go run cmd/api/main.go > "$ABS_ROOT/media.log" 2>&1) &

# 5. Start Chat Service (Port 8085)
echo "Starting Chat Service on :8085..."
(cd backend-new/services/chat-service && DB_PATH="$DB_PATH" go run cmd/api/main.go > "$ABS_ROOT/chat.log" 2>&1) &

# 6. Start Group Service (Port 8086)
echo "Starting Group Service on :8086..."
(cd backend-new/services/group-service && DB_PATH="$DB_PATH" go run cmd/api/main.go > "$ABS_ROOT/group.log" 2>&1) &

# Wait a bit for services to initialize
sleep 5

# 7. Start API Gateway (Port 8080)
echo "Starting API Gateway on :8080..."
(cd backend-new/api-gateway && \
 AUTH_SERVICE_URL=http://localhost:8081 \
 USER_SERVICE_URL=http://localhost:8082 \
 POST_SERVICE_URL=http://localhost:8083 \
 MEDIA_SERVICE_URL=http://localhost:8084 \
 CHAT_SERVICE_URL=http://localhost:8085 \
 GROUP_SERVICE_URL=http://localhost:8086 \
 go run cmd/api/main.go > "$ABS_ROOT/gateway.log" 2>&1) &

# 8. Start Frontend (Port 3000)
echo "Starting Frontend on :3000..."
(cd ui-frontend && npm run dev > "$ABS_ROOT/frontend.log" 2>&1) &

echo "All services started!"
echo "Backend Gateway: http://localhost:8080"
echo "Frontend: http://localhost:3000"
echo "Logs are being written to $ABS_ROOT/*.log"
