# Audit Service Test Kılavuzu

## 🚀 Hızlı Test (Otomatik Script)

```bash
cd /home/ahmet/historian/go-services/audit
./test_audit_service.sh
```

Bu script:
- ✅ Postgres ve NATS'ı başlatır
- ✅ Audit Service'i build eder ve çalıştırır
- ✅ Test audit eventleri gönderir
- ✅ Chain integrity'yi doğrular
- ✅ Veritabanındaki logları gösterir

---

## 📋 Manuel Test Adımları

### 1. Altyapıyı Başlat

```bash
cd /home/ahmet/historian/ops
docker-compose up -d postgres nats
```

### 2. Audit Service'i Çalıştır

```bash
cd /home/ahmet/historian/go-services/audit

# Environment variables
export DB_URL="postgres://postgres:postgres@localhost:5432/historian?sslmode=disable"
export NATS_URL="nats://localhost:4222"
export PORT="8082"

# Build ve run
go run ./cmd/server
```

### 3. Test Event'leri Gönder

Başka bir terminal'de:

```bash
# NATS CLI kur (eğer yoksa)
go install github.com/nats-io/natscli/nats@latest

# Login event gönder
echo '{"actor":"admin","action":"login","ip":"127.0.0.1"}' | \
  nats pub sys.auth.login --server=localhost:4222

# Audit event gönder
echo '{"actor":"admin","action":"changed_setpoint","device":"PLC-001"}' | \
  nats pub sys.audit.setpoint --server=localhost:4222
```

### 4. Chain Integrity'yi Doğrula

```bash
curl http://localhost:8082/api/v1/audit/verify | jq
```

Beklenen çıktı:
```json
{
  "valid": true
}
```

### 5. Veritabanını İncele

```bash
docker exec -it ops-postgres-1 psql -U postgres -d historian

# SQL sorguları
SELECT * FROM audit_logs ORDER BY timestamp DESC LIMIT 10;

# Hash chain'i kontrol et
SELECT 
  id, 
  timestamp, 
  actor, 
  action, 
  LEFT(prev_hash, 8) as prev, 
  LEFT(curr_hash, 8) as curr 
FROM audit_logs 
ORDER BY timestamp;
```

---

## 🧪 Test Senaryoları

### Senaryo 1: Temel İşlevsellik
1. ✅ Service başlatma
2. ✅ NATS event'i alma
3. ✅ Veritabanına yazma
4. ✅ Hash hesaplama
5. ✅ Verification endpoint

### Senaryo 2: Concurrent Writes (Race Condition)
```bash
# Go test ile
DB_URL="postgres://postgres:postgres@localhost:5432/historian?sslmode=disable" \
  go test ./internal/repository -v -run TestPostgresRepository_Integration
```

### Senaryo 3: Chain Tampering Detection
```bash
# Veritabanında bir hash'i bozalım
docker exec -it ops-postgres-1 psql -U postgres -d historian -c \
  "UPDATE audit_logs SET curr_hash = 'tampered' WHERE id = (SELECT id FROM audit_logs LIMIT 1);"

# Verify endpoint'i çağır - "valid": false dönmeli
curl http://localhost:8082/api/v1/audit/verify | jq
```

---

## 🔍 Debugging

### Logları İzle
```bash
# Audit Service logs
# Service çalışırken terminal'de görünür

# NATS logs
docker logs -f ops-nats-1

# Postgres logs
docker logs -f ops-postgres-1
```

### NATS Stream'leri Kontrol Et
```bash
nats stream ls --server=localhost:4222
nats stream info AUDIT_EVENTS --server=localhost:4222
nats consumer ls AUDIT_EVENTS --server=localhost:4222
```

---

## 🧹 Temizlik

```bash
# Service'i durdur (Ctrl+C)

# Altyapıyı durdur
cd /home/ahmet/historian/ops
docker-compose down

# Verileri de sil (opsiyonel)
docker-compose down -v
```

---

## ✅ Başarı Kriterleri

- [ ] Service başarıyla başlıyor
- [ ] NATS event'leri alınıyor
- [ ] Veritabanına log yazılıyor
- [ ] Hash chain doğru hesaplanıyor
- [ ] Verification endpoint `valid: true` dönüyor
- [ ] Concurrent write'lar race condition yaratmıyor
- [ ] Tampered data tespit ediliyor

---

## 📊 Performans Testi (Opsiyonel)

```bash
# 1000 event gönder
for i in {1..1000}; do
  echo "{\"actor\":\"user$i\",\"action\":\"test\",\"index\":$i}" | \
    nats pub sys.audit.test --server=localhost:4222
done

# Veritabanında kaç log var?
docker exec -it ops-postgres-1 psql -U postgres -d historian -c \
  "SELECT COUNT(*) FROM audit_logs;"
```
