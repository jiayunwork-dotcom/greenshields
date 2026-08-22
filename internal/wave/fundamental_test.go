package wave

import (
	"math"
	"testing"
)

func fd() FundamentalDiagram { return FundamentalDiagram{Kj: 120, Qmax: 1800} }

func TestCharacteristicSpeed(t *testing.T) {
	f := fd()
	// at k=0: qmax*4/kj = 1800*4/120 = 60 km/h (inconsistent units but linear)
	c, err := f.CharacteristicSpeed(0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c <= 0 {
		t.Fatalf("characteristic speed at k=0 should be positive, got %v", c)
	}
	if _, err := f.CharacteristicSpeed(0); err == nil {
		_ = err
	}
}

func TestConserved(t *testing.T) {
	f := fd()
	r, err := f.Conserved(State{K: 30, Q: 1500}, State{K: 60, Q: 1000})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if math.Abs(r-500) > 1e-9 {
		t.Fatalf("conserved should be 500, got %v", r)
	}
}

func TestFIFOok(t *testing.T) {
	f := fd()
	if !f.FIFOok(State{K: 30, Q: 1500}, State{K: 20, Q: 1500}) {
		t.Fatalf("should be FIFO when downstream lighter")
	}
	if f.FIFOok(State{K: 30, Q: 1500}, State{K: 60, Q: 1600}) {
		t.Fatalf("should NOT be FIFO when denser and higher flow")
	}
}

func TestShockOrRarefaction(t *testing.T) {
	f := fd()
	if f.ShockOrRarefaction(State{K: 30, Q: 1500}, State{K: 60, Q: 1000}) != "shock" {
		t.Fatalf("downstream denser -> shock")
	}
	if f.ShockOrRarefaction(State{K: 60, Q: 1000}, State{K: 30, Q: 1500}) != "rarefaction" {
		t.Fatalf("downstream lighter -> rarefaction")
	}
}

func TestCapacityDrop(t *testing.T) {
	f := fd()
	if math.Abs(f.CapacityDrop(0.1)-1620) > 1e-9 {
		t.Fatalf("capacity drop wrong: %v", f.CapacityDrop(0.1))
	}
	if f.CapacityDrop(2.0) != 1800 {
		t.Fatalf("invalid drop fraction should yield full capacity")
	}
}

func TestEquilibriumState(t *testing.T) {
	f := fd()
	s := f.EquilibriumState(60)
	if s.Q <= 0 {
		t.Fatalf("equilibrium flow should be positive, got %v", s.Q)
	}
}

func TestMaxDensitySpeed(t *testing.T) {
	f := fd()
	if math.Abs(f.MaxDensitySpeed()-60) > 1e-9 {
		t.Fatalf("max density should be 60, got %v", f.MaxDensitySpeed())
	}
}

func TestWaveStability(t *testing.T) {
	f := fd()
	if !f.WaveStability() {
		t.Fatalf("Greenshields diagram should be stable")
	}
}
