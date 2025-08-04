# Build stage
FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Install git and build dependencies
RUN apk add --no-cache git gcc musl-dev

# Copy go mod and sum files
COPY go.mod go.sum ./

# Copy vendor directory
COPY vendor/ ./vendor/

# Copy source code
COPY . .

# Build the application with vendored dependencies
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -o main ./cmd/api

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/main .

# Copy config files
COPY config.json .

# Copy firebase credentials if it exists (will be created during deployment)
COPY firebase-credentials.json* ./

# Expose port
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
