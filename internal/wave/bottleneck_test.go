package wave

import (
	"math"
	"testing"
)

func TestBottleneckActive(t *testing.T) {
	b := Bottleneck{Arrival: 2000, Capacity: 1800, Reduction: 0.2}
	if !b.Active() {
		t.Fatalf("should be active when arrival exceeds reduced capacity")
	}
	if math.Abs(b.EffectiveCapacity()-1440) > 1e-9 {
		t.Fatalf("effective capacity should be 1440, got %v", b.EffectiveCapacity())
	}
}

func TestBottleneckGrowthRate(t *testing.T) {
	b := Bottleneck{Arrival: 2000, Capacity: 1800, Reduction: 0.2}
	if math.Abs(b.QueueGrowthRate()-560) > 1e-9 {
		t.Fatalf("growth rate should be 560 veh/h, got %v", b.QueueGrowthRate())
	}
	b2 := Bottleneck{Arrival: 1000, Capacity: 1800, Reduction: 0.2}
	if b2.QueueGrowthRate() != 0 {
		t.Fatalf("no growth under capacity")
	}
}

func TestBottleneckQueueAt(t *testing.T) {
	b := Bottleneck{Arrival: 2000, Capacity: 1800, Reduction: 0.2}
	q := b.QueueAt(10, 3600)
	if q <= 10 {
		t.Fatalf("queue should grow, got %v", q)
	}
}

func TestBottleneckBackWaveSpeed(t *testing.T) {
	b := Bottleneck{Arrival: 2000, Capacity: 1800, Reduction: 0.2}
	w := b.BackWaveSpeed()
	if w >= 0 {
		t.Fatalf("queue should propagate upstream (negative speed), got %v", w)
	}
}

func TestBottleneckRecoveryTime(t *testing.T) {
	b := Bottleneck{Arrival: 2000, Capacity: 1800, Reduction: 0.2}
	rt := b.RecoveryTime(100)
	if rt <= 0 {
		t.Fatalf("recovery time should be positive, got %v", rt)
	}
	if b.RecoveryTime(0) != 0 {
		t.Fatalf("zero queue -> 0")
	}
}

func TestSpillbackDistance(t *testing.T) {
	b := Bottleneck{Arrival: 2000, Capacity: 1800, Reduction: 0.2}
	d := b.SpillbackDistance(100, 120)
	if math.Abs(d-100.0/120.0) > 1e-9 {
		t.Fatalf("spillback distance wrong: %v", d)
	}
	if b.SpillbackDistance(100, 0) != 0 {
		t.Fatalf("zero jam density -> 0")
	}
}
