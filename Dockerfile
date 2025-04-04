# Build Stage
FROM golang:1.20-alpine AS builder

# Set working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to download dependencies first.
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source code.
COPY . .

# Build the binary. Adjust the command if main package path is different.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ipchecker ./cmd/ipchecker

# Run Stage: Use a minimal image
FROM alpine:latest

# Install any necessary dependencies (for example, CA certificates)
RUN apk --no-cache add ca-certificates

# Set working directory in the runtime container
WORKDIR /root/

# Copy the binary from the builder stage.
COPY --from=builder /app/ipchecker .

# Expose the ports used by the HTTP and gRPC servers.
EXPOSE 8080
EXPOSE 50051

# Set the entrypoint command to run the binary.
CMD ["./ipchecker"]
