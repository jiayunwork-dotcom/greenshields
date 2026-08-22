package wave

// Bottleneck models a fixed-capacity reduction (e.g. a lane drop or a red
// phase) at which the downstream capacity drops from qCap to qBottleneck. The
// resulting queue grows at the rate of the difference between arrival flow and
// the bottleneck capacity.
type Bottleneck struct {
	Arrival   float64 // arriving flow upstream (veh/h)
	Capacity  float64 // downstream capacity (veh/h)
	Reduction float64 // fractional capacity reduction at the bottleneck (0..1)
}

// Active reports whether the bottleneck is the binding constraint: arrival
// exceeds the reduced capacity so a queue forms.
func (b Bottleneck) Active() bool {
	return b.Arrival > b.EffectiveCapacity()
}

// EffectiveCapacity returns the post-reduction capacity at the bottleneck.
func (b Bottleneck) EffectiveCapacity() float64 {
	if b.Reduction < 0 {
		return b.Capacity
	}
	return b.Capacity * (1 - b.Reduction)
}

// QueueGrowthRate returns the rate (veh/h) at which the queue length increases
// while the bottleneck is active. Non-negative; 0 when not active.
func (b Bottleneck) QueueGrowthRate() float64 {
	eff := b.EffectiveCapacity()
	if b.Arrival <= eff {
		return 0
	}
	return b.Arrival - eff
}

// QueueAt returns the standing queue (vehicles) after duration (s) of the
// bottleneck being active, starting from an initial queue q0.
func (b Bottleneck) QueueAt(q0, duration float64) float64 {
	growth := b.QueueGrowthRate()
	q := q0 + growth*duration/3600.0
	if q < 0 {
		return 0
	}
	return q
}

// BackWaveSpeed is the speed (km/h, signed) at which the queue discharge
// interface moves upstream. It is the kinematic-wave speed between the upstream
// arriving state and the jam state at the bottleneck. A negative value means
// the queue extends backward.
func (b Bottleneck) BackWaveSpeed() float64 {
	up := State{K: 0, Q: b.Arrival}
	// jam state: zero speed, flow = effective capacity (discharging)
	down := State{K: 1, Q: b.EffectiveCapacity()}
	w, err := ShockSpeed(up, down)
	if err != nil {
		return 0
	}
	return w * 3.6 // m/s -> km/h
}

// RecoveryTime returns the time (s) needed to clear an initial queue q0 once
// the bottleneck is removed (arrival <= capacity). Zero if not applicable.
func (b Bottleneck) RecoveryTime(q0 float64) float64 {
	excess := b.Arrival // once cleared, discharge at arrival rate
	if excess <= 0 || q0 <= 0 {
		return 0
	}
	return q0 / (excess / 3600.0)
}

// TotalDelay integrates the queue over a duration (s) assuming a constant
// arrival rate and a fixed downstream capacity. Vehicles in queue accumulate
// delay equal to their waiting time; this is the area under the queue curve.
func (b Bottleneck) TotalDelay(q0, duration float64) float64 {
	if duration <= 0 {
		return 0
	}
	rate := b.QueueGrowthRate() / 3600.0 // veh/s
	if rate <= 0 {
		return 0
	}
	// queue grows linearly from q0; area = q0*d + 0.5*rate*d^2
	return q0*duration + 0.5*rate*duration*duration
}

// SpillbackDistance estimates how far upstream a queue extends given a jam
// density kJam (veh/km) and the current queue length (vehicles).
func (b Bottleneck) SpillbackDistance(queue, kJam float64) float64 {
	if kJam <= 0 {
		return 0
	}
	return queue / kJam // km
}
