#!/bin/bash
set -e

# Directory containing proto files
PROTO_DIR="./crates/historian-core/src/proto"
# Output directory for Python generated code
OUT_DIR="./services/sim/proto"

# Ensure output directory exists
mkdir -p "$OUT_DIR"

# Generate Python code
# We use grpcio-tools to generate both proto and grpc stubs (though we might not strictly need grpc for NATS)
# If grpcio-tools is not available in host, advice user to install or use docker.
# For now, we assume standard protoc is available or use python -m grpc_tools.protoc

echo "Generating Python Protobuf definitions..."

# Check if python3-protoc is installed or use pip version
if command -v protoc &> /dev/null; then
    protoc -I="$PROTO_DIR" --python_out="$OUT_DIR" "$PROTO_DIR"/*.proto
else
    # Fallback to python module if available
    python3 -m grpc_tools.protoc -I="$PROTO_DIR" --python_out="$OUT_DIR" "$PROTO_DIR"/*.proto
fi

# Add __init__.py to make it a package
touch "$OUT_DIR/__init__.py"

echo "Done."
