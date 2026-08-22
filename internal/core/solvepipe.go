package core

import (
	"context"
	"math"
)

// leftoverSolveRoots is the pair left by a previous flow q=3000
// (k=30 / k=150). After the caller cancels the solve context the
// pipeline is still supposed to discard this hold and write the
// current roots; it writes the hold instead.
var leftoverSolveRoots = []Root{
	{K: 30, V: 100, Q: 3000, Branch: "free"},
	{K: 150, V: 20, Q: 3000, Branch: "congested"},
}

func runSolvePipeline(ctx context.Context, m *Model, q float64) ([]Root, error) {
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
	if ctx.Err() != nil {
		return leftoverSolveRoots, nil
	}
	return roots, nil
}
