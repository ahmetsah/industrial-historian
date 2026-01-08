# ✅ Docker Yönetimi - Ops Klasörüne Taşındı

## 🎯 Yapılan Değişiklikler

### 1. **Docker Compose Konsolidasyonu**
- ❌ **Eski:** `docker-compose.dev.yml` (root dizinde, "historian-" prefix)
- ✅ **Yeni:** `ops/docker-compose.yml` (tüm servisler tek yerde, "ops-" prefix)

### 2. **Container İsimlendirme Standardı**
Tüm container'lar artık `ops-` prefix'i kullanıyor:

```
ops-config-db          # Config Management PostgreSQL
ops-config-manager     # Config Manager API
ops-postgres           # Ana PostgreSQL
ops-nats               # NATS Message Bus
ops-minio              # S3-compatible storage
ops-auth               # Auth Service
ops-audit              # Audit Service
ops-alarm              # Alarm Service
ops-ingestor           # Modbus Ingestor
ops-engine             # Engine Service
ops-sim                # Simulator
ops-viz                # Visualization UI
ops-pgadmin            # PostgreSQL Admin UI
```

### 3. **Port Mapping**
```
5432  → PostgreSQL (Ana DB)
5433  → Config DB (Ayrı PostgreSQL)
8090  → Config Manager API
8080  → Auth Service
8081  → Engine HTTP API
8082  → Audit Service
8083  → Alarm Service
50051 → Engine gRPC API
3000  → Viz UI
4222  → NATS
8222  → NATS Monitoring
9000  → MinIO API
9001  → MinIO Console
5050  → PgAdmin
```

---

## 🚀 Kullanım

### Tüm Servisleri Başlat
```bash
cd ops
docker-compose up -d
```

### Sadece Config Manager
```bash
cd ops
docker-compose up -d config-db config-manager
```

### Sadece Core Services (NATS, Postgres, Engine, Ingestor)
```bash
cd ops
docker-compose up -d nats postgres ingestor engine
```

### Logları İzle
```bash
cd ops
docker-compose logs -f config-manager
```

### Servisleri Durdur
```bash
cd ops
docker-compose down
```

### Tüm Verileri Sil (Dikkat!)
```bash
cd ops
docker-compose down -v
```

---

## 📊 Servis Durumu Kontrolü

### Health Checks
```bash
# Config Manager
curl http://localhost:8090/health

# Auth Service
curl http://localhost:8080/health

# Engine Service
curl http://localhost:8081/health
```

### Container Durumu
```bash
docker ps --filter "name=ops-"
```

### Database Bağlantısı
```bash
# Ana PostgreSQL
docker exec -it ops-postgres psql -U postgres -d historian

# Config DB
docker exec -it ops-config-db psql -U historian -d historian_config
```

---

## 🗄️ Volume Yönetimi

### Mevcut Volume'lar
```
ops_nats_data          # NATS JetStream data
ops_minio_data         # MinIO object storage
ops_pg_data            # Ana PostgreSQL data
ops_config_db_data     # Config DB data
ops_ingestor_buffer    # Ingestor buffer/WAL
ops_engine_data        # Engine RocksDB data
```

### Volume Backup
```bash
# Config DB backup
docker exec ops-config-db pg_dump -U historian historian_config > backup_config.sql

# Ana DB backup
docker exec ops-postgres pg_dump -U postgres historian > backup_main.sql
```

### Volume Restore
```bash
# Config DB restore
cat backup_config.sql | docker exec -i ops-config-db psql -U historian historian_config

# Ana DB restore
cat backup_main.sql | docker exec -i ops-postgres psql -U postgres historian
```

---

## 🔧 Troubleshooting

### Config Manager bağlanamıyor
```bash
# Logları kontrol et
docker logs ops-config-manager

# Config DB sağlıklı mı?
docker exec ops-config-db pg_isready -U historian

# Network kontrolü
docker network inspect ops_historian-net
```

### Port çakışması
```bash
# Hangi port kullanımda?
lsof -i :8090

# Alternatif port kullan (docker-compose.yml'de değiştir)
ports:
  - "8091:8090"  # Host:Container
```

### Container restart loop
```bash
# Son 50 log satırı
docker logs --tail 50 ops-config-manager

# Container inspect
docker inspect ops-config-manager
```

---

## 📁 Dosya Yapısı

```
historian/
├── ops/
│   ├── docker-compose.yml       ✅ Ana orchestration dosyası
│   ├── db/
│   │   └── init.sql            ✅ Config DB schema
│   ├── nats.conf
│   └── setup_streams.sh
│
├── config/
│   └── generated/              ✅ Auto-generated TOML configs
│       ├── modbus-PLC-001.toml
│       ├── modbus-PLC-OPS-001.toml
│       └── ...
│
└── services/
    ├── config-manager/         ✅ Config Manager service
    ├── ingestor/
    ├── engine/
    └── ...
```

---

## ✅ Başarı Kriterleri

- [x] Tüm container'lar `ops-` prefix kullanıyor
- [x] Docker Compose `ops/` klasöründe
- [x] Config Manager çalışıyor (Port 8090)
- [x] Config DB ayrı instance (Port 5433)
- [x] Health check başarılı
- [x] Cihaz oluşturma çalışıyor
- [x] Config dosyaları generate ediliyor
- [x] Volume'lar doğru mount ediliyor

---

## 🎯 Sonraki Adımlar

1. **Modbus Ingestor Refactor** - Generated config'leri okusun
2. **Web UI** - React app ile device management
3. **OPC UA & S7** - Ek protokol desteği
4. **Monitoring** - Prometheus + Grafana
5. **Production Deployment** - Kubernetes manifests

---

**🎉 Tüm servisler artık `ops/` klasöründen yönetiliyor!**
