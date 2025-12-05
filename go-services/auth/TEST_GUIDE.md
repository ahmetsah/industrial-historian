# Auth Service Test Kılavuzu

## 🚀 Hızlı Test

### Otomatik Test Script (Önerilen)

```bash
cd /home/ahmet/historian/go-services/auth
python3 test_auth.py
```

---

## 📋 Manuel Test Adımları

### 1. Auth Service'i Başlat

**Seçenek A: Docker Compose ile (Önerilen)**
```bash
cd /home/ahmet/historian/ops
docker-compose up -d postgres nats auth
```

**Seçenek B: Lokal olarak**
```bash
cd /home/ahmet/historian/go-services/auth

# Environment variables
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/historian?sslmode=disable"
export NATS_URL="nats://localhost:4222"
export PORT="8080"
export PRIVATE_KEY_PATH="private.pem"

# Admin kullanıcısı oluştur (ilk kez)
go run main.go -seed-admin -admin-user admin -admin-pass admin123

# Service'i çalıştır
go run main.go
```

### 2. Test Senaryoları

#### A. Login Testi
```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

**Beklenen Çıktı:**
```json
{
  "token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...",
  "role": "ADMIN"
}
```

#### B. Re-Authentication Testi (FDA 21 CFR Part 11)
```bash
# Önce login olup token al
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' | jq -r '.token')

# Re-auth yap
curl -X POST http://localhost:8080/api/v1/re-auth \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"password":"admin123"}'
```

#### C. Service Account Oluşturma (ADMIN only)
```bash
# Admin token ile
curl -X POST http://localhost:8080/api/v1/service-accounts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"username":"plc-service"}'
```

**Beklenen Çıktı:**
```json
{
  "token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...",
  "message": "Service account created"
}
```

#### D. RBAC Testi
```bash
# OPERATOR kullanıcısı oluştur (database'de)
docker exec ops-postgres-1 psql -U postgres -d historian -c \
  "INSERT INTO users (username, password_hash, role) VALUES ('operator', '\$2a\$10\$...', 'OPERATOR');"

# Operator olarak login
OPERATOR_TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"username":"operator","password":"operator123"}' | jq -r '.token')

# Admin endpoint'e erişmeye çalış (403 dönmeli)
curl -X POST http://localhost:8080/api/v1/service-accounts \
  -H "Authorization: Bearer $OPERATOR_TOKEN" \
  -d '{"username":"test"}'
```

---

## 🧪 Test Senaryoları

### Senaryo 1: Temel Kimlik Doğrulama
- ✅ Başarılı login
- ✅ Hatalı şifre (401)
- ✅ Olmayan kullanıcı (401)
- ✅ JWT token oluşturma
- ✅ Token içinde role bilgisi

### Senaryo 2: FDA 21 CFR Part 11 Uyumluluğu
- ✅ Re-authentication zorunluluğu
- ✅ Electronic signature (password tekrar girme)
- ✅ NATS'a audit event gönderme
- ✅ Yeni token oluşturma

### Senaryo 3: RBAC (Role-Based Access Control)
- ✅ ADMIN rolü tüm endpoint'lere erişebilir
- ✅ OPERATOR rolü sadece okuma yapabilir
- ✅ SERVICE rolü sadece veri yazabilir
- ✅ Yetkisiz erişim 403 döner

### Senaryo 4: Service Accounts
- ✅ Long-lived JWT (10 yıl)
- ✅ Sadece ADMIN oluşturabilir
- ✅ SERVICE rolü ile sınırlı

---

## 🔍 Debugging

### Logları İzle
```bash
# Auth service logs
docker-compose logs -f auth

# Database'deki kullanıcıları gör
docker exec ops-postgres-1 psql -U postgres -d historian -c "SELECT * FROM users;"

# NATS event'lerini izle
docker run --rm --network ops_historian-net natsio/nats-box \
  nats sub sys.auth.login --server nats://nats:4222
```

### JWT Token'ı Decode Et
```bash
# jwt.io kullan ya da:
echo $TOKEN | cut -d. -f2 | base64 -d | jq
```

### Database Schema
```sql
-- Users tablosu
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL CHECK (role IN ('ADMIN', 'OPERATOR', 'SERVICE'))
);
```

---

## 🐛 Sık Karşılaşılan Sorunlar

### 1. "Failed to connect to database"
```bash
# Postgres çalışıyor mu?
docker-compose ps postgres

# Bağlantı string doğru mu?
echo $DATABASE_URL
```

### 2. "Failed to connect to NATS"
```bash
# NATS çalışıyor mu?
docker-compose ps nats

# NATS'a bağlanabiliyor muyuz?
docker run --rm --network ops_historian-net natsio/nats-box \
  nats server ping --server nats://nats:4222
```

### 3. "Invalid token"
- Token expire olmuş olabilir (default 24 saat)
- Private key doğru mu? (`private.pem`)
- Token format'ı doğru mu? (`Bearer <token>`)

### 4. "403 Forbidden"
- Kullanıcının rolü endpoint için yeterli mi?
- RBAC middleware doğru çalışıyor mu?

---

## 📊 Performans Testi

```bash
# Apache Bench ile
ab -n 1000 -c 10 -p login.json -T application/json \
  http://localhost:8080/api/v1/login

# login.json:
# {"username":"admin","password":"admin123"}
```

---

## ✅ Başarı Kriterleri

- [ ] Login endpoint çalışıyor
- [ ] JWT token oluşturuluyor
- [ ] Re-authentication çalışıyor
- [ ] NATS'a audit event gönderiliyor
- [ ] RBAC middleware çalışıyor
- [ ] Service account oluşturuluyor
- [ ] Yetkisiz erişim engelleniy or
- [ ] Token validation çalışıyor

---

## 🔐 Güvenlik Notları

1. **Production'da:**
   - `admin123` gibi default şifreler kullanmayın
   - Private key'i güvenli bir yerde saklayın
   - HTTPS kullanın
   - Token expiration süresini ayarlayın

2. **FDA Uyumluluğu:**
   - Re-authentication zorunlu
   - Tüm login/re-auth event'leri audit edilmeli
   - Electronic signature (password) saklanmamalı

3. **RBAC:**
   - En az yetki prensibi
   - Role'ler database'de tanımlı
   - Middleware her endpoint'te kontrol ediyor
