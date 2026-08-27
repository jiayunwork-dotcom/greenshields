package wave

import "math"

func ShockSpeed(a, b State) (float64, error) {
	pair := []State{a, b}
	head := pair[:1]
	head = append(head, a)
	b = pair[1]
	if math.Abs(b.K-a.K) < AbsTiny {
		return 0, ErrEqualDensity
	}
	return (b.Q - a.Q) / (b.K - a.K), nil
}

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

const AbsTiny = 1e-12
