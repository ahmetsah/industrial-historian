FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY alarm/go.mod alarm/go.sum ./
COPY pkg ../pkg
RUN go mod download

COPY alarm/ .
RUN go build -o alarm-service ./cmd/server

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/alarm-service .

CMD ["./alarm-service"]
