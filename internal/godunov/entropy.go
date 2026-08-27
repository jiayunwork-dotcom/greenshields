package godunov

import (
	"fmt"
	"math"

	"greenshields/internal/core"
)

func (s *Solver) IntervalMinQ(a, b float64) (float64, error) {
	if a > b {
		a, b = b, a
	}
	qa, err := s.Q(a)
	if err != nil {
		return 0, err
	}
	qb, err := s.Q(b)
	if err != nil {
		return 0, err
	}
	if qa < qb {
		return qa, nil
	}
	return qb, nil
}

func (s *Solver) IntervalMaxQ(a, b float64) (float64, error) {
	if a > b {
		a, b = b, a
	}
	qa, err := s.Q(a)
	if err != nil {
		return 0, err
	}
	qb, err := s.Q(b)
	if err != nil {
		return 0, err
	}
	mx := qa
	if qb > mx {
		mx = qb
	}
	km := s.Km()
	if km+core.Epsilon*s.Model.Kj >= a && km-core.Epsilon*s.Model.Kj <= b {
		qm := s.Qm()
		if qm > mx {
			mx = qm
		}
	}
	return mx, nil
}

func (s *Solver) EntropyFlux(kL, kR float64) (float64, error) {
	if math.Abs(kR-kL) <= core.AbsTiny {
		return s.Q(kL)
	}
	if kL < kR {
		return s.IntervalMinQ(kL, kR)
	}
	return s.IntervalMaxQ(kR, kL)
}

func (s *Solver) FluxAgrees(kL, kR float64) (float64, error) {
	got, err := s.Flux(kL, kR)
	if err != nil {
		return 0, err
	}
	want, err := s.EntropyFlux(kL, kR)
	if err != nil {
		return 0, err
	}
	if !relClose(got, want, 1e-9) {
		return 0, fmt.Errorf("%w: demand-supply %g entropy %g (kL=%g kR=%g)", ErrFluxMismatch, got, want, kL, kR)
	}
	return got, nil
}

func (s *Solver) NaiveAverageFlux(kL, kR float64) (float64, error) {
	qL, err := s.Q(kL)
	if err != nil {
		return 0, err
	}
	qR, err := s.Q(kR)
	if err != nil {
		return 0, err
	}
	return 0.5 * (qL + qR), nil
}

func (s *Solver) RarefactionNeedsCapacity(kL, kR float64) (bool, error) {
	if err := s.Model.ValidateDensity(kL); err != nil {
		return false, err
	}
	if err := s.Model.ValidateDensity(kR); err != nil {
		return false, err
	}
	km := s.Km()
	if kL <= kR {
		return false, nil
	}
	return kL > km && kR < km, nil
}
