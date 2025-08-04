# Build stage
FROM golang:1.21-alpine AS builder
# Install git and build dependencies
RUN apk add --no-cache gcc git musl-dev

# Set working directory
WORKDIR /app

RUN go mod download

# Build the application with vendored dependencies
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -a -installsuffix cgo -o main ./cmd/api

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/main .

# Expose port
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
