package core

import "math"

func BPRFreeFlow() float64 { return 1.0 }

const (
	BPRAlpha = 0.15
	BPRBeta  = 4.0
)

func BPRTime(x float64) (float64, error) {
	if x < 0 {
		return 0, errBadSignal
	}
	t0 := BPRFreeFlow()
	return t0 * (1.0 + BPRAlpha*math.Pow(x, BPRBeta)), nil
}

func MarginalDelay(x float64) (float64, error) {
	if x < 0 {
		return 0, errBadSignal
	}
	t0 := BPRFreeFlow()
	return t0 * BPRAlpha * BPRBeta * math.Pow(x, BPRBeta-1), nil
}

func VehicleHoursTraveled(q, capacity, hours float64) (float64, error) {
	if capacity <= 0 || hours <= 0 {
		return 0, errBadSignal
	}
	x := q / capacity
	t, err := BPRTime(x)
	if err != nil {
		return 0, err
	}
	return q * hours * t / 60.0, nil
}

func CongestionCost(q, capacity, hours, vot float64) (float64, error) {
	base, err := VehicleHoursTraveled(q, capacity, hours)
	if err != nil {
		return 0, err
	}
	free := q * hours * BPRFreeFlow() / 60.0
	extra := base - free
	if extra < 0 {
		extra = 0
	}
	return extra * vot, nil
}

func RampMeteringRate(mainFlow, capacity, targetOcc, measuredOcc, maxRate float64) float64 {
	if mainFlow+maxRate <= capacity {
		return 0
	}
	delta := 0.5 * (targetOcc - measuredOcc)
	r := maxRate + delta
	if r < 0 {
		r = 0
	}
	if r > maxRate {
		r = maxRate
	}
	return r
}

func LevelOfService(x float64) (string, string) {
	switch {
	case x < 0.6:
		return "A", "free flow"
	case x < 0.7:
		return "B", "reasonably free"
	case x < 0.8:
		return "C", "stable, restricted"
	case x < 0.9:
		return "D", "approaching unstable"
	case x < 1.0:
		return "E", "at capacity"
	default:
		return "F", "forced / breakdown"
	}
}

func DelaySlope(q, capacity float64) (float64, error) {
	if capacity <= 0 {
		return 0, errZeroCapacity
	}
	x := q / capacity
	if x >= 1 {
		return 0, nil
	}
	d1, err := MarginalDelay(x)
	if err != nil {
		return 0, err
	}
	return d1, nil
}

func SpillbackProbability(queue, storageCapacity float64) float64 {
	if storageCapacity <= 0 {
		return 0
	}
	if queue <= 0 {
		return 0
	}
	p := queue / storageCapacity
	if p > 1 {
		return 1
	}
	return p
}
