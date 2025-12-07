import asyncio
import os
import signal
import sys
import logging
from nats.aio.client import Client as NATS
from nats.js.api import StreamConfig
import common_pb2

# Configure logging
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(name)s - %(levelname)s - %(message)s')
logger = logging.getLogger("digital-twin")

# Configuration
NATS_URL = os.getenv("NATS_URL", "nats://localhost:4222")
SENSOR_SUBJECT = "enterprise.site.area.line.reactor.temp"
PREDICTION_SUBJECT = "enterprise.site.area.line.reactor.temp.predicted"
ANOMALY_SUBJECT = "sys.analytics.anomaly"

from reactor import CSTRModel
from detector import AnomalyDetector
import time

async def main():
    nc = NATS()
    model = CSTRModel()
    logger.info("CSTR Model initialized")
    
    # Initialize Detector
    detector = AnomalyDetector(SENSOR_SUBJECT, window_size=50, z_score_threshold=3.0)
    logger.info("Anomaly Detector initialized")
    
    async def disconnected_cb():
        logger.warning("Got disconnected from NATS...")

    async def reconnected_cb():
        logger.info("Got reconnected to NATS...")

    async def error_cb(e):
        logger.error(f"There was an error with NATS: {e}")

    async def closed_cb():
        logger.info("Connection to NATS is closed")

    try:
        await nc.connect(
            servers=[NATS_URL],
            reconnected_cb=reconnected_cb,
            disconnected_cb=disconnected_cb,
            error_cb=error_cb,
            closed_cb=closed_cb,
        )
        logger.info(f"Connected to NATS at {NATS_URL}")
        
        js = nc.jetstream()

        # Subscribe to sensor data
        # We use a durable consumer to ensure we don't miss messages if we restart,
        # but for a real-time twin, we might prioritize latest data. 
        # Using simple subscribe for now.
        async def message_handler(msg):
            try:
                data = common_pb2.SensorData()
                data.ParseFromString(msg.data)
                
                # We expect the sensor_id to be something like "enterprise.site.area.line.reactor.temp"
                # If we get data, we treat it as the input 'Tc' (Cooling Temp) or just feed it to the model 
                # For this demo, we assume the input value is the Cooling Temp (Tc) that drives the process.
                
                loop = asyncio.get_running_loop()
                # Run solver in executor to avoid blocking the event loop
                start_time = time.time()
                pred_T, pred_Ca = await loop.run_in_executor(None, model.solve_step, data.value)
                solve_time = (time.time() - start_time) * 1000

                if pred_T is not None:
                    logger.info(f"Solved in {solve_time:.2f}ms. Input Tc={data.value} -> Pred T={pred_T:.2f}")

                    # Publish Prediction
                    pred_msg = common_pb2.SensorData()
                    pred_msg.sensor_id = PREDICTION_SUBJECT
                    pred_msg.value = pred_T
                    pred_msg.timestamp_ms = int(time.time() * 1000)
                    pred_msg.quality = 1
                    
                    await nc.publish(PREDICTION_SUBJECT, pred_msg.SerializeToString())
                    
                    # Anomaly Detection
                    # NOTE: We are comparing Actual (Input Tc) vs Predicted (Output T) which is chemically invalid 
                    # but functionally demonstrates the pipeline for this story.
                    anomaly = detector.check(actual=data.value, predicted=pred_T)
                    if anomaly:
                        logger.info(f"Publishing Anomaly: {anomaly.severity} Residual={anomaly.residual:.2f}")
                        await nc.publish(ANOMALY_SUBJECT, anomaly.SerializeToString())
                    
                else:
                    logger.warning("Model solve returned None")

            except Exception as e:
                logger.error(f"Failed to process message: {e}")

        await nc.subscribe(SENSOR_SUBJECT, cb=message_handler)
        logger.info(f"Subscribed to {SENSOR_SUBJECT}")

        # Keep alive
        stop_event = asyncio.Event()

        def signal_handler():
            stop_event.set()

        loop = asyncio.get_running_loop()
        for sig in (signal.SIGINT, signal.SIGTERM):
            loop.add_signal_handler(sig, signal_handler)

        await stop_event.wait()

    except Exception as e:
        logger.error(f"Error in main loop: {e}")
    finally:
        await nc.drain()
        logger.info("NATS connection drained")

if __name__ == "__main__":
    try:
        asyncio.run(main())
    except KeyboardInterrupt:
        pass
