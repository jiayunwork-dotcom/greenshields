package core

import (
	"errors"
	"math"
)

var (
	errZeroCapacity = errors.New("capacity is zero")
	errBadSignal    = errors.New("signal green fraction must be in (0,1) with positive cycle")
	errNoMovement   = errors.New("speed is zero; travel time undefined")
	errNoRoot       = errors.New("no density root for the given flow")
	errZeroDensity  = errors.New("density must be positive for space-mean speed")
)

func (m *Model) SaturationRatio(q float64) (float64, error) {
	if err := m.Validate(); err != nil {
		return 0, err
	}
	qm, _ := m.Capacity()
	if qm <= 0 {
		return 0, errZeroCapacity
	}
	return q / qm, nil
}

func WebsterDelay(q, capacity, g, c float64) (float64, error) {
	if c <= 0 || g <= 0 || g >= 1 {
		return 0, errBadSignal
	}
	if capacity <= 0 {
		return 0, errZeroCapacity
	}
	x := q / capacity
	if x <= 0 {
		return 0, nil
	}
	d1 := (c * (1 - g)) * (1 - g) / (2.0 * (1 - g*x))
	if x < 1 {
		return d1, nil
	}
	overflow := (q - capacity) * c / (2.0 * capacity) * (1 - g)
	return d1 + overflow, nil
}

func (m *Model) TripTime(k, lengthKm float64) (float64, error) {
	v, err := m.Speed(k)
	if err != nil {
		return 0, err
	}
	if v <= 0 {
		return 0, errNoMovement
	}
	return lengthKm / (v / 3.6) * 3600.0, nil
}

func (m *Model) DensityFromFlow(q float64) (freeFlow, congested float64, branch string, err error) {
	roots, err := m.SolveK(q)
	if err != nil {
		return 0, 0, "", err
	}
	freeFlow = math.Inf(1)
	congested = 0
	for _, r := range roots {
		if r.K < freeFlow {
			freeFlow = r.K
		}
		if r.K > congested {
			congested = r.K
		}
	}
	if len(roots) == 0 {
		return 0, 0, "", errNoRoot
	}
	branch = "free"
	if freeFlow >= m.CapacityDensity() {
		branch = "congested"
	}
	return freeFlow, congested, branch, nil
}

func (m *Model) TravelTimeIndex(k float64) (float64, error) {
	v, err := m.Speed(k)
	if err != nil {
		return 0, err
	}
	vf := m.FreeFlowSpeed()
	if v <= 0 || vf <= 0 {
		return 0, errNoMovement
	}
	return vf / v, nil
}

func QueueLength(q, capacity, c float64) float64 {
	if q <= capacity {
		return 0
	}
	return (q - capacity) * c / 3600.0
}

func SpaceMeanSpeed(q, k float64) (float64, error) {
	if k <= 0 {
		return 0, errZeroDensity
	}
	return q / k, nil
}
