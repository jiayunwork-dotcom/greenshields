package wave

// Spillback describes the propagation of a jam from a downstream bottleneck
// back into the upstream stream. It records the two states and the resulting
// wave speed.
type Spillback struct {
	Upstream   State
	Downstream State
	WaveSpeed  float64
	Direction  string
}

// ClassifySpillback builds a spillback descriptor from an upstream state and a
// downstream (more congested) state. A negative wave speed means the
// disturbance travels upstream, which is exactly the congestion spillback
// behaviour seen when a downstream jam propagates backward into free-flowing
// traffic.
func ClassifySpillback(up, down State) (Spillback, error) {
	w, err := ShockSpeed(up, down)
	if err != nil {
		return Spillback{}, err
	}
	return Spillback{
		Upstream:   up,
		Downstream: down,
		WaveSpeed:  w,
		Direction:  Direction(w),
	}, nil
}

// IsUpstream reports whether the spillback propagates toward the upstream
// boundary (a backward-moving jam).
func (s Spillback) IsUpstream() bool {
	return s.WaveSpeed < 0
}
