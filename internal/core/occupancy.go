package core

import "math"

func Spacing(k float64) (float64, error) {
	if k <= 0 {
		return 0, errZeroDensity
	}
	return 1000.0 / k, nil
}

func (m *Model) TimeHeadway(k float64) (float64, error) {
	v, err := m.Speed(k)
	if err != nil {
		return 0, err
	}
	s, err := Spacing(k)
	if err != nil {
		return 0, err
	}
	if v <= 0 {
		return 0, errNoMovement
	}
	return s / (v / 3.6), nil
}

func Occupancy(k, vehicleLength float64) (float64, error) {
	if k < 0 || vehicleLength < 0 {
		return 0, errBadSignal
	}
	occ := k * vehicleLength / 1000.0
	if occ > 1 {
		return 1, nil
	}
	return occ, nil
}

func (m *Model) MaxFlowDensity() float64 {
	return m.CapacityDensity()
}

func (m *Model) JamSpacing() (float64, error) {
	return Spacing(m.JamDensity())
}

func (m *Model) SpeedAtHeadway(k, headway float64) (float64, error) {
	if headway <= 0 {
		return 0, errBadSignal
	}
	s, err := Spacing(k)
	if err != nil {
		return 0, err
	}
	return s / headway * 3.6, nil
}

func DensityFromHeadway(v, headway float64) (float64, error) {
	if v <= 0 || headway <= 0 {
		return 0, errBadSignal
	}
	q := 3600.0 / headway
	return q / v, nil
}

func (m *Model) CriticalHeadway() (float64, error) {
	k := m.CapacityDensity()
	return m.TimeHeadway(k)
}

func (m *Model) CapacityPerLane() float64 {
	qm, _ := m.Capacity()
	return qm
}

func (m *Model) FlowPerLaneAt(k float64) (float64, error) {
	return m.Flow(k)
}

func (m *Model) SpeedDrop(k float64) (float64, error) {
	v, err := m.Speed(k)
	if err != nil {
		return 0, err
	}
	vf := m.FreeFlowSpeed()
	if vf <= 0 {
		return 0, errNoMovement
	}
	drop := (vf - v) / vf
	if drop < 0 {
		return 0, nil
	}
	return drop, nil
}

func (m *Model) Uniformity(k1, k2 float64, n int) (float64, error) {
	if n < 2 || k2 <= k1 {
		return 0, errBadSignal
	}
	sum := 0.0
	sumSq := 0.0
	for i := 0; i < n; i++ {
		k := k1 + (k2-k1)*float64(i)/float64(n-1)
		v, err := m.Speed(k)
		if err != nil {
			return 0, err
		}
		sum += v
		sumSq += v * v
	}
	mean := sum / float64(n)
	if mean == 0 {
		return 0, errNoMovement
	}
	variance := sumSq/float64(n) - mean*mean
	if variance < 0 {
		variance = 0
	}
	return math.Sqrt(variance) / mean, nil
}
