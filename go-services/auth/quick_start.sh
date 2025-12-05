#!/bin/bash
set -e

echo "🔐 Auth Service Quick Start"
echo "============================"

cd /home/ahmet/historian/ops

# 1. Start services
echo -e "\n📦 Starting Postgres, NATS, and Auth..."
docker-compose up -d postgres nats auth

# Wait for services
echo "⏳ Waiting for services to be ready..."
sleep 8

# 2. Check if auth service is running
echo -e "\n✅ Checking Auth Service status..."
docker-compose ps auth

# 3. Seed admin user
echo -e "\n👤 Creating admin user..."
docker-compose exec auth ./auth-server -seed-admin -admin-user admin -admin-pass admin123 || echo "Admin user may already exist"

# 4. Test login
echo -e "\n🔑 Testing login..."
sleep 2
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' | python3 -c "import sys, json; print(json.load(sys.stdin).get('token', 'FAILED')[:50])")

if [ "$TOKEN" != "FAILED" ]; then
    echo "✅ Login successful!"
    echo "   Token: $TOKEN..."
else
    echo "❌ Login failed"
    echo "   Checking logs..."
    docker-compose logs --tail=20 auth
fi

echo -e "\n✅ Auth Service is ready!"
echo -e "\n💡 Next steps:"
echo "  1. Run full tests: cd /home/ahmet/historian/go-services/auth && python3 test_auth.py"
echo "  2. View logs: docker-compose logs -f auth"
echo "  3. Check users: docker exec ops-postgres-1 psql -U postgres -d historian -c 'SELECT * FROM users;'"
