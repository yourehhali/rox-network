FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install build dependencies
# sqlite3 requires gcc and musl-dev
RUN apk add --no-cache gcc musl-dev

# Copy the entire backend source code
# We copy everything because of go.work dependency
COPY . .

# Build argument to specify which service to build
ARG SERVICE_NAME

# Build the binary
# We find the main.go file. Most services are in services/<name>/cmd/api/main.go
# api-gateway is in api-gateway/cmd/api/main.go
# verification is in verification/main.go (if used)
# We will use a shell script logic or just rely on the path convention passed in ARG
# Convention: SERVICE_NAME should be the directory name in services/ or api-gateway

RUN if [ "$SERVICE_NAME" = "api-gateway" ]; then \
        go build -o /app/binary ./api-gateway/cmd/api/main.go; \
    else \
        go build -o /app/binary ./services/$SERVICE_NAME/cmd/api/main.go; \
    fi

# Final stage
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies (sqlite3 libraries if dynamically linked, though go-sqlite3 usually is static with CGO_ENABLED=1)
# But alpine might need libc compatibility if not built purely static.
# We'll install nothing for now, check if it works.

# Copy the binary
COPY --from=builder /app/binary /app/service

# Copy database schema (required by some services)
COPY --from=builder /app/database/schema.sql /app/database/schema.sql

# Create uploads directory
RUN mkdir -p /app/uploads

# Expose port (can be overridden)
EXPOSE 8080

CMD ["/app/service"]
