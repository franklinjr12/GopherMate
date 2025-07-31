# for building this
# docker build --no-cache -t gophermatebackend .
# for running
# docker run -d -p 8080:8080 -p 8443:8443 -p 5432:5432 -p 6543:6543 gophermatebackend

# Start from the official Golang image for building
FROM golang:1.23 AS builder

WORKDIR /app

# Copy go.mod and go.sum first for dependency caching
COPY go.mod go.sum .env ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go app
RUN go build -o build/app ./cmd/main.go

# Use a minimal base image for running
FROM debian:bookworm-slim
WORKDIR /app

# Install ca-certificates for HTTPS connections
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# Copy the built binary from the builder
COPY --from=builder /app/build/app ./app
COPY --from=builder /app/.env ./

# Expose HTTP, HTTPS, and PostgreSQL ports
EXPOSE 8080 8443 5432 6543

# Set environment variable for port (default 8080)
ENV PORT=8080

# Force IPv4 for network connections
ENV GODEBUG=netdns=go

# Command to run the app
CMD ["./app"]
