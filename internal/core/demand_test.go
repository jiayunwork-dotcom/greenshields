package core

import (
	"math"
	"testing"
)

func gard() *Model {
	m, _ := New(120, 120)
	return m
}

func TestSaturationRatio(t *testing.T) {
	m := gard()
	r, err := m.SaturationRatio(60)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	qm, _ := m.Capacity()
	if math.Abs(r-60/qm) > 1e-9 {
		t.Fatalf("saturation ratio wrong: %v", r)
	}
}

func TestWebsterDelayFree(t *testing.T) {
	d, err := WebsterDelay(600, 1800, 0.5, 120)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d <= 0 {
		t.Fatalf("delay should be positive, got %v", d)
	}
}

func TestWebsterDelayBadSignal(t *testing.T) {
	if _, err := WebsterDelay(600, 1800, 0, 120); err == nil {
		t.Fatalf("zero green fraction should error")
	}
	if _, err := WebsterDelay(600, 1800, 0.5, 0); err == nil {
		t.Fatalf("zero cycle should error")
	}
}

func TestTripTime(t *testing.T) {
	m := gard()
	tt, err := m.TripTime(30, 1.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tt <= 0 {
		t.Fatalf("trip time should be positive, got %v", tt)
	}
	if _, err := m.TripTime(120, 1.0); err == nil {
	}
}

func TestDensityFromFlow(t *testing.T) {
	m := gard()
	free, cong, branch, err := m.DensityFromFlow(1500)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if free >= cong {
		t.Fatalf("free-flow density should be below congested, got %v vs %v", free, cong)
	}
	if branch != "free" {
		t.Fatalf("low flow should be free branch, got %v", branch)
	}
}

func TestTravelTimeIndex(t *testing.T) {
	m := gard()
	idx, err := m.TravelTimeIndex(30)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if idx <= 1 {
		t.Fatalf("congestion index should be > 1 away from free flow, got %v", idx)
	}
}

func TestQueueLength(t *testing.T) {
	q := QueueLength(2000, 1800, 3600)
	if q <= 0 {
		t.Fatalf("queue should be positive when oversaturated, got %v", q)
	}
	if QueueLength(1000, 1800, 3600) != 0 {
		t.Fatalf("under capacity queue should be 0")
	}
}

func TestSpaceMeanSpeed(t *testing.T) {
	s, err := SpaceMeanSpeed(1500, 30)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if math.Abs(s-50.0) > 1e-9 {
		t.Fatalf("space-mean speed should be 50, got %v", s)
	}
	if _, err := SpaceMeanSpeed(1500, 0); err == nil {
		t.Fatalf("zero density should error")
	}
}
