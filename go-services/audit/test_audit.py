#!/usr/bin/env python3
"""
Simple test script for Audit Service
Sends test events via NATS and verifies chain integrity
"""

import json
import subprocess
import time
import urllib.request

def send_nats_event(subject, data):
    """Send event to NATS using nats-box container"""
    cmd = [
        "docker", "run", "--rm",
        "--network", "ops_historian-net",
        "natsio/nats-box",
        "nats", "pub", subject,
        "--server", "nats://nats:4222",
        json.dumps(data)
    ]
    result = subprocess.run(cmd, capture_output=True, text=True)
    if result.returncode == 0:
        print(f"✅ Sent {subject}: {data}")
    else:
        print(f"❌ Failed to send {subject}: {result.stderr}")

def verify_chain():
    """Verify audit chain integrity"""
    try:
        with urllib.request.urlopen("http://localhost:8082/api/v1/audit/verify") as response:
            result = json.loads(response.read())
            if result.get("valid"):
                print("✅ Chain integrity verified!")
                return True
            else:
                print(f"❌ Chain broken at: {result.get('broken_id')}")
                return False
    except Exception as e:
        print(f"❌ Verification failed: {e}")
        return False

def query_database():
    """Query audit logs from database"""
    cmd = [
        "docker", "exec", "ops-postgres-1",
        "psql", "-U", "postgres", "-d", "historian",
        "-c", "SELECT id, timestamp, actor, action, LEFT(prev_hash, 8) as prev, LEFT(curr_hash, 8) as curr FROM audit_logs ORDER BY timestamp;"
    ]
    result = subprocess.run(cmd, capture_output=True, text=True)
    print("\n📊 Audit Logs:")
    print(result.stdout)

def main():
    print("🧪 Audit Service Test")
    print("=" * 50)
    
    # 1. Initial verification
    print("\n1️⃣ Initial chain verification (should be empty):")
    verify_chain()
    
    # 2. Send test events
    print("\n2️⃣ Sending test events...")
    
    send_nats_event("sys.auth.login", {
        "actor": "admin",
        "action": "login",
        "ip": "127.0.0.1"
    })
    
    send_nats_event("sys.audit.setpoint", {
        "actor": "admin",
        "action": "changed_setpoint",
        "device": "PLC-001",
        "old_value": 50,
        "new_value": 75
    })
    
    send_nats_event("sys.audit.alarm", {
        "actor": "operator",
        "action": "acknowledged_alarm",
        "alarm_id": "ALM-123"
    })
    
    # Wait for processing
    print("\n⏳ Waiting for events to be processed...")
    time.sleep(3)
    
    # 3. Verify chain after events
    print("\n3️⃣ Verifying chain integrity after events:")
    verify_chain()
    
    # 4. Query database
    print("\n4️⃣ Querying database:")
    query_database()
    
    print("\n✅ Test completed!")
    print("\n💡 Useful commands:")
    print("  - View logs: docker-compose logs -f audit")
    print("  - Stop services: cd ops && docker-compose down")

if __name__ == "__main__":
    main()
