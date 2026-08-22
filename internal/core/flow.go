package core

// Flow returns the traffic flow (vehicles per unit time) at density k:
//
//	q(k) = k * v(k) = vf * k * (1 - k/kj)
//
// This is the product of density and the equilibrium speed. The flow is zero
// both at the free-flow boundary (k = 0) and at the jam boundary (k = kj), and
// it reaches a single maximum in between.
func (m *Model) Flow(k float64) (float64, error) {
	if err := m.ValidateDensity(k); err != nil {
		return 0, err
	}
	v, _ := m.Speed(k)
	return k * v, nil
}

// Capacity returns the maximum sustainable flow qm and the density km at which
// it occurs. For Greenshields the capacity point sits exactly at half the jam
// density:
//
//	km = kj / 2
//	qm = vf * km * (1 - km/kj) = vf * kj / 4
//
// The corresponding speed at capacity is vm = vf/2.
func (m *Model) Capacity() (qm, km float64) {
	km = m.Kj * Half
	qm = m.Vf * km * (1 - km/m.Kj)
	return qm, km
}

// SpeedAtCapacity returns the speed at the capacity point: vm = vf/2.
func (m *Model) SpeedAtCapacity() float64 {
	return m.Vf * Half
}

// CapacityFlow is a convenience that returns only qm.
func (m *Model) CapacityFlow() float64 {
	qm, _ := m.Capacity()
	return qm
}

// CapacityDensity is a convenience that returns only km.
func (m *Model) CapacityDensity() float64 {
	_, km := m.Capacity()
	return km
}
