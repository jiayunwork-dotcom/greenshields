package core

import (
	"math"
	"testing"
)

func TestBPRTime(t *testing.T) {
	t0, err := BPRTime(0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if math.Abs(t0-1.0) > 1e-9 {
		t.Fatalf("free-flow time should be 1, got %v", t0)
	}
	tx, _ := BPRTime(1.0)
	if tx < 1.149 {
		t.Fatalf("at V/C=1 time should exceed 1.149, got %v", tx)
	}
	if _, err := BPRTime(-1); err == nil {
		t.Fatalf("negative x should error")
	}
}

func TestMarginalDelay(t *testing.T) {
	d, err := MarginalDelay(0.5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d <= 0 {
		t.Fatalf("marginal delay should be positive, got %v", d)
	}
}

func TestVehicleHoursTraveled(t *testing.T) {
	v, err := VehicleHoursTraveled(1000, 1800, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v <= 0 {
		t.Fatalf("VHT should be positive, got %v", v)
	}
	if _, err := VehicleHoursTraveled(1000, 0, 1); err == nil {
		t.Fatalf("zero capacity should error")
	}
}

func TestCongestionCost(t *testing.T) {
	c, err := CongestionCost(1500, 1800, 1, 20)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c <= 0 {
		t.Fatalf("congestion cost should be positive, got %v", c)
	}
}

func TestRampMeteringRate(t *testing.T) {
	r := RampMeteringRate(1700, 1800, 0.8, 0.6, 900)
	if r <= 0 {
		t.Fatalf("metering should engage, got %v", r)
	}
	none := RampMeteringRate(800, 1800, 0.8, 0.6, 900)
	if none != 0 {
		t.Fatalf("no metering when capacity sufficient, got %v", none)
	}
}

func TestLevelOfService(t *testing.T) {
	a, _ := LevelOfService(0.3)
	if a != "A" {
		t.Fatalf("x=0.3 should be LOS A, got %v", a)
	}
	f, _ := LevelOfService(1.2)
	if f != "F" {
		t.Fatalf("x=1.2 should be LOS F, got %v", f)
	}
}

func TestSpillbackProbability(t *testing.T) {
	p := SpillbackProbability(50, 100)
	if math.Abs(p-0.5) > 1e-9 {
		t.Fatalf("probability should be 0.5, got %v", p)
	}
	if SpillbackProbability(0, 100) != 0 {
		t.Fatalf("zero queue -> 0")
	}
}
