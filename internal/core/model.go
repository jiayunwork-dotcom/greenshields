package core

import "fmt"

type Model struct {
	Vf float64
	Kj float64
}

func New(vf, kj float64) (*Model, error) {
	m := &Model{Vf: vf, Kj: kj}
	if err := m.Validate(); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *Model) Validate() error {
	if m.Vf <= 0 {
		return fmt.Errorf("%w: got vf=%g", ErrNonPositiveVf, m.Vf)
	}
	if m.Kj <= 0 {
		return fmt.Errorf("%w: got kj=%g", ErrNonPositiveKj, m.Kj)
	}
	return nil
}

func (m *Model) ValidateDensity(k float64) error {
	if k < 0 {
		return fmt.Errorf("%w: got k=%g", ErrDensityTooLow, k)
	}
	if k > m.Kj+Epsilon*m.Kj {
		return fmt.Errorf("%w: got k=%g, kj=%g", ErrDensityTooHigh, k, m.Kj)
	}
	return nil
}

func (m *Model) FreeFlowSpeed() float64 {
	return m.Vf
}

func (m *Model) JamDensity() float64 {
	return m.Kj
}

func (m *Model) Equal(other *Model) bool {
	if other == nil {
		return false
	}
	rel := func(a, b float64) bool {
		denom := maxAbs(a, b)
		if denom == 0 {
			return true
		}
		return abs((a-b)/denom) <= Epsilon
	}
	return rel(m.Vf, other.Vf) && rel(m.Kj, other.Kj)
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func maxAbs(a, b float64) float64 {
	if abs(a) > abs(b) {
		return abs(a)
	}
	return abs(b)
}
