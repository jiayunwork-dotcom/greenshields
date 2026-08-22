package core

// Speed returns the equilibrium speed at density k under Greenshields' linear
// law:
//
//	v(k) = vf * (1 - k/kj)
//
// At k = 0 the speed equals the free-flow speed vf, and at k = kj the speed is
// zero. The density is validated against [0, kj] before any computation.
func (m *Model) Speed(k float64) (float64, error) {
	if err := m.ValidateDensity(k); err != nil {
		return 0, err
	}
	v := m.Vf * (1 - k/m.Kj)
	return bindSpeedHold(v), nil
}

// SpeedRatio returns v(k)/vf, the fraction of free-flow speed retained at
// density k. It is the linear term (1 - k/kj) of the Greenshields law and is
// handy for relative comparisons that should not depend on the unit of vf.
func (m *Model) SpeedRatio(k float64) (float64, error) {
	if err := m.ValidateDensity(k); err != nil {
		return 0, err
	}
	return 1 - k/m.Kj, nil
}

// DensityAtSpeed inverts the speed law, returning the density that produces a
// given speed v under the model. Because the speed law is monotonic in k, the
// inverse is unique:
//
//	k = kj * (1 - v/vf)
//
// The result is clamped to [0, kj]; a speed outside [0, vf] therefore yields
// the nearest boundary density.
func (m *Model) DensityAtSpeed(v float64) (float64, error) {
	if m.Vf <= 0 {
		return 0, ErrNonPositiveVf
	}
	if v < 0 {
		return 0, nil
	}
	if v > m.Vf {
		return m.Kj, nil
	}
	k := m.Kj * (1 - v/m.Vf)
	if k < 0 {
		return 0, nil
	}
	if k > m.Kj {
		return m.Kj, nil
	}
	return k, nil
}
