#!/usr/bin/env python3
"""
Ingestor Service Test Script - Real Modbus Device
Tests connection to real Modbus slave at 172.29.80.1:5020
"""

import subprocess
import time
import json
from datetime import datetime

class Colors:
    GREEN = '\033[0;32m'
    BLUE = '\033[0;34m'
    RED = '\033[0;31m'
    YELLOW = '\033[1;33m'
    NC = '\033[0m'

def print_header(text):
    print(f"\n{Colors.BLUE}{'='*60}{Colors.NC}")
    print(f"{Colors.BLUE}{text}{Colors.NC}")
    print(f"{Colors.BLUE}{'='*60}{Colors.NC}")

def print_success(text):
    print(f"{Colors.GREEN}✅ {text}{Colors.NC}")

def print_error(text):
    print(f"{Colors.RED}❌ {text}{Colors.NC}")

def print_info(text):
    print(f"{Colors.YELLOW}ℹ️  {text}{Colors.NC}")

def check_modbus_connection():
    """Check if Modbus device is reachable"""
    print_header("Checking Modbus Device Connection")
    
    try:
        result = subprocess.run(
            ["nc", "-zv", "172.29.80.1", "5020"],
            capture_output=True,
            text=True,
            timeout=5
        )
        
        if result.returncode == 0 or "succeeded" in result.stderr.lower():
            print_success("Modbus device is reachable at 172.29.80.1:5020")
            return True
        else:
            print_error("Cannot connect to Modbus device")
            print_info("   Make sure the device is powered on and network is configured")
            return False
    except subprocess.TimeoutExpired:
        print_error("Connection timeout")
        return False
    except Exception as e:
        print_error(f"Failed to check connection: {e}")
        return False

def check_nats_connection():
    """Check if NATS is reachable"""
    print_header("Checking NATS Connection")
    
    try:
        result = subprocess.run(
            ["nc", "-zv", "172.29.80.1", "4222"],
            capture_output=True,
            text=True,
            timeout=5
        )
        
        if result.returncode == 0 or "succeeded" in result.stderr.lower():
            print_success("NATS is reachable at 172.29.80.1:4222")
            return True
        else:
            print_error("Cannot connect to NATS")
            print_info("   Start NATS: cd ops && docker-compose up -d nats")
            return False
    except Exception as e:
        print_error(f"Failed to check NATS: {e}")
        return False

def check_ingestor_running():
    """Check if ingestor is running"""
    print_header("Checking Ingestor Status")
    
    try:
        # Check for both cargo and the binary
        result = subprocess.run(
            ["pgrep", "-f", "ingestor"],
            capture_output=True,
            text=True
        )
        
        if result.returncode == 0:
            pids = result.stdout.strip().split('\n')
            print_success(f"Ingestor is running")
            print_info(f"   PID(s): {', '.join(pids)}")
            return True
        else:
            print_error("Ingestor is not running")
            print_info("   Start: cd /home/ahmet/historian && cargo run -p ingestor")
            return False
    except Exception as e:
        print_error(f"Failed to check ingestor: {e}")
        return False

def subscribe_to_nats(duration=15):
    """Subscribe to NATS data.raw subject"""
    print_header(f"Subscribing to NATS 'data.raw' for {duration} seconds")
    print_info("Waiting for sensor data from Modbus...")
    
    try:
        # Use local NATS if available, otherwise use docker
        result = subprocess.run([
            "timeout", str(duration),
            "nats", "sub", "data.raw",
            "--server", "nats://172.29.80.1:4222"
        ], capture_output=True, text=True, timeout=duration+5)
        
        if result.stdout:
            lines = result.stdout.strip().split('\n')
            message_count = len([l for l in lines if 'data.raw' in l or '{' in l])
            
            if message_count > 0:
                print_success(f"Received {message_count} messages")
                print("\n📊 Sample messages:")
                
                # Show first few messages
                shown = 0
                for line in lines:
                    if '{' in line and shown < 3:
                        try:
                            # Try to parse and pretty print JSON
                            data = json.loads(line)
                            print(f"\n   Message {shown + 1}:")
                            for key, value in data.items():
                                print(f"      {key}: {value}")
                            shown += 1
                        except:
                            print(f"   {line[:100]}")
                            shown += 1
                
                return True
            else:
                print_error("No messages received")
                print_info("   Check if Ingestor is reading from Modbus")
                return False
        else:
            print_error("No output from NATS subscription")
            return False
            
    except subprocess.TimeoutExpired:
        print_info("Subscription timeout (no messages received)")
        return False
    except FileNotFoundError:
        print_error("'nats' CLI not found")
        print_info("   Install: go install github.com/nats-io/natscli/nats@latest")
        print_info("   Or use Docker: docker run --rm --network host natsio/nats-box nats sub data.raw --server nats://172.29.80.1:4222")
        return False
    except Exception as e:
        print_error(f"Failed to subscribe: {e}")
        return False

def test_modbus_read():
    """Test direct Modbus read (requires pymodbus)"""
    print_header("Testing Direct Modbus Read")
    
    try:
        import sys
        result = subprocess.run([
            sys.executable, "-c",
            """
from pymodbus.client import ModbusTcpClient
import json

client = ModbusTcpClient('172.29.80.1', port=5020)
if client.connect():
    result = client.read_holding_registers(0, 10, slave=1)
    if not result.isError():
        print(json.dumps({
            'success': True,
            'registers': result.registers,
            'addresses': list(range(10))
        }))
    else:
        print(json.dumps({'success': False, 'error': 'Read error'}))
    client.close()
else:
    print(json.dumps({'success': False, 'error': 'Connection failed'}))
"""
        ], capture_output=True, text=True, timeout=10)
        
        if result.returncode == 0 and result.stdout:
            data = json.loads(result.stdout)
            if data.get('success'):
                print_success("Successfully read from Modbus device")
                print("\n   📊 Register values:")
                for i, val in enumerate(data['registers']):
                    print(f"      Address {i}: {val}")
                return True
            else:
                print_error(f"Modbus read failed: {data.get('error')}")
                return False
        else:
            print_error("Failed to read from Modbus")
            if result.stderr:
                print_info(f"   Error: {result.stderr[:200]}")
            return False
            
    except ImportError:
        print_info("pymodbus not installed (optional)")
        print_info("   Install: pip install pymodbus")
        return None
    except Exception as e:
        print_error(f"Modbus read test failed: {e}")
        return None

def check_config_file():
    """Check if config file exists and is correct"""
    print_header("Checking Configuration")
    
    config_path = "/home/ahmet/historian/config/default.toml"
    
    try:
        with open(config_path, 'r') as f:
            content = f.read()
            
        checks = {
            '172.29.80.1': 'Modbus IP address',
            '5020': 'Modbus port',
            'data.raw': 'NATS subject',
            'adres_0': 'Register definitions'
        }
        
        all_ok = True
        for check, desc in checks.items():
            if check in content:
                print_success(f"{desc} configured")
            else:
                print_error(f"{desc} missing")
                all_ok = False
        
        return all_ok
        
    except FileNotFoundError:
        print_error(f"Config file not found: {config_path}")
        return False
    except Exception as e:
        print_error(f"Failed to check config: {e}")
        return False

def main():
    print(f"\n{Colors.BLUE}🔧 Ingestor Service Test - Real Modbus Device{Colors.NC}")
    print(f"{Colors.BLUE}{'='*60}{Colors.NC}")
    print(f"Time: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
    print(f"\n{Colors.YELLOW}Device: 172.29.80.1:5020 (Unit ID: 1){Colors.NC}")
    print(f"{Colors.YELLOW}NATS: 172.29.80.1:4222{Colors.NC}")
    
    # Test 1: Config file
    config_ok = check_config_file()
    time.sleep(1)
    
    # Test 2: Modbus connection
    modbus_ok = check_modbus_connection()
    time.sleep(1)
    
    # Test 3: NATS connection
    nats_ok = check_nats_connection()
    time.sleep(1)
    
    # Test 4: Direct Modbus read (optional)
    modbus_read_ok = test_modbus_read()
    time.sleep(1)
    
    # Test 5: Ingestor status
    ingestor_ok = check_ingestor_running()
    time.sleep(1)
    
    # Test 6: NATS subscription (if ingestor is running)
    if ingestor_ok and nats_ok:
        print_info("\nWaiting 5 seconds for data to flow...")
        time.sleep(5)
        subscribe_to_nats(15)
    else:
        print_info("\nSkipping NATS subscription (Ingestor not running)")
    
    # Summary
    print_header("Test Summary")
    
    if config_ok and modbus_ok and nats_ok and ingestor_ok:
        print_success("All systems operational!")
        print_info("\n📊 Data Flow:")
        print("   Modbus (172.29.80.1:5020) → Ingestor → NATS (data.raw)")
        print("\n   Registers being read:")
        print("   - adres_0 to adres_9 (addresses 0-9)")
        print("\n   Calculated tags:")
        print("   - Delta_P = Pressure - 10")
        print("   - Efficiency = (Temperature / Pressure) * 100")
    else:
        print_error("Some components need attention")
        
    print(f"\n{Colors.YELLOW}💡 Useful Commands:{Colors.NC}")
    print("  - Start Ingestor: cd /home/ahmet/historian && cargo run -p ingestor")
    print("  - View NATS messages: nats sub data.raw --server nats://172.29.80.1:4222")
    print("  - Check config: cat /home/ahmet/historian/config/default.toml")
    print("  - Debug mode: RUST_LOG=debug cargo run -p ingestor")

if __name__ == "__main__":
    main()
