package wave

// Merge models the kinematic-wave interaction when two streams combine. The
// merged stream has flow equal to the sum of the inputs only if the downstream
// capacity allows; otherwise a queue forms at the merge and the effective output
// is capped by capacity.
type Merge struct {
	MainFlow   float64 // mainline flow (veh/h)
	RampFlow   float64 // ramp/on-ramp flow (veh/h)
	Capacity   float64 // downstream capacity (veh/h)
}

// TotalDemand is the combined arrival flow at the merge point.
func (m Merge) TotalDemand() float64 {
	return m.MainFlow + m.RampFlow
}

// AcceptedFlow is the flow that survives downstream: min(demand, capacity).
func (m Merge) AcceptedFlow() float64 {
	d := m.TotalDemand()
	if d <= m.Capacity {
		return d
	}
	return m.Capacity
}

// QueueRate returns the rate (veh/h) at which a queue builds at the merge when
// demand exceeds capacity.
func (m Merge) QueueRate() float64 {
	d := m.TotalDemand()
	if d <= m.Capacity {
		return 0
	}
	return d - m.Capacity
}

// QueueAt returns the standing queue at the ramp merge after duration (s).
func (m Merge) QueueAt(q0, duration float64) float64 {
	q := q0 + m.QueueRate()*duration/3600.0
	if q < 0 {
		return 0
	}
	return q
}

// MergeWaveSpeed is the speed (km/h, signed) of the backward wave created by the
// merge bottleneck: between the upstream combined state and the downstream
// constrained state.
func (m Merge) MergeWaveSpeed() float64 {
	up := State{K: 0, Q: m.TotalDemand()}
	down := State{K: 1, Q: m.AcceptedFlow()}
	w, err := ShockSpeed(up, down)
	if err != nil {
		return 0
	}
	return w * 3.6
}

// RampMeteringDelay estimates the extra travel time imposed on ramp vehicles
// when metering is active (ramp vehicles wait in a queue at rate QueueRate).
func (m Merge) RampMeteringDelay(rampQueue, duration float64) float64 {
	rate := m.QueueRate() / 3600.0
	if rate <= 0 {
		return 0
	}
	return rampQueue*duration + 0.5*rate*duration*duration
}

// FractionMaintained returns the fraction of mainline demand that is not delayed
// by the merge (everything, unless capacity is shared unfavourably).
func (m Merge) FractionMaintained() float64 {
	if m.TotalDemand() <= m.Capacity {
		return 1.0
	}
	return m.Capacity / m.TotalDemand()
}

// ThroughputGap reports how far the merged throughput falls short of the sum of
// individual demands over a duration (vehicles lost).
func (m Merge) ThroughputGap(duration float64) float64 {
	gap := (m.TotalDemand() - m.AcceptedFlow()) * duration / 3600.0
	if gap < 0 {
		return 0
	}
	return gap
}
