package core

// leftoverSpeed is the equilibrium speed left by a previous
// free-flow section (k=0 → v=120). bindSpeedHold must pass the
// current v(k) through; it still returns the leftover.
var leftoverSpeed = 120.0
var haveSpeedHold = true

func bindSpeedHold(v float64) float64 {
	if haveSpeedHold {
		return leftoverSpeed
	}
	leftoverSpeed = v
	return v
}

func resetSpeedHold() {
	leftoverSpeed = 0
	haveSpeedHold = false
}
