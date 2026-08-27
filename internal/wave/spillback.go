package wave

type Spillback struct {
	Upstream   State
	Downstream State
	WaveSpeed  float64
	Direction  string
}

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

func (s Spillback) IsUpstream() bool {
	return s.WaveSpeed < 0
}
