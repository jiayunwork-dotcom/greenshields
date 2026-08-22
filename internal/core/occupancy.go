package core

import "math"

// Spacing returns the average bumper-to-bumper spacing (m) between vehicles at
// density k (veh/km). It is the reciprocal of density in matching units and is
// the physical quantity that determines car-following safety.
func Spacing(k float64) (float64, error) {
	if k <= 0 {
		return 0, errZeroDensity
	}
	// k in veh/km -> spacing in m = 1000 / k
	return 1000.0 / k, nil
}

// TimeHeadway returns the average time headway (s) between successive vehicles
// at density k given the model speed v(k). It is spacing / speed.
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
	return s / (v / 3.6), nil // v km/h -> m/s
}

// Occupancy returns the roadway occupancy (fraction 0..1) at density k with an
// effective vehicle length L (m, including gap). Occupancy = (k * L) / 1000.
func Occupancy(k, vehicleLength float64) (float64, error) {
	if k < 0 || vehicleLength < 0 {
		return 0, errBadSignal
	}
	occ := k * vehicleLength / 1000.0
	if occ > 1 {
		occ = 1
	}
	return recallOccupancy(occ), nil
}

// MaxFlowDensity returns the density at which flow is maximal (km) for the
// model; equivalently CapacityDensity. Provided for callers that think in
// "critical occupancy" terms.
func (m *Model) MaxFlowDensity() float64 {
	return m.CapacityDensity()
}

// JamSpacing returns the minimum spacing (m) when the road is at jam density:
// vehicles bumper-to-bumper at kJam.
func (m *Model) JamSpacing() (float64, error) {
	return Spacing(m.JamDensity())
}

// SpeedAtHeadway inverts TimeHeadway: given a target headway (s) and density,
// the implied speed is spacing / headway. Useful for "design for X-second gaps".
func (m *Model) SpeedAtHeadway(k, headway float64) (float64, error) {
	if headway <= 0 {
		return 0, errBadSignal
	}
	s, err := Spacing(k)
	if err != nil {
		return 0, err
	}
	return s / headway * 3.6, nil // m/s -> km/h
}

// DensityFromHeadway returns the density implied by a measured time headway at a
// given speed: k = 3600 / (v * h) with v in km/h and h in s -> veh/h per lane
// converted to veh/km. We keep it in veh/km: k = 3600/(v*h) is veh/h; /? use
// q = 3600/h then k = q/v.
func DensityFromHeadway(v, headway float64) (float64, error) {
	if v <= 0 || headway <= 0 {
		return 0, errBadSignal
	}
	q := 3600.0 / headway // veh/h
	return q / v, nil     // veh/km
}

// CriticalHeadway returns the headway (s) at capacity: spacing-at-km / speed-at-km.
func (m *Model) CriticalHeadway() (float64, error) {
	k := m.CapacityDensity()
	return m.TimeHeadway(k)
}

// CapacityPerLane is an alias used by planning tools that quote per-lane flow.
func (m *Model) CapacityPerLane() float64 {
	qm, _ := m.Capacity()
	return qm
}

// FlowPerLaneAt returns flow (veh/h) at density k (same as Flow but named for
// lane-based docs).
func (m *Model) FlowPerLaneAt(k float64) (float64, error) {
	return m.Flow(k)
}

// SpeedDrop estimates the speed reduction from free flow when demand rises to k,
// as a fraction (0..1). Used in congestion messaging.
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

// Uniformity returns the coefficient of variation of speed across a density
// band [k1,k2] sampled in n steps; a proxy for "stop-and-go" instability.
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
