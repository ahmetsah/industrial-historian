# 🚀 Quick Start Guide - Config Manager Prototype

## Prerequisites

- Docker & Docker Compose
- Go 1.21+ (for local development)
- curl & jq (for testing)

---

## 🎯 Quick Start (5 Minutes)

### Step 1: Start Services

```bash
# From project root
cd ops
docker-compose up -d config-db config-manager

# Wait for services to be healthy
docker-compose ps
```

**Expected Output:**
```
NAME                  STATUS
ops-config-db         Up (healthy)
ops-config-manager    Up
```

### Step 2: Verify Health

```bash
curl http://localhost:8090/health | jq '.'
```

**Expected Response:**
```json
{
  "status": "healthy",
  "service": "config-manager",
  "version": "2.0.0"
}
```

### Step 3: Create Your First Device

```bash
curl -X POST http://localhost:8090/api/v1/devices/modbus \
  -H "Content-Type: application/json" \
  -d '{
    "name": "PLC-001",
    "description": "Test PLC",
    "ip": "192.168.1.10",
    "port": 502,
    "unit_id": 1,
    "poll_interval_ms": 1000,
    "registers": [
      {
        "address": 0,
        "name": "Factory1.Line1.PLC001.Temp.T001",
        "data_type": "Float32",
        "unit": "°C"
      }
    ]
  }' | jq '.'
```

**Expected Response:**
```json
{
  "device": {
    "id": "uuid-here",
    "name": "PLC-001",
    ...
  },
  "config": {
    "id": "config-uuid",
    "file_path": "/config/generated/modbus-PLC-001.toml",
    "hash": "sha256-hash"
  },
  "message": "Device created and config generated successfully"
}
```

### Step 4: Check Generated Config

```bash
cat config/generated/modbus-PLC-001.toml
```

**You should see:**
```toml
# Modbus Device: PLC-001
# Generated at: ...

[[modbus_devices]]
ip = "192.168.1.10"
port = 502
unit_id = 1
...
```

### Step 5: Run Full Test Suite

```bash
./scripts/test_api.sh
```

---

## 📊 API Endpoints

### Health Check
```bash
GET http://localhost:8090/health
```

### Device Management

**List All Devices**
```bash
GET http://localhost:8090/api/v1/devices
GET http://localhost:8090/api/v1/devices?protocol=modbus
```

**Get Device**
```bash
GET http://localhost:8090/api/v1/devices/{id}
```

**Delete Device**
```bash
DELETE http://localhost:8090/api/v1/devices/{id}
```

### Modbus Devices

**List Modbus Devices**
```bash
GET http://localhost:8090/api/v1/devices/modbus
```

**Create Modbus Device**
```bash
POST http://localhost:8090/api/v1/devices/modbus
Content-Type: application/json

{
  "name": "PLC-001",
  "description": "Description",
  "ip": "192.168.1.10",
  "port": 502,
  "unit_id": 1,
  "poll_interval_ms": 1000,
  "registers": [
    {
      "address": 0,
      "name": "TagName",
      "data_type": "Float32",
      "unit": "°C"
    }
  ]
}
```

**Get Modbus Device**
```bash
GET http://localhost:8090/api/v1/devices/modbus/{id}
```

**Update Modbus Device**
```bash
PUT http://localhost:8090/api/v1/devices/modbus/{id}
Content-Type: application/json

{
  "name": "PLC-001-Updated",
  ...
}
```

### Config Generation

**Generate Config**
```bash
POST http://localhost:8090/api/v1/config/generate/{device_id}
```

**Get Latest Config**
```bash
GET http://localhost:8090/api/v1/config/latest/{device_id}
```

---

## 🐛 Troubleshooting

### Database Connection Error

```bash
# Check if Postgres is running
docker-compose -f docker-compose.dev.yml logs postgres

# Restart services
docker-compose -f docker-compose.dev.yml restart
```

### Config Files Not Generated

```bash
# Check permissions
ls -la config/generated/

# Create directory if missing
mkdir -p config/generated
```

### Port Already in Use

```bash
# Check what's using port 8090
lsof -i :8090

# Or change port in docker-compose.dev.yml
```

---

## 🧪 Development Mode

### Run Locally (Without Docker)

```bash
# 1. Start only Postgres
docker-compose -f docker-compose.dev.yml up -d postgres

# 2. Set environment variables
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=historian
export DB_PASSWORD=password
export DB_NAME=historian_config
export CONFIG_DIR=./config/generated
export PORT=8090

# 3. Run config manager
cd services/config-manager
go run cmd/server/main.go
```

### Watch Logs

```bash
# All services
docker-compose -f docker-compose.dev.yml logs -f

# Specific service
docker-compose -f docker-compose.dev.yml logs -f config-manager
```

### Database Access

```bash
# Connect to Postgres
docker exec -it historian-postgres psql -U historian -d historian_config

# List devices
SELECT * FROM devices;

# List Modbus devices with registers
SELECT d.name, m.ip, m.port, COUNT(r.id) as register_count
FROM devices d
JOIN modbus_devices m ON d.id = m.id
LEFT JOIN modbus_registers r ON m.id = r.device_id
GROUP BY d.name, m.ip, m.port;
```

---

## 📁 File Structure

```
historian/
├── config/
│   └── generated/              # Auto-generated config files
│       ├── modbus-PLC-001.toml
│       └── modbus-PLC-002.toml
│
├── services/
│   └── config-manager/
│       ├── cmd/server/main.go  # Entry point
│       ├── internal/
│       │   ├── api/            # HTTP handlers
│       │   ├── models/         # Data models
│       │   ├── repository/     # Database layer
│       │   └── generator/      # Config templates
│       └── Dockerfile
│
├── ops/
│   └── db/
│       └── init.sql            # Database schema
│
├── scripts/
│   └── test_api.sh             # API tests
│
└── docker-compose.dev.yml      # Development setup
```

---

## ✅ Success Criteria

You've successfully set up the prototype when:

- [ ] Health check returns 200
- [ ] Can create Modbus device via API
- [ ] Config file is auto-generated in `config/generated/`
- [ ] Can list devices
- [ ] Can update device (config regenerates)
- [ ] Can delete device
- [ ] All tests in `test_api.sh` pass

---

## 🎯 Next Steps

1. **Web UI** - Create React app for device management
2. **Modbus Ingestor** - Refactor to read from generated configs
3. **Hot Reload** - Add file watcher for config changes
4. **OPC UA** - Add OPC UA protocol support
5. **S7** - Add Siemens S7 protocol support

---

## 📞 Need Help?

Check logs:
```bash
docker-compose -f docker-compose.dev.yml logs -f config-manager
```

Reset everything:
```bash
docker-compose -f docker-compose.dev.yml down -v
docker-compose -f docker-compose.dev.yml up -d
```

---

**🎉 Happy Coding!**
