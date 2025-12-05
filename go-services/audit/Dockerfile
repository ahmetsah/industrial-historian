# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /build

# Copy go mod files
COPY audit/go.mod audit/go.sum ./
COPY pkg ../pkg
RUN go mod download

# Copy source
COPY audit/ .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o audit-server ./cmd/server

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /build/audit-server .

EXPOSE 8082

CMD ["./audit-server"]
