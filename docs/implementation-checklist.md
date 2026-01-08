# Config Manager Implementation Checklist

## Phase 1: API Layer (2 days)

### Day 1: Core API
- [ ] `internal/api/handlers.go` - REST handlers
  - [ ] GET /api/v1/devices
  - [ ] POST /api/v1/devices/modbus
  - [ ] GET /api/v1/devices/:id
  - [ ] PUT /api/v1/devices/:id
  - [ ] DELETE /api/v1/devices/:id
  
- [ ] `internal/api/middleware.go` - CORS, logging, auth
- [ ] `cmd/server/main.go` - Main application entry
- [ ] Database connection setup
- [ ] Basic error handling

### Day 2: Config Generation
- [ ] POST /api/v1/config/generate/:id
- [ ] POST /api/v1/config/deploy/:id
- [ ] File watcher integration
- [ ] Docker restart trigger
- [ ] Integration tests

## Phase 2: Modbus Ingestor Refactor (2 days)

### Day 3: Restructure
- [ ] Move `services/ingestor` → `services/ingestor-modbus`
- [ ] Update Dockerfile
- [ ] Config hot-reload support
- [ ] NATS subject naming: `data.modbus.{device_name}`
- [ ] Health check endpoint

### Day 4: Testing
- [ ] Unit tests
- [ ] Integration tests with Config Manager
- [ ] Docker Compose setup
- [ ] End-to-end test

## Phase 3: Web UI (3 days)

### Day 5-6: React App
- [ ] Create `web/config-ui` project
- [ ] Device list page
- [ ] Modbus config form
- [ ] API integration
- [ ] Real-time status updates

### Day 7: Polish
- [ ] Validation
- [ ] Error handling
- [ ] Loading states
- [ ] Responsive design

## Phase 4: Additional Protocols (1 week)

### OPC UA Ingestor (3 days)
- [ ] Create `services/ingestor-opc`
- [ ] OPC UA client implementation
- [ ] Config parser
- [ ] Tests

### S7 Ingestor (3 days)
- [ ] Create `services/ingestor-s7`
- [ ] S7 protocol implementation
- [ ] Config parser
- [ ] Tests

### UI Forms (1 day)
- [ ] OPC config form
- [ ] S7 config form

## Phase 5: Production Ready (1 week)

### Kubernetes
- [ ] Deployment manifests
- [ ] Service definitions
- [ ] ConfigMaps
- [ ] Secrets management

### Monitoring
- [ ] Prometheus metrics
- [ ] Grafana dashboards
- [ ] Alert rules

### Documentation
- [ ] API documentation (Swagger)
- [ ] Deployment guide
- [ ] User manual

---

## Quick Start Commands

```bash
# 1. Start database
docker-compose up -d postgres

# 2. Run migrations
psql -h localhost -U historian -d historian_config -f ops/db/init.sql

# 3. Start config manager
cd services/config-manager
go run cmd/server/main.go

# 4. Start modbus ingestor
cd services/ingestor-modbus
cargo run

# 5. Start web UI
cd web/config-ui
npm start
```

---

## Success Criteria

- [ ] Can create Modbus device via UI
- [ ] Config file auto-generated
- [ ] Ingestor picks up new config
- [ ] Data flows to NATS
- [ ] Engine receives data
- [ ] Visible in Viz dashboard

---

## Estimated Timeline

| Phase | Duration | Status |
|-------|----------|--------|
| Phase 1: API | 2 days | 🔨 Next |
| Phase 2: Modbus Refactor | 2 days | ⏳ Pending |
| Phase 3: Web UI | 3 days | ⏳ Pending |
| Phase 4: Protocols | 1 week | ⏳ Pending |
| Phase 5: Production | 1 week | ⏳ Pending |
| **Total** | **~3 weeks** | |

---

## Current Status: 30% Complete ✅

✅ Architecture design
✅ Database schema
✅ Models & Repository
✅ Config Generator
🔨 API Layer (next)
⏳ Modbus Ingestor refactor
⏳ Web UI
⏳ Additional protocols
