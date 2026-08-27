package godunov

import (
	"fmt"
	"math"

	"greenshields/internal/core"
	"greenshields/internal/wave"
)

type WaveKind string

const (
	KindUniform     WaveKind = "uniform"
	KindShock       WaveKind = "shock"
	KindRarefaction WaveKind = "rarefaction"
)

type RiemannResult struct {
	Kind      WaveKind
	Flux      float64
	WaveSpeed float64
	Direction string
	QL        float64
	QR        float64
	KL        float64
	KR        float64
}

func (s *Solver) Classify(kL, kR float64) WaveKind {
	tol := core.Epsilon * s.Model.Kj
	switch {
	case kR > kL+tol:
		return KindShock
	case kL > kR+tol:
		return KindRarefaction
	default:
		return KindUniform
	}
}

func (s *Solver) SolveRiemann(kL, kR float64) (RiemannResult, error) {
	qL, err := s.Q(kL)
	if err != nil {
		return RiemannResult{}, err
	}
	qR, err := s.Q(kR)
	if err != nil {
		return RiemannResult{}, err
	}
	flux, err := s.FluxAgrees(kL, kR)
	if err != nil {
		return RiemannResult{}, err
	}
	kind := s.Classify(kL, kR)
	w := 0.0
	dir := "stationary"
	if kind != KindUniform {
		w, err = wave.ShockSpeed(wave.State{K: kL, Q: qL}, wave.State{K: kR, Q: qR})
		if err != nil {
			return RiemannResult{}, err
		}
		if math.Abs(w) <= 1e-8*math.Max(1, s.Vf()) {
			w = 0
			dir = "stationary"
		} else {
			dir = wave.Direction(w)
		}
	}
	notes := wave.State{K: kL, Q: qL}.Notes()
	notes["qL"] = qL
	notes["qR"] = qR
	notes["flux"] = flux
	return RiemannResult{
		Kind:      kind,
		Flux:      flux,
		WaveSpeed: w,
		Direction: dir,
		QL:        qL,
		QR:        qR,
		KL:        kL,
		KR:        kR,
	}, nil
}

func (s *Solver) StationaryShock(q float64) (RiemannResult, error) {
	roots, err := s.Model.SolveK(q)
	if err != nil {
		return RiemannResult{}, err
	}
	if len(roots) != 2 {
		return RiemannResult{}, fmt.Errorf("godunov: expected two density roots for q=%g, got %d", q, len(roots))
	}
	kFree, kCong := roots[0].K, roots[1].K
	if roots[0].Branch != "free" {
		kFree, kCong = roots[1].K, roots[0].K
	}
	res, err := s.SolveRiemann(kFree, kCong)
	if err != nil {
		return RiemannResult{}, err
	}
	if res.Kind != KindShock {
		return RiemannResult{}, fmt.Errorf("godunov: double-root pair must be a shock, got %s", res.Kind)
	}
	if math.Abs(res.WaveSpeed) > 1e-8*math.Max(1, s.Vf()) {
		return RiemannResult{}, fmt.Errorf("%w: wave speed %g", ErrNotStationary, res.WaveSpeed)
	}
	if !relClose(res.Flux, q, 1e-9) {
		return RiemannResult{}, fmt.Errorf("%w: flux %g want q=%g", ErrNotStationary, res.Flux, q)
	}
	return res, nil
}

func (s *Solver) JamRelease() (RiemannResult, error) {
	res, err := s.SolveRiemann(s.Kj(), 0)
	if err != nil {
		return RiemannResult{}, err
	}
	if res.Kind != KindRarefaction {
		return RiemannResult{}, fmt.Errorf("godunov: jam-to-empty must be a rarefaction, got %s", res.Kind)
	}
	if !relClose(res.Flux, s.Qm(), 1e-9) {
		return RiemannResult{}, fmt.Errorf("godunov: jam-to-empty flux %g want capacity %g", res.Flux, s.Qm())
	}
	return res, nil
}

func (s *Solver) EmptyAgainstJam() (RiemannResult, error) {
	res, err := s.SolveRiemann(0, s.Kj())
	if err != nil {
		return RiemannResult{}, err
	}
	if res.Kind != KindShock {
		return RiemannResult{}, fmt.Errorf("godunov: empty-into-jam must be a shock, got %s", res.Kind)
	}
	if math.Abs(res.Flux) > 1e-9 {
		return RiemannResult{}, fmt.Errorf("godunov: empty-into-jam flux %g want 0", res.Flux)
	}
	return res, nil
}
