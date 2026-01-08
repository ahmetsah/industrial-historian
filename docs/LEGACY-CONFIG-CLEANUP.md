# 🗂️ Eski Config Dosyaları - Gerekli mi?

## 📊 Durum Analizi

### Eski Mimari (Artık Kullanılmıyor)
```
config/default.toml         → Eski monolitik ingestor için
config/modbus_devices.csv   → 100+ cihaz için CSV örneği
```

### Yeni Mimari (Aktif)
```
config/generated/           → Config Manager tarafından oluşturulan
├── modbus-PLC-001.toml
├── modbus-PLC-002.toml
└── modbus-PLC-*.toml
```

---

## ✅ Karar: SİLİNEBİLİR (Ama Yedek Alın)

### `config/default.toml`
**Durum:** ❌ Artık kullanılmıyor
**Neden:**
- Eski monolitik ingestor için tasarlanmıştı
- Yeni mikroservis mimarisi `config/generated/` kullanıyor
- Docker Compose'da artık mount edilmiyor

**Öneri:**
```bash
# Yedek al
mv config/default.toml config/default.toml.backup

# Veya tamamen sil
rm config/default.toml
```

### `config/modbus_devices.csv`
**Durum:** ❌ Artık kullanılmıyor
**Neden:**
- 100+ cihaz için CSV yaklaşımı örneğiydi
- Yeni sistemde PostgreSQL kullanıyoruz
- Config Manager API ile yönetiliyor

**Öneri:**
```bash
# Referans olarak docs/ altına taşı
mv config/modbus_devices.csv docs/examples/modbus_devices.csv.example

# Veya sil
rm config/modbus_devices.csv
```

---

## 📁 Yeni Config Yapısı

### Aktif Dizin
```
config/
└── generated/              ✅ KULLANILIYOR
    ├── modbus-PLC-001.toml
    ├── modbus-PLC-002.toml
    └── modbus-PLC-*.toml
```

### Oluşturma Yöntemi
```bash
# Config Manager API ile
curl -X POST http://localhost:8090/api/v1/devices/modbus \
  -d '{"name":"PLC-NEW",...}'

# Otomatik olarak oluşturulur:
# config/generated/modbus-PLC-NEW.toml
```

---

## 🔄 Migration (Eski → Yeni)

Eğer `default.toml`'deki cihazları yeni sisteme taşımak isterseniz:

### Adım 1: Eski Config'i Analiz Et
```bash
cat config/default.toml
```

### Adım 2: Her Cihaz için API Call
```bash
# Örnek: default.toml'deki ilk cihaz
curl -X POST http://localhost:8090/api/v1/devices/modbus \
  -H "Content-Type: application/json" \
  -d '{
    "name": "PLC-MAIN",
    "ip": "172.29.80.1",
    "port": 5020,
    "unit_id": 1,
    "poll_interval_ms": 3000,
    "registers": [
      {"address": 0, "name": "Factory1.Line1.PLC1.Analog.adres_0", "data_type": "Int16"},
      {"address": 1, "name": "Factory1.Line1.PLC1.Analog.adres_1", "data_type": "Int16"},
      {"address": 2, "name": "Factory1.Line1.PLC1.Analog.adres_2", "data_type": "Int16"}
    ]
  }'
```

### Adım 3: Ingestor Ekle
```bash
./scripts/add_modbus_ingestor.sh PLC-MAIN
```

### Adım 4: Eski Dosyayı Yedekle
```bash
mv config/default.toml config/default.toml.migrated
```

---

## 🧹 Temizlik Önerileri

### Güvenli Temizlik (Önerilen)
```bash
# Yedek dizini oluştur
mkdir -p config/backup

# Eski dosyaları yedekle
mv config/default.toml config/backup/
mv config/modbus_devices.csv config/backup/

# .gitignore'a ekle
echo "config/backup/" >> .gitignore
```

### Agresif Temizlik
```bash
# Tamamen sil (DİKKAT!)
rm config/default.toml
rm config/modbus_devices.csv
```

### Referans Olarak Sakla
```bash
# docs/examples/ altına taşı
mkdir -p docs/examples
mv config/default.toml docs/examples/default.toml.example
mv config/modbus_devices.csv docs/examples/modbus_devices.csv.example
```

---

## 📋 Kontrol Listesi

Eski dosyaları silmeden önce kontrol edin:

- [ ] Yeni mikroservis mimarisi çalışıyor mu?
  ```bash
  docker ps | grep ingestor-modbus
  ```

- [ ] Generated config'ler oluşturuluyor mu?
  ```bash
  ls -la config/generated/
  ```

- [ ] Config Manager API çalışıyor mu?
  ```bash
  curl http://localhost:8090/health
  ```

- [ ] Eski config'te önemli veri var mı?
  ```bash
  cat config/default.toml
  ```

- [ ] Yedek alındı mı?
  ```bash
  ls -la config/backup/
  ```

---

## 🎯 Önerilen Aksiyon

### Seçenek 1: Güvenli (Önerilen)
```bash
# Yedek al ve referans olarak sakla
mkdir -p docs/examples
mv config/default.toml docs/examples/default.toml.example
mv config/modbus_devices.csv docs/examples/modbus_devices.csv.example

# .gitignore güncelle
echo "config/generated/*.toml" >> .gitignore
```

### Seçenek 2: Temiz Başlangıç
```bash
# Yedek al
mkdir -p config/backup
mv config/default.toml config/backup/
mv config/modbus_devices.csv config/backup/

# Sadece generated/ kullan
ls config/generated/
```

---

## 📊 Karşılaştırma

| Özellik | Eski (default.toml) | Yeni (generated/*.toml) |
|---------|---------------------|-------------------------|
| **Yönetim** | Manuel düzenleme | API + Web UI |
| **Validasyon** | Yok | PostgreSQL constraints |
| **Versiyonlama** | Git | Database + hash |
| **Ölçeklenebilirlik** | Zor (tek dosya) | Kolay (cihaz başına) |
| **Hot-reload** | Yok | Gelecekte eklenecek |
| **Audit trail** | Yok | PostgreSQL'de |
| **Backup** | Git | Database backup |

---

## 🔍 Hangi Dosyalar Gerekli?

### ✅ Gerekli (Saklanmalı)
```
config/
└── generated/              ← Config Manager tarafından oluşturulan
    └── *.toml
```

### ❌ Gerekli Değil (Silinebilir/Taşınabilir)
```
config/
├── default.toml            ← Eski monolitik sistem
└── modbus_devices.csv      ← CSV örneği
```

### 📚 Referans (docs/examples/ altında)
```
docs/
└── examples/
    ├── default.toml.example
    └── modbus_devices.csv.example
```

---

## 🚀 Hemen Yapılacaklar

```bash
# 1. Yedek al
mkdir -p docs/examples
cp config/default.toml docs/examples/default.toml.example
cp config/modbus_devices.csv docs/examples/modbus_devices.csv.example

# 2. Eski dosyaları kaldır
rm config/default.toml
rm config/modbus_devices.csv

# 3. .gitignore güncelle
cat >> .gitignore <<EOF

# Generated configs (auto-created by Config Manager)
config/generated/*.toml

# Backup configs
config/backup/
EOF

# 4. Commit
git add .
git commit -m "chore: remove legacy config files, use generated configs"
```

---

## ✅ Sonuç

**EVET, bu dosyalar artık gerekli değil!**

- ✅ `config/default.toml` → Silinebilir (örnek olarak docs/examples/'a taşı)
- ✅ `config/modbus_devices.csv` → Silinebilir (örnek olarak docs/examples/'a taşı)
- ✅ Yeni sistem `config/generated/` kullanıyor
- ✅ Config Manager API ile yönetiliyor
- ✅ PostgreSQL'de saklanıyor

**Önerim:** Referans olarak `docs/examples/` altına taşıyın, sonra silin.
