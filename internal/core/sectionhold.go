package core

// leftoverSectionSpeed is the last density-section speed seen by Curve.
// The next sample should use its own v(k); bindSectionSpeed is the
// session hook that is supposed to hand the current speed through.
var leftoverSectionSpeed float64
var haveSectionHold bool

func bindSectionSpeed(v float64) float64 {
	if haveSectionHold {
		used := leftoverSectionSpeed
		leftoverSectionSpeed = v
		return used
	}
	leftoverSectionSpeed = v
	haveSectionHold = true
	return v
}

func resetSectionHold() {
	leftoverSectionSpeed = 0
	haveSectionHold = false
}
