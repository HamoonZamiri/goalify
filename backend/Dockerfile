# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/goalify ./cmd/main.go
# Install goose to run the migrations
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Copy the binary from builder
COPY --from=builder /app/bin/goalify .

# Copy migrations (needed if running migrations in-app)
COPY --from=builder /app/db/migrations ./db/migrations

COPY --from=builder /go/bin/goose /usr/local/bin/goose

# Expose port (fly.io uses PORT env variable, defaults to 8080)
EXPOSE 8080

# Run the application
CMD ["./goalify"]
