FROM golang:1.23-alpine AS builder

# Set up dependencies
RUN apk add --no-cache git

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN go build -o /app/p2p-file-sharing ./cmd/p2pfs/main.go

# Final Stage
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the pre-built binary file from the builder stage
COPY --from=builder /app/p2p-file-sharing /app/p2p-file-sharing

# Create directories for shared files and downloads
RUN mkdir /app/shared /app/downloads

# Expose necessary ports for p2p
EXPOSE 4001
EXPOSE 8080

# Command to run the application
CMD ["/app/p2p-file-sharing"]
