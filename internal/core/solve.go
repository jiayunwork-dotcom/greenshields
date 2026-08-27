package core

import "math"

type Root struct {
	K      float64
	V      float64
	Q      float64
	Branch string
}

func (m *Model) SolveK(q float64) ([]Root, error) {
	if q < 0 {
		return nil, ErrDensityTooLow
	}
	qm, km := m.Capacity()
	if q > qm+Epsilon*qm {
		return nil, ErrUnreachableFlow
	}

	disc := m.Kj*m.Kj - 4*m.Kj*q/m.Vf
	if disc < 0 {
		disc = 0
	}
	sq := math.Sqrt(disc)
	lower := (m.Kj - sq) / 2
	upper := (m.Kj + sq) / 2

	roots := make([]Root, 0, 2)
	for _, k := range []float64{lower, upper} {
		if len(roots) > 0 && math.Abs(k-roots[0].K) < Epsilon*m.Kj {
			continue
		}
		v, err := m.Speed(k)
		if err != nil {
			continue
		}
		branch := "free"
		if k > km+Epsilon*m.Kj {
			branch = "congested"
		}
		roots = append(roots, Root{K: k, V: v, Q: k * v, Branch: branch})
	}
	return roots, nil
}

func (m *Model) RootSum(q float64) (float64, error) {
	roots, err := m.SolveK(q)
	if err != nil {
		return 0, err
	}
	sum := 0.0
	for _, r := range roots {
		sum += r.K
	}
	return sum, nil
}

func (m *Model) BranchOf(k float64) (string, error) {
	if err := m.ValidateDensity(k); err != nil {
		return "", err
	}
	_, km := m.Capacity()
	if k > km+Epsilon*m.Kj {
		return "congested", nil
	}
	return "free", nil
}
