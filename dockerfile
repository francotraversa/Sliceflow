# Stage 1: Build the Go binary
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install git (for go modules) and build tools
RUN apk add --no-cache git build-base
RUN apk add --no-cache tzdata

# Copy go.mod and go.sum first (better cache)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the app
COPY . .

# Build binary
RUN go build -o server ./cmd/api

# Stage 2: Run the app
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server .

# Run the binary
CMD ["./server"]
