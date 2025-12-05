package core

import "testing"

func TestEvaluate(t *testing.T) {
	tests := []struct {
		name       string
		def        AlarmDefinition
		value      float64
		shouldFire bool
	}{
		{
			name:       "High Alarm Trigger",
			def:        AlarmDefinition{Type: "High", Threshold: 100},
			value:      101,
			shouldFire: true,
		},
		{
			name:       "High Alarm Clear",
			def:        AlarmDefinition{Type: "High", Threshold: 100},
			value:      99,
			shouldFire: false,
		},
		{
			name:       "Low Alarm Trigger",
			def:        AlarmDefinition{Type: "Low", Threshold: 50},
			value:      49,
			shouldFire: true,
		},
		{
			name:       "Low Alarm Clear",
			def:        AlarmDefinition{Type: "Low", Threshold: 50},
			value:      51,
			shouldFire: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fired := Evaluate(&tt.def, tt.value)
			if fired != tt.shouldFire {
				t.Errorf("Expected fired=%v, got %v", tt.shouldFire, fired)
			}
		})
	}
}
