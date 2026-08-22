package core

import (
	"math"
	"testing"
)

// approx compares two floats within a relative tolerance.
func approx(a, b, tol float64) bool {
	denom := math.Max(math.Abs(a), math.Abs(b))
	if denom == 0 {
		return true
	}
	return math.Abs(a-b)/denom <= tol
}

// TestValidateErrors checks the simple, single-rule validation failures:
// non-positive vf, non-positive kj, density above jam, and density below zero.
// (Simple difficulty: a straight rule that must reject bad input.)
func TestValidateErrors(t *testing.T) {
	// Model-parameter errors: New must reject a non-positive vf or kj.
	modelCases := []struct {
		vf, kj float64
	}{
		{-1, 100}, // vf <= 0
		{100, -1}, // kj <= 0
	}
	for _, c := range modelCases {
		if _, err := New(c.vf, c.kj); err == nil {
			t.Errorf("New(%g,%g): expected parameter error, got nil", c.vf, c.kj)
		}
	}

	// Density errors: ValidateDensity must reject k outside [0, kj].
	m, _ := New(100, 100)
	densityCases := []float64{200, -5}
	for _, k := range densityCases {
		if err := m.ValidateDensity(k); err == nil {
			t.Errorf("ValidateDensity(%g): expected error for out-of-range density", k)
		}
	}

	// A genuinely valid configuration must not error.
	if _, err := New(120, 180); err != nil {
		t.Errorf("New(120,180) unexpected error: %v", err)
	}
	if err := m.ValidateDensity(90); err != nil {
		t.Errorf("ValidateDensity(90): unexpected error: %v", err)
	}
}

// TestCapacityQm checks qm = vf*kj/4, km = kj/2, and vm = vf/2.
func TestCapacityQm(t *testing.T) {
	m, _ := New(120, 180)
	qm, km := m.Capacity()
	if !approx(qm, 120*180/4, 1e-9) {
		t.Errorf("qm = %g, want %g", qm, float64(120*180/4))
	}
	if !approx(km, 90, 1e-9) {
		t.Errorf("km = %g, want 90", km)
	}
	if !approx(m.SpeedAtCapacity(), 60, 1e-9) {
		t.Errorf("vm = %g, want 60", m.SpeedAtCapacity())
	}
}

// TestZeroDensity checks the free-flow boundary: k=0 -> v=vf, q=0.
func TestZeroDensity(t *testing.T) {
	m, _ := New(120, 180)
	v, err := m.Speed(0)
	if err != nil {
		t.Fatalf("Speed(0) error: %v", err)
	}
	if v != 120 {
		t.Errorf("v(0) = %g, want 120 (free-flow speed)", v)
	}
	q, _ := m.Flow(0)
	if q != 0 {
		t.Errorf("q(0) = %g, want 0", q)
	}
}

// TestJamDensity checks the jam boundary: k=kj -> v=0, q=0.
func TestJamDensity(t *testing.T) {
	m, _ := New(120, 180)
	v, err := m.Speed(180)
	if err != nil {
		t.Fatalf("Speed(kj) error: %v", err)
	}
	if v != 0 {
		t.Errorf("v(kj) = %g, want 0", v)
	}
	q, _ := m.Flow(180)
	if q != 0 {
		t.Errorf("q(kj) = %g, want 0", q)
	}
}

// TestSolveKDoubleRoots checks that a sub-capacity flow yields two labelled
// roots and that each root reproduces the requested flow.
func TestSolveKDoubleRoots(t *testing.T) {
	m, _ := New(120, 180)
	q := 2700.0 // = qm/2
	roots, err := m.SolveK(q)
	if err != nil {
		t.Fatalf("SolveK(%g) error: %v", q, err)
	}
	if len(roots) != 2 {
		t.Fatalf("len(roots) = %d, want 2", len(roots))
	}
	if roots[0].Branch != "free" {
		t.Errorf("root[0].Branch = %q, want free", roots[0].Branch)
	}
	if roots[1].Branch != "congested" {
		t.Errorf("root[1].Branch = %q, want congested", roots[1].Branch)
	}
	for _, r := range roots {
		back, _ := m.Flow(r.K)
		if !approx(back, q, 1e-9) {
			t.Errorf("Flow(%g) = %g, want %g (root must reproduce q)", r.K, back, q)
		}
	}
}

// TestSolveKSumEqualsKj verifies the identity k_free + k_congested = kj.
func TestSolveKSumEqualsKj(t *testing.T) {
	m, _ := New(120, 180)
	roots, err := m.SolveK(2700)
	if err != nil {
		t.Fatalf("SolveK error: %v", err)
	}
	sum := roots[0].K + roots[1].K
	if !approx(sum, 180, 1e-9) {
		t.Errorf("k_free + k_congested = %g, want 180", sum)
	}
}

// TestCrossRuleDoubleVf checks that doubling only vf doubles qm while kj is
// unchanged.
func TestCrossRuleDoubleVf(t *testing.T) {
	m, _ := New(120, 180)
	d := ModelWithDoubleVf(m)
	if d.Kj != 180 {
		t.Errorf("doubling vf changed kj to %g, want 180", d.Kj)
	}
	if !approx(CapacityRatio(m, d), 2, 1e-9) {
		t.Errorf("capacity ratio after doubling vf = %g, want 2", CapacityRatio(m, d))
	}
}

// TestCrossRuleDoubleKj checks that doubling only kj doubles qm while vf is
// unchanged.
func TestCrossRuleDoubleKj(t *testing.T) {
	m, _ := New(120, 180)
	d := ModelWithDoubleKj(m)
	if d.Vf != 120 {
		t.Errorf("doubling kj changed vf to %g, want 120", d.Vf)
	}
	if !approx(CapacityRatio(m, d), 2, 1e-9) {
		t.Errorf("capacity ratio after doubling kj = %g, want 2", CapacityRatio(m, d))
	}
}

// TestCurvePeakAtHalf checks that the sampled curve peaks at the capacity
// density and that the peak equals qm.
func TestCurvePeakAtHalf(t *testing.T) {
	m, _ := New(120, 180)
	pts := m.Curve(101)
	maxQ, atK := MaxFlow(pts)
	qm, km := m.Capacity()
	if !approx(atK, km, 1e-9) {
		t.Errorf("curve peaks at k=%g, want %g", atK, km)
	}
	if !approx(maxQ, qm, 1e-9) {
		t.Errorf("curve peak q=%g, want %g", maxQ, qm)
	}
}

// TestSolveKNoRealRoots checks that a flow above capacity yields an error.
func TestSolveKNoRealRoots(t *testing.T) {
	m, _ := New(120, 180)
	if _, err := m.SolveK(10000); err == nil {
		t.Errorf("SolveK(10000) above capacity: expected error, got nil")
	}
}
