#!/bin/bash
set -e

echo "🚀 Starting Local CI Checks..."

# 1. Rust Checks
echo -e "\n🦀 Running Rust Checks..."
echo "  - Formatting..."
cargo fmt --all -- --check
echo "  - Clippy..."
cargo clippy --workspace --all-targets --all-features -- -D warnings
echo "  - Tests..."
cargo test --workspace
echo "  - Build..."
cargo build --workspace --release

# 2. Go Checks
echo -e "\n🐹 Running Go Checks..."

# Ensure protos are generated
echo "  - Generating Protos..."
./scripts/gen-go-proto.sh

SERVICES=("auth" "audit" "alarm")
for svc in "${SERVICES[@]}"; do
    echo "  - Checking $svc..."
    cd "go-services/$svc"
    echo "    - Vet..."
    go vet ./...
    echo "    - Test..."
    go test ./...
    echo "    - Build..."
    go build ./...
    cd ../..
done

# 3. Frontend Checks
echo -e "\n⚛️ Running Frontend Checks..."
cd viz
echo "  - Installing dependencies..."
npm ci
echo "  - Linting..."
npm run lint
echo "  - Building..."
npm run build
cd ..

echo -e "\n✅ All CI checks passed successfully!"
