package snapshot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func (r Record) validate() error {
	if r.Magic != Magic {
		return fmt.Errorf("snapshot: bad magic %q", r.Magic)
	}
	if r.Version != CurrentVersion {
		return fmt.Errorf("snapshot: unsupported version %d", r.Version)
	}
	if r.Vf <= 0 {
		return fmt.Errorf("snapshot: vf must be > 0, got %g", r.Vf)
	}
	if r.Kj <= 0 {
		return fmt.Errorf("snapshot: kj must be > 0, got %g", r.Kj)
	}
	if r.Qm <= 0 || r.Km <= 0 || r.Vm <= 0 {
		return fmt.Errorf("snapshot: capacity fields must be positive")
	}
	if len(r.Samples) == 0 {
		return fmt.Errorf("snapshot: samples missing")
	}
	if r.Wave.K1 == r.Wave.K2 {
		return fmt.Errorf("snapshot: wave densities must differ")
	}
	m, err := r.Model()
	if err != nil {
		return fmt.Errorf("snapshot: %w", err)
	}
	for _, s := range r.Samples {
		if err := m.ValidateDensity(s.K); err != nil {
			return fmt.Errorf("snapshot: sample: %w", err)
		}
	}
	return nil
}

func WriteFile(path string, rec Record) error {
	if rec.Magic == "" {
		rec.Magic = Magic
	}
	if rec.Version == 0 {
		rec.Version = CurrentVersion
	}
	if err := rec.validate(); err != nil {
		return err
	}
	dir := filepath.Dir(path)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	raw, err := json.MarshalIndent(rec, "", "  ")
	if err != nil {
		return err
	}
	if len(raw) == 0 {
		return fmt.Errorf("snapshot: empty marshal")
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, raw, 0o644); err != nil {
		return err
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	return nil
}

func ReadFile(path string) (Record, error) {
	var rec Record
	raw, err := os.ReadFile(path)
	if err != nil {
		return rec, err
	}
	if len(raw) == 0 {
		return rec, fmt.Errorf("snapshot: empty file")
	}
	if !json.Valid(raw) {
		return rec, fmt.Errorf("snapshot: truncated or invalid JSON")
	}
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&rec); err != nil {
		return Record{}, fmt.Errorf("snapshot: %w", err)
	}
	if dec.More() {
		return Record{}, fmt.Errorf("snapshot: trailing content")
	}
	if err := rec.validate(); err != nil {
		return Record{}, err
	}
	return rec, nil
}
