#!/bin/bash
# Yeni Modbus Ingestor Ekleme Script'i
# Kullanım: ./add_modbus_ingestor.sh DEVICE_NAME

set -e

DEVICE_NAME=$1

if [ -z "$DEVICE_NAME" ]; then
    echo "❌ Hata: Cihaz adı belirtilmedi!"
    echo "Kullanım: ./add_modbus_ingestor.sh DEVICE_NAME"
    echo "Örnek: ./add_modbus_ingestor.sh PLC-002"
    exit 1
fi

# Küçük harf ve tire yerine alt çizgi
DEVICE_NAME_LOWER=$(echo $DEVICE_NAME | tr '[:upper:]' '[:lower:]' | tr '-' '')

echo "🚀 Yeni Modbus Ingestor Ekleniyor: $DEVICE_NAME"
echo "================================================"

# 1. Config dosyasının varlığını kontrol et
CONFIG_FILE="config/generated/modbus-${DEVICE_NAME}.toml"
if [ ! -f "$CONFIG_FILE" ]; then
    echo "❌ Hata: Config dosyası bulunamadı: $CONFIG_FILE"
    echo ""
    echo "Önce Config Manager API ile cihazı oluşturun:"
    echo "curl -X POST http://localhost:8090/api/v1/devices/modbus \\"
    echo "  -H 'Content-Type: application/json' \\"
    echo "  -d '{\"name\":\"$DEVICE_NAME\", ...}'"
    exit 1
fi

echo "✅ Config dosyası bulundu: $CONFIG_FILE"

# 2. Docker Compose'a servis ekle
echo ""
echo "📝 Docker Compose'a servis ekleniyor..."

COMPOSE_FILE="ops/docker-compose.yml"

# Servis zaten var mı kontrol et
if grep -q "ingestor-modbus-${DEVICE_NAME_LOWER}:" "$COMPOSE_FILE"; then
    echo "⚠️  Servis zaten mevcut: ingestor-modbus-${DEVICE_NAME_LOWER}"
else
    # Servis tanımını ekle
    cat >> "$COMPOSE_FILE" <<EOF

  # Modbus Ingestor: $DEVICE_NAME
  ingestor-modbus-${DEVICE_NAME_LOWER}:
    build:
      context: ..
      dockerfile: services/ingestor/Dockerfile
    container_name: ops-ingestor-modbus-${DEVICE_NAME_LOWER}
    environment:
      RUST_LOG: info
      CONFIG_FILE: /config/modbus-${DEVICE_NAME}
      NATS_URL: nats://nats:4222
    volumes:
      - ../config/generated:/config:ro
      - ingestor_buffer_${DEVICE_NAME_LOWER}:/data/buffer
    depends_on:
      - nats
      - config-manager
    networks:
      - historian-net
    restart: unless-stopped
    extra_hosts:
      - "host.docker.internal:host-gateway"
EOF
    echo "✅ Servis eklendi: ingestor-modbus-${DEVICE_NAME_LOWER}"
fi

# 3. Volume ekle
echo ""
echo "💾 Volume ekleniyor..."

if grep -q "ingestor_buffer_${DEVICE_NAME_LOWER}:" "$COMPOSE_FILE"; then
    echo "⚠️  Volume zaten mevcut: ingestor_buffer_${DEVICE_NAME_LOWER}"
else
    # volumes: satırını bul ve altına ekle
    sed -i "/^volumes:/a\\  ingestor_buffer_${DEVICE_NAME_LOWER}:  # $DEVICE_NAME buffer" "$COMPOSE_FILE"
    echo "✅ Volume eklendi: ingestor_buffer_${DEVICE_NAME_LOWER}"
fi

# 4. Docker Compose'u yeniden yükle ve servisi başlat
echo ""
echo "🐳 Docker servisi başlatılıyor..."
cd ops
docker-compose up -d ingestor-modbus-${DEVICE_NAME_LOWER}

# 5. Durum kontrolü
echo ""
echo "📊 Servis durumu:"
docker-compose ps ingestor-modbus-${DEVICE_NAME_LOWER}

echo ""
echo "📋 Logları görmek için:"
echo "docker logs -f ops-ingestor-modbus-${DEVICE_NAME_LOWER}"

echo ""
echo "🎉 Tamamlandı! Yeni ingestor çalışıyor: $DEVICE_NAME"
