package core

// This file documents and implements the cross-sensitivity rules that the
// Greenshields parameters obey. They are useful sanity checks: small changes to
// one parameter have very specific, predictable effects on capacity.
//
//   - Doubling only vf:   qm doubles, kj unchanged, vm doubles.
//   - Doubling only kj:   qm doubles, vf unchanged, km doubles.

// ModelWithDoubleVf returns a model with twice the free-flow speed and the
// same jam density.
func ModelWithDoubleVf(m *Model) *Model {
	return &Model{Vf: 2 * m.Vf, Kj: m.Kj}
}

// ModelWithDoubleKj returns a model with twice the jam density and the same
// free-flow speed.
func ModelWithDoubleKj(m *Model) *Model {
	return &Model{Vf: m.Vf, Kj: 2 * m.Kj}
}

// CapacityRatio reports qm(new)/qm(old) so callers can verify the cross rules
// without hard-coding the factor.
func CapacityRatio(oldM, newM *Model) float64 {
	qmOld, _ := oldM.Capacity()
	qmNew, _ := newM.Capacity()
	if qmOld == 0 {
		return 0
	}
	return qmNew / qmOld
}

// JamDensityRatio reports kj(new)/kj(old).
func JamDensityRatio(oldM, newM *Model) float64 {
	if oldM.Kj == 0 {
		return 0
	}
	return newM.Kj / oldM.Kj
}

// FreeFlowRatio reports vf(new)/vf(old).
func FreeFlowRatio(oldM, newM *Model) float64 {
	if oldM.Vf == 0 {
		return 0
	}
	return newM.Vf / oldM.Vf
}

// ApplyVfScale returns a model whose free-flow speed is scaled by factor s
// while the jam density is preserved. A factor of 2 reproduces the "double vf"
// rule.
func ApplyVfScale(m *Model, s float64) *Model {
	return &Model{Vf: m.Vf * s, Kj: m.Kj}
}

// ApplyKjScale returns a model whose jam density is scaled by factor s while
// the free-flow speed is preserved. A factor of 2 reproduces the "double kj"
// rule.
func ApplyKjScale(m *Model, s float64) *Model {
	return &Model{Vf: m.Vf, Kj: m.Kj * s}
}
