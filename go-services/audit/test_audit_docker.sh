#!/bin/bash
set -e

echo "🧪 Audit Service Test (Docker Compose)"
echo "======================================="

cd /home/ahmet/historian/ops

# 1. Start all services
echo -e "\n📦 Starting services (Postgres, NATS, Audit)..."
docker-compose up -d postgres nats audit

# Wait for services
echo "⏳ Waiting for services to be ready..."
sleep 8

# 2. Check if audit service is running
echo -e "\n✅ Checking Audit Service status..."
docker-compose ps audit

# 3. Check logs
echo -e "\n📋 Audit Service logs (last 20 lines):"
docker-compose logs --tail=20 audit

# 4. Test verification endpoint
echo -e "\n🔍 Testing verification endpoint..."
VERIFY_RESULT=$(curl -s http://localhost:8082/api/v1/audit/verify)
echo "Result: $VERIFY_RESULT"

# 5. Send test events using docker exec
echo -e "\n📤 Sending test audit events..."

# Login event
echo '{"actor":"admin","action":"login","ip":"127.0.0.1"}' | \
  docker exec -i ops-nats-1 nats pub sys.auth.login

# Audit events
echo '{"actor":"admin","action":"changed_setpoint","device":"PLC-001","old_value":50,"new_value":75}' | \
  docker exec -i ops-nats-1 nats pub sys.audit.setpoint

echo '{"actor":"operator","action":"acknowledged_alarm","alarm_id":"ALM-123"}' | \
  docker exec -i ops-nats-1 nats pub sys.audit.alarm

# Wait for processing
sleep 3

# 6. Verify chain integrity again
echo -e "\n🔐 Verifying chain integrity after events..."
VERIFY_RESULT=$(curl -s http://localhost:8082/api/v1/audit/verify)
echo "Result: $VERIFY_RESULT"

if echo "$VERIFY_RESULT" | grep -q '"valid":true'; then
    echo -e "\n✅ Chain integrity verified!"
else
    echo -e "\n❌ Chain integrity check failed!"
fi

# 7. Query database
echo -e "\n📊 Audit logs in database:"
docker exec ops-postgres-1 psql -U postgres -d historian -c \
  "SELECT id, timestamp, actor, action, LEFT(prev_hash, 8) as prev, LEFT(curr_hash, 8) as curr FROM audit_logs ORDER BY timestamp;"

echo -e "\n✅ Test completed!"
echo -e "\n💡 Useful commands:"
echo "  - View logs: docker-compose logs -f audit"
echo "  - Stop services: docker-compose down"
echo "  - Restart audit: docker-compose restart audit"
