# 🚀 Yeni Modbus Ingestor Ekleme Kılavuzu

## 📋 Genel Bakış

Her yeni Modbus cihazı için ayrı bir ingestor instance'ı çalıştırılır. Bu sayede:
- ✅ Her cihaz bağımsız çalışır
- ✅ Bir cihazın problemi diğerlerini etkilemez
- ✅ Cihaza özel buffer ve restart
- ✅ Kolay scale ve yönetim

---

## 🎯 Yöntem 1: Otomatik (Script ile) - ÖNERİLEN

### Adım 1: Config Manager ile Cihaz Oluştur

```bash
curl -X POST http://localhost:8090/api/v1/devices/modbus \
  -H "Content-Type: application/json" \
  -d '{
    "name": "PLC-002",
    "description": "İkinci üretim hattı PLC",
    "ip": "192.168.1.20",
    "port": 502,
    "unit_id": 1,
    "poll_interval_ms": 1000,
    "registers": [
      {
        "address": 0,
        "name": "Factory1.Line2.PLC002.Temp.T001",
        "data_type": "Float32",
        "unit": "°C"
      },
      {
        "address": 2,
        "name": "Factory1.Line2.PLC002.Pressure.P001",
        "data_type": "Int16",
        "unit": "bar"
      }
    ]
  }'
```

**Sonuç:** `config/generated/modbus-PLC-002.toml` oluşturulur

### Adım 2: Script ile Ingestor Ekle

```bash
./scripts/add_modbus_ingestor.sh PLC-002
```

**Script otomatik olarak:**
1. ✅ Config dosyasını kontrol eder
2. ✅ Docker Compose'a servisi ekler
3. ✅ Volume oluşturur
4. ✅ Container'ı başlatır
5. ✅ Durum raporunu gösterir

### Adım 3: Kontrol Et

```bash
# Container durumu
docker ps | grep PLC-002

# Logları izle
docker logs -f ops-ingestor-modbus-plc002

# NATS'e veri gidiyor mu?
docker exec ops-nats-1 nats sub "data.modbus.PLC-002"
```

---

## 🛠️ Yöntem 2: Manuel

### Adım 1: Config Oluştur (Yukarıdaki gibi)

### Adım 2: Docker Compose'a Ekle

`ops/docker-compose.yml` dosyasını açın ve ekleyin:

```yaml
  # Modbus Ingestor: PLC-002
  ingestor-modbus-plc002:
    build:
      context: ..
      dockerfile: services/ingestor/Dockerfile
    container_name: ops-ingestor-modbus-plc002
    environment:
      RUST_LOG: info
      CONFIG_FILE: /config/modbus-PLC-002
      NATS_URL: nats://nats:4222
    volumes:
      - ../config/generated:/config:ro
      - ingestor_buffer_plc002:/data/buffer
    depends_on:
      - nats
      - config-manager
    networks:
      - historian-net
    restart: unless-stopped
    extra_hosts:
      - "host.docker.internal:host-gateway"
```

### Adım 3: Volume Ekle

`volumes:` bölümüne ekleyin:

```yaml
volumes:
  # ... diğer volume'lar
  ingestor_buffer_plc002:  # PLC-002 buffer
```

### Adım 4: Başlat

```bash
cd ops
docker-compose up -d ingestor-modbus-plc002
```

---

## 📊 Çoklu Cihaz Örneği

### 3 Cihaz Senaryosu

```bash
# 1. PLC-001 (Ana hat)
curl -X POST http://localhost:8090/api/v1/devices/modbus -d '{"name":"PLC-001",...}'
./scripts/add_modbus_ingestor.sh PLC-001

# 2. PLC-002 (İkinci hat)
curl -X POST http://localhost:8090/api/v1/devices/modbus -d '{"name":"PLC-002",...}'
./scripts/add_modbus_ingestor.sh PLC-002

# 3. PLC-003 (Paketleme)
curl -X POST http://localhost:8090/api/v1/devices/modbus -d '{"name":"PLC-003",...}'
./scripts/add_modbus_ingestor.sh PLC-003
```

**Sonuç:**
```
ops-ingestor-modbus-plc001  → data.modbus.PLC-001
ops-ingestor-modbus-plc002  → data.modbus.PLC-002
ops-ingestor-modbus-plc003  → data.modbus.PLC-003
```

---

## 🔄 Veri Akışı

```
┌─────────────────────────────────────────────────────────┐
│              Config Manager API                          │
│           http://localhost:8090                          │
└────────────────┬────────────────────────────────────────┘
                 │ POST /api/v1/devices/modbus
                 ↓
┌─────────────────────────────────────────────────────────┐
│         PostgreSQL (historian DB)                        │
│         ├─ devices                                       │
│         ├─ modbus_devices                                │
│         ├─ modbus_registers                              │
│         └─ config_generations                            │
└────────────────┬────────────────────────────────────────┘
                 │ Generates
                 ↓
┌─────────────────────────────────────────────────────────┐
│         config/generated/                                │
│         ├─ modbus-PLC-001.toml                          │
│         ├─ modbus-PLC-002.toml                          │
│         └─ modbus-PLC-003.toml                          │
└────────────────┬────────────────────────────────────────┘
                 │ Mounted to (read-only)
                 ↓
┌─────────────────────────────────────────────────────────┐
│  ops-ingestor-modbus-plc001                             │
│  ops-ingestor-modbus-plc002                             │
│  ops-ingestor-modbus-plc003                             │
└────────────────┬────────────────────────────────────────┘
                 │ Publishes to
                 ↓
┌─────────────────────────────────────────────────────────┐
│                    NATS JetStream                        │
│         ├─ data.modbus.PLC-001                          │
│         ├─ data.modbus.PLC-002                          │
│         └─ data.modbus.PLC-003                          │
└────────────────┬────────────────────────────────────────┘
                 │ Subscribes (data.>)
                 ↓
┌─────────────────────────────────────────────────────────┐
│              ops-engine                                  │
│         (RocksDB + S3 Storage)                          │
└─────────────────────────────────────────────────────────┘
```

---

## 🎛️ Yönetim Komutları

### Tüm Ingestor'ları Listele
```bash
docker ps --filter "name=ops-ingestor-modbus"
```

### Belirli Bir Ingestor'u Restart Et
```bash
cd ops
docker-compose restart ingestor-modbus-plc002
```

### Logları İzle
```bash
# Tek ingestor
docker logs -f ops-ingestor-modbus-plc002

# Tüm ingestor'lar
docker-compose logs -f | grep "ingestor-modbus"
```

### Resource Kullanımı
```bash
docker stats --filter "name=ops-ingestor-modbus"
```

### Ingestor'u Durdur
```bash
cd ops
docker-compose stop ingestor-modbus-plc002
```

### Ingestor'u Kaldır
```bash
cd ops
docker-compose rm -s -f ingestor-modbus-plc002

# Volume'u da sil (DİKKAT: Veri kaybı!)
docker volume rm ops_ingestor_buffer_plc002
```

---

## 🔍 Troubleshooting

### Config dosyası bulunamıyor
```bash
# Config dosyasını kontrol et
ls -la config/generated/modbus-*.toml

# Yoksa yeniden generate et
DEVICE_ID=$(docker exec ops-postgres-1 psql -U postgres -d historian -t -c "SELECT id FROM devices WHERE name='PLC-002';" | tr -d ' \n')
curl -X POST http://localhost:8090/api/v1/config/generate/$DEVICE_ID
```

### NATS bağlantı hatası
```bash
# NATS çalışıyor mu?
docker ps | grep nats

# Config'te NATS URL doğru mu?
cat config/generated/modbus-PLC-002.toml | grep -A 2 "\[nats\]"

# Doğru format:
# url = "nats://nats:4222"
```

### Modbus bağlantı hatası
```bash
# PLC IP'si erişilebilir mi?
ping 192.168.1.20

# Port açık mı?
nc -zv 192.168.1.20 502

# Logları kontrol et
docker logs ops-ingestor-modbus-plc002 | grep -i error
```

### Container başlamıyor
```bash
# Build hatası var mı?
cd ops
docker-compose build ingestor-modbus-plc002

# Yeniden başlat
docker-compose up -d ingestor-modbus-plc002

# Detaylı loglar
docker logs --tail 100 ops-ingestor-modbus-plc002
```

---

## 📈 Performans ve Limitler

### Tek Sunucuda Kaç Ingestor?
- **Önerilen:** 10-20 cihaz
- **Maksimum:** 50-100 cihaz (donanıma bağlı)

### Resource İhtiyacı (Cihaz Başına)
- **CPU:** ~5-10% (1 core)
- **RAM:** ~50-100 MB
- **Disk:** ~10 MB (buffer için)
- **Network:** ~1-10 KB/s

### Ölçekleme Stratejisi
```
1-10 cihaz:    Tek sunucu
10-50 cihaz:   Tek sunucu (güçlü donanım)
50+ cihaz:     Çoklu sunucu (Kubernetes)
```

---

## 🎯 Best Practices

### 1. İsimlendirme
```
✅ İyi:  PLC-001, PLC-MAIN, PLC-LINE-A
❌ Kötü: plc1, test, device
```

### 2. Config Yönetimi
```bash
# Config'leri git'e commit etmeyin
echo "config/generated/*.toml" >> .gitignore

# Backup alın
tar -czf configs-backup-$(date +%Y%m%d).tar.gz config/generated/
```

### 3. Monitoring
```bash
# Healthcheck script
for container in $(docker ps --filter "name=ops-ingestor-modbus" --format "{{.Names}}"); do
    echo "Checking $container..."
    docker logs --tail 5 $container | grep -q "Connected to NATS" && echo "✅ OK" || echo "❌ FAIL"
done
```

### 4. Deployment
```bash
# Yeni cihaz eklerken
1. Önce test ortamında dene
2. Config'i doğrula
3. Logları izle
4. NATS'e veri geldiğini kontrol et
5. Production'a geç
```

---

## 📚 Örnek Senaryolar

### Senaryo 1: Fabrika Genişlemesi
```bash
# 5 yeni hat ekleniyor
for i in {11..15}; do
    curl -X POST http://localhost:8090/api/v1/devices/modbus \
      -d "{\"name\":\"PLC-LINE-$i\", \"ip\":\"192.168.1.$((i+10))\", ...}"
    ./scripts/add_modbus_ingestor.sh PLC-LINE-$i
done
```

### Senaryo 2: Cihaz Değişimi
```bash
# Eski cihazı durdur
cd ops
docker-compose stop ingestor-modbus-plcold

# Yeni cihaz ekle (aynı isimle)
curl -X POST http://localhost:8090/api/v1/config/generate/DEVICE_ID

# Restart
docker-compose restart ingestor-modbus-plcold
```

### Senaryo 3: Toplu Restart
```bash
# Tüm ingestor'ları restart et
cd ops
docker-compose restart $(docker-compose ps --services | grep "ingestor-modbus")
```

---

## ✅ Checklist

Yeni ingestor eklerken:

- [ ] Config Manager ile cihaz oluşturuldu
- [ ] Config dosyası generate edildi
- [ ] Config dosyası doğrulandı (NATS URL, IP, port)
- [ ] Docker Compose'a servis eklendi
- [ ] Volume eklendi
- [ ] Container başlatıldı
- [ ] Loglar kontrol edildi
- [ ] NATS'e bağlantı başarılı
- [ ] Modbus bağlantısı başarılı
- [ ] Veri akışı doğrulandı

---

**🎉 Artık istediğiniz kadar Modbus cihazı ekleyebilirsiniz!**
