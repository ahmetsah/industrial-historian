# 100+ Modbus Slave Yönetimi Kılavuzu

## 📊 Ölçeklenebilirlik Analizi

### Mevcut Mimari Kapasitesi

| Slave Sayısı | Durum | Gerekli Değişiklik |
|--------------|-------|-------------------|
| 1-10 | ✅ Sorunsuz | Yok |
| 10-30 | ✅ İyi | Buffer size artırımı |
| 30-50 | ⚠️ Dikkat | Buffer + Config yönetimi |
| 50-100 | ⚠️ Sınırda | Connection pooling önerilir |
| 100+ | ❌ Yeniden tasarım | Mikroservis mimarisi |

---

## 🎯 Yaklaşım 1: Basit (10-30 Slave)

### ✅ Yapılan Değişiklik
```rust
// main.rs - Dinamik buffer size
let buffer_size = (100 * device_count).max(1000);
```

### Kullanım
```toml
# config/default.toml
[[modbus_devices]]
ip = "192.168.1.10"
# ... (30 cihaza kadar manuel eklenebilir)
```

### Avantajlar
- ✅ Kolay yönetim
- ✅ Kod değişikliği minimal
- ✅ Debug kolay

### Dezavantajlar
- ❌ Config dosyası şişer
- ❌ 30+ cihazda yönetim zorlaşır

---

## 🎯 Yaklaşım 2: CSV Config (30-100 Slave)

### Config Yapısı
```csv
# config/modbus_devices.csv
ip,port,unit_id,poll_ms,addr,name,type
192.168.1.10,502,1,1000,0,F1.L1.PLC01.T001,Float32
192.168.1.11,502,1,1000,0,F1.L1.PLC02.T001,Float32
# ... 100 satır
```

### CSV Parser Ekle

**Dosya:** `services/ingestor/src/config.rs`

```rust
use csv::ReaderBuilder;
use std::fs::File;

#[derive(Debug, Deserialize)]
struct CsvDevice {
    ip: String,
    port: u16,
    unit_id: u8,
    poll_ms: u64,
    addr: u16,
    name: String,
    #[serde(rename = "type")]
    data_type: String,
}

impl Settings {
    pub fn from_csv(path: &str) -> Result<Self, ConfigError> {
        let file = File::open(path)?;
        let mut reader = ReaderBuilder::new()
            .has_headers(true)
            .from_reader(file);
        
        let mut devices_map: HashMap<String, ModbusConfig> = HashMap::new();
        
        for result in reader.deserialize() {
            let record: CsvDevice = result?;
            let key = format!("{}:{}", record.ip, record.port);
            
            let device = devices_map.entry(key.clone()).or_insert_with(|| {
                ModbusConfig {
                    ip: record.ip.clone(),
                    port: record.port,
                    unit_id: record.unit_id,
                    poll_interval_ms: record.poll_ms,
                    registers: Vec::new(),
                }
            });
            
            device.registers.push(RegisterConfig {
                address: record.addr,
                name: record.name,
                data_type: record.data_type,
            });
        }
        
        Ok(Settings {
            modbus_devices: devices_map.into_values().collect(),
            // ... diğer alanlar
        })
    }
}
```

### Kullanım
```rust
// main.rs
let settings = if Path::new("config/modbus_devices.csv").exists() {
    Settings::from_csv("config/modbus_devices.csv")?
} else {
    Settings::new()?
};
```

### Avantajlar
- ✅ Excel ile düzenlenebilir
- ✅ Toplu import/export kolay
- ✅ 100 cihaz rahatça yönetilebilir

### Dezavantajlar
- ⚠️ Ekstra parsing kodu gerekli
- ⚠️ TOML'den farklı format

---

## 🎯 Yaklaşım 3: Connection Pool (50-100+ Slave)

### Sorun
100 cihaz = 100 açık TCP connection = Yüksek kaynak kullanımı

### Çözüm: Semaphore ile Sınırlama

```rust
// Max 20 concurrent connection
let semaphore = Arc::new(Semaphore::new(20));

for device in devices {
    let sem = semaphore.clone();
    tokio::spawn(async move {
        loop {
            let _permit = sem.acquire().await.unwrap();
            // Poll device
            poll_once(&device).await;
            // Permit auto-release
        }
    });
}
```

### Performans Kazancı
- **Öncesi:** 100 connection × 4KB = 400KB RAM
- **Sonrası:** 20 connection × 4KB = 80KB RAM
- **Kazanç:** %80 RAM tasarrufu

### Avantajlar
- ✅ Düşük kaynak kullanımı
- ✅ Network stack'e daha az yük
- ✅ 100+ cihaz destekler

### Dezavantajlar
- ⚠️ Polling latency artar (sıra bekler)
- ⚠️ Kod karmaşıklığı artar

---

## 🎯 Yaklaşım 4: Mikroservis (100+ Slave)

### Mimari
```
Ingestor 1 (PLC 1-20)  ──┐
Ingestor 2 (PLC 21-40) ──┼──> NATS ──> Engine
Ingestor 3 (PLC 41-60) ──┤
Ingestor 4 (PLC 61-80) ──┤
Ingestor 5 (PLC 81-100)──┘
```

### Docker Compose
```yaml
services:
  ingestor-group1:
    image: historian-ingestor
    environment:
      - CONFIG_FILE=/config/group1.toml
    volumes:
      - ./config/group1.toml:/config/group1.toml

  ingestor-group2:
    image: historian-ingestor
    environment:
      - CONFIG_FILE=/config/group2.toml
    volumes:
      - ./config/group2.toml:/config/group2.toml
  
  # ... 5 instance toplam
```

### Avantajlar
- ✅ Tam izolasyon (bir grup düşerse diğerleri çalışır)
- ✅ Bağımsız ölçekleme
- ✅ Kolay deployment

### Dezavantajlar
- ❌ Daha fazla resource (5× container)
- ❌ Orchestration gerekli (K8s/Compose)

---

## 📊 Performans Hesaplamaları

### Senaryo: 100 Slave, Her biri 10 register, 1 saniye poll

**Veri Akışı:**
- 100 slave × 10 register = 1000 tag
- 1000 tag × 1 poll/saniye = **1000 mesaj/saniye**
- 1000 mesaj × 100 byte = **100 KB/saniye**

**Kaynak Kullanımı (Yaklaşım 1):**
- RAM: ~50MB (100 task + buffer)
- CPU: ~10% (single core)
- Network: 100 KB/s

**Kaynak Kullanımı (Yaklaşım 3 - Pool):**
- RAM: ~20MB (20 connection + buffer)
- CPU: ~8% (daha az context switch)
- Network: 100 KB/s (aynı)

---

## 🛠️ Önerilen Strateji (100 Slave için)

### Aşama 1: Hızlı Başlangıç (1 gün)
✅ **Yapıldı:** Dinamik buffer size
```rust
let buffer_size = (100 * device_count).max(1000);
```

### Aşama 2: Config Yönetimi (2 gün)
🔨 **Yapılacak:** CSV parser ekle
- `config.rs`'e `from_csv()` metodu
- Excel template oluştur
- Test et

### Aşama 3: Optimizasyon (3 gün)
🔨 **Yapılacak:** Connection pooling
- `modbus_pool.rs` modülü
- Semaphore ile 20 concurrent limit
- Benchmark yap

### Aşama 4: Production (1 hafta)
🔨 **Yapılacak:** Monitoring ekle
- Prometheus metrics (cihaz başına)
- Grafana dashboard
- Alert rules

---

## 🧪 Test Planı

### 1. Stress Test
```bash
# 100 Modbus simulator başlat
for i in {1..100}; do
  python scripts/modbus_simulator.py --port $((5000 + i)) &
done

# Ingestor başlat
cargo run --release

# Metrics izle
watch -n 1 'ps aux | grep ingestor'
```

### 2. Beklenen Sonuçlar
- ✅ CPU < 20%
- ✅ RAM < 100MB
- ✅ Mesaj kaybı yok
- ✅ Tüm cihazlar bağlanıyor

---

## 📋 Karar Matrisi

| Slave Sayısı | Önerilen Yaklaşım | Uygulama Süresi |
|--------------|-------------------|-----------------|
| 1-10 | Mevcut (TOML) | ✅ Hazır |
| 10-30 | TOML + Dinamik buffer | ✅ Hazır |
| 30-50 | CSV Config | 🔨 2 gün |
| 50-100 | CSV + Connection Pool | 🔨 5 gün |
| 100+ | Mikroservis | 🔨 2 hafta |

---

## 🚀 Hemen Yapılabilecekler

### Seçenek A: Basit (30 Slave'e kadar)
```bash
# Hiçbir şey yapma, mevcut sistem yeterli!
# Sadece config/default.toml'e cihaz ekle
```

### Seçenek B: CSV (100 Slave)
```bash
# 1. CSV dosyası oluştur
vim config/modbus_devices.csv

# 2. CSV parser ekle (kod örnekleri yukarıda)

# 3. Test et
cargo run
```

### Seçenek C: Mikroservis (100+ Slave)
```bash
# 1. Config'i 5 gruba böl
# 2. Docker Compose güncelle
# 3. Deploy et
docker-compose up -d --scale ingestor=5
```

---

## ❓ Hangi Yaklaşımı Seçmeliyim?

**Sorular:**
1. Kaç slave planlıyorsunuz? → **100**
2. Hepsi aynı ağda mı? → **Evet/Hayır**
3. Kritik sistem mi? (Downtime kabul edilemez) → **Evet/Hayır**
4. Geliştirme süresi kısıtlı mı? → **Evet/Hayır**

**Cevaplarınıza göre:**
- Evet, Evet, Hayır, Evet → **Yaklaşım 2 (CSV)**
- Evet, Evet, Evet, Hayır → **Yaklaşım 3 (Pool)**
- Evet, Hayır, Evet, Hayır → **Yaklaşım 4 (Mikroservis)**

---

## 📞 Sonuç

**Kısa Cevap:** Evet, 100 slave ekleyebilirsiniz!

**Uzun Cevap:** 
- ✅ **10-30 slave:** Şu anki sistem hazır
- ⚠️ **30-100 slave:** CSV config + connection pool önerilir
- 🔨 **100+ slave:** Mikroservis mimarisi gerekli

**Bir sonraki adım ne olsun?**
1. CSV parser'ı implement edeyim mi?
2. Connection pool örneği hazırlayayım mı?
3. Mikroservis setup'ı göstereyim mi?

Hangisini istersiniz? 🚀
