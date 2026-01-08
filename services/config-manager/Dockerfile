FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o config-manager ./cmd/server

# Final stage
FROM alpine:latest

# Install ca-certificates and docker-cli for container orchestration
RUN apk --no-cache add ca-certificates docker-cli

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/config-manager .

# Create config directory
RUN mkdir -p /config/generated

EXPOSE 8090

CMD ["./config-manager"]
