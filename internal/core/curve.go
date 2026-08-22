package core

// Point is a single sample of the fundamental diagram.
type Point struct {
	K float64
	V float64
	Q float64
}

// Curve samples both q(k) and v(k) from density 0 up to the jam density in
// steps points. The endpoints are pinned exactly: the first point is k = 0
// (v = vf, q = 0) and the final point is k = kj (v = 0, q = 0). This makes the
// curve suitable for plotting and for downstream wave calculations that need a
// dense, ordered set of states.
func (m *Model) Curve(steps int) []Point {
	if steps < 2 {
		steps = DefaultSteps
	}
	pts := make([]Point, 0, steps)
	for i := 0; i < steps; i++ {
		k := m.Kj * float64(i) / float64(steps-1)
		v, _ := m.Speed(k)
		q, _ := m.Flow(k)
		pts = append(pts, Point{K: k, V: v, Q: q})
	}
	// Pin the jam endpoint exactly regardless of floating point rounding.
	pts[steps-1] = Point{K: m.Kj, V: 0, Q: 0}
	return pts
}

// MaxFlow returns the largest flow present in a curve. For a Greenshields
// curve this equals the capacity qm and occurs at k = km.
func MaxFlow(pts []Point) (maxQ, atK float64) {
	for _, p := range pts {
		if p.Q > maxQ {
			maxQ = p.Q
			atK = p.K
		}
	}
	return maxQ, atK
}

// SampleAt returns the flow and speed at a given density by linear
// interpolation between the two surrounding curve samples. It is useful when a
// caller already has a curve in hand and wants a quick lookup without building
// a fresh model.
func SampleAt(pts []Point, k float64) (v, q float64, ok bool) {
	if len(pts) < 2 {
		return 0, 0, false
	}
	if k < pts[0].K || k > pts[len(pts)-1].K {
		return 0, 0, false
	}
	for i := 1; i < len(pts); i++ {
		if k <= pts[i].K {
			a, b := pts[i-1], pts[i]
			if b.K == a.K {
				return a.V, a.Q, true
			}
			t := (k - a.K) / (b.K - a.K)
			return a.V + t*(b.V-a.V), a.Q + t*(b.Q-a.Q), true
		}
	}
	last := pts[len(pts)-1]
	return last.V, last.Q, true
}
