import unittest
import numpy as np
from detector import AnomalyDetector
class TestAnomalyDetector(unittest.TestCase):
    def test_detection(self):
        detector = AnomalyDetector("test_tag")
        
        # Feed normal data (mean=0, std=1)
        # We need at least 20 points
        for _ in range(50):
            # Perfect prediction: residual = 0 + noise
            detector.check(10.0 + np.random.normal(0, 0.1), 10.0)
            
        # No anomaly yet (residuals are small)
        event = detector.check(10.0, 10.0)
        self.assertIsNone(event)
        
        # Introduce anomaly: huge residual
        # Actual=20, Pred=10 -> Residual=10.
        # Check against small noise std dev (~0.1) -> Z Score ~ 100
        event = detector.check(20.0, 10.0)
        
        self.assertIsNotNone(event)
        self.assertEqual(event.severity, "CRITICAL")
        self.assertAlmostEqual(event.residual, 10.0, delta=0.1)

if __name__ == '__main__':
    unittest.main()
