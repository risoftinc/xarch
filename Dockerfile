# Build stage
FROM golang:1.24.6-alpine AS builder

# Set working directory
WORKDIR /app

# Install git and ca-certificates (needed for go mod download)
RUN apk add --no-cache git ca-certificates

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN adduser -D -s /bin/sh appuser

# Set working directory
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Copy config files
COPY --from=builder /app/config ./config

# Change ownership to appuser
RUN chown -R appuser:appuser /root/

# Switch to non-root user
USER appuser

# Expose ports (HTTP and gRPC)
EXPOSE ${HOST_API} ${HOST_GRPC}

# Run the application
CMD ["./main"]
