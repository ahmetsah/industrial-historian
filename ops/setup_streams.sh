#!/bin/sh
set -e

echo "Waiting for NATS server..."
sleep 5

echo "Resetting EVENTS stream..."

# 1. Varsa eski stream'i SİL (Hata verirse yoksay '|| true')
# -f: Force (Soru sorma)
nats stream delete EVENTS -f --server nats:4222 || true

echo "Creating EVENTS stream config..."

# 2. JSON Konfigürasyonu oluştur
cat <<EOF > /tmp/events_stream.json
{
  "name": "EVENTS",
  "subjects": ["data.>"],
  "retention": "limits",
  "max_consumers": -1,
  "max_msgs": -1,
  "max_bytes": -1,
  "max_age": 31536000000000000,
  "storage": "file",
  "discard": "old",
  "num_replicas": 1,
  "duplicate_window": 120000000000,
  "no_ack": false
}
EOF

echo "Applying configuration..."

# 3. Yeni ayarlarla stream'i EKLE
nats stream add --config /tmp/events_stream.json --server nats:4222

echo "Stream EVENTS created successfully."