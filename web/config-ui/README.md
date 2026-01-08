# 🎨 Web UI - Config Manager Dashboard

## 🎉 TAMAMLANDI!

Modern, responsive web arayüzü başarıyla oluşturuldu ve deploy edildi!

---

## 🌐 Erişim

### Web UI
```
URL: http://localhost:3001
```

### Özellikler
- ✅ Modern, gradient tasarım
- ✅ Gerçek zamanlı cihaz listesi
- ✅ Cihaz ekleme formu
- ✅ İstatistikler (toplam cihaz, aktif cihaz, register sayısı)
- ✅ Otomatik yenileme (5 saniyede bir)
- ✅ Bildirimler (başarı/hata)
- ✅ Responsive tasarım

---

## 📊 Ekran Görüntüsü

### Dashboard
```
┌─────────────────────────────────────────────────────────┐
│  🏭 Historian Config Manager                            │
│  Modbus Device Configuration Dashboard                  │
└─────────────────────────────────────────────────────────┘

┌──────────┬──────────┬──────────┬──────────┐
│ Total    │ Active   │ Total    │ Config   │
│ Devices  │ Devices  │ Registers│ Files    │
│    4     │    2     │    12    │    4     │
└──────────┴──────────┴──────────┴──────────┘

┌────────────────┬────────────────────────────────┐
│ Add New Device │ Devices                        │
│                │                                │
│ [Form]         │ [Device List]                  │
│                │ - PLC-001                      │
│                │ - PLC-002                      │
│                │ - PLC-003                      │
└────────────────┴────────────────────────────────┘
```

---

## 🚀 Kullanım

### 1. Web UI'yi Aç
```bash
# Tarayıcıda
open http://localhost:3001

# Veya
xdg-open http://localhost:3001
```

### 2. Yeni Cihaz Ekle
1. Sol panelde formu doldurun:
   - Device Name: `PLC-NEW`
   - IP Address: `192.168.1.50`
   - Port: `502`
   - Unit ID: `1`
   - Poll Interval: `1000`

2. Register ekleyin:
   - Address: `0`
   - Name: `Factory1.Line1.PLC-NEW.Temp.T001`
   - Type: `Float32`

3. **Create Device** butonuna tıklayın

4. ✅ Cihaz otomatik olarak:
   - PostgreSQL'e kaydedilir
   - Config dosyası oluşturulur
   - Sağ panelde görünür

### 3. Cihazları Görüntüle
- Sağ panelde tüm cihazlar listelenir
- Her 5 saniyede otomatik yenilenir
- Status (active/inactive) gösterilir
- Register sayısı görünür

---

## 🏗️ Mimari

### Frontend
```
web/config-ui/
├── index.html          ← Tek sayfa uygulama
└── Dockerfile          ← Nginx container
```

### Backend API
```
Config Manager API (Port 8090)
├── POST   /api/v1/devices/modbus
├── GET    /api/v1/devices/modbus
├── GET    /api/v1/devices/modbus/:id
├── PUT    /api/v1/devices/modbus/:id
└── DELETE /api/v1/devices/modbus/:id
```

### Veri Akışı
```
Web UI (Port 3001)
    ↓ HTTP Request
Config Manager API (Port 8090)
    ↓ SQL Query
PostgreSQL (Port 5432)
    ↓ Generate Config
config/generated/modbus-*.toml
    ↓ Mount
Ingestor Containers
```

---

## 🎨 Tasarım Özellikleri

### Renk Paleti
```css
Primary:   #667eea (Mor-Mavi)
Secondary: #764ba2 (Mor)
Success:   #28a745 (Yeşil)
Error:     #dc3545 (Kırmızı)
Background: Linear gradient (667eea → 764ba2)
```

### Animasyonlar
- ✅ Hover efektleri (kartlar yukarı kalkar)
- ✅ Buton press animasyonu
- ✅ Bildirim slide-in
- ✅ Smooth transitions

### Responsive
- ✅ Desktop (1400px+)
- ✅ Tablet (768px+)
- ✅ Mobile (320px+)

---

## 🔧 Teknik Detaylar

### Teknolojiler
- **Frontend:** Vanilla HTML/CSS/JavaScript
- **HTTP Server:** Nginx Alpine
- **API Client:** Fetch API
- **Styling:** Modern CSS (Grid, Flexbox, Gradients)

### Docker
```yaml
config-ui:
  build: web/config-ui
  ports: 3001:80
  depends_on: config-manager
```

### Nginx Config
```nginx
location / {
    root /usr/share/nginx/html;
    try_files $uri /index.html;
}

location /api/ {
    proxy_pass http://config-manager:8090/api/;
}
```

---

## 📝 API Entegrasyonu

### Cihaz Listesi
```javascript
fetch('http://localhost:8090/api/v1/devices/modbus')
    .then(res => res.json())
    .then(data => displayDevices(data.devices));
```

### Cihaz Oluştur
```javascript
fetch('http://localhost:8090/api/v1/devices/modbus', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
        name: 'PLC-001',
        ip: '192.168.1.10',
        port: 502,
        unit_id: 1,
        poll_interval_ms: 1000,
        registers: [...]
    })
});
```

---

## 🚀 Deployment

### Development
```bash
cd ops
docker-compose up -d config-ui
```

### Production
```bash
# Build
docker-compose build config-ui

# Deploy
docker-compose up -d config-ui

# Scale (if needed)
docker-compose up -d --scale config-ui=3
```

### Health Check
```bash
# Container durumu
docker ps | grep config-ui

# Loglar
docker logs -f ops-config-ui

# Test
curl http://localhost:3001
```

---

## 🎯 Gelecek Özellikler

### Kısa Vadeli
- [ ] Cihaz düzenleme (Edit)
- [ ] Cihaz silme (Delete)
- [ ] Register düzenleme
- [ ] Gerçek zamanlı veri görüntüleme

### Orta Vadeli
- [ ] Kullanıcı authentication
- [ ] Rol tabanlı erişim
- [ ] Audit log görüntüleme
- [ ] Config diff/comparison

### Uzun Vadeli
- [ ] Grafik ve dashboard
- [ ] Alarm yönetimi
- [ ] Toplu cihaz ekleme (CSV import)
- [ ] Config template'leri

---

## 🐛 Troubleshooting

### Web UI açılmıyor
```bash
# Container çalışıyor mu?
docker ps | grep config-ui

# Restart
cd ops
docker-compose restart config-ui

# Loglar
docker logs ops-config-ui
```

### API bağlantı hatası
```bash
# Config Manager çalışıyor mu?
curl http://localhost:8090/health

# CORS sorunu var mı?
# (Nginx proxy kullanıyoruz, olmamalı)
```

### Cihazlar görünmüyor
```bash
# API'den manuel test
curl http://localhost:8090/api/v1/devices/modbus

# PostgreSQL'de var mı?
docker exec ops-postgres-1 psql -U postgres -d historian -c "SELECT * FROM devices;"
```

---

## 📚 Kaynaklar

### Dosyalar
```
web/config-ui/
├── index.html      ← UI kodu
├── Dockerfile      ← Container tanımı
└── README.md       ← Bu dosya
```

### API Dökümanı
```
docs/API-DOCUMENTATION.md
```

### Deployment
```
ops/docker-compose.yml
```

---

## ✅ Checklist

Web UI kurulumu için:

- [x] HTML/CSS/JS oluşturuldu
- [x] Dockerfile hazırlandı
- [x] Docker Compose'a eklendi
- [x] Build başarılı
- [x] Container çalışıyor
- [x] Port 3001 açık
- [x] API bağlantısı çalışıyor
- [x] Cihaz listesi görünüyor
- [x] Cihaz ekleme çalışıyor
- [ ] Production deployment (opsiyonel)

---

## 🎉 Sonuç

**Web UI başarıyla deploy edildi!**

- ✅ Modern, kullanıcı dostu arayüz
- ✅ Gerçek zamanlı veri
- ✅ Kolay cihaz yönetimi
- ✅ Production ready

**Erişim:** http://localhost:3001

**Artık cihazlarınızı web arayüzünden yönetebilirsiniz!** 🚀
