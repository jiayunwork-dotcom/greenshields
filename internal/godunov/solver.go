package godunov

import (
	"math"

	"greenshields/internal/core"
)

type Solver struct {
	Model *core.Model
}

func New(vf, kj float64) (*Solver, error) {
	m, err := core.New(vf, kj)
	if err != nil {
		return nil, err
	}
	return &Solver{Model: m}, nil
}

func FromModel(m *core.Model) (*Solver, error) {
	if m == nil {
		return nil, core.ErrNonPositiveVf
	}
	if err := m.Validate(); err != nil {
		return nil, err
	}
	return &Solver{Model: m}, nil
}

func (s *Solver) Q(k float64) (float64, error) {
	return s.Model.Flow(k)
}

func (s *Solver) V(k float64) (float64, error) {
	return s.Model.Speed(k)
}

func (s *Solver) Km() float64 {
	_, km := s.Model.Capacity()
	return km
}

func (s *Solver) Qm() float64 {
	qm, _ := s.Model.Capacity()
	return qm
}

func (s *Solver) Kj() float64 {
	return s.Model.Kj
}

func (s *Solver) Vf() float64 {
	return s.Model.Vf
}

func (s *Solver) Characteristic(k float64) (float64, error) {
	if err := s.Model.ValidateDensity(k); err != nil {
		return 0, err
	}
	return s.Model.Vf * (1 - 2*k/s.Model.Kj), nil
}

func (s *Solver) Sonic() bool {
	c, err := s.Characteristic(s.Km())
	if err != nil {
		return false
	}
	return math.Abs(c) <= core.AbsTiny
}

func (s *Solver) Branch(k float64) (string, error) {
	return s.Model.BranchOf(k)
}

func (s *Solver) Equilibrium(k float64) (float64, float64, error) {
	v, err := s.V(k)
	if err != nil {
		return 0, 0, err
	}
	q, err := s.Q(k)
	if err != nil {
		return 0, 0, err
	}
	return v, q, nil
}

func (s *Solver) ClampDensity(k float64) float64 {
	if k < 0 {
		return 0
	}
	if k > s.Model.Kj {
		return s.Model.Kj
	}
	return k
}

func relClose(a, b, tol float64) bool {
	scale := math.Max(1, math.Max(math.Abs(a), math.Abs(b)))
	return math.Abs(a-b) <= tol*scale
}
