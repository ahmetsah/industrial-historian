#!/usr/bin/env python3
"""
Simple test script for Alarm Service
Tests Alarm Definition, Triggering, Acknowledge, Shelving, and Clearing
"""

import json
import subprocess
import time
import urllib.request
import urllib.error

ALARM_API_URL = "http://localhost:8083/api/v1/alarms"

def send_nats_sensor_data(sensor_id, value):
    """Send sensor data to NATS using nats-box container"""
    subject = f"enterprise.site1.area1.line1.device1.{sensor_id}"
    data = {
        "sensor_id": sensor_id,
        "value": value,
        "timestamp_ms": int(time.time() * 1000),
        "quality": 1
    }
    
    # Using protobuf encoding would be ideal, but for simplicity we might need to check 
    # if the ingestor or alarm service expects JSON or Proto.
    # The Alarm Service NatsTransport unmarshals Proto: `proto.Unmarshal(msg.Data, &sensorData)`
    # So we cannot just send JSON string via NATS if the service expects Proto bytes.
    # However, the `test_audit.py` sent JSON. Let's check `audit` service.
    # Audit service might be using JSON.
    # Alarm service definitely uses Proto.
    
    # If we use `nats pub`, we send raw bytes.
    # If the service expects Proto, we need to send Proto-encoded bytes.
    # Generating Proto bytes in Python without compiled protos is hard.
    
    # Alternative: Use the Ingestor? 
    # The Ingestor takes Modbus/OPC-UA and publishes Proto to NATS.
    # But we want to test Alarm Service in isolation or integration.
    
    # If we want to test Alarm Service integration, we should publish what it expects.
    # It expects `historian.v1.SensorData` proto.
    
    # If I cannot easily generate proto bytes here, maybe I can use a small Go utility?
    # Or maybe I can assume the user has `protoc` or similar?
    
    # Wait, `nats-box` might not support proto encoding on the fly.
    
    # Let's look at `test_ingestor.py` to see how it tests?
    # Or maybe `test_engine.py`?
    pass

def create_definition(tag, threshold, type_, priority):
    """Create an alarm definition via HTTP API"""
    data = {
        "tag": tag,
        "threshold": threshold,
        "type": type_,
        "priority": priority
    }
    req = urllib.request.Request(
        f"{ALARM_API_URL}/definitions",
        data=json.dumps(data).encode('utf-8'),
        headers={'Content-Type': 'application/json'},
        method='POST'
    )
    try:
        with urllib.request.urlopen(req) as response:
            if response.status == 201:
                print(f"✅ Created definition for {tag}")
                return json.loads(response.read())
    except urllib.error.HTTPError as e:
        print(f"❌ Failed to create definition: {e}")
        return None

def get_active_alarms():
    """Get active alarms via HTTP API"""
    try:
        with urllib.request.urlopen(f"{ALARM_API_URL}/active") as response:
            return json.loads(response.read())
    except Exception as e:
        print(f"❌ Failed to get active alarms: {e}")
        return []

def ack_alarm(alarm_id):
    """Acknowledge alarm via HTTP API"""
    req = urllib.request.Request(
        f"{ALARM_API_URL}/{alarm_id}/ack",
        method='POST'
    )
    try:
        with urllib.request.urlopen(req) as response:
            if response.status == 200:
                print(f"✅ Acknowledged alarm {alarm_id}")
                return True
    except urllib.error.HTTPError as e:
        print(f"❌ Failed to ack alarm: {e}")
        return False

def shelve_alarm(alarm_id, duration_sec):
    """Shelve alarm via HTTP API"""
    data = {"duration_seconds": duration_sec}
    req = urllib.request.Request(
        f"{ALARM_API_URL}/{alarm_id}/shelve",
        data=json.dumps(data).encode('utf-8'),
        headers={'Content-Type': 'application/json'},
        method='POST'
    )
    try:
        with urllib.request.urlopen(req) as response:
            if response.status == 200:
                print(f"✅ Shelved alarm {alarm_id}")
                return True
    except urllib.error.HTTPError as e:
        print(f"❌ Failed to shelve alarm: {e}")
        return False

# ... (rest of the script)
