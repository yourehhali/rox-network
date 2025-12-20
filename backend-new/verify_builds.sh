#!/bin/bash
echo "Verifying builds..."

echo "Building pkg..."
cd pkg && go build ./... && cd ..

echo "Building api-gateway..."
cd api-gateway && go build -o bin/api-gateway cmd/api/main.go && cd ..

echo "Building auth-service..."
cd services/auth-service && go build -o bin/auth-service cmd/api/main.go && cd ../..

echo "Building user-service..."
cd services/user-service && go build -o bin/user-service cmd/api/main.go && cd ../..

echo "Building post-service..."
cd services/post-service && go build -o bin/post-service cmd/api/main.go && cd ../..

echo "Building media-service..."
cd services/media-service && go build -o bin/media-service cmd/api/main.go && cd ../..

echo "Building chat-service..."
cd services/chat-service && go build -o bin/chat-service cmd/api/main.go && cd ../..

echo "Building group-service..."
cd services/group-service && go build -o bin/group-service cmd/api/main.go && cd ../..

echo "All builds completed successfully!"
