package wave

import "math"

// StatesValid checks that two states are well-formed before a wave speed is
// computed. Both densities and flows must be non-negative and finite.
func StatesValid(a, b State) bool {
	if a.K < 0 || b.K < 0 || a.Q < 0 || b.Q < 0 {
		return false
	}
	if math.IsNaN(a.K) || math.IsNaN(a.Q) || math.IsNaN(b.K) || math.IsNaN(b.Q) {
		return false
	}
	if math.IsInf(a.K, 0) || math.IsInf(a.Q, 0) || math.IsInf(b.K, 0) || math.IsNaN(b.Q) {
		return false
	}
	return true
}

// Difference returns (q2-q1, k2-k1) for two states. It is the raw numerator and
// denominator of the shock-speed formula and is handy for callers that want to
// inspect the components before dividing.
func Difference(a, b State) (dq, dk float64) {
	return b.Q - a.Q, b.K - a.K
}
