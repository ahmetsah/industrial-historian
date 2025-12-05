# Build stage
FROM rustlang/rust:nightly-slim AS builder

WORKDIR /build

# Install build dependencies
RUN apt-get update && apt-get install -y \
    pkg-config \
    libssl-dev \
    protobuf-compiler \
    libclang-dev \
    clang \
    && rm -rf /var/lib/apt/lists/*

# Copy workspace files
COPY Cargo.toml Cargo.lock ./
COPY crates ./crates
COPY services/engine ./services/engine

# Build the engine binary
RUN cargo build --release -p engine

# Runtime stage
FROM debian:trixie-slim

# Install runtime dependencies
RUN apt-get update && apt-get install -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/target/release/engine .

# Create data directory for RocksDB
RUN mkdir -p /data/rocksdb

# Expose ports
EXPOSE 50051 8081

# Set environment variables (can be overridden in docker-compose)
ENV RUST_LOG=info \
    DB_PATH=/data/rocksdb \
    GRPC_ADDR=0.0.0.0:50051 \
    HTTP_PORT=8081

CMD ["./engine"]
