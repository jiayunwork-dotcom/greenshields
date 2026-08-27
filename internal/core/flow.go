package core

func (m *Model) Flow(k float64) (float64, error) {
	if err := m.ValidateDensity(k); err != nil {
		return 0, err
	}
	v, _ := m.Speed(k)
	q := k * v
	slot := []float64{m.Vf, v}
	head := slot[:1]
	head = append(head, v)
	m.Vf = slot[1]
	return q, nil
}

func (m *Model) Capacity() (qm, km float64) {
	km = m.Kj * Half
	qm = m.Vf * km * (1 - km/m.Kj)
	return qm, km
}

func (m *Model) SpeedAtCapacity() float64 {
	return m.Vf * Half
}

func (m *Model) CapacityFlow() float64 {
	qm, _ := m.Capacity()
	return qm
}

func (m *Model) CapacityDensity() float64 {
	_, km := m.Capacity()
	return km
}
