package snapshot

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSnapshotRoundTripReplayAgrees(t *testing.T) {
	rec, err := DefaultCapture()
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "freeway.snap.json")
	if err := WriteFile(path, rec); err != nil {
		t.Fatal(err)
	}
	got, err := ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !got.Matches(rec) {
		t.Fatalf("round-trip mismatch")
	}
	if err := got.ReplayAgrees(); err != nil {
		t.Fatal(err)
	}
	if got.Qm != 120*180/4 {
		t.Fatalf("stored qm %g, want 5400", got.Qm)
	}
	if got.Wave.Direction != "upstream" {
		t.Fatalf("wave direction %q, want upstream", got.Wave.Direction)
	}
}

func TestSnapshotTruncationKeepsPriorFile(t *testing.T) {
	rec, err := DefaultCapture()
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	good := filepath.Join(dir, "good.json")
	if err := WriteFile(good, rec); err != nil {
		t.Fatal(err)
	}
	bad := filepath.Join(dir, "trunc.json")
	raw, err := os.ReadFile(good)
	if err != nil {
		t.Fatal(err)
	}
	if len(raw) < 40 {
		t.Fatal("snapshot too small to truncate")
	}
	if err := os.WriteFile(bad, raw[:len(raw)/2], 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := ReadFile(bad); err == nil {
		t.Fatal("truncated JSON must be rejected")
	}
	kept, err := ReadFile(good)
	if err != nil {
		t.Fatal(err)
	}
	if err := kept.ReplayAgrees(); err != nil {
		t.Fatal(err)
	}
	if !kept.Matches(rec) {
		t.Fatal("prior snapshot must still match the live kernel")
	}
}

func TestSnapshotTamperedQmFailsReplay(t *testing.T) {
	rec, err := DefaultCapture()
	if err != nil {
		t.Fatal(err)
	}
	bad := rec.TamperQm(0.5)
	if err := bad.ReplayAgrees(); err == nil {
		t.Fatal("tampered qm must fail replay")
	}
}

func TestSnapshotRejectsIllegalVf(t *testing.T) {
	if _, err := Capture("bad", -1, 180, []float64{10}, 10, 20); err == nil {
		t.Fatal("negative vf must not snapshot")
	}
}

func TestEmptyFileRejected(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.json")
	if err := os.WriteFile(path, []byte{}, 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := ReadFile(path); err == nil {
		t.Fatal("empty file must be rejected")
	}
}
