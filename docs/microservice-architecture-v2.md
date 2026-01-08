# Industrial Historian - Mikroservis Mimarisi v2.0

## 📐 Mimari Genel Bakış

```
┌─────────────────────────────────────────────────────────────┐
│                    Web Management UI                         │
│              (React - Config Management)                     │
└────────────────────┬────────────────────────────────────────┘
                     │ HTTP/REST
                     ↓
┌─────────────────────────────────────────────────────────────┐
│              Config Management Service (Go)                  │
│    - Device CRUD API                                        │
│    - Config validation                                       │
│    - Template management                                     │
└────────────────────┬────────────────────────────────────────┘
                     │ PostgreSQL
                     ↓
┌─────────────────────────────────────────────────────────────┐
│                  Device Configurations DB                    │
│    - modbus_devices                                         │
│    - opc_devices                                            │
│    - s7_devices                                             │
└─────────────────────────────────────────────────────────────┘
                     │
                     │ Config Files (Generated)
                     ↓
┌──────────────┬──────────────┬──────────────┬───────────────┐
│   Modbus     │     OPC      │    S7Net     │   MQTT        │
│  Ingestor    │  Ingestor    │  Ingestor    │  Ingestor     │
│  (Rust)      │  (Rust)      │  (Rust)      │  (Rust)       │
└──────┬───────┴──────┬───────┴──────┬───────┴──────┬────────┘
       │              │              │              │
       └──────────────┴──────────────┴──────────────┘
                      │
                      ↓ NATS JetStream
       ┌──────────────────────────────────┐
       │      Message Bus (NATS)          │
       │   Subject: data.{protocol}.{id}  │
       └──────────────┬───────────────────┘
                      ↓
       ┌──────────────────────────────────┐
       │        Engine Service             │
       │    (Storage + Query + Export)     │
       └──────────────────────────────────┘
```

## 🎯 Tasarım Prensipleri

1. **Protocol Agnostic:** Her protokol için ayrı ingestor
2. **Config as Data:** Tüm konfigürasyon DB'de
3. **Web-First:** UI üzerinden tam kontrol
4. **Hot Reload:** Restart gerektirmeden config değişikliği
5. **Scalable:** Her ingestor bağımsız ölçeklenebilir

---

## 📁 Yeni Proje Yapısı

```
historian/
├── services/
│   ├── ingestor-modbus/        # Modbus Protocol Ingestor
│   │   ├── Dockerfile
│   │   ├── Cargo.toml
│   │   └── src/
│   │       ├── main.rs
│   │       ├── adapter.rs
│   │       └── config.rs
│   │
│   ├── ingestor-opc/           # OPC UA Protocol Ingestor
│   │   ├── Dockerfile
│   │   ├── Cargo.toml
│   │   └── src/
│   │       ├── main.rs
│   │       ├── adapter.rs
│   │       └── config.rs
│   │
│   ├── ingestor-s7/            # Siemens S7 Protocol Ingestor
│   │   ├── Dockerfile
│   │   ├── Cargo.toml
│   │   └── src/
│   │       ├── main.rs
│   │       ├── adapter.rs
│   │       └── config.rs
│   │
│   ├── config-manager/         # Config Management Service (Go)
│   │   ├── Dockerfile
│   │   ├── go.mod
│   │   └── internal/
│   │       ├── api/            # REST API handlers
│   │       ├── models/         # Device models
│   │       ├── repository/     # DB access
│   │       └── generator/      # Config file generator
│   │
│   └── engine/                 # Existing Engine Service
│
├── web/                        # Management Web UI
│   ├── config-ui/
│   │   ├── package.json
│   │   └── src/
│   │       ├── pages/
│   │       │   ├── DeviceList.tsx
│   │       │   ├── ModbusConfig.tsx
│   │       │   ├── OPCConfig.tsx
│   │       │   └── S7Config.tsx
│   │       └── api/
│   │           └── configApi.ts
│   │
│   └── viz/                    # Existing Visualization UI
│
├── config/
│   ├── generated/              # Auto-generated configs
│   │   ├── modbus-instance-1.toml
│   │   ├── modbus-instance-2.toml
│   │   ├── opc-instance-1.toml
│   │   └── s7-instance-1.toml
│   │
│   └── templates/              # Config templates
│       ├── modbus.template.toml
│       ├── opc.template.toml
│       └── s7.template.toml
│
└── ops/
    ├── docker-compose.yml      # Orchestration
    └── k8s/                    # Kubernetes manifests
        ├── ingestor-modbus-deployment.yaml
        ├── ingestor-opc-deployment.yaml
        └── config-manager-deployment.yaml
```

---

## 🗄️ Database Schema

```sql
-- Device Configurations Database

-- Protocol types
CREATE TYPE protocol_type AS ENUM ('modbus', 'opc', 's7', 'mqtt');
CREATE TYPE device_status AS ENUM ('active', 'inactive', 'error');

-- Main devices table
CREATE TABLE devices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    protocol protocol_type NOT NULL,
    status device_status DEFAULT 'inactive',
    config JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    created_by VARCHAR(255),
    UNIQUE(name)
);

-- Modbus-specific config
CREATE TABLE modbus_devices (
    id UUID PRIMARY KEY REFERENCES devices(id) ON DELETE CASCADE,
    ip VARCHAR(45) NOT NULL,
    port INTEGER DEFAULT 502,
    unit_id INTEGER NOT NULL,
    poll_interval_ms INTEGER DEFAULT 1000,
    timeout_ms INTEGER DEFAULT 5000,
    retry_count INTEGER DEFAULT 3
);

CREATE TABLE modbus_registers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id UUID REFERENCES modbus_devices(id) ON DELETE CASCADE,
    address INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL,
    data_type VARCHAR(50) NOT NULL,
    scale_factor FLOAT DEFAULT 1.0,
    offset FLOAT DEFAULT 0.0,
    unit VARCHAR(50),
    description TEXT
);

-- OPC UA-specific config
CREATE TABLE opc_devices (
    id UUID PRIMARY KEY REFERENCES devices(id) ON DELETE CASCADE,
    endpoint_url VARCHAR(512) NOT NULL,
    security_mode VARCHAR(50) DEFAULT 'None',
    security_policy VARCHAR(50) DEFAULT 'None',
    username VARCHAR(255),
    password_encrypted TEXT,
    poll_interval_ms INTEGER DEFAULT 1000
);

CREATE TABLE opc_nodes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id UUID REFERENCES opc_devices(id) ON DELETE CASCADE,
    node_id VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    data_type VARCHAR(50),
    unit VARCHAR(50),
    description TEXT
);

-- S7-specific config
CREATE TABLE s7_devices (
    id UUID PRIMARY KEY REFERENCES devices(id) ON DELETE CASCADE,
    ip VARCHAR(45) NOT NULL,
    rack INTEGER DEFAULT 0,
    slot INTEGER DEFAULT 1,
    poll_interval_ms INTEGER DEFAULT 1000
);

CREATE TABLE s7_data_blocks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id UUID REFERENCES s7_devices(id) ON DELETE CASCADE,
    db_number INTEGER NOT NULL,
    start_address INTEGER NOT NULL,
    length INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL,
    data_type VARCHAR(50) NOT NULL,
    unit VARCHAR(50),
    description TEXT
);

-- Config generation tracking
CREATE TABLE config_generations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id UUID REFERENCES devices(id) ON DELETE CASCADE,
    config_hash VARCHAR(64) NOT NULL,
    file_path VARCHAR(512) NOT NULL,
    generated_at TIMESTAMP DEFAULT NOW(),
    deployed_at TIMESTAMP,
    status VARCHAR(50) DEFAULT 'pending'
);

-- Indexes
CREATE INDEX idx_devices_protocol ON devices(protocol);
CREATE INDEX idx_devices_status ON devices(status);
CREATE INDEX idx_modbus_ip_port ON modbus_devices(ip, port);
CREATE INDEX idx_opc_endpoint ON opc_devices(endpoint_url);
```

---

## 🔧 Config Manager Service API

### REST Endpoints

```
# Device Management
GET    /api/v1/devices                    # List all devices
GET    /api/v1/devices/:id                # Get device details
POST   /api/v1/devices                    # Create device
PUT    /api/v1/devices/:id                # Update device
DELETE /api/v1/devices/:id                # Delete device

# Protocol-specific
GET    /api/v1/devices/modbus             # List Modbus devices
POST   /api/v1/devices/modbus             # Create Modbus device
GET    /api/v1/devices/opc                # List OPC devices
POST   /api/v1/devices/opc                # Create OPC device
GET    /api/v1/devices/s7                 # List S7 devices
POST   /api/v1/devices/s7                 # Create S7 device

# Config Generation
POST   /api/v1/config/generate/:id        # Generate config file
POST   /api/v1/config/deploy/:id          # Deploy config (trigger restart)
GET    /api/v1/config/validate            # Validate config JSON

# Templates
GET    /api/v1/templates/:protocol        # Get config template
POST   /api/v1/templates/:protocol        # Save custom template

# Health & Status
GET    /api/v1/health                     # Service health
GET    /api/v1/devices/:id/status         # Device connection status
```

---

## 🐳 Docker Compose Orchestration

```yaml
version: '3.8'

services:
  # Config Management
  config-manager:
    build: ./services/config-manager
    ports:
      - "8090:8090"
    environment:
      - DATABASE_URL=postgres://historian:password@postgres:5432/historian_config
      - CONFIG_DIR=/config/generated
    volumes:
      - ./config/generated:/config/generated
    depends_on:
      - postgres

  # Modbus Ingestors (Auto-scaled)
  ingestor-modbus:
    build: ./services/ingestor-modbus
    environment:
      - NATS_URL=nats://nats:4222
      - CONFIG_FILE=/config/modbus-${INSTANCE_ID}.toml
      - INSTANCE_ID=${INSTANCE_ID:-1}
    volumes:
      - ./config/generated:/config:ro
    depends_on:
      - nats
      - config-manager
    deploy:
      replicas: 3  # 3 Modbus ingestor instance

  # OPC UA Ingestor
  ingestor-opc:
    build: ./services/ingestor-opc
    environment:
      - NATS_URL=nats://nats:4222
      - CONFIG_FILE=/config/opc-${INSTANCE_ID}.toml
    volumes:
      - ./config/generated:/config:ro
    depends_on:
      - nats

  # S7 Ingestor
  ingestor-s7:
    build: ./services/ingestor-s7
    environment:
      - NATS_URL=nats://nats:4222
      - CONFIG_FILE=/config/s7-${INSTANCE_ID}.toml
    volumes:
      - ./config/generated:/config:ro
    depends_on:
      - nats

  # Config UI
  config-ui:
    build: ./web/config-ui
    ports:
      - "3001:3000"
    environment:
      - REACT_APP_API_URL=http://localhost:8090
    depends_on:
      - config-manager

  # Infrastructure
  postgres:
    image: postgres:15-alpine
    environment:
      - POSTGRES_DB=historian_config
      - POSTGRES_USER=historian
      - POSTGRES_PASSWORD=password
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./ops/db/init.sql:/docker-entrypoint-initdb.d/init.sql

  nats:
    image: nats:latest
    ports:
      - "4222:4222"
    command: "-js -sd /data"
    volumes:
      - nats_data:/data

  engine:
    build: ./services/engine
    # ... existing engine config

volumes:
  postgres_data:
  nats_data:
```

---

## 🔄 Config Hot Reload Mechanism

```rust
// services/ingestor-modbus/src/main.rs

use notify::{Watcher, RecursiveMode, watcher};
use std::sync::mpsc::channel;
use std::time::Duration;

#[tokio::main]
async fn main() -> Result<()> {
    let config_path = env::var("CONFIG_FILE")?;
    
    // Initial load
    let mut settings = Settings::load(&config_path)?;
    
    // Start adapters
    let (reload_tx, mut reload_rx) = mpsc::channel(1);
    spawn_adapters(&settings, reload_tx.clone()).await;
    
    // Watch config file for changes
    let (tx, rx) = channel();
    let mut watcher = watcher(tx, Duration::from_secs(2))?;
    watcher.watch(&config_path, RecursiveMode::NonRecursive)?;
    
    tokio::spawn(async move {
        for event in rx {
            match event {
                Ok(notify::DebouncedEvent::Write(_)) => {
                    info!("Config file changed, reloading...");
                    if let Ok(new_settings) = Settings::load(&config_path) {
                        reload_tx.send(new_settings).await.ok();
                    }
                }
                _ => {}
            }
        }
    });
    
    // Handle reloads
    while let Some(new_settings) = reload_rx.recv().await {
        info!("Applying new configuration");
        settings = new_settings;
        // Gracefully restart adapters
        spawn_adapters(&settings, reload_tx.clone()).await;
    }
    
    Ok(())
}
```

---

## 📊 Implementation Roadmap

### Phase 1: Foundation (Week 1)
- [ ] Database schema setup
- [ ] Config Manager service skeleton
- [ ] Modbus ingestor refactoring
- [ ] Docker Compose setup

### Phase 2: Config Management (Week 2)
- [ ] REST API implementation
- [ ] Config generation logic
- [ ] File watcher / hot reload
- [ ] Validation layer

### Phase 3: Web UI (Week 3)
- [ ] React app setup
- [ ] Device list page
- [ ] Modbus config form
- [ ] API integration

### Phase 4: Additional Protocols (Week 4)
- [ ] OPC UA ingestor
- [ ] S7 ingestor
- [ ] Protocol-specific UI forms

### Phase 5: Production Ready (Week 5)
- [ ] Kubernetes manifests
- [ ] Monitoring & metrics
- [ ] Backup & restore
- [ ] Documentation

---

## 🎯 Next Steps

1. **Database Setup**
2. **Config Manager Service**
3. **Modbus Ingestor Refactor**
4. **Web UI Prototype**

Hangi adımdan başlayalım?
