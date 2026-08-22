package core

import "context"

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
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return runSolvePipeline(ctx, m, q)
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
