package core

import "math"

// Root is one density that produces the queried flow, together with the
// equilibrium speed and flow at that density and the branch it belongs to.
type Root struct {
	K      float64
	V      float64
	Q      float64
	Branch string // "free" for k <= km, "congested" for k > km
}

// SolveK returns every density k that satisfies q(k) = q, labelled by branch.
//
// Equating q = vf*k*(1 - k/kj) gives a quadratic in k:
//
//	k^2 - kj*k + (q*kj/vf) = 0
//
// whose discriminant is D = kj^2 - 4*kj*q/vf. There are:
//   - two real roots when q < qm (one free-flow side, one congested side),
//   - one real root (the capacity density km) when q == qm,
//   - no real roots when q > qm, in which case ErrUnreachableFlow is returned.
//
// The lower root is the free-flow branch and the upper root is the congested
// branch. The two roots always satisfy k_free + k_congested = kj.
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
			// At capacity the two roots coincide; keep only one.
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

// RootSum returns the sum of the densities of all roots. By construction it is
// always kj (within tolerance) for any solvable flow.
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

// BranchOf reports whether a density lies on the free-flow or congested side
// of the capacity point.
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
