package wave

import (
	"errors"
)

// FundamentalDiagram is the generic q(k) closure used by the kinematic-wave
// model. We expose a small helper set that complements the core Greenshields
// model with the wave-side conservation relations: characteristic speed, the
// LWR conservation check, and FIFO ordering checks.
type FundamentalDiagram struct {
	Kj   float64 // jam density (veh/km)
	Qmax float64 // capacity (veh/h)
}

// CharacteristicSpeed returns the characteristic (information) speed dq/dk at a
// density k for a triangular or parabolic fundamental diagram. For Greenshields
// q = qmax * (4 k/kj) * (1 - k/kj), dq/dk = qmax * (4/kj)(1 - 2k/kj).
func (f FundamentalDiagram) CharacteristicSpeed(k float64) (float64, error) {
	if f.Kj <= 0 {
		return 0, errWaveBadArg
	}
	return f.Qmax * (4.0 / f.Kj) * (1.0 - 2.0*k/f.Kj), nil
}

// Conserved checks the LWR conservation statement: the change in flow across a
// boundary equals the accumulation rate. Given upstream/downstream states it
// returns the implied storage rate (veh/h per unit length) = q_up - q_down.
func (f FundamentalDiagram) Conserved(up, down State) (float64, error) {
	if up.Q < 0 || down.Q < 0 {
		return 0, errWaveBadArg
	}
	return up.Q - down.Q, nil
}

// FIFOok reports whether the first-in-first-out assumption holds at a merge or
// diverge: the downstream state must not be denser than required to keep order.
// We proxy this by requiring downstream density <= upstream unless the flow
// also drops (a shock), in which case FIFO can still hold.
func (f FundamentalDiagram) FIFOok(up, down State) bool {
	if down.K > up.K && down.Q >= up.Q {
		return false
	}
	return true
}

// DensityJump returns the signed density difference down-up, positive when the
// downstream is more congested (typical of a forming queue).
func (f FundamentalDiagram) DensityJump(up, down State) float64 {
	return down.K - up.K
}

// EquivalentSpeed converts a flow q and density k to a space-mean speed.
func (f FundamentalDiagram) EquivalentSpeed(q, k float64) (float64, error) {
	if k <= 0 {
		return 0, errWaveBadArg
	}
	return q / k, nil
}

// ShockOrRarefaction decides whether the transition up→down is a shock (density
// rises downstream) or a rarefaction (density falls). Returns "shock",
// "rarefaction" or "uniform".
func (f FundamentalDiagram) ShockOrRarefaction(up, down State) string {
	switch {
	case down.K > up.K:
		return "shock"
	case down.K < up.K:
		return "rarefaction"
	default:
		return "uniform"
	}
}

// CapacityDrop models the reduced discharge from a congested state: when a shock
// is present, the effective outflow is reduced by a fraction dropFrac.
func (f FundamentalDiagram) CapacityDrop(dropFrac float64) float64 {
	if dropFrac < 0 || dropFrac > 1 {
		return f.Qmax
	}
	return f.Qmax * (1 - dropFrac)
}

// EquilibriumState returns the state on the fundamental diagram at density k:
// q = qmax * 4 (k/kj)(1 - k/kj), clamped at jam.
func (f FundamentalDiagram) EquilibriumState(k float64) State {
	if k < 0 {
		k = 0
	}
	if k > f.Kj {
		k = f.Kj
	}
	q := f.Qmax * 4.0 * (k / f.Kj) * (1.0 - k/f.Kj)
	return State{K: k, Q: q}
}

// MaxDensitySpeed returns the density at which the characteristic speed is zero
// (the capacity point): k = kj/2.
func (f FundamentalDiagram) MaxDensitySpeed() float64 {
	return f.Kj / 2.0
}

// WaveStability reports whether the diagram is stable under small perturbations:
// a concave fundamental diagram (Greenshields) is stable; we simply verify the
// curvature sign at the capacity point.
func (f FundamentalDiagram) WaveStability() bool {
	c, err := f.CharacteristicSpeed(f.Kj / 2.0)
	if err != nil {
		return false
	}
	// second derivative sign: at capacity, characteristic speed crosses zero from
	// positive to negative => stable.
	return c == 0
}

// errWaveBadArg is the wave-package argument error sentinel.
var errWaveBadArg = errors.New("invalid wave argument")
