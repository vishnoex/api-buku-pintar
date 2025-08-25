# Build stage
FROM golang:1.24-alpine AS builder

# Install git and build dependencies
RUN apk add --no-cache gcc git musl-dev

# Set working directory
WORKDIR /app

# Copy go mod and sum files first for better caching
COPY go.mod go.sum ./

# Download dependencies with proper flags and error handling
RUN go mod download -x || (echo "Failed to download dependencies" && exit 1)

# Copy source code
COPY . .

# Verify dependencies are correct
RUN go mod verify

# Build the application with vendored dependencies
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -a -installsuffix cgo -o main ./cmd/api

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS and other necessary packages
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/main .

# Copy config files
COPY config.json .
COPY firebase-credentials.json .

# Expose port
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
