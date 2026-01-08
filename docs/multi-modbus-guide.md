# Multi-Modbus Client Implementation Guide

## Yaklaşım 1: Tek Ingestor, Çoklu Adapter (Önerilen)

### Adım 1: Config Yapısını Güncelle

**Dosya:** `services/ingestor/src/config.rs`

```rust
#[derive(Debug, Deserialize)]
pub struct Settings {
    pub modbus_devices: Vec<ModbusConfig>,  // ← Tekil yerine çoğul
    pub nats: NatsConfig,
    pub buffer: BufferConfig,
    #[serde(default)]
    pub calculated_tags: Vec<CalculatedTagConfig>,
}
```

### Adım 2: Main.rs'i Güncelle

**Dosya:** `services/ingestor/src/main.rs`

```rust
#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    dotenv().ok();
    tracing_subscriber::fmt::init();
    info!("Starting Ingestor Service");

    let settings = Settings::new()?;
    let buffer = HybridBuffer::new(
        settings.buffer.memory_capacity, 
        settings.buffer.disk_path
    );
    let mut publisher = Publisher::new(buffer, settings.nats.url.clone());

    // Tek kanal - tüm adapter'lar buraya yazar
    let (tx_raw, mut rx_raw) = mpsc::channel::<SensorData>(1000);
    let (tx_pub, rx_pub) = mpsc::channel::<SensorData>(1000);

    // Publisher task
    tokio::spawn(async move {
        publisher.run(rx_pub).await;
    });

    // Engine task
    let mut engine = Engine::new(settings.calculated_tags);
    tokio::spawn(async move {
        while let Some(data) = rx_raw.recv().await {
            let results = engine.process(data);
            for res in results {
                if let Err(e) = tx_pub.send(res).await {
                    error!("Failed to send to publisher: {}", e);
                    break;
                }
            }
        }
    });

    // ✨ HER MODBUS CİHAZI İÇİN AYRI ADAPTER BAŞLAT
    for modbus_config in settings.modbus_devices {
        let tx_clone = tx_raw.clone();
        let device_name = modbus_config.ip.clone();
        
        tokio::spawn(async move {
            info!("Starting Modbus adapter for {}", device_name);
            let mut adapter = ModbusAdapter::new(modbus_config, tx_clone);
            
            if let Err(e) = adapter.connect().await {
                error!("Initial connection failed for {}: {}", device_name, e);
            }
            adapter.poll_loop().await;
        });
    }

    info!("Ingestor running with {} Modbus devices", settings.modbus_devices.len());
    tokio::signal::ctrl_c().await?;
    info!("Shutting down...");
    Ok(())
}
```

### Adım 3: Config Dosyası Örneği

**Dosya:** `config/default.toml`

```toml
# Cihaz 1: Ana Hat PLC
[[modbus_devices]]
ip = "192.168.1.10"
port = 502
unit_id = 1
poll_interval_ms = 1000

[[modbus_devices.registers]]
address = 0
name = "Factory1.Line1.Mixer.Temp.T001"
data_type = "Float32"

[[modbus_devices.registers]]
address = 2
name = "Factory1.Line1.Mixer.Pressure.P001"
data_type = "Int16"

# Cihaz 2: Yardımcı Hat PLC
[[modbus_devices]]
ip = "192.168.1.20"
port = 502
unit_id = 1
poll_interval_ms = 2000

[[modbus_devices.registers]]
address = 0
name = "Factory1.Line2.Pump.Speed.S001"
data_type = "UInt16"

# Cihaz 3: Kalite Kontrol PLC
[[modbus_devices]]
ip = "192.168.1.30"
port = 502
unit_id = 2
poll_interval_ms = 5000

[[modbus_devices.registers]]
address = 100
name = "Factory1.QC.Scale.Weight.W001"
data_type = "Float32"

[nats]
url = "nats://localhost:4222"
subject = "data.raw"

[buffer]
memory_capacity = 10000
disk_path = "ops/data/ingestor_wal/buffer.wal"
```

---

## Yaklaşım 2: Her Cihaz İçin Ayrı Ingestor Instance (Mikroservis)

### Docker Compose ile Çoklu Instance

**Dosya:** `docker-compose.yml`

```yaml
services:
  ingestor-plc1:
    build: ./services/ingestor
    environment:
      - CONFIG_FILE=/config/plc1.toml
    volumes:
      - ./config/plc1.toml:/config/plc1.toml
    depends_on:
      - nats

  ingestor-plc2:
    build: ./services/ingestor
    environment:
      - CONFIG_FILE=/config/plc2.toml
    volumes:
      - ./config/plc2.toml:/config/plc2.toml
    depends_on:
      - nats

  ingestor-plc3:
    build: ./services/ingestor
    environment:
      - CONFIG_FILE=/config/plc3.toml
    volumes:
      - ./config/plc3.toml:/config/plc3.toml
    depends_on:
      - nats
```

**Her config dosyası tek bir cihaz içerir:**

`config/plc1.toml`:
```toml
[modbus]
ip = "192.168.1.10"
port = 502
unit_id = 1
poll_interval_ms = 1000

[[modbus.registers]]
address = 0
name = "Factory1.Line1.Mixer.Temp.T001"
data_type = "Float32"
```

---

## 🎯 Önerilen Yaklaşım: Yaklaşım 1

**Neden?**
- ✅ Daha az kaynak kullanımı (tek process)
- ✅ Merkezi konfigürasyon yönetimi
- ✅ Kolay debug ve monitoring
- ✅ Shared buffer ve publisher (verimli)

**Ne zaman Yaklaşım 2?**
- Cihazlar farklı ağlarda
- Bağımsız ölçekleme gerekli
- Fault isolation kritik

---

## 📝 Uygulama Adımları

1. **Config struct'ı güncelle** (`config.rs`)
2. **Main.rs'i güncelle** (loop ile adapter spawn)
3. **Config dosyasını güncelle** (array formatına)
4. **Test et:**
   ```bash
   cd services/ingestor
   cargo run
   ```

---

## 🧪 Test Senaryosu

```bash
# 1. Modbus Simulator başlat (3 farklı port)
python scripts/modbus_simulator.py --port 5020
python scripts/modbus_simulator.py --port 5021
python scripts/modbus_simulator.py --port 5022

# 2. Config'i güncelle
# 3. Ingestor'u başlat
cargo run --bin ingestor

# 4. NATS'i izle
nats sub "data.>"
```

Beklenen çıktı:
```
[data.raw] Factory1.Line1.Mixer.Temp.T001: 23.5
[data.raw] Factory1.Line2.Pump.Speed.S001: 1450
[data.raw] Factory1.QC.Scale.Weight.W001: 125.3
```

---

## ⚠️ Dikkat Edilmesi Gerekenler

1. **Tag Naming:** Her cihaz için unique prefix kullanın
   - ✅ `PLC1.Line1.Temp.T001`
   - ❌ `Line1.Temp.T001` (conflict riski)

2. **Poll Interval:** Cihaz yüküne göre ayarlayın
   - Kritik: 100-500ms
   - Normal: 1000-2000ms
   - Yavaş: 5000ms+

3. **Channel Buffer Size:** Cihaz sayısına göre artırın
   ```rust
   let (tx_raw, rx_raw) = mpsc::channel::<SensorData>(
       100 * settings.modbus_devices.len()
   );
   ```

4. **Error Handling:** Her adapter bağımsız fail olmalı
   - Bir cihaz düşerse diğerleri çalışmaya devam etmeli

---

## 🚀 Hızlı Başlangıç

İsterseniz size hazır kod değişikliklerini uygulayayım:
- [ ] `config.rs` güncelleme
- [ ] `main.rs` güncelleme  
- [ ] Örnek multi-device config dosyası
- [ ] Test script'i

Devam edelim mi?
