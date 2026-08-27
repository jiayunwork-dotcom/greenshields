package wave

import (
	"math"
	"testing"
)

func TestMergeDemand(t *testing.T) {
	m := Merge{MainFlow: 1500, RampFlow: 400, Capacity: 1800}
	if math.Abs(m.TotalDemand()-1900) > 1e-9 {
		t.Fatalf("demand should be 1900, got %v", m.TotalDemand())
	}
	if math.Abs(m.AcceptedFlow()-1800) > 1e-9 {
		t.Fatalf("accepted should be capped at 1800, got %v", m.AcceptedFlow())
	}
	if math.Abs(m.QueueRate()-100) > 1e-9 {
		t.Fatalf("queue rate should be 100, got %v", m.QueueRate())
	}
}

func TestMergeNoQueue(t *testing.T) {
	m := Merge{MainFlow: 1000, RampFlow: 200, Capacity: 1800}
	if m.QueueRate() != 0 {
		t.Fatalf("no queue under capacity")
	}
	if m.QueueAt(0, 3600) != 0 {
		t.Fatalf("no queue should accumulate")
	}
}

func TestMergeWaveSpeed(t *testing.T) {
	m := Merge{MainFlow: 1500, RampFlow: 400, Capacity: 1800}
	w := m.MergeWaveSpeed()
	if w >= 0 {
		t.Fatalf("merge wave should be backward (negative), got %v", w)
	}
}

func TestRampMeteringDelay(t *testing.T) {
	m := Merge{MainFlow: 1500, RampFlow: 400, Capacity: 1800}
	d := m.RampMeteringDelay(10, 3600)
	if d <= 0 {
		t.Fatalf("delay should be positive, got %v", d)
	}
}

func TestFractionMaintained(t *testing.T) {
	m := Merge{MainFlow: 1500, RampFlow: 400, Capacity: 1800}
	frac := m.FractionMaintained()
	if math.Abs(frac-1800.0/1900.0) > 1e-9 {
		t.Fatalf("fraction maintained wrong: %v", frac)
	}
}

func TestThroughputGap(t *testing.T) {
	m := Merge{MainFlow: 1500, RampFlow: 400, Capacity: 1800}
	g := m.ThroughputGap(3600)
	if math.Abs(g-100.0) > 1e-9 {
		t.Fatalf("throughput gap should be 100 veh, got %v", g)
	}
}
