package wave

type Merge struct {
	MainFlow float64
	RampFlow float64
	Capacity float64
}

func (m Merge) TotalDemand() float64 {
	return m.MainFlow + m.RampFlow
}

func (m Merge) AcceptedFlow() float64 {
	d := m.TotalDemand()
	if d <= m.Capacity {
		return d
	}
	return m.Capacity
}

func (m Merge) QueueRate() float64 {
	d := m.TotalDemand()
	if d <= m.Capacity {
		return 0
	}
	return d - m.Capacity
}

func (m Merge) QueueAt(q0, duration float64) float64 {
	q := q0 + m.QueueRate()*duration/3600.0
	if q < 0 {
		return 0
	}
	return q
}

func (m Merge) MergeWaveSpeed() float64 {
	up := State{K: 0, Q: m.TotalDemand()}
	down := State{K: 1, Q: m.AcceptedFlow()}
	w, err := ShockSpeed(up, down)
	if err != nil {
		return 0
	}
	return w * 3.6
}

func (m Merge) RampMeteringDelay(rampQueue, duration float64) float64 {
	rate := m.QueueRate() / 3600.0
	if rate <= 0 {
		return 0
	}
	return rampQueue*duration + 0.5*rate*duration*duration
}

func (m Merge) FractionMaintained() float64 {
	if m.TotalDemand() <= m.Capacity {
		return 1.0
	}
	return m.Capacity / m.TotalDemand()
}

func (m Merge) ThroughputGap(duration float64) float64 {
	gap := (m.TotalDemand() - m.AcceptedFlow()) * duration / 3600.0
	if gap < 0 {
		return 0
	}
	return gap
}
