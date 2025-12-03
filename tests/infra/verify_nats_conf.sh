#!/bin/bash
set -e

CONFIG_FILE="ops/nats.conf"

echo "Verifying NATS configuration..."

if [ ! -f "$CONFIG_FILE" ]; then
    echo "❌ Error: $CONFIG_FILE does not exist."
    exit 1
fi

if ! grep -q "jetstream" "$CONFIG_FILE"; then
    echo "❌ Error: JetStream not enabled in $CONFIG_FILE."
    exit 1
fi

if ! grep -q 'store_dir: "/data/jetstream"' "$CONFIG_FILE"; then
    echo "❌ Error: store_dir not configured correctly in $CONFIG_FILE."
    exit 1
fi

echo "✅ NATS configuration verification passed."
exit 0
