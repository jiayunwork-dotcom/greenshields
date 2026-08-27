package core

import (
	"testing"
)

func TestArrivalDeparture(t *testing.T) {
	curve, maxQ, total := ArrivalDeparture(2000, 1800, 3600)
	if len(curve) == 0 {
		t.Fatalf("curve should be non-empty")
	}
	if maxQ <= 0 {
		t.Fatalf("max queue should be positive under oversaturation, got %v", maxQ)
	}
	if total <= 0 {
		t.Fatalf("total delay should be positive, got %v", total)
	}
}

func TestArrivalDepartureUnderCapacity(t *testing.T) {
	curve, maxQ, total := ArrivalDeparture(1000, 1800, 3600)
	if maxQ > 1e-6 {
		t.Fatalf("no queue expected under capacity, got %v", maxQ)
	}
	if total > 1e-6 {
		t.Fatalf("no delay expected under capacity, got %v", total)
	}
	_ = curve
}

func TestTravelTimeFromNCurve(t *testing.T) {
	curve, _, _ := ArrivalDeparture(2000, 1800, 3600)
	dt := TravelTimeFromNCurve(curve, 100)
	if dt < 0 {
		t.Fatalf("travel time should be non-negative, got %v", dt)
	}
	if TravelTimeFromNCurve(nil, 100) != 0 {
		t.Fatalf("empty curve should yield 0")
	}
}

func TestShockObserver(t *testing.T) {
	o := ShockObserver{}
	if o.Update(30, 1500) {
		t.Fatalf("first sample should not be a shock")
	}
	if !o.Update(60, 1000) {
		t.Fatalf("backward shock expected")
	}
	if o.Update(60, 1500) {
		t.Fatalf("no shock when flow rises")
	}
}

func TestThroughputLoss(t *testing.T) {
	loss := ThroughputLoss(1800, 1500, 3600)
	if loss <= 0 {
		t.Fatalf("throughput loss should be positive, got %v", loss)
	}
	if ThroughputLoss(1800, 1800, 3600) != 0 {
		t.Fatalf("no loss when at capacity")
	}
}
