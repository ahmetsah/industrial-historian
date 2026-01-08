# 🚀 Historian Web UI & Deployment Guide

Bu doküman, Web UI'a eklenen **Dinamik Deployment**, **API-Driven Configuration** ve **Gelişmiş Cihaz Yönetimi** özelliklerini detaylandırır.

---

## ✨ Yeni Özellikler

### 1. 🚀 Deploy Butonu
Web UI üzerinden cihazları tek tıkla Docker container olarak ayağa kaldırabilirsiniz.
- **Fonksiyon:**
    - **İlk Deploy:** Cihaz için henüz bir container yoksa, otomatik olarak oluşturur ve başlatır.
    - **Hot Reload:** Cihaz çalışıyorsa, container'ı **restart** ederek yeni konfigürasyonu yükler.
- **Teknoloji:** Backend, Docker socket üzerinden orkestrasyon yapar.

### 2. 🗑️ Gelişmiş Silme (Cleanup)
- **Database:** Cihaz kaydını siler.
- **Docker:** Çalışan ilgili Ingestor container'ını durdurur ve siler (`docker rm -f`).
- **Sonuç:** Sistemde yetim container kalmaz.

### 3. 🌐 API-Driven Configuration
- Ingestor, `CONFIG_URL` ortam değişkenini kullanarak Config Manager API'sından JSON config çeker.
- Fiziksel dosya yönetimine gerek kalmaz.

### 4. ✏️ Cihaz Düzenleme (Edit)
Mevcut cihazların güncellenmesi artık çok kolay.
- **Edit Modu:** Kalem (✏️) ikonuna tıklayınca form cihaz bilgileriyle dolar.
- **Register Yönetimi:** Form üzerinde dinamik olarak register ekleyebilir, silebilir ve düzenleyebilirsiniz.
- **Update:** Değişiklikler `PUT` isteği ile sunucuya gönderilir.

### 5. 📊 Detaylı Durum Takibi (Status)
Cihazların durumu artık iki ayrı gösterge ile takip edilebilir:

#### A. 📦 Deployment Status
Cihazın Docker Container durumunu gösterir.
- **Deployed:** Container başarıyla oluşturuldu ve çalışıyor (Running).
- **Not Deployed:** Container durduruldu veya henüz oluşturulmadı.

#### B. 🔌 Connection Status
Config Manager'ın hedef cihaza (PLC/Modbus Device) erişip erişemediğini gösterir.
- **Connected (Mavi):** Hedef IP ve Port'a başarılı TCP bağlantısı kurulabiliyor.
- **Disconnected (Kırmızı):** Hedef IP'ye ulaşılamıyor veya bağlantı reddedildi.
- **Idle (Sarı):** Cihaz henüz deploy edilmediği için bağlantı kontrolü yapılmıyor.

---

## 🛠️ Mimari Değişiklikler

### Config Manager (Go)
- **Docker Integration:** `docker-cli` ve socket mount ile container yönetimi.
- **Status Checks:** 
    - `docker ps` ile deployment kontrolü.
    - `net.Dial` ile TCP bağlantı kontrolü.
- **Dynamic Config API:** Ingestor'lar için JSON konfigürasyon sunar.

### Ingestor (Rust)
- **Reqwest Client:** Konfigürasyonu HTTP üzerinden asenkron olarak çeker.
- **Hot Reload:** Container restart edildiğinde yeni config ile başlar.

---

## 🧪 Kullanım Senaryosu

1. **Ekle:** "Add New Device" formu ile cihazı oluşturun. Status: **Not Deployed / Idle**.
2. **Deploy (🚀):** Roket ikonuna basın. 
    - Container başlar -> **Deployed**.
    - Bağlantı sağlanırsa -> **Connected**. Bağlantı yoksa -> **Disconnected**.
3. **Düzenle (✏️):** Register eklemek veya IP düzeltmek için kullanın.
4. **Uygula:** Değişikliğin hemen yansıması için tekrar Deploy (🚀) butonuna basın.
5. **Sil (🗑️):** Cihazı ve container'ı kalıcı olarak siler.
