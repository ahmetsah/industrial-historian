# 🔧 PgAdmin ile PostgreSQL Bağlantısı

## 🌐 Web Arayüzüne Erişim

### 1. PgAdmin'i Aç
```
URL: http://localhost:5050
```

### 2. Giriş Bilgileri
```
Email:    admin@historian.com
Password: admin
```

---

## 🔌 PostgreSQL Server Ekleme

### Adım 1: Add New Server
1. Sol menüde **Servers** üzerine sağ tıklayın
2. **Register** → **Server** seçin

### Adım 2: General Tab
```
Name: Historian PostgreSQL
```

### Adım 3: Connection Tab
```
Host name/address: postgres
Port:              5432
Maintenance DB:    historian
Username:          postgres
Password:          postgres
```

✅ **Save password** kutusunu işaretleyin

### Adım 4: Save
**Save** butonuna tıklayın

---

## 📊 Database'leri Görüntüleme

Bağlantı kurulduktan sonra:

```
Servers
  └─ Historian PostgreSQL
      └─ Databases
          └─ historian
              ├─ Schemas
              │   └─ public
              │       ├─ Tables (10 tablo)
              │       │   ├─ devices
              │       │   ├─ modbus_devices
              │       │   ├─ modbus_registers
              │       │   ├─ opc_devices
              │       │   ├─ opc_nodes
              │       │   ├─ s7_devices
              │       │   ├─ s7_data_blocks
              │       │   ├─ config_generations
              │       │   ├─ deployment_history
              │       │   └─ config_templates
              │       └─ Views (2 view)
              │           ├─ v_devices_complete
              │           └─ v_latest_configs
```

---

## 🔍 Örnek Sorgular

### Tüm Cihazları Listele
```sql
SELECT 
    d.name,
    d.protocol,
    d.status,
    d.enabled,
    d.created_at
FROM devices d
ORDER BY d.created_at DESC;
```

### Modbus Cihaz Detayları
```sql
SELECT 
    d.name,
    m.ip,
    m.port,
    m.unit_id,
    m.poll_interval_ms,
    COUNT(r.id) as register_count
FROM devices d
JOIN modbus_devices m ON d.id = m.id
LEFT JOIN modbus_registers r ON m.id = r.device_id
GROUP BY d.name, m.ip, m.port, m.unit_id, m.poll_interval_ms;
```

### Son Oluşturulan Config'ler
```sql
SELECT 
    d.name,
    cg.file_path,
    cg.config_hash,
    cg.generated_at,
    cg.status
FROM config_generations cg
JOIN devices d ON cg.device_id = d.id
ORDER BY cg.generated_at DESC
LIMIT 10;
```

### Complete Device View (Hazır View)
```sql
SELECT * FROM v_devices_complete;
```

---

## 🛠️ Veri Düzenleme

### Yeni Kayıt Ekle
1. İlgili tabloya sağ tıklayın
2. **View/Edit Data** → **All Rows**
3. Üstteki toolbar'dan **Add Row** (+) butonuna tıklayın
4. Verileri girin
5. **Save** (💾) butonuna tıklayın

### Kayıt Güncelle
1. Tabloda satıra çift tıklayın
2. Değeri değiştirin
3. **Save** butonuna tıklayın

### Kayıt Sil
1. Satırı seçin
2. **Delete** butonuna tıklayın

---

## 📈 Grafik ve Analiz

### Query Tool
1. Database'e sağ tıklayın
2. **Query Tool** seçin
3. SQL sorgunuzu yazın
4. **Execute** (▶) butonuna tıklayın

### Export Data
1. Sorgu sonucunda **Download** butonuna tıklayın
2. Format seçin (CSV, JSON, etc.)

---

## 🔒 Güvenlik Notları

### Development (Şu anki)
```
Email:    admin@historian.com
Password: admin
Host:     localhost:5050
```
⚠️ **Sadece development için!**

### Production
```yaml
# docker-compose.yml
pgadmin:
  environment:
    PGADMIN_DEFAULT_EMAIL: your-email@company.com
    PGADMIN_DEFAULT_PASSWORD: strong-password-here
  # Reverse proxy arkasında çalıştırın
  # HTTPS kullanın
```

---

## 🚀 Hızlı Erişim

### Tarayıcı Bookmark
```
http://localhost:5050
```

### Docker Container
```bash
# PgAdmin container'ını kontrol et
docker ps | grep pgadmin

# Logları izle
docker logs -f ops-pgadmin
```

---

## 🐛 Troubleshooting

### PgAdmin açılmıyor
```bash
# Container çalışıyor mu?
docker ps | grep pgadmin

# Restart
cd ops
docker-compose restart pgadmin

# Logları kontrol et
docker logs ops-pgadmin
```

### Bağlantı hatası
```bash
# PostgreSQL çalışıyor mu?
docker ps | grep postgres

# Network kontrolü
docker network inspect ops_historian-net
```

### Şifre hatası
```
Username: postgres
Password: postgres
Database: historian
```

---

## 📚 Faydalı Linkler

- **PgAdmin Docs:** https://www.pgadmin.org/docs/
- **PostgreSQL Docs:** https://www.postgresql.org/docs/

---

**🎉 Artık PgAdmin ile veritabanınızı görsel olarak yönetebilirsiniz!**
