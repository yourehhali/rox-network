#!/bin/bash

echo "Stopping all services..."

# 1. Kill Frontend (Next.js) by Port and specific command
echo "Stopping Frontend..."
# Kill process on port 3000 (Next.js)
if lsof -ti:3000 >/dev/null; then
    echo "Killing process on port 3000..."
    lsof -ti:3000 | xargs kill -9
fi
# Fallback for npm/next just in case
pkill -f "npm run dev"
pkill -f "next-server"

# 2. Kill Go Backend Services
echo "Stopping Go Backend Services..."

# We primarily rely on ports 8080-8086 to identify our services.
# This avoids killing random processes like browser "Media Service" helpers
# which might match generic names like "media-service".

PORTS=(8080 8081 8082 8083 8084 8085 8086)
for PORT in "${PORTS[@]}"; do
    if lsof -ti:$PORT >/dev/null; then
        echo "Clearing port $PORT..."
        lsof -ti:$PORT | xargs kill -9
    fi
done

# Cleanup any lingering go run processes that might strictly match our path
# Be very specific to avoid collateral damage
pkill -f "go run cmd/api/main.go"

echo "All application instances have been shut down."
