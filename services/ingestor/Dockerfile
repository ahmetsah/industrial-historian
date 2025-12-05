# Build stage
FROM rustlang/rust:nightly-slim AS builder

WORKDIR /build

# Install build dependencies
RUN apt-get update && apt-get install -y \
    pkg-config \
    libssl-dev \
    protobuf-compiler \
    && rm -rf /var/lib/apt/lists/*

# Copy workspace files
COPY Cargo.toml Cargo.lock ./
COPY crates ./crates
COPY services/ingestor ./services/ingestor

# Build the ingestor binary with nightly features
RUN cargo build --release -p ingestor

# Runtime stage
FROM debian:trixie-slim

# Install runtime dependencies
RUN apt-get update && apt-get install -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/target/release/ingestor .

# Create config directory
RUN mkdir -p /app/config /app/ops/data/ingestor_wal

# Set environment variables (can be overridden in docker-compose)
ENV RUST_LOG=info

CMD ["./ingestor"]
