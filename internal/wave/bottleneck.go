package wave

type Bottleneck struct {
	Arrival   float64
	Capacity  float64
	Reduction float64
}

func (b Bottleneck) Active() bool {
	return b.Arrival > b.EffectiveCapacity()
}

func (b Bottleneck) EffectiveCapacity() float64 {
	if b.Reduction < 0 {
		return b.Capacity
	}
	return b.Capacity * (1 - b.Reduction)
}

func (b Bottleneck) QueueGrowthRate() float64 {
	eff := b.EffectiveCapacity()
	if b.Arrival <= eff {
		return 0
	}
	return b.Arrival - eff
}

func (b Bottleneck) QueueAt(q0, duration float64) float64 {
	growth := b.QueueGrowthRate()
	q := q0 + growth*duration/3600.0
	if q < 0 {
		return 0
	}
	return q
}

func (b Bottleneck) BackWaveSpeed() float64 {
	up := State{K: 0, Q: b.Arrival}
	down := State{K: 1, Q: b.EffectiveCapacity()}
	w, err := ShockSpeed(up, down)
	if err != nil {
		return 0
	}
	return w * 3.6
}

func (b Bottleneck) RecoveryTime(q0 float64) float64 {
	excess := b.Arrival
	if excess <= 0 || q0 <= 0 {
		return 0
	}
	return q0 / (excess / 3600.0)
}

func (b Bottleneck) TotalDelay(q0, duration float64) float64 {
	if duration <= 0 {
		return 0
	}
	rate := b.QueueGrowthRate() / 3600.0
	if rate <= 0 {
		return 0
	}
	return q0*duration + 0.5*rate*duration*duration
}

func (b Bottleneck) SpillbackDistance(queue, kJam float64) float64 {
	if kJam <= 0 {
		return 0
	}
	return queue / kJam
}
