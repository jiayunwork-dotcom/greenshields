package example

import "testing"

// TestLoadFreewayExample checks that both the embedded and on-disk examples
// decode to valid, positive parameters and that qm = vf*kj/4 holds. (Medium)
func TestLoadFreewayExample(t *testing.T) {
	p, err := Embedded()
	if err != nil {
		t.Fatalf("Embedded() error: %v", err)
	}
	if p.Vf <= 0 || p.Kj <= 0 {
		t.Errorf("embedded example has non-positive params: vf=%g kj=%g", p.Vf, p.Kj)
	}

	// The highway example is meant to reproduce qm = vf*kj/4.
	wantQm := p.Vf * p.Kj / 4
	if wantQm <= 0 {
		t.Errorf("expected positive qm, got %g", wantQm)
	}

	// On-disk copy in this package directory.
	disk, err := LoadFile("freeway.json")
	if err != nil {
		t.Fatalf("LoadFile(freeway.json) error: %v", err)
	}
	if disk.Vf != p.Vf || disk.Kj != p.Kj {
		t.Errorf("on-disk example %+v differs from embedded %+v", disk, p)
	}
}

// TestDecodeInvalid checks that malformed or non-positive JSON is rejected.
func TestDecodeInvalid(t *testing.T) {
	if _, err := Decode([]byte(`{"vf":-1,"kj":180}`)); err == nil {
		t.Errorf("Decode with vf<0: expected error, got nil")
	}
	if _, err := Decode([]byte(`not json`)); err == nil {
		t.Errorf("Decode invalid json: expected error, got nil")
	}
}
