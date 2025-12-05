# Docker Deployment Guide

## 🚀 Quick Start

### Start All Services
```bash
cd ops
docker-compose up -d
```

### Check Status
```bash
docker-compose ps
```

### View Logs
```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f ingestor
docker-compose logs -f engine
```

### Stop Services
```bash
docker-compose down
```

---

## 📦 Services

### Infrastructure
- **NATS JetStream** (4222, 8222) - Message broker
- **PostgreSQL** (5432) - Auth & Audit database
- **MinIO** (9000, 9001) - S3-compatible storage
- **PgAdmin** (5050) - Database management

### Application Services
- **Auth** (8080) - Authentication & JWT
- **Audit** (8082) - Immutable audit logs
- **Ingestor** - Modbus data collection
- **Engine** (50051, 8081) - Time-series storage & query

---

## 🔧 Configuration

### Environment Variables

Services can be configured via environment variables in `docker-compose.yml`:

**Ingestor:**
```yaml
environment:
  RUST_LOG: info
  NATS_URL: nats://nats:4222
```

**Engine:**
```yaml
environment:
  RUST_LOG: info
  NATS_URL: nats://nats:4222
  NATS_SUBJECT: data.raw
  DB_PATH: /data/rocksdb
  GRPC_ADDR: 0.0.0.0:50051
  HTTP_PORT: 8081
```

### Volumes

Data is persisted in Docker volumes:
- `nats_data` - NATS JetStream data
- `pg_data` - PostgreSQL database
- `minio_data` - MinIO object storage
- `ingestor_buffer` - Ingestor WAL buffer
- `engine_data` - RocksDB time-series data

---

## 🏗️ Building

### Build All Services
```bash
docker-compose build
```

### Build Specific Service
```bash
docker-compose build ingestor
docker-compose build engine
```

### Rebuild Without Cache
```bash
docker-compose build --no-cache ingestor
```

---

## 🧪 Testing

### Test Ingestor
```bash
# Check if running
docker-compose ps ingestor

# View logs
docker-compose logs -f ingestor

# Run test script
cd ../services/ingestor
python3 test_ingestor.py
```

### Test Engine
```bash
# Check if running
docker-compose ps engine

# View logs
docker-compose logs -f engine

# Test HTTP export
curl 'http://localhost:8081/api/v1/export?sensor_id=adres_0&start_ts=0&end_ts=9999999999999'

# Test gRPC (requires grpcurl)
grpcurl -plaintext localhost:50051 list

# Run test script
cd ../services/engine
python3 test_engine.py
```

---

## 🔍 Troubleshooting

### Service Won't Start

```bash
# Check logs
docker-compose logs service-name

# Restart service
docker-compose restart service-name

# Rebuild and restart
docker-compose up -d --build service-name
```

### Port Already in Use

```bash
# Find process using port
sudo lsof -i :8080

# Kill process
sudo kill -9 PID
```

### Network Issues

```bash
# Recreate network
docker-compose down
docker network prune
docker-compose up -d
```

### Volume Issues

```bash
# Remove all volumes (WARNING: deletes data)
docker-compose down -v

# Remove specific volume
docker volume rm ops_engine_data
```

---

## 📊 Monitoring

### Resource Usage
```bash
# All containers
docker stats

# Specific container
docker stats ops-ingestor ops-engine
```

### Database Size
```bash
# RocksDB (Engine)
docker exec ops-engine du -sh /data/rocksdb

# PostgreSQL
docker exec ops-postgres-1 psql -U postgres -d historian -c "SELECT pg_size_pretty(pg_database_size('historian'));"
```

### NATS Stream Info
```bash
docker run --rm --network ops_historian-net natsio/nats-box \
  nats stream info EVENTS --server nats://nats:4222
```

---

## 🔄 Updates

### Update Service
```bash
# Pull latest code
git pull

# Rebuild and restart
docker-compose up -d --build service-name
```

### Update All Services
```bash
git pull
docker-compose down
docker-compose build
docker-compose up -d
```

---

## 🗑️ Cleanup

### Stop and Remove Containers
```bash
docker-compose down
```

### Remove Containers and Volumes
```bash
docker-compose down -v
```

### Remove Images
```bash
docker-compose down --rmi all
```

### Full Cleanup
```bash
docker-compose down -v --rmi all
docker system prune -a
```

---

## 🌐 Network

All services run on the `historian-net` bridge network:

```bash
# Inspect network
docker network inspect ops_historian-net

# List connected containers
docker network inspect ops_historian-net | grep Name
```

---

## 📝 Notes

### Modbus Connection
Ingestor connects to Modbus device on host machine (172.29.80.1:5020).
This is configured via `extra_hosts` in docker-compose.yml.

### Data Persistence
All data is stored in Docker volumes and persists across restarts.
To completely reset, use `docker-compose down -v`.

### Restart Policy
All services have `restart: always` policy and will automatically
restart on failure or system reboot.

---

## 🎯 Production Checklist

- [ ] Update default passwords (PostgreSQL, MinIO, PgAdmin)
- [ ] Configure proper resource limits (CPU, Memory)
- [ ] Set up log rotation
- [ ] Configure backup strategy for volumes
- [ ] Enable TLS/SSL for external endpoints
- [ ] Set up monitoring (Prometheus, Grafana)
- [ ] Configure firewall rules
- [ ] Test disaster recovery procedures
