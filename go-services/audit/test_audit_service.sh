#!/bin/bash
set -e

echo "🧪 Audit Service Test Script"
echo "=============================="

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# 1. Start infrastructure
echo -e "\n${BLUE}1. Starting infrastructure (Postgres + NATS)...${NC}"
cd /home/ahmet/historian/ops
docker-compose up -d postgres nats

# Wait for services
echo -e "${BLUE}Waiting for Postgres...${NC}"
sleep 5

# 2. Build Audit Service
echo -e "\n${BLUE}2. Building Audit Service...${NC}"
cd /home/ahmet/historian/go-services/audit
go build -o audit-server ./cmd/server

# 3. Run Audit Service
echo -e "\n${BLUE}3. Starting Audit Service...${NC}"
export DB_URL="postgres://postgres:postgres@localhost:5432/historian?sslmode=disable"
export NATS_URL="nats://localhost:4222"
export PORT="8082"

./audit-server &
AUDIT_PID=$!
echo "Audit Service started with PID: $AUDIT_PID"

# Wait for service to start
sleep 3

# 4. Test verification endpoint (should be valid initially)
echo -e "\n${BLUE}4. Testing verification endpoint (empty chain)...${NC}"
VERIFY_RESULT=$(curl -s http://localhost:8082/api/v1/audit/verify)
echo "Result: $VERIFY_RESULT"

# 5. Publish test audit events via NATS
echo -e "\n${BLUE}5. Publishing test audit events...${NC}"

# Install nats CLI if not present
if ! command -v nats &> /dev/null; then
    echo "Installing nats CLI..."
    go install github.com/nats-io/natscli/nats@latest
fi

# Publish some test events
echo "Publishing login event..."
echo '{"actor":"admin","action":"login","ip":"127.0.0.1"}' | nats pub sys.auth.login --server=localhost:4222

echo "Publishing audit events..."
echo '{"actor":"admin","action":"changed_setpoint","device":"PLC-001","old_value":50,"new_value":75}' | nats pub sys.audit.setpoint --server=localhost:4222
echo '{"actor":"operator","action":"acknowledged_alarm","alarm_id":"ALM-123"}' | nats pub sys.audit.alarm --server=localhost:4222

# Wait for processing
sleep 2

# 6. Verify chain integrity
echo -e "\n${BLUE}6. Verifying chain integrity...${NC}"
VERIFY_RESULT=$(curl -s http://localhost:8082/api/v1/audit/verify)
echo "Result: $VERIFY_RESULT"

if echo "$VERIFY_RESULT" | grep -q '"valid":true'; then
    echo -e "${GREEN}✅ Chain integrity verified!${NC}"
else
    echo -e "${RED}❌ Chain integrity check failed!${NC}"
fi

# 7. Query database to see logs
echo -e "\n${BLUE}7. Querying audit logs from database...${NC}"
docker exec -it ops-postgres-1 psql -U postgres -d historian -c "SELECT id, timestamp, actor, action, prev_hash, curr_hash FROM audit_logs ORDER BY timestamp;"

# 8. Cleanup
echo -e "\n${BLUE}8. Cleanup...${NC}"
echo "Stopping Audit Service (PID: $AUDIT_PID)..."
kill $AUDIT_PID 2>/dev/null || true

echo -e "\n${GREEN}✅ Test completed!${NC}"
echo -e "${BLUE}To stop infrastructure: cd /home/ahmet/historian/ops && docker-compose down${NC}"
