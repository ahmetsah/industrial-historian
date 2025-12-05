# Ingestor Service Test Kılavuzu

## 🚀 Hızlı Test

### Otomatik Test Script

```bash
cd /home/ahmet/historian/services/ingestor
python3 test_ingestor.py
```

---

## 📋 Manuel Test Adımları

### 1. NATS'ı Başlat

```bash
cd /home/ahmet/historian/ops
docker-compose up -d nats
```

### 2. Ingestor'ı Başlat

```bash
cd /home/ahmet/historian
cargo run -p ingestor
```

**Beklenen Çıktı:**
```
INFO Starting Ingestor Service
INFO Configuration loaded
INFO Calculation Engine started
INFO Ingestor running. Press Ctrl+C to stop.
```

### 3. NATS Event'lerini İzle

Başka bir terminal'de:

```bash
docker run --rm --network ops_historian-net natsio/nats-box \
  nats sub 'enterprise.>' --server nats://nats:4222
```

### 4. Modbus Simulator Başlat (Opsiyonel)

Eğer gerçek bir Modbus cihazınız yoksa:

```bash
# Python Modbus simulator kur
pip install pymodbus

# Simulator'ı başlat
pymodbus.simulator --modbus-server tcp --modbus-port 5502
```

---

## 🧪 Test Senaryoları

### Senaryo 1: Temel Veri Akışı
1. ✅ Modbus'tan veri okuma
2. ✅ Calculation Engine'de işleme
3. ✅ NATS'a publish etme
4. ✅ Hybrid Buffer (memory + disk)

### Senaryo 2: Calculated Tags
Ingestor, ham sensör verilerini okuyup hesaplanmış tag'ler oluşturur:

**Örnek Konfigürasyon** (`config/ingestor.toml`):
```toml
[[calculated_tags]]
name = "total_flow"
expression = "flow_1 + flow_2"
dependencies = ["flow_1", "flow_2"]

[[calculated_tags]]
name = "efficiency"
expression = "output / input * 100"
dependencies = ["output", "input"]
```

### Senaryo 3: Store & Forward
Ağ kesintisinde veri kaybı olmamalı:

1. NATS'ı durdur: `docker-compose stop nats`
2. Ingestor çalışmaya devam etmeli
3. Veriler disk'e yazılmalı
4. NATS'ı başlat: `docker-compose start nats`
5. Biriken veriler gönderilmeli

---

## 🔍 Debugging

### Ingestor Loglarını İzle

```bash
# Cargo ile çalıştırıyorsanız
cd /home/ahmet/historian
RUST_LOG=debug cargo run -p ingestor
```

### NATS Stream Durumunu Kontrol Et

```bash
docker run --rm --network ops_historian-net natsio/nats-box \
  nats stream info EVENTS --server nats://nats:4222
```

**Beklenen Çıktı:**
```
Information for Stream EVENTS

Configuration:
             Subjects: enterprise.>
             Storage: file
```

### NATS'a Manuel Event Gönder

```bash
docker run --rm --network ops_historian-net natsio/nats-box \
  nats pub enterprise.site1.area1.line1.sensor1 \
  '{"sensor_id":"sensor1","value":42.5,"timestamp":"2025-12-04T10:00:00Z"}' \
  --server nats://nats:4222
```

### Modbus Bağlantısını Test Et

```bash
# Modbus TCP test (Python)
python3 << EOF
from pymodbus.client import ModbusTcpClient

client = ModbusTcpClient('localhost', port=5502)
client.connect()
result = client.read_holding_registers(0, 10, slave=1)
print(f"Registers: {result.registers}")
client.close()
EOF
```

---

## 📊 Konfigürasyon

### Environment Variables

`.env` dosyası:
```bash
NATS_URL=nats://localhost:4222
MODBUS_HOST=localhost
MODBUS_PORT=5502
MODBUS_SLAVE_ID=1
BUFFER_MEMORY_CAPACITY=10000
BUFFER_DISK_PATH=/tmp/ingestor_buffer
```

### Ingestor Config

`config/ingestor.toml`:
```toml
[nats]
url = "nats://localhost:4222"

[modbus]
host = "localhost"
port = 5502
slave_id = 1
poll_interval_ms = 1000

[modbus.registers]
# Holding registers to read
holding = [
    { address = 0, count = 10, tag_prefix = "temp" },
    { address = 100, count = 5, tag_prefix = "pressure" }
]

[buffer]
memory_capacity = 10000
disk_path = "/tmp/ingestor_buffer"

[[calculated_tags]]
name = "avg_temp"
expression = "(temp_0 + temp_1 + temp_2) / 3"
dependencies = ["temp_0", "temp_1", "temp_2"]
```

---

## 🐛 Sık Karşılaşılan Sorunlar

### 1. "Failed to connect to Modbus"
```
ERROR Initial connection failed: Connection refused
```

**Çözüm:**
- Modbus cihazı/simulator çalışıyor mu?
- Host ve port doğru mu?
- Firewall engeli var mı?

```bash
# Port'u kontrol et
nc -zv localhost 5502
```

### 2. "Failed to publish to NATS"
```
ERROR Failed to send to publisher: channel closed
```

**Çözüm:**
- NATS çalışıyor mu?
- NATS URL doğru mu?

```bash
# NATS'ı kontrol et
docker-compose ps nats
```

### 3. "Buffer disk full"
```
ERROR Failed to write to disk buffer: No space left
```

**Çözüm:**
- Disk alanını kontrol et
- Buffer path'i değiştir
- Eski buffer dosyalarını temizle

```bash
# Buffer dosyalarını temizle
rm -rf /tmp/ingestor_buffer/*
```

### 4. "Calculation engine error"
```
ERROR Failed to evaluate expression: Unknown variable 'sensor_x'
```

**Çözüm:**
- Calculated tag dependencies doğru mu?
- Tüm bağımlı sensörler okunuyor mu?
- Expression syntax'ı doğru mu?

---

## 📈 Performans Testi

### Throughput Test

```bash
# NATS message rate'i izle
docker run --rm --network ops_historian-net natsio/nats-box \
  nats bench enterprise.test --pub 10 --msgs 1000 --server nats://nats:4222
```

### Memory Usage

```bash
# Ingestor memory kullanımı
ps aux | grep ingestor
```

### Buffer Performance

```bash
# Buffer dosya boyutu
du -sh /tmp/ingestor_buffer/
```

---

## ✅ Başarı Kriterleri

- [ ] Ingestor başarıyla başlıyor
- [ ] Modbus'a bağlanıyor
- [ ] Sensör verileri okunuyor
- [ ] Calculated tags hesaplanıyor
- [ ] NATS'a event'ler gönderiliyor
- [ ] Buffer çalışıyor (memory + disk)
- [ ] Store & Forward çalışıyor
- [ ] Ağ kesintisinde veri kaybı yok

---

## 🔄 Veri Akışı

```
┌─────────────┐
│   Modbus    │ (PLC/RTU)
│   TCP 5502  │
└──────┬──────┘
       │
       ▼
┌─────────────────────┐
│  Modbus Adapter     │
│  - Read registers   │
│  - Convert to       │
│    SensorData       │
└──────┬──────────────┘
       │
       ▼
┌─────────────────────┐
│ Calculation Engine  │
│  - Evaluate exprs   │
│  - Create calc tags │
└──────┬──────────────┘
       │
       ▼
┌─────────────────────┐
│  Hybrid Buffer      │
│  - Memory (10k)     │
│  - Disk (overflow)  │
└──────┬──────────────┘
       │
       ▼
┌─────────────────────┐
│   NATS Publisher    │
│  - Publish events   │
│  - Retry on fail    │
└──────┬──────────────┘
       │
       ▼
┌─────────────────────┐
│   NATS JetStream    │
│   (enterprise.>)    │
└─────────────────────┘
```

---

## 🎯 İleri Seviye Testler

### 1. Failover Test
```bash
# NATS'ı durdur
docker-compose stop nats

# 1 dakika bekle (veriler buffer'a yazılmalı)
sleep 60

# NATS'ı başlat
docker-compose start nats

# Biriken veriler gönderilmeli
docker run --rm --network ops_historian-net natsio/nats-box \
  nats stream info EVENTS --server nats://nats:4222
```

### 2. Load Test
```bash
# Çok sayıda register oku
# config/ingestor.toml'da register sayısını artır
[modbus.registers]
holding = [
    { address = 0, count = 100, tag_prefix = "sensor" }
]
```

### 3. Expression Test
```bash
# Karmaşık hesaplamalar
[[calculated_tags]]
name = "complex_calc"
expression = "sqrt(sensor_1^2 + sensor_2^2) * 1.5 + offset"
dependencies = ["sensor_1", "sensor_2", "offset"]
```
