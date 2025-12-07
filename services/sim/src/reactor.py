from gekko import GEKKO
import logging

logger = logging.getLogger("process-model")

class CSTRModel:
    def __init__(self):
        self.m = GEKKO(remote=False)
        
        # Time discretization (for dynamic simulation)
        # We simulate 1 step ahead
        self.m.time = [0, 0.1] 

        # Parameters / Inputs (Manipulated Variables)
        # Tc: Cooling Jacket Temperature (K)
        self.Tc = self.m.MV(value=300)
        self.Tc.STATUS = 0 # Measuring, not optimizing (it's an input)
        self.Tc.FSTATUS = 0 # Receives measurement

        # Variables / Outputs (Controlled Variables)
        # T: Reactor Temperature (K)
        self.T = self.m.CV(value=350) 
        self.T.STATUS = 1 # We want to simulate this
        self.T.FSTATUS = 1 # Feedback status

        # Ca: Concentration (mol/L)
        self.Ca = self.m.CV(value=0.8)
        self.Ca.STATUS = 1
        self.Ca.FSTATUS = 1

        # Constants
        self.q = 100.0  # Volumetric Flow Rate (L/min)
        self.V = 100.0  # Volume (L)
        self.rho = 1000.0 # Density (g/L)
        self.Cp = 0.239 # Heat Capacity (J/g K)
        self.mdelH = 50000.0 # Heat of Reaction (J/mol)
        self.E_R = 8750.0 # Activation Energy / R (K)
        self.k0 = 7.2e10 # Pre-exponential factor (1/min)
        self.UA = 50000.0 # Heat Transfer Coefficient (J/min K)
        self.Tf = 350.0 # Feed Temperature (K)
        self.Caf = 1.0 # Feed Concentration (mol/L)

        # Equations
        # Arrhenius Rate Law
        # k = k0 * exp(-E/RT)
        k = self.m.Intermediate(self.k0 * self.m.exp(-self.E_R / self.T))
        
        # Mass Balance: accumulation = in - out - reaction
        # V * dCa/dt = q*(Caf - Ca) - V*k*Ca
        self.m.Equation(self.V * self.Ca.dt() == self.q * (self.Caf - self.Ca) - self.V * k * self.Ca)

        # Energy Balance: accumulation = in - out + reaction - cooling
        # V*rho*Cp * dT/dt = q*rho*Cp*(Tf - T) + V*mdelH*k*Ca + UA*(Tc - T)
        self.m.Equation(self.V * self.rho * self.Cp * self.T.dt() == \
                        self.q * self.rho * self.Cp * (self.Tf - self.T) \
                        + self.V * self.mdelH * k * self.Ca \
                        + self.UA * (self.Tc - self.T))

        # Configuration
        self.m.options.IMODE = 4 # Dynamic Simulation
        self.m.options.NODES = 3 # Collocation nodes
        self.m.options.SOLVER = 3 # IPOPT

    def solve_step(self, cooling_temp_input: float):
        """
        Advances the simulation by one time step given the input.
        """
        try:
            # Update Input (MV)
            self.Tc.MEAS = cooling_temp_input
            
            # Solve
            # disp=False to suppress output
            self.m.solve(disp=False)
            
            # Extract Results (last point in horizon)
            predicted_T = self.T.value[-1]
            predicted_Ca = self.Ca.value[-1]
            
            return predicted_T, predicted_Ca
            
        except Exception as e:
            logger.error(f"Solver failed: {e}")
            return None, None
