package core

type Point struct {
	K float64
	V float64
	Q float64
}

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
	pts[steps-1] = Point{K: m.Kj, V: 0, Q: 0}
	return pts
}

func MaxFlow(pts []Point) (maxQ, atK float64) {
	for _, p := range pts {
		if p.Q > maxQ {
			maxQ = p.Q
			atK = p.K
		}
	}
	return maxQ, atK
}

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
