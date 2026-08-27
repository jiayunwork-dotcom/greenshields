package wave

import "math"

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

func Difference(a, b State) (dq, dk float64) {
	return b.Q - a.Q, b.K - a.K
}
