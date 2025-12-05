#!/bin/bash

echo "📡 NATS Data Monitor"
echo "===================="
echo "Subscribing to: data.raw"
echo "Server: nats://172.29.80.1:4222"
echo ""
echo "Press Ctrl+C to stop"
echo ""

# Check if nats CLI is available
if command -v nats &> /dev/null; then
    nats sub data.raw --server nats://172.29.80.1:4222
else
    echo "❌ 'nats' CLI not found"
    echo ""
    echo "Install with:"
    echo "  go install github.com/nats-io/natscli/nats@latest"
    echo ""
    echo "Or use Docker:"
    docker run --rm --network host natsio/nats-box \
        nats sub data.raw --server nats://172.29.80.1:4222
fi
