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

# 3. Python Checks (Sim)
echo -e "\n🐍 Running Python Checks (Sim)..."
cd services/sim
echo "  - Running Tests..."

# Create/Activate venv to ensure dependencies are available
if [ ! -f ".venv/bin/activate" ]; then
    echo "  - Creating python venv..."
    rm -rf .venv # Clean up potential broken install
    python3 -m venv .venv || {
        echo "❌ Failed to create venv. Please install python3-venv (and python3-pip)."
        echo "   Ubuntu/Debian: sudo apt install python3-venv python3-pip"
        exit 1
    }
fi
source .venv/bin/activate

echo "  - Installing requirements..."
pip install -q -r requirements.txt


# Add src to PYTHONPATH so tests can import modules
export PYTHONPATH=$PYTHONPATH:$(pwd)/src:$(pwd)/proto

cd src
python3 -m unittest discover -p "test_*.py"
cd .. # Back to services/sim
deactivate
cd ../.. # Back to root

# 4. Frontend Checks
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
