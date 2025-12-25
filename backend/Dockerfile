# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/goalify ./cmd/main.go

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Copy the binary from builder
COPY --from=builder /app/bin/goalify .

# Copy migrations (needed for running migrations on startup)
COPY --from=builder /app/internal/db/migrations ./internal/db/migrations

# Expose port (fly.io uses PORT env variable, defaults to 8080)
EXPOSE 8080

# Run the application
CMD ["./goalify"]
