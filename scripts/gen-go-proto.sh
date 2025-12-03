#!/bin/bash
set -e

# Ensure protoc is in PATH
export PATH="$HOME/.local/bin:$PATH"
export PATH="$HOME/go/bin:$PATH"
if ! command -v protoc &> /dev/null; then
    echo "❌ Error: protoc is not installed or not in PATH."
    exit 1
fi

PROTO_ROOT="crates/historian-core/src/proto"
GO_OUT_DIR="go-services/pkg/proto"

echo "Generating Go code from protos..."

mkdir -p "$GO_OUT_DIR"

protoc -I="$PROTO_ROOT" \
  --go_out="$GO_OUT_DIR" --go_opt=paths=source_relative \
  --go-grpc_out="$GO_OUT_DIR" --go-grpc_opt=paths=source_relative \
  "$PROTO_ROOT"/*.proto

echo "✅ Go code generation complete."
