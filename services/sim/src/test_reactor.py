import unittest
from reactor import CSTRModel

class TestReactorModel(unittest.TestCase):
    def test_solve_step(self):
        model = CSTRModel()
        
        # Initial solve
        T, Ca = model.solve_step(300.0)
        print(f"Initial: T={T}, Ca={Ca}")
        
        self.assertIsNotNone(T)
        self.assertIsNotNone(Ca)
        self.assertTrue(300 < T < 400) # Reasonable range for reactor temp
        
        # Process perturbation
        for _ in range(5):
             T, Ca = model.solve_step(310.0) # Increase cooling temp
             print(f"Step: T={T}, Ca={Ca}")
             self.assertIsNotNone(T)

if __name__ == '__main__':
    unittest.main()
