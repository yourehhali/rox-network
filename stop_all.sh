#!/bin/bash

echo "Stopping all services..."

# Kill Frontend (Next.js) on port 3000
if lsof -ti:3000 >/dev/null; then
    echo "Stopping Frontend (Next.js) on port 3000..."
    lsof -ti:3000 | xargs kill -9
else
    echo "Frontend not running on port 3000."
fi

# Kill Frontend (Next.js) on port 3001
if lsof -ti:3001 >/dev/null; then
    echo "Stopping Frontend (Next.js) on port 3001..."
    lsof -ti:3001 | xargs kill -9
else
    echo "Frontend not running on port 3001."
fi

# Kill Backend Services by Name
echo "Stopping Go Backend Services..."
pkill -9 -f "go run cmd/api/main.go"
pkill -9 -f "api-gateway"

# Ensure Backend Ports are Clear
PORTS=(8080 8081 8082 8083 8084 8085 8086)
for PORT in "${PORTS[@]}"; do
    if lsof -ti:$PORT >/dev/null; then
        echo "Clearing port $PORT..."
        lsof -ti:$PORT | xargs kill -9
    fi
done

echo "All application instances have been shut down."
