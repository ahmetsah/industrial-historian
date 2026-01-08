# ✅ PostgreSQL Birleştirme - Tek Database

## 🎯 Yapılan Değişiklik

### Öncesi (❌ Karmaşık)
```
ops-postgres       → historian DB (Ana veriler)
ops-config-db      → historian_config DB (Config veriler)
```
**Sorun:** 2 ayrı PostgreSQL container, 2 ayrı volume, karmaşık yönetim

### Sonrası (✅ Basit)
```
ops-postgres       → historian DB
                     ├─ Ana tablolar (time-series data)
                     └─ Config tablolar (devices, modbus_devices, etc.)
```
**Avantaj:** Tek container, tek volume, kolay yönetim, JOIN'ler mümkün

---

## 📊 Database Yapısı

### Ana PostgreSQL (ops-postgres)
**Database:** `historian`

**Tablolar:**
```sql
-- Config Management (10 tablo)
devices                 -- Ana cihaz tablosu
modbus_devices          -- Modbus cihaz detayları
modbus_registers        -- Modbus register'lar
opc_devices             -- OPC UA cihaz detayları
opc_nodes               -- OPC UA node'lar
s7_devices              -- Siemens S7 cihaz detayları
s7_data_blocks          -- S7 data block'ları
config_generations      -- Oluşturulan config dosyaları
deployment_history      -- Deployment geçmişi
config_templates        -- Config şablonları

-- Ana Historian Tabloları (varsa)
-- time_series_data
-- sensor_metadata
-- etc.
```

---

## 🔧 Bağlantı Bilgileri

### Config Manager Service
```yaml
environment:
  - DB_HOST=postgres          # ✅ Tek PostgreSQL
  - DB_PORT=5432
  - DB_USER=postgres
  - DB_PASSWORD=postgres
  - DB_NAME=historian         # ✅ Aynı database
```

### Manuel Bağlantı
```bash
# Container içinden
docker exec -it ops-postgres-1 psql -U postgres -d historian

# Host'tan
psql -h localhost -p 5432 -U postgres -d historian
```

---

## 📋 Sorgular

### Config Cihazlarını Listele
```sql
SELECT 
    d.name,
    d.protocol,
    d.status,
    CASE 
        WHEN d.protocol = 'modbus' THEN m.ip || ':' || m.port
        WHEN d.protocol = 'opc' THEN o.endpoint_url
        WHEN d.protocol = 's7' THEN s.ip
    END as connection
FROM devices d
LEFT JOIN modbus_devices m ON d.id = m.id
LEFT JOIN opc_devices o ON d.id = o.id
LEFT JOIN s7_devices s ON d.id = s.id;
```

### Modbus Cihaz Detayları
```sql
SELECT 
    d.name,
    m.ip,
    m.port,
    m.unit_id,
    COUNT(r.id) as register_count
FROM devices d
JOIN modbus_devices m ON d.id = m.id
LEFT JOIN modbus_registers r ON m.id = r.device_id
GROUP BY d.name, m.ip, m.port, m.unit_id;
```

### Son Oluşturulan Config'ler
```sql
SELECT 
    d.name,
    cg.file_path,
    cg.generated_at,
    cg.status
FROM config_generations cg
JOIN devices d ON cg.device_id = d.id
ORDER BY cg.generated_at DESC
LIMIT 10;
```

---

## 🚀 Kullanım

### Servisleri Başlat
```bash
cd ops
docker-compose up -d postgres config-manager
```

### Database Kontrolü
```bash
# Tabloları listele
docker exec ops-postgres-1 psql -U postgres -d historian -c "\dt"

# Cihazları listele
docker exec ops-postgres-1 psql -U postgres -d historian -c "SELECT * FROM devices;"
```

### Backup
```bash
# Tüm database (hem ana hem config tabloları)
docker exec ops-postgres-1 pg_dump -U postgres historian > backup_full.sql

# Sadece config tabloları
docker exec ops-postgres-1 pg_dump -U postgres historian \
  -t devices -t modbus_devices -t modbus_registers \
  -t opc_devices -t opc_nodes -t s7_devices -t s7_data_blocks \
  -t config_generations -t deployment_history -t config_templates \
  > backup_config_only.sql
```

### Restore
```bash
# Tüm database
cat backup_full.sql | docker exec -i ops-postgres-1 psql -U postgres historian

# Sadece config tabloları
cat backup_config_only.sql | docker exec -i ops-postgres-1 psql -U postgres historian
```

---

## ✅ Avantajlar

### 1. **Basitlik**
- ✅ Tek PostgreSQL container
- ✅ Tek volume (`ops_pg_data`)
- ✅ Tek port (5432)
- ✅ Tek backup/restore

### 2. **Performans**
- ✅ JOIN'ler aynı database içinde (hızlı)
- ✅ Transaction'lar tek DB'de (ACID garantisi)
- ✅ Connection pool paylaşımı

### 3. **Yönetim**
- ✅ Tek yerde monitoring
- ✅ Tek yerde backup
- ✅ Tek yerde migration
- ✅ Kolay troubleshooting

### 4. **Maliyet**
- ✅ Daha az RAM kullanımı
- ✅ Daha az disk I/O
- ✅ Daha az network overhead

---

## 📊 Kaynak Kullanımı

### Öncesi (2 PostgreSQL)
```
ops-postgres:    ~100MB RAM
ops-config-db:   ~100MB RAM
Total:           ~200MB RAM
```

### Sonrası (1 PostgreSQL)
```
ops-postgres:    ~120MB RAM
Total:           ~120MB RAM
Tasarruf:        ~80MB RAM (40%)
```

---

## 🔄 Migration (Eski Veriler Varsa)

Eğer `ops-config-db`'de veri varsa:

```bash
# 1. Eski DB'den export
docker exec ops-config-db pg_dump -U historian historian_config > old_config.sql

# 2. Ana DB'ye import
cat old_config.sql | docker exec -i ops-postgres-1 psql -U postgres historian

# 3. Eski container'ı kaldır
docker rm ops-config-db
docker volume rm ops_config_db_data
```

---

## 🎯 Sonraki Adımlar

1. ✅ **Tek PostgreSQL** - Tamamlandı
2. ✅ **Config tabloları** - Oluşturuldu
3. ✅ **Config Manager bağlantısı** - Güncellendi
4. ⏳ **Modbus Ingestor refactor** - Config'leri okusun
5. ⏳ **Web UI** - Device management

---

## 📝 Notlar

- **Schema:** `ops/db/init.sql` dosyasında
- **Auto-load:** PostgreSQL ilk başlatıldığında otomatik yüklenir
- **Seed data:** `PLC-001` örnek cihazı otomatik oluşturulur
- **Views:** `v_devices_complete` ve `v_latest_configs` hazır

---

**🎉 Artık tek bir PostgreSQL ile hem ana veriler hem de config veriler yönetiliyor!**
