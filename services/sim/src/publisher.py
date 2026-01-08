#!/usr/bin/env python3
"""
Test Data Publisher - Generates simulated sensor data for testing
Publishes to data.{sensor_id} topics that Engine subscribes to
"""

import asyncio
import os
import math
import random
import time
import logging
from nats.aio.client import Client as NATS
import common_pb2

logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(name)s - %(levelname)s - %(message)s')
logger = logging.getLogger("data-publisher")

NATS_URL = os.getenv("NATS_URL", "nats://localhost:4222")
PUBLISH_INTERVAL = float(os.getenv("PUBLISH_INTERVAL", "1.0"))  # seconds

# Simulated sensors following factory/line/device/type/tag naming convention
SENSORS = [
    {"id": "Factory1.Line1.PLC001.Temperature.T001", "base": 25, "amplitude": 5, "noise": 0.5},
    {"id": "Factory1.Line1.PLC001.Temperature.T002", "base": 30, "amplitude": 8, "noise": 0.8},
    {"id": "Factory1.Line1.PLC001.Pressure.P001", "base": 101.3, "amplitude": 2, "noise": 0.1},
    {"id": "Factory1.Line1.PLC002.Flow.F001", "base": 50, "amplitude": 10, "noise": 1.0},
    {"id": "Factory1.Line2.PLC001.Temperature.T001", "base": 80, "amplitude": 15, "noise": 2.0},
    {"id": "Factory1.Line2.PLC001.Level.L001", "base": 75, "amplitude": 20, "noise": 3.0},
]

def generate_value(sensor: dict, t: float) -> float:
    """Generate a simulated sensor value with sinusoidal pattern and noise"""
    base = sensor["base"]
    amplitude = sensor["amplitude"]
    noise = sensor["noise"]
    
    # Create a sinusoidal pattern with some variation
    period = 60 + random.uniform(-10, 10)  # ~60 second period
    value = base + amplitude * math.sin(2 * math.pi * t / period)
    value += random.gauss(0, noise)
    
    return round(value, 3)

async def main():
    nc = NATS()
    
    try:
        await nc.connect(servers=[NATS_URL])
        logger.info(f"Connected to NATS at {NATS_URL}")
        
        start_time = time.time()
        
        while True:
            t = time.time() - start_time
            timestamp_ms = int(time.time() * 1000)
            
            for sensor in SENSORS:
                value = generate_value(sensor, t)
                
                # Create protobuf message
                msg = common_pb2.SensorData()
                msg.sensor_id = sensor["id"]
                msg.value = value
                msg.timestamp_ms = timestamp_ms
                msg.quality = 1
                
                # Publish to data.{sensor_id} topic
                subject = f"data.{sensor['id']}"
                await nc.publish(subject, msg.SerializeToString())
                
            logger.info(f"Published {len(SENSORS)} sensor readings at t={t:.1f}s")
            await asyncio.sleep(PUBLISH_INTERVAL)
            
    except KeyboardInterrupt:
        logger.info("Shutting down...")
    except Exception as e:
        logger.error(f"Error: {e}")
    finally:
        await nc.drain()
        logger.info("NATS connection closed")

if __name__ == "__main__":
    asyncio.run(main())
