#!/usr/bin/env python3
"""
Test script for Alarm Service
Tests Alarm Definition, Triggering, Acknowledge, Shelving, and Clearing.
Uses a helper Go program to publish Protobuf messages to NATS.
"""

import json
import subprocess
import time
import urllib.request
import urllib.error
import sys
import os

ALARM_API_URL = "http://localhost:8083/api/v1/alarms"
PUBLISHER_CMD = ["go", "run", "./cmd/test-publisher/main.go"]

def send_sensor_data(sensor_id, value):
    """Send sensor data to NATS using Go utility"""
    cmd = PUBLISHER_CMD + [
        "-sensor", sensor_id,
        "-value", str(value)
    ]
    result = subprocess.run(cmd, capture_output=True, text=True)
    if result.returncode == 0:
        print(f"✅ Sent {sensor_id}: {value}")
    else:
        print(f"❌ Failed to send data: {result.stderr}")
        sys.exit(1)

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

def main():
    print("🧪 Alarm Service Test")
    print("=" * 50)
    
    # Check if we are in the right directory
    if not os.path.exists("cmd/test-publisher/main.go"):
        print("❌ Please run this script from 'go-services/alarm' directory")
        sys.exit(1)

    tag = f"test_sensor_{int(time.time())}"
    threshold = 80.0
    
    # 1. Create Definition
    print("\n1️⃣ Creating Alarm Definition...")
    def_resp = create_definition(tag, threshold, "High", "Critical")
    if not def_resp:
        sys.exit(1)
        
    # 2. Trigger Alarm
    print("\n2️⃣ Triggering Alarm (Value: 85 > 80)...")
    send_sensor_data(tag, 85.0)
    time.sleep(1) # Wait for processing
    
    alarms = get_active_alarms()
    active_alarm = None
    for a in alarms:
        if a['definition_id'] == def_resp['id']:
            active_alarm = a
            break
            
    if active_alarm:
        print(f"✅ Alarm Triggered: ID={active_alarm['id']}, State={active_alarm['state']}")
        if active_alarm['state'] != "UnackActive":
            print(f"❌ Expected UnackActive, got {active_alarm['state']}")
    else:
        print("❌ Alarm NOT triggered")
        sys.exit(1)
        
    # 3. Acknowledge Alarm
    print("\n3️⃣ Acknowledging Alarm...")
    if ack_alarm(active_alarm['id']):
        time.sleep(1)
        alarms = get_active_alarms()
        for a in alarms:
            if a['id'] == active_alarm['id']:
                print(f"✅ Alarm State: {a['state']}")
                if a['state'] != "AckActive":
                    print(f"❌ Expected AckActive, got {a['state']}")
                break
    
    # 4. Clear Alarm
    print("\n4️⃣ Clearing Alarm (Value: 70 < 80)...")
    send_sensor_data(tag, 70.0)
    time.sleep(1)
    
    alarms = get_active_alarms()
    found = False
    for a in alarms:
        if a['definition_id'] == def_resp['id']:
            found = True
            break
    
    if not found:
        print("✅ Alarm Cleared (Removed from active list)")
    else:
        print("❌ Alarm still active (Should be cleared)")
        
    # 5. Test Shelving
    print("\n5️⃣ Testing Shelving...")
    # Trigger again
    send_sensor_data(tag, 90.0)
    time.sleep(1)
    alarms = get_active_alarms()
    new_alarm = None
    for a in alarms:
        if a['definition_id'] == def_resp['id']:
            new_alarm = a
            break
            
    if new_alarm:
        print(f"✅ Alarm Re-triggered: ID={new_alarm['id']}")
        shelve_alarm(new_alarm['id'], 10)
        time.sleep(1)
        
        alarms = get_active_alarms()
        for a in alarms:
            if a['id'] == new_alarm['id']:
                print(f"✅ Alarm State: {a['state']}")
                if a['state'] != "Shelved":
                    print(f"❌ Expected Shelved, got {a['state']}")
                break
    else:
        print("❌ Failed to re-trigger alarm")

    # 6. Integration: Verify Audit Log
    print("\n6️⃣ Verifying Audit Log (Integration)...")
    
    # Wait a bit for Audit Service to process the event
    time.sleep(2)
    
    # Query Postgres directly via docker exec
    # We look for action 'alarm_AckActive' which corresponds to the Acknowledge step
    cmd = [
        "docker", "exec", "ops-postgres",
        "psql", "-U", "postgres", "-d", "historian",
        "-t", "-c", "SELECT count(*) FROM audit_logs WHERE action = 'alarm_AckActive';"
    ]
    
    try:
        result = subprocess.run(cmd, capture_output=True, text=True)
        if result.returncode == 0:
            count = int(result.stdout.strip())
            if count > 0:
                print(f"✅ Audit Log Verified: Found {count} 'alarm_AckActive' entries.")
            else:
                print("❌ Audit Log Verification Failed: No 'alarm_AckActive' entries found.")
                # Debug: Show all logs
                debug_cmd = [
                    "docker", "exec", "ops-postgres",
                    "psql", "-U", "postgres", "-d", "historian",
                    "-c", "SELECT id, action, timestamp FROM audit_logs ORDER BY timestamp DESC LIMIT 5;"
                ]
                debug_res = subprocess.run(debug_cmd, capture_output=True, text=True)
                print("Recent Logs:")
                print(debug_res.stdout)
        else:
             print(f"❌ Failed to query database: {result.stderr}")
             
    except Exception as e:
        print(f"⚠️ Could not verify Audit Log: {e}")

    print("\n✅ Test Completed!")

if __name__ == "__main__":
    main()
