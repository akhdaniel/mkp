# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the migration tool
RUN CGO_ENABLED=0 GOOS=linux go build -o migrate cmd/migrate/main.go

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o server cmd/server/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates bash

WORKDIR /app

# Copy the binaries from builder
COPY --from=builder /app/server .
COPY --from=builder /app/migrate .

# Copy migrations
COPY --from=builder /app/internal/database/migrations ./internal/database/migrations

# Copy entrypoint script
COPY docker-entrypoint.sh .
RUN chmod +x docker-entrypoint.sh

# Expose port
EXPOSE 8080

# Run the entrypoint script
ENTRYPOINT ["./docker-entrypoint.sh"]