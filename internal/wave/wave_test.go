package wave

import (
	"math"
	"testing"

	"greenshields/internal/core"
)

func mustFlow(m *core.Model, k float64) float64 {
	q, err := m.Flow(k)
	if err != nil {
		panic(err)
	}
	return q
}

func TestWaveShockFormulaSign(t *testing.T) {
	vf, kj := 120.0, 180.0
	m, _ := core.New(vf, kj)

	k1 := kj / 4.0
	k2 := kj
	q1 := mustFlow(m, k1)
	q2 := mustFlow(m, k2)

	w, err := ShockSpeed(State{K: k1, Q: q1}, State{K: k2, Q: q2})
	if err != nil {
		t.Fatalf("ShockSpeed error: %v", err)
	}
	want := -vf / 4.0
	if math.Abs(w-want) > 1e-6 {
		t.Errorf("wave speed = %g, want %g (must use (q2-q1)/(k2-k1), not (v2-v1)/...)", w, want)
	}
}

func TestWaveNegativeSpillback(t *testing.T) {
	vf, kj := 120.0, 180.0
	m, _ := core.New(vf, kj)

	up := State{K: kj / 4.0, Q: mustFlow(m, kj/4.0)}
	down := State{K: kj, Q: 0}

	sp, err := ClassifySpillback(up, down)
	if err != nil {
		t.Fatalf("ClassifySpillback error: %v", err)
	}
	if !sp.IsUpstream() {
		t.Errorf("spillback wave speed = %g, want < 0 (upstream propagation)", sp.WaveSpeed)
	}
	if sp.Direction != "upstream" {
		t.Errorf("direction = %q, want upstream", sp.Direction)
	}
}

func TestWaveZeroDenominator(t *testing.T) {
	if _, err := ShockSpeed(State{K: 50, Q: 1000}, State{K: 50, Q: 1500}); err == nil {
		t.Errorf("expected ErrEqualDensity for equal densities, got nil")
	}
}
