package snapshot

import (
	"fmt"

	"greenshields/internal/core"
	"greenshields/internal/wave"
)

const (
	Magic          = "GSFD"
	CurrentVersion = 1
	Tol            = 1e-9
)

type Sample struct {
	K    float64 `json:"k"`
	V    float64 `json:"v"`
	Q    float64 `json:"q"`
	Side string  `json:"side"`
}

type WavePair struct {
	K1        float64 `json:"k1"`
	Q1        float64 `json:"q1"`
	K2        float64 `json:"k2"`
	Q2        float64 `json:"q2"`
	W         float64 `json:"w"`
	Direction string  `json:"direction"`
}

type Record struct {
	Magic   string   `json:"magic"`
	Version int      `json:"version"`
	Name    string   `json:"name"`
	Vf      float64  `json:"vf"`
	Kj      float64  `json:"kj"`
	Km      float64  `json:"km"`
	Qm      float64  `json:"qm"`
	Vm      float64  `json:"vm"`
	Samples []Sample `json:"samples"`
	Wave    WavePair `json:"wave"`
}

func Capture(name string, vf, kj float64, densities []float64, k1, k2 float64) (Record, error) {
	m, err := core.New(vf, kj)
	if err != nil {
		return Record{}, fmt.Errorf("snapshot: %w", err)
	}
	qm, km := m.Capacity()
	rec := Record{
		Magic:   Magic,
		Version: CurrentVersion,
		Name:    name,
		Vf:      vf,
		Kj:      kj,
		Km:      km,
		Qm:      qm,
		Vm:      m.SpeedAtCapacity(),
	}
	for _, k := range densities {
		v, err := m.Speed(k)
		if err != nil {
			return Record{}, fmt.Errorf("snapshot: sample k=%g: %w", k, err)
		}
		q, err := m.Flow(k)
		if err != nil {
			return Record{}, fmt.Errorf("snapshot: sample k=%g: %w", k, err)
		}
		side, err := m.BranchOf(k)
		if err != nil {
			return Record{}, err
		}
		rec.Samples = append(rec.Samples, Sample{K: k, V: v, Q: q, Side: side})
	}
	if err := m.ValidateDensity(k1); err != nil {
		return Record{}, fmt.Errorf("snapshot: wave k1: %w", err)
	}
	if err := m.ValidateDensity(k2); err != nil {
		return Record{}, fmt.Errorf("snapshot: wave k2: %w", err)
	}
	q1, err := m.Flow(k1)
	if err != nil {
		return Record{}, err
	}
	q2, err := m.Flow(k2)
	if err != nil {
		return Record{}, err
	}
	w, err := wave.ShockSpeed(wave.State{K: k1, Q: q1}, wave.State{K: k2, Q: q2})
	if err != nil {
		return Record{}, fmt.Errorf("snapshot: %w", err)
	}
	rec.Wave = WavePair{
		K1:        k1,
		Q1:        q1,
		K2:        k2,
		Q2:        q2,
		W:         w,
		Direction: wave.Direction(w),
	}
	if err := rec.validate(); err != nil {
		return Record{}, err
	}
	return rec, nil
}

func DefaultCapture() (Record, error) {
	return Capture("urban-freeway", 120, 180, []float64{0, 45, 90, 162, 180}, 45, 162)
}

func (r Record) Model() (*core.Model, error) {
	return core.New(r.Vf, r.Kj)
}
