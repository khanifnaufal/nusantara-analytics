# Stage 1: Build the Go binary
FROM golang:1.26.4 AS builder

WORKDIR /app

# Copy dependency files and download them
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o main .

# Stage 2: Runtime image
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Expose port 8080 (the default fallback port)
EXPOSE 8080

# Command to run the application
CMD ["./main"]
