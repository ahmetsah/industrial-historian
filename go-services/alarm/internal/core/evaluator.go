package core

func Evaluate(def *AlarmDefinition, value float64) bool {
	switch def.Type {
	case "High":
		return value > def.Threshold
	case "Low":
		return value < def.Threshold
	default:
		return false
	}
}
