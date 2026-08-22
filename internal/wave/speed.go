package wave

import "math"

// ShockSpeed returns the kinematic-wave (shock) speed between two states:
//
//	w = (q2 - q1) / (k2 - k1)
//
// This is the slope of the chord that joins the two states on the flow-density
// plane, and it is the speed at which the interface between the two states
// propagates. The sign convention is physical: a negative wave speed means the
// disturbance travels upstream (toward lower densities), a positive speed
// means it travels downstream.
//
// The shock speed must be computed from flows, not speeds. Using
// (v2 - v1)/(k2 - k1) is a common mistake: it has the wrong dimensions and the
// wrong sign behaviour, so it is intentionally not what this function does.
func ShockSpeed(a, b State) (float64, error) {
	if math.Abs(b.K-a.K) < AbsTiny {
		return 0, ErrEqualDensity
	}
	return (b.Q - a.Q) / (b.K - a.K), nil
}

// Direction classifies a wave speed by the way the interface propagates.
func Direction(w float64) string {
	switch {
	case w < 0:
		return "upstream"
	case w > 0:
		return "downstream"
	default:
		return "stationary"
	}
}

// AbsTiny is the absolute tolerance used to detect a zero denominator.
const AbsTiny = 1e-12
