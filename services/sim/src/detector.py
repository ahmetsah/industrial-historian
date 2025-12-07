import numpy as np
from collections import deque
import logging
import analytics_pb2
import time

logger = logging.getLogger("anomaly-detector")

class AnomalyDetector:
    def __init__(self, source_tag: str, window_size: int = 100, z_score_threshold: float = 3.0):
        self.source_tag = source_tag
        self.window_size = window_size
        self.z_score_threshold = z_score_threshold
        self.residuals = deque(maxlen=window_size)
    
    def check(self, actual: float, predicted: float) -> analytics_pb2.AnomalyEvent | None:
        """
        Checks for anomaly given current actual and predicted values.
        Returns AnomalyEvent if anomaly detected, else None.
        """
        residual = abs(actual - predicted)
        self.residuals.append(residual)
        
        # Need minimum samples to detect statistically
        if len(self.residuals) < 20:
            return None
            
        window = np.array(self.residuals)
        mean = np.mean(window)
        std = np.std(window)
        
        # Avoid division by zero
        if std == 0:
            return None
            
        z_score = (residual - mean) / std
        
        if abs(z_score) > self.z_score_threshold:
            logger.warning(f"Anomaly detected! Z-Score={z_score:.2f}, Residual={residual:.4f}, Mean={mean:.4f}, Std={std:.4f}")
            
            event = analytics_pb2.AnomalyEvent()
            event.source_tag = self.source_tag
            event.actual = actual
            event.predicted = predicted
            event.residual = residual
            event.timestamp_ms = int(time.time() * 1000)
            event.severity = "CRITICAL" if abs(z_score) > 5.0 else "WARNING"
            
            return event
            
        return None
