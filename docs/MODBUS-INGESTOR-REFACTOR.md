# ✅ Modbus Ingestor Refactor - Tamamlandı!

## 🎯 Yapılan Değişiklikler

### **1. Config Loading Sistemi**
- ✅ `CONFIG_FILE` environment variable desteği
- ✅ Generated config dosyalarını okuma
- ✅ NATS subject konfigürasyonu
- ✅ Backward compatibility (eski config'ler de çalışır)

### **2. Publisher Güncellemesi**
- ✅ Dinamik NATS subject
- ✅ Her cihaz kendi subject'ine publish eder
- ✅ `data.modbus.PLC-001` formatı

### **3. Docker Compose Yapısı**
- ✅ Her cihaz için ayrı ingestor instance
- ✅ Generated config mount (`/config:ro`)
- ✅ Cihaza özel buffer volume
- ✅ Template ve örnek eklendi

---

## 📊 Yeni Mimari

### Öncesi (Monolitik)
```
ops-ingestor  → config/default.toml
              → Tüm cihazlar tek instance'da
              → Tek buffer
```

### Sonrası (Mikroservis)
```
ops-ingestor-modbus-plc-001  → config/generated/modbus-PLC-001.toml
                             → data.modbus.PLC-001
                             → ingestor_buffer_plc001

ops-ingestor-modbus-plc-002  → config/generated/modbus-PLC-002.toml
                             → data.modbus.PLC-002
                             → ingestor_buffer_plc002
```

**Avantajlar:**
- ✅ Bağımsız restart (bir cihaz diğerini etkilemez)
- ✅ Cihaza özel buffer (veri izolasyonu)
- ✅ Kolay scale (yeni cihaz = yeni service)
- ✅ Cihaza özel NATS subject (filtreleme kolaylaşır)

---

## 🚀 Kullanım

### 1. Config Manager ile Cihaz Oluştur
```bash
curl -X POST http://localhost:8090/api/v1/devices/modbus \
  -H "Content-Type: application/json" \
  -d '{
    "name": "PLC-002",
    "ip": "192.168.1.20",
    "port": 502,
    "unit_id": 1,
    "poll_interval_ms": 1000,
    "registers": [
      {
        "address": 1,
        "name": "Factory1.Line2.PLC002.Temp.T001",
        "data_type": "Float32",
        "unit": "°C"
      }
    ]
  }'
```

**Sonuç:** `config/generated/modbus-PLC-002.toml` oluşturulur

### 2. Docker Compose'a Servis Ekle
```yaml
# ops/docker-compose.yml
ingestor-modbus-plc-002:
  build:
    context: ..
    dockerfile: services/ingestor/Dockerfile
  container_name: ops-ingestor-modbus-plc-002
  environment:
    RUST_LOG: info
    CONFIG_FILE: /config/modbus-PLC-002
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

### 3. Volume Ekle
```yaml
volumes:
  ingestor_buffer_plc002:  # Yeni cihaz için
```

### 4. Başlat
```bash
cd ops
docker-compose up -d ingestor-modbus-plc-002
```

---

## 📁 Config Dosyası Formatı

### Generated Config (modbus-PLC-001.toml)
```toml
[[modbus_devices]]
ip = "192.168.1.10"
port = 502
unit_id = 1
poll_interval_ms = 1000
timeout_ms = 5000
retry_count = 3

[[modbus_devices.registers]]
address = 0
name = "Factory1.Line1.PLC001.Temp.T001"
data_type = "Float32"
register_type = "holding"
scale_factor = 1
offset = 0
unit = "°C"

[nats]
url = "${NATS_URL}"
subject = "data.modbus.PLC-001"  # ← Cihaza özel subject

[buffer]
memory_capacity = 10000
disk_path = "/data/buffer/PLC-001.wal"
```

---

## 🔄 Veri Akışı

```
Modbus Device (192.168.1.10:502)
    ↓
ops-ingestor-modbus-plc-001
    ↓ (reads config)
config/generated/modbus-PLC-001.toml
    ↓ (publishes to)
NATS: data.modbus.PLC-001
    ↓ (subscribes)
ops-engine (data.>)
    ↓
RocksDB + S3
```

---

## 🧪 Test

### 1. Config Dosyasını Kontrol Et
```bash
cat config/generated/modbus-PLC-001.toml
```

### 2. Ingestor Loglarını İzle
```bash
docker logs -f ops-ingestor-modbus-plc-001
```

**Beklenen çıktı:**
```
INFO ingestor: Loaded config from /config/modbus-PLC-001
INFO ingestor: Starting Modbus adapter for 192.168.1.10:502
INFO ingestor::publisher: Connected to NATS at nats://nats:4222
INFO ingestor::modbus: Polling 1 registers from unit 1
```

### 3. NATS Subject'i Kontrol Et
```bash
# NATS CLI ile
docker exec ops-nats-1 nats sub "data.modbus.>"

# Veya
docker exec ops-nats-1 nats stream info DATA
```

### 4. Engine'de Veriyi Kontrol Et
```bash
docker logs -f ops-engine | grep "PLC-001"
```

---

## 🎯 Otomatik Deployment (Gelecek)

### Script ile Otomatik Servis Ekleme
```bash
#!/bin/bash
# scripts/add_ingestor.sh

DEVICE_NAME=$1
DEVICE_NAME_LOWER=$(echo $DEVICE_NAME | tr '[:upper:]' '[:lower:]' | tr '-' '')

# 1. Docker Compose'a ekle
cat >> ops/docker-compose.yml <<EOF

  ingestor-modbus-${DEVICE_NAME_LOWER}:
    build:
      context: ..
      dockerfile: services/ingestor/Dockerfile
    container_name: ops-ingestor-modbus-${DEVICE_NAME_LOWER}
    environment:
      RUST_LOG: info
      CONFIG_FILE: /config/modbus-${DEVICE_NAME}
    volumes:
      - ../config/generated:/config:ro
      - ingestor_buffer_${DEVICE_NAME_LOWER}:/data/buffer
    depends_on:
      - nats
      - config-manager
    networks:
      - historian-net
    restart: unless-stopped
    extra_hosts:
      - "host.docker.internal:host-gateway"
EOF

# 2. Volume ekle
sed -i "/^volumes:/a\\  ingestor_buffer_${DEVICE_NAME_LOWER}:" ops/docker-compose.yml

# 3. Başlat
cd ops && docker-compose up -d ingestor-modbus-${DEVICE_NAME_LOWER}
```

**Kullanım:**
```bash
./scripts/add_ingestor.sh PLC-003
```

---

## 📊 Monitoring

### Container Durumu
```bash
docker ps --filter "name=ops-ingestor-modbus"
```

### Resource Kullanımı
```bash
docker stats --filter "name=ops-ingestor-modbus"
```

### Tüm Ingestor Logları
```bash
docker-compose logs -f --tail=100 | grep ingestor-modbus
```

---

## ✅ Başarı Kriterleri

- [x] CONFIG_FILE env var desteği
- [x] Generated config okuma
- [x] NATS subject konfigürasyonu
- [x] Docker Compose template
- [x] Örnek servis (PLC-001)
- [x] Cihaza özel buffer
- [x] Bağımsız restart
- [ ] Otomatik deployment script (opsiyonel)
- [ ] Hot-reload (gelecek)

---

## 🎯 Sonraki Adımlar

1. ✅ **Modbus Ingestor Refactor** - Tamamlandı!
2. ⏳ **Web UI** - Device management + Auto-deploy
3. ⏳ **Hot Reload** - Config değişikliklerini otomatik algılama
4. ⏳ **OPC UA & S7** - Diğer protokoller

---

**🎉 Artık her Modbus cihazı için ayrı, bağımsız ingestor instance çalışıyor!**
