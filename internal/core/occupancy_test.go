package core

import (
	"math"
	"testing"
)

func mk() *Model {
	m, _ := New(120, 120)
	return m
}

func TestSpacing(t *testing.T) {
	s, err := Spacing(50)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if math.Abs(s-20.0) > 1e-9 {
		t.Fatalf("spacing should be 20 m, got %v", s)
	}
	if _, err := Spacing(0); err == nil {
		t.Fatalf("zero density should error")
	}
}

func TestTimeHeadway(t *testing.T) {
	m := mk()
	h, err := m.TimeHeadway(30)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h <= 0 {
		t.Fatalf("headway should be positive, got %v", h)
	}
}

func TestOccupancy(t *testing.T) {
	occ, err := Occupancy(120, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if math.Abs(occ-0.6) > 1e-9 {
		t.Fatalf("occupancy should be 0.6, got %v", occ)
	}
}

func TestJamSpacing(t *testing.T) {
	m := mk()
	s, err := m.JamSpacing()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s <= 0 {
		t.Fatalf("jam spacing should be positive, got %v", s)
	}
}

func TestSpeedAtHeadway(t *testing.T) {
	m := mk()
	v, err := m.SpeedAtHeadway(30, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v <= 0 {
		t.Fatalf("speed should be positive, got %v", v)
	}
}

func TestDensityFromHeadway(t *testing.T) {
	k, err := DensityFromHeadway(60, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// q = 3600/2 = 1800 veh/h, k = 1800/60 = 30 veh/km
	if math.Abs(k-30.0) > 1e-9 {
		t.Fatalf("density should be 30, got %v", k)
	}
}

func TestCriticalHeadway(t *testing.T) {
	m := mk()
	h, err := m.CriticalHeadway()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h <= 0 {
		t.Fatalf("critical headway should be positive, got %v", h)
	}
}

func TestSpeedDrop(t *testing.T) {
	m := mk()
	d, err := m.SpeedDrop(90)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d <= 0 || d >= 1 {
		t.Fatalf("speed drop should be in (0,1), got %v", d)
	}
}

func TestUniformity(t *testing.T) {
	m := mk()
	u, err := m.Uniformity(10, 100, 20)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u < 0 {
		t.Fatalf("uniformity cannot be negative, got %v", u)
	}
	if _, err := m.Uniformity(10, 5, 20); err == nil {
		t.Fatalf("k2<=k1 should error")
	}
}
