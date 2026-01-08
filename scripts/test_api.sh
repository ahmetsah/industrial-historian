#!/bin/bash

# Config Manager API Test Script
# Usage: ./test_api.sh

API_URL="http://localhost:8090/api/v1"

echo "🧪 Config Manager API Tests"
echo "============================"
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test 1: Health Check
echo -e "${YELLOW}Test 1: Health Check${NC}"
response=$(curl -s -w "\n%{http_code}" ${API_URL%/api/v1}/health)
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

if [ "$http_code" = "200" ]; then
    echo -e "${GREEN}✓ PASS${NC}"
    echo "$body" | jq '.'
else
    echo -e "${RED}✗ FAIL (HTTP $http_code)${NC}"
fi
echo ""

# Test 2: List Devices (Empty)
echo -e "${YELLOW}Test 2: List Devices${NC}"
response=$(curl -s -w "\n%{http_code}" $API_URL/devices)
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

if [ "$http_code" = "200" ]; then
    echo -e "${GREEN}✓ PASS${NC}"
    echo "$body" | jq '.'
else
    echo -e "${RED}✗ FAIL (HTTP $http_code)${NC}"
fi
echo ""

# Test 3: Create Modbus Device
echo -e "${YELLOW}Test 3: Create Modbus Device${NC}"
response=$(curl -s -w "\n%{http_code}" -X POST $API_URL/devices/modbus \
  -H "Content-Type: application/json" \
  -d '{
    "name": "PLC-001",
    "description": "Main production line PLC",
    "ip": "192.168.1.10",
    "port": 502,
    "unit_id": 1,
    "poll_interval_ms": 1000,
    "registers": [
      {
        "address": 0,
        "name": "Factory1.Line1.PLC001.Temp.T001",
        "data_type": "Float32",
        "scale_factor": 1.0,
        "offset": 0.0,
        "unit": "°C",
        "description": "Reactor temperature"
      },
      {
        "address": 2,
        "name": "Factory1.Line1.PLC001.Pressure.P001",
        "data_type": "Int16",
        "scale_factor": 0.1,
        "offset": 0.0,
        "unit": "bar",
        "description": "Reactor pressure"
      }
    ]
  }')

http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

if [ "$http_code" = "201" ]; then
    echo -e "${GREEN}✓ PASS${NC}"
    echo "$body" | jq '.'
    
    # Extract device ID for next tests
    DEVICE_ID=$(echo "$body" | jq -r '.device.id')
    echo ""
    echo "Device ID: $DEVICE_ID"
else
    echo -e "${RED}✗ FAIL (HTTP $http_code)${NC}"
    echo "$body"
fi
echo ""

# Test 4: List Modbus Devices
echo -e "${YELLOW}Test 4: List Modbus Devices${NC}"
response=$(curl -s -w "\n%{http_code}" $API_URL/devices/modbus)
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

if [ "$http_code" = "200" ]; then
    echo -e "${GREEN}✓ PASS${NC}"
    echo "$body" | jq '.'
else
    echo -e "${RED}✗ FAIL (HTTP $http_code)${NC}"
fi
echo ""

# Test 5: Get Specific Modbus Device
if [ ! -z "$DEVICE_ID" ]; then
    echo -e "${YELLOW}Test 5: Get Modbus Device${NC}"
    response=$(curl -s -w "\n%{http_code}" $API_URL/devices/modbus/$DEVICE_ID)
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')

    if [ "$http_code" = "200" ]; then
        echo -e "${GREEN}✓ PASS${NC}"
        echo "$body" | jq '.'
    else
        echo -e "${RED}✗ FAIL (HTTP $http_code)${NC}"
    fi
    echo ""
fi

# Test 6: Get Latest Config
if [ ! -z "$DEVICE_ID" ]; then
    echo -e "${YELLOW}Test 6: Get Latest Config${NC}"
    response=$(curl -s -w "\n%{http_code}" $API_URL/config/latest/$DEVICE_ID)
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')

    if [ "$http_code" = "200" ]; then
        echo -e "${GREEN}✓ PASS${NC}"
        echo "$body" | jq '.'
        
        # Show generated config file
        CONFIG_FILE=$(echo "$body" | jq -r '.file_path')
        if [ -f "$CONFIG_FILE" ]; then
            echo ""
            echo "Generated Config File:"
            echo "======================"
            cat "$CONFIG_FILE"
        fi
    else
        echo -e "${RED}✗ FAIL (HTTP $http_code)${NC}"
    fi
    echo ""
fi

# Test 7: Create Second Device
echo -e "${YELLOW}Test 7: Create Second Modbus Device${NC}"
response=$(curl -s -w "\n%{http_code}" -X POST $API_URL/devices/modbus \
  -H "Content-Type: application/json" \
  -d '{
    "name": "PLC-002",
    "description": "Secondary line PLC",
    "ip": "192.168.1.20",
    "port": 502,
    "unit_id": 1,
    "poll_interval_ms": 2000,
    "registers": [
      {
        "address": 0,
        "name": "Factory1.Line2.PLC002.Speed.S001",
        "data_type": "UInt16",
        "unit": "RPM"
      }
    ]
  }')

http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

if [ "$http_code" = "201" ]; then
    echo -e "${GREEN}✓ PASS${NC}"
    echo "$body" | jq '.device.name, .config.file_path'
else
    echo -e "${RED}✗ FAIL (HTTP $http_code)${NC}"
fi
echo ""

# Summary
echo "============================"
echo -e "${GREEN}✅ Tests Complete${NC}"
echo ""
echo "Check generated config files in: ./config/generated/"
echo ""
