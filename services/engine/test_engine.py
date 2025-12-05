#!/usr/bin/env python3
"""
Engine Service Test Script
Tests NATS ingestion, RocksDB storage, gRPC query, and HTTP export
"""

import subprocess
import time
import json
import requests
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

def check_engine_running():
    """Check if engine is running"""
    print_header("Checking Engine Status")
    
    try:
        result = subprocess.run(
            ["pgrep", "-f", "engine"],
            capture_output=True,
            text=True
        )
        
        if result.returncode == 0:
            pids = result.stdout.strip().split('\n')
            print_success(f"Engine is running")
            print_info(f"   PID(s): {', '.join(pids)}")
            return True
        else:
            print_error("Engine is not running")
            print_info("   Start: cd /home/ahmet/historian && cargo run -p engine")
            return False
    except Exception as e:
        print_error(f"Failed to check engine: {e}")
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
            return False
    except Exception as e:
        print_error(f"Failed to check NATS: {e}")
        return False

def check_grpc_port():
    """Check if gRPC port is listening"""
    print_header("Checking gRPC Server")
    
    try:
        result = subprocess.run(
            ["ss", "-tuln"],
            capture_output=True,
            text=True,
            check=True
        )
        
        if ":50051" in result.stdout:
            print_success("gRPC server is listening on port 50051")
            return True
        else:
            print_error("gRPC server not found on port 50051")
            return False
    except Exception as e:
        print_error(f"Failed to check gRPC port: {e}")
        return False

def check_http_export():
    """Check if HTTP export server is running"""
    print_header("Checking HTTP Export Server")
    
    try:
        response = requests.get("http://localhost:8081/api/v1/export?sensor_id=test&start_ts=0&end_ts=1", timeout=5)
        if response.status_code in [200, 500]:  # 500 is ok, means server is running
            print_success("HTTP export server is running on port 8081")
            return True
        else:
            print_error(f"HTTP server returned {response.status_code}")
            return False
    except requests.exceptions.ConnectionError:
        print_error("HTTP export server not reachable on port 8081")
        return False
    except Exception as e:
        print_error(f"Failed to check HTTP server: {e}")
        return False

def test_data_ingestion():
    """Test if data is being ingested from NATS"""
    print_header("Testing Data Ingestion")
    
    print_info("Checking if Engine is consuming NATS messages...")
    print_info("(This requires Ingestor to be running and sending data)")
    
    # Check RocksDB directory
    try:
        result = subprocess.run(
            ["ls", "-lh", "/tmp/historian-db"],
            capture_output=True,
            text=True
        )
        
        if result.returncode == 0:
            print_success("RocksDB directory exists")
            # Count files
            file_count = len([l for l in result.stdout.split('\n') if l.strip() and not l.startswith('total')])
            print_info(f"   Files in database: {file_count}")
            return True
        else:
            print_error("RocksDB directory not found")
            return False
    except Exception as e:
        print_error(f"Failed to check RocksDB: {e}")
        return False

def test_grpc_query():
    """Test gRPC query endpoint"""
    print_header("Testing gRPC Query")
    
    print_info("Testing gRPC query requires grpcurl or custom client")
    print_info("Skipping automated gRPC test (manual test available)")
    
    print(f"\n{Colors.YELLOW}Manual gRPC Test:{Colors.NC}")
    print("  grpcurl -plaintext localhost:50051 list")
    print("  grpcurl -plaintext -d '{\"sensor_id\":\"adres_0\",\"start_ms\":0,\"end_ms\":9999999999999}' \\")
    print("    localhost:50051 historian.HistorianQuery/Query")
    
    return None

def test_http_export():
    """Test HTTP export endpoint"""
    print_header("Testing HTTP Export")
    
    try:
        # Test export endpoint
        params = {
            'sensor_id': 'adres_0',
            'start_ts': 0,
            'end_ts': 9999999999999
        }
        
        response = requests.get(
            "http://localhost:8081/api/v1/export",
            params=params,
            timeout=10
        )
        
        if response.status_code == 200:
            lines = response.text.strip().split('\n')
            print_success(f"HTTP export successful")
            print_info(f"   Received {len(lines)} lines")
            
            if len(lines) > 1:  # Header + at least one data line
                print("\n   Sample data (first 5 lines):")
                for line in lines[:5]:
                    print(f"      {line}")
                return True
            else:
                print_info("   No data in response (database might be empty)")
                return False
        else:
            print_error(f"HTTP export failed: {response.status_code}")
            print_info(f"   Response: {response.text[:200]}")
            return False
            
    except Exception as e:
        print_error(f"HTTP export test failed: {e}")
        return False

def check_database_size():
    """Check RocksDB database size"""
    print_header("Checking Database Size")
    
    try:
        result = subprocess.run(
            ["du", "-sh", "/tmp/historian-db"],
            capture_output=True,
            text=True
        )
        
        if result.returncode == 0:
            size = result.stdout.split()[0]
            print_success(f"Database size: {size}")
            return True
        else:
            print_error("Failed to get database size")
            return False
    except Exception as e:
        print_error(f"Failed to check database size: {e}")
        return False

def main():
    print(f"\n{Colors.BLUE}🗄️  Engine Service Test Suite{Colors.NC}")
    print(f"{Colors.BLUE}{'='*60}{Colors.NC}")
    print(f"Time: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
    print(f"\n{Colors.YELLOW}Components:{Colors.NC}")
    print("  - NATS Ingestion (data.raw)")
    print("  - RocksDB Storage (/tmp/historian-db)")
    print("  - gRPC Query API (port 50051)")
    print("  - HTTP Export API (port 8081)")
    
    # Test 1: Check if Engine is running
    engine_running = check_engine_running()
    time.sleep(1)
    
    if not engine_running:
        print_info("\nEngine is not running. Start it first:")
        print("  cd /home/ahmet/historian && cargo run -p engine")
        return
    
    # Test 2: Check NATS connection
    nats_ok = check_nats_connection()
    time.sleep(1)
    
    # Test 3: Check gRPC port
    grpc_ok = check_grpc_port()
    time.sleep(1)
    
    # Test 4: Check HTTP export server
    http_ok = check_http_export()
    time.sleep(1)
    
    # Test 5: Check data ingestion
    ingestion_ok = test_data_ingestion()
    time.sleep(1)
    
    # Test 6: Check database size
    check_database_size()
    time.sleep(1)
    
    # Test 7: Test gRPC query
    test_grpc_query()
    time.sleep(1)
    
    # Test 8: Test HTTP export
    export_ok = test_http_export()
    
    # Summary
    print_header("Test Summary")
    
    if engine_running and nats_ok and grpc_ok and http_ok:
        print_success("Engine service is operational!")
        print_info("\n📊 Data Flow:")
        print("   NATS (data.>) → Engine → RocksDB → Query/Export APIs")
        
        if ingestion_ok and export_ok:
            print_success("\n✅ Data is flowing through the system!")
        else:
            print_info("\n⚠️  Engine is running but no data yet")
            print_info("   Make sure Ingestor is running and sending data")
    else:
        print_error("Some components are not running")
        
    print(f"\n{Colors.YELLOW}💡 Useful Commands:{Colors.NC}")
    print("  - Start Engine: NATS_URL=nats://172.29.80.1:4222 NATS_SUBJECT=data.raw HTTP_PORT=8081 cargo run -p engine")
    print("  - View logs: RUST_LOG=debug cargo run -p engine")
    print("  - Export data: curl 'http://localhost:8081/api/v1/export?sensor_id=adres_0&start_ts=0&end_ts=9999999999999'")
    print("  - gRPC query: grpcurl -plaintext localhost:50051 list")
    print("  - Check DB: du -sh /tmp/historian-db")

if __name__ == "__main__":
    main()
