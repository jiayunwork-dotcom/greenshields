package godunov

import (
	"math"
	"testing"
)

func mustSolver(t *testing.T) *Solver {
	t.Helper()
	s, err := New(120, 180)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return s
}

func TestGodunovFluxEqualsMinDemandSupply(t *testing.T) {
	s := mustSolver(t)
	pairs := [][2]float64{
		{20, 40},
		{40, 140},
		{90, 90},
		{160, 20},
		{0, 180},
		{180, 0},
	}
	for _, p := range pairs {
		d, err := s.Demand(p[0])
		if err != nil {
			t.Fatalf("Demand(%g): %v", p[0], err)
		}
		sup, err := s.Supply(p[1])
		if err != nil {
			t.Fatalf("Supply(%g): %v", p[1], err)
		}
		want := d
		if sup < want {
			want = sup
		}
		got, err := s.Flux(p[0], p[1])
		if err != nil {
			t.Fatalf("Flux(%g,%g): %v", p[0], p[1], err)
		}
		if math.Abs(got-want) > 1e-9*math.Max(1, math.Abs(want)) {
			t.Errorf("Flux(%g,%g)=%g, want min(D,S)=%g", p[0], p[1], got, want)
		}
	}
}

func TestGodunovEntropyAgreesDemandSupply(t *testing.T) {
	s := mustSolver(t)
	ks := []float64{0, 30, 45, 90, 135, 162, 180}
	for _, kL := range ks {
		for _, kR := range ks {
			got, err := s.FluxAgrees(kL, kR)
			if err != nil {
				t.Fatalf("FluxAgrees(%g,%g): %v", kL, kR, err)
			}
			naive, err := s.NaiveAverageFlux(kL, kR)
			if err != nil {
				t.Fatalf("NaiveAverageFlux: %v", err)
			}
			needCap, err := s.RarefactionNeedsCapacity(kL, kR)
			if err != nil {
				t.Fatalf("RarefactionNeedsCapacity: %v", err)
			}
			if needCap && math.Abs(naive-s.Qm()) > 1 && math.Abs(got-s.Qm()) > 1e-6*s.Qm() {
				t.Errorf("rarefaction across peak: flux=%g want qm=%g (naive average %g is not entropy)", got, s.Qm(), naive)
			}
		}
	}
}

func TestGodunovStationaryShockZeroSpeed(t *testing.T) {
	s := mustSolver(t)
	q := 2700.0
	res, err := s.StationaryShock(q)
	if err != nil {
		t.Fatalf("StationaryShock: %v", err)
	}
	if math.Abs(res.WaveSpeed) > 1e-8 {
		t.Errorf("wave speed = %g, want 0", res.WaveSpeed)
	}
	if res.Direction != "stationary" {
		t.Errorf("direction = %q, want stationary", res.Direction)
	}
	dx := 0.1
	dt, err := s.TimeStep(dx, res.KL, res.KR, 0.5)
	if err != nil {
		t.Fatalf("TimeStep: %v", err)
	}
	if err := s.HoldStationary(q, dx, dt); err != nil {
		t.Fatalf("HoldStationary: %v", err)
	}
}

func TestGodunovRarefactionAcrossPeakUsesCapacity(t *testing.T) {
	s := mustSolver(t)
	res, err := s.JamRelease()
	if err != nil {
		t.Fatalf("JamRelease: %v", err)
	}
	if math.Abs(res.Flux-s.Qm()) > 1e-9*s.Qm() {
		t.Errorf("jam-to-empty flux = %g, want qm=%g", res.Flux, s.Qm())
	}
	avg, err := s.NaiveAverageFlux(s.Kj(), 0)
	if err != nil {
		t.Fatalf("NaiveAverageFlux: %v", err)
	}
	if math.Abs(avg) > 1e-9 {
		t.Fatalf("sanity: jam and empty both have q=0, naive average should be 0, got %g", avg)
	}
	blocked, err := s.EmptyAgainstJam()
	if err != nil {
		t.Fatalf("EmptyAgainstJam: %v", err)
	}
	if math.Abs(blocked.Flux) > 1e-9 {
		t.Errorf("empty-into-jam flux = %g, want 0", blocked.Flux)
	}
	dx := 0.25
	dt, err := s.TimeStep(dx, s.Kj(), 0, 0.4)
	if err != nil {
		t.Fatalf("TimeStep: %v", err)
	}
	nL, nR, err := s.StepPair(s.Kj(), 0, dx, dt)
	if err != nil {
		t.Fatalf("StepPair: %v", err)
	}
	if err := MassConserved(s.Kj(), 0, nL, nR); err != nil {
		t.Fatalf("mass: %v", err)
	}
}

func TestGodunovRejectsIllegalDensity(t *testing.T) {
	s := mustSolver(t)
	if _, err := s.Flux(-1, 10); err == nil {
		t.Fatal("negative density must fail")
	}
	if _, err := s.Flux(10, 400); err == nil {
		t.Fatal("k > kj must fail")
	}
	if _, err := New(-1, 180); err == nil {
		t.Fatal("negative vf must fail")
	}
}
