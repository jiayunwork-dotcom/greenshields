package core

import (
	"math"
	"testing"
)

func TestTripGeneration(t *testing.T) {
	trips, err := TripGeneration(1000, 0.3, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if math.Abs(trips-300) > 1e-9 {
		t.Fatalf("trips should be 300, got %v", trips)
	}
	capped, _ := TripGeneration(100000, 0.3, 100)
	if math.Abs(capped-100) > 1e-9 {
		t.Fatalf("should cap at 100, got %v", capped)
	}
}

func TestLogitChoice(t *testing.T) {
	p, err := LogitChoice(1.0, []float64{0.5, 0.2}, 1.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p <= 0 || p >= 1 {
		t.Fatalf("probability must be in (0,1), got %v", p)
	}
	if _, err := LogitChoice(1.0, []float64{}, 0); err == nil {
		t.Fatalf("theta=0 should error")
	}
}

func TestModeSplit(t *testing.T) {
	m1, m2, err := ModeSplit(1000, 0.5, 1.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if math.Abs(m1+m2-1000) > 1e-9 {
		t.Fatalf("split must sum to demand, got %v+%v", m1, m2)
	}
}

func TestGrowthFactor(t *testing.T) {
	out, err := GrowthFactor([]float64{10, 20, 30}, 1.1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if math.Abs(out[2]-33) > 1e-9 {
		t.Fatalf("third should be 33, got %v", out[2])
	}
}

func TestAttractionConstrained(t *testing.T) {
	out, f, err := AttractionConstrained([]float64{10, 20, 30}, 60)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if math.Abs(f-1.0) > 1e-9 {
		t.Fatalf("factor should be 1, got %v", f)
	}
	sum := 0.0
	for _, v := range out {
		sum += v
	}
	if math.Abs(sum-60) > 1e-9 {
		t.Fatalf("should sum to 60, got %v", sum)
	}
}

func TestAverageTripLength(t *testing.T) {
	avg, err := AverageTripLength([]float64{1, 5, 10}, []float64{10, 10, 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if math.Abs(avg-5.333333) > 1e-6 {
		t.Fatalf("average should be 16/3, got %v", avg)
	}
}
