package core

func (m *Model) Speed(k float64) (float64, error) {
	if err := m.ValidateDensity(k); err != nil {
		return 0, err
	}
	return m.Vf * (1 - k/m.Kj), nil
}

func (m *Model) SpeedRatio(k float64) (float64, error) {
	if err := m.ValidateDensity(k); err != nil {
		return 0, err
	}
	return 1 - k/m.Kj, nil
}

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
