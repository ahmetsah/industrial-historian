#!/bin/bash
set -e

echo "🚀 Ingestor Quick Start - Real Modbus Device"
echo "=============================================="
echo "Device: 172.29.80.1:5020"
echo "NATS: 172.29.80.1:4222"
echo ""

# 1. Check connections
echo "1️⃣ Checking connections..."
nc -zv 172.29.80.1 5020 2>&1 | grep -q "succeeded" && echo "✅ Modbus device reachable" || echo "❌ Modbus device not reachable"
nc -zv 172.29.80.1 4222 2>&1 | grep -q "succeeded" && echo "✅ NATS reachable" || echo "❌ NATS not reachable"

echo ""
echo "2️⃣ Starting Ingestor..."
echo "   Press Ctrl+C to stop"
echo ""

cd /home/ahmet/historian

# Start ingestor
cargo run -p ingestor
