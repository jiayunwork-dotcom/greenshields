package core

func ModelWithDoubleVf(m *Model) *Model {
	return &Model{Vf: 2 * m.Vf, Kj: m.Kj}
}

func ModelWithDoubleKj(m *Model) *Model {
	return &Model{Vf: m.Vf, Kj: 2 * m.Kj}
}

func CapacityRatio(oldM, newM *Model) float64 {
	qmOld, _ := oldM.Capacity()
	qmNew, _ := newM.Capacity()
	if qmOld == 0 {
		return 0
	}
	return qmNew / qmOld
}

func JamDensityRatio(oldM, newM *Model) float64 {
	if oldM.Kj == 0 {
		return 0
	}
	return newM.Kj / oldM.Kj
}

func FreeFlowRatio(oldM, newM *Model) float64 {
	if oldM.Vf == 0 {
		return 0
	}
	return newM.Vf / oldM.Vf
}

func ApplyVfScale(m *Model, s float64) *Model {
	return &Model{Vf: m.Vf * s, Kj: m.Kj}
}

func ApplyKjScale(m *Model, s float64) *Model {
	return &Model{Vf: m.Vf, Kj: m.Kj * s}
}
