#!/bin/bash
set -e

COMPOSE_FILE="ops/docker-compose.yml"

echo "Verifying Docker Compose configuration..."

if [ ! -f "$COMPOSE_FILE" ]; then
    echo "❌ Error: $COMPOSE_FILE does not exist."
    exit 1
fi

if ! grep -q "nats:" "$COMPOSE_FILE"; then
    echo "❌ Error: nats service not defined in $COMPOSE_FILE."
    exit 1
fi

if ! grep -q "minio:" "$COMPOSE_FILE"; then
    echo "❌ Error: minio service not defined in $COMPOSE_FILE."
    exit 1
fi

if ! grep -q "postgres:" "$COMPOSE_FILE"; then
    echo "❌ Error: postgres service not defined in $COMPOSE_FILE."
    exit 1
fi

if ! grep -q "volumes:" "$COMPOSE_FILE"; then
    echo "❌ Error: volumes not defined in $COMPOSE_FILE."
    exit 1
fi

echo "✅ Docker Compose verification passed."
exit 0
