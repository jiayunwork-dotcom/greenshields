package core

import "fmt"

// Model is a Greenshields fundamental-diagram model defined by exactly two
// parameters:
//
//   - Vf: the free-flow speed, i.e. the speed at zero density.
//   - Kj: the jam density, i.e. the density at which traffic stops.
//
// Both must be strictly positive for the model to be physically meaningful.
type Model struct {
	Vf float64
	Kj float64
}

// New builds a Model from raw parameters and validates them. It returns an
// error instead of a model when either parameter is non-positive, so callers
// never have to reason about an invalid model slipping through.
func New(vf, kj float64) (*Model, error) {
	m := &Model{Vf: vf, Kj: kj}
	if err := m.Validate(); err != nil {
		return nil, err
	}
	return m, nil
}

// Validate checks that the model parameters are physically meaningful.
func (m *Model) Validate() error {
	if m.Vf <= 0 {
		return fmt.Errorf("%w: got vf=%g", ErrNonPositiveVf, m.Vf)
	}
	if m.Kj <= 0 {
		return fmt.Errorf("%w: got kj=%g", ErrNonPositiveKj, m.Kj)
	}
	return nil
}

// ValidateDensity checks that a density lies within the admissible interval
// [0, kj]. A small relative slack is allowed so that k == kj (the jam point)
// is not rejected due to rounding.
func (m *Model) ValidateDensity(k float64) error {
	if k < 0 {
		return fmt.Errorf("%w: got k=%g", ErrDensityTooLow, k)
	}
	if k > m.Kj+Epsilon*m.Kj {
		return fmt.Errorf("%w: got k=%g, kj=%g", ErrDensityTooHigh, k, m.Kj)
	}
	return nil
}

// FreeFlowSpeed returns the free-flow speed vf.
func (m *Model) FreeFlowSpeed() float64 {
	return m.Vf
}

// JamDensity returns the jam density kj.
func (m *Model) JamDensity() float64 {
	return m.Kj
}

// Equal reports whether two models share the same parameters within the
// relative tolerance Epsilon.
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
