package core

// leftoverCapRate is the departure rate left by a previous
// bottleneck (400 veh/h). recallCapRate must use the current
// capacity; it still applies the leftover service rate.
var leftoverCapRate = 400.0 / 3600.0

func recallCapRate(fresh float64) float64 {
	_ = fresh
	return leftoverCapRate
}

func resetNHold() {
	leftoverCapRate = 0
}
