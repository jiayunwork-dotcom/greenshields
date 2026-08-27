package godunov

import (
	"fmt"
	"math"

	"greenshields/internal/core"
	"greenshields/internal/wave"
)

func (s *Solver) MaxWaveSpeed(kL, kR float64) (float64, error) {
	cL, err := s.Characteristic(kL)
	if err != nil {
		return 0, err
	}
	cR, err := s.Characteristic(kR)
	if err != nil {
		return 0, err
	}
	a := math.Abs(cL)
	if math.Abs(cR) > a {
		a = math.Abs(cR)
	}
	if math.Abs(kR-kL) > core.AbsTiny {
		qL, err := s.Q(kL)
		if err != nil {
			return 0, err
		}
		qR, err := s.Q(kR)
		if err != nil {
			return 0, err
		}
		w, err := wave.ShockSpeed(wave.State{K: kL, Q: qL}, wave.State{K: kR, Q: qR})
		if err == nil && math.Abs(w) > a {
			a = math.Abs(w)
		}
	}
	return a, nil
}

func (s *Solver) TimeStep(dx, kL, kR, cfl float64) (float64, error) {
	if dx <= 0 || cfl <= 0 || cfl > 1 {
		return 0, ErrBadCFL
	}
	a, err := s.MaxWaveSpeed(kL, kR)
	if err != nil {
		return 0, err
	}
	if a < core.AbsTiny {
		return 0, ErrZeroWave
	}
	return cfl * dx / a, nil
}

func (s *Solver) StepPair(kL, kR, dx, dt float64) (float64, float64, error) {
	if dx <= 0 || dt <= 0 {
		return 0, 0, ErrBadCFL
	}
	a, err := s.MaxWaveSpeed(kL, kR)
	if err != nil {
		return 0, 0, err
	}
	if a > 0 && dt > dx/a+core.Epsilon {
		return 0, 0, fmt.Errorf("%w: dt=%g exceeds dx/a=%g", ErrBadCFL, dt, dx/a)
	}
	fStar, err := s.Flux(kL, kR)
	if err != nil {
		return 0, 0, err
	}
	qL, err := s.Q(kL)
	if err != nil {
		return 0, 0, err
	}
	qR, err := s.Q(kR)
	if err != nil {
		return 0, 0, err
	}
	nL := kL - (dt/dx)*(fStar-qL)
	nR := kR - (dt/dx)*(qR-fStar)
	nL = s.ClampDensity(nL)
	nR = s.ClampDensity(nR)
	if err := MassConserved(kL, kR, nL, nR); err != nil {
		return 0, 0, err
	}
	return nL, nR, nil
}

func MassConserved(kL, kR, nL, nR float64) error {
	old := kL + kR
	now := nL + nR
	if !relClose(old, now, 1e-9) {
		return fmt.Errorf("%w: before %g after %g", ErrMassLeak, old, now)
	}
	return nil
}

func (s *Solver) HoldStationary(q, dx, dt float64) error {
	res, err := s.StationaryShock(q)
	if err != nil {
		return err
	}
	nL, nR, err := s.StepPair(res.KL, res.KR, dx, dt)
	if err != nil {
		return err
	}
	if !relClose(nL, res.KL, 1e-8) || !relClose(nR, res.KR, 1e-8) {
		return fmt.Errorf("%w: densities moved from (%g,%g) to (%g,%g)", ErrNotStationary, res.KL, res.KR, nL, nR)
	}
	return nil
}
