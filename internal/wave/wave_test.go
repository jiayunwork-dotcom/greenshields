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

// TestWaveShockFormulaSign is the high-difficulty check. It verifies the wave
// speed uses (q2-q1)/(k2-k1) by computing a spillback analytically. A common
// mistake is to use (v2-v1)/(k2-k1); this test pins the correct value
// (deterministic: w = -vf/4) so that wrong formula is caught. (High difficulty)
func TestWaveShockFormulaSign(t *testing.T) {
	vf, kj := 120.0, 180.0
	m, _ := core.New(vf, kj)

	k1 := kj / 4.0 // free-flow upstream
	k2 := kj       // jammed downstream
	q1 := mustFlow(m, k1)
	q2 := mustFlow(m, k2)

	w, err := ShockSpeed(State{K: k1, Q: q1}, State{K: k2, Q: q2})
	if err != nil {
		t.Fatalf("ShockSpeed error: %v", err)
	}
	// Analytic: w = (0 - q1) / (kj - kj/4) = -vf/4.
	want := -vf / 4.0
	if math.Abs(w-want) > 1e-6 {
		t.Errorf("wave speed = %g, want %g (must use (q2-q1)/(k2-k1), not (v2-v1)/...)", w, want)
	}
}

// TestWaveNegativeSpillback checks that an upstream free-flow / downstream jam
// transition yields a negative (upstream-propagating) wave speed. (Medium)
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

// TestWaveZeroDenominator checks that identical densities are rejected.
func TestWaveZeroDenominator(t *testing.T) {
	if _, err := ShockSpeed(State{K: 50, Q: 1000}, State{K: 50, Q: 1500}); err == nil {
		t.Errorf("expected ErrEqualDensity for equal densities, got nil")
	}
}
