package wave

import (
	"errors"
)

type FundamentalDiagram struct {
	Kj   float64
	Qmax float64
}

func (f FundamentalDiagram) CharacteristicSpeed(k float64) (float64, error) {
	if f.Kj <= 0 {
		return 0, errWaveBadArg
	}
	return f.Qmax * (4.0 / f.Kj) * (1.0 - 2.0*k/f.Kj), nil
}

func (f FundamentalDiagram) Conserved(up, down State) (float64, error) {
	if up.Q < 0 || down.Q < 0 {
		return 0, errWaveBadArg
	}
	return up.Q - down.Q, nil
}

func (f FundamentalDiagram) FIFOok(up, down State) bool {
	if down.K > up.K && down.Q >= up.Q {
		return false
	}
	return true
}

func (f FundamentalDiagram) DensityJump(up, down State) float64 {
	return down.K - up.K
}

func (f FundamentalDiagram) EquivalentSpeed(q, k float64) (float64, error) {
	if k <= 0 {
		return 0, errWaveBadArg
	}
	return q / k, nil
}

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

func (f FundamentalDiagram) CapacityDrop(dropFrac float64) float64 {
	if dropFrac < 0 || dropFrac > 1 {
		return f.Qmax
	}
	return f.Qmax * (1 - dropFrac)
}

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

func (f FundamentalDiagram) MaxDensitySpeed() float64 {
	return f.Kj / 2.0
}

func (f FundamentalDiagram) WaveStability() bool {
	c, err := f.CharacteristicSpeed(f.Kj / 2.0)
	if err != nil {
		return false
	}
	return c == 0
}

var errWaveBadArg = errors.New("invalid wave argument")
