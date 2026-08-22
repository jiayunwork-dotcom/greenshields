package core

import "math"

// BPRFreeFlow is the free-flow travel time (minutes) used by the Bureau of
// Public Roads volume-delay function. The BPR model is a complementary
// macroscopic delay model layered on top of the Greenshields fundamental
// diagram; it is kept here because both describe speed-flow conformity.
func BPRFreeFlow() float64 { return 1.0 }

// BPRAlpha and BPRBeta are the standard BPR shape constants t = t0 (1 + α (x)^β)
// with x = V/C. Defaults α=0.15, β=4 are the classic highway values.
const (
	BPRAlpha = 0.15
	BPRBeta  = 4.0
)

// BPRTime returns the travel time (minutes) on a link under volume-to-capacity
// ratio x using the BPR function. At x=0 it equals the free-flow time; it rises
// steeply as x→1.
func BPRTime(x float64) (float64, error) {
	if x < 0 {
		return 0, errBadSignal
	}
	t0 := BPRFreeFlow()
	return t0 * (1.0 + BPRAlpha*math.Pow(x, BPRBeta)), nil
}

// MarginalDelay returns the extra travel time per additional vehicle (minutes)
// as the derivative of the BPR function with respect to x, scaled by the
// free-flow time.
func MarginalDelay(x float64) (float64, error) {
	if x < 0 {
		return 0, errBadSignal
	}
	t0 := BPRFreeFlow()
	return t0 * BPRAlpha * BPRBeta * math.Pow(x, BPRBeta-1), nil
}

// VehicleHoursTraveled sums total travel time (vehicle-hours) for a demand Q
// (veh/h) across a link of free-flow time t0 (min) over a period of hours.
func VehicleHoursTraveled(q, capacity, hours float64) (float64, error) {
	if capacity <= 0 || hours <= 0 {
		return 0, errBadSignal
	}
	x := q / capacity
	t, err := BPRTime(x)
	if err != nil {
		return 0, err
	}
	return q * hours * t / 60.0, nil // minutes -> hours
}

// CongestionCost monetises the additional delay caused by congestion for a
// demand Q over `hours` at a value of time vot (currency per vehicle-hour).
// It compares the BPR travel time against free-flow and prices the difference.
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

// RampMeteringRate returns the ALINEA-style metered rate (veh/h) that keeps the
// merge from exceeding capacity: it releases the ramp at a rate proportional to
// mainline occupancy error from a target. Returns 0 when no control needed.
func RampMeteringRate(mainFlow, capacity, targetOcc, measuredOcc, maxRate float64) float64 {
	if mainFlow+maxRate <= capacity {
		return 0 // no metering needed
	}
	// ALINEA: r_next = r_prev + Kp*(occ_target - occ_measured), clamped >=0
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

// LevelOfService classifies a volume-to-capacity ratio into the standard LOS
// bands A–F used in practice. Returns the letter and a human note.
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

// DelaySlope estimates the slope of the macroscopic delay-vs-demand curve at a
// given demand, useful for identifying the knee of the curve near capacity.
func DelaySlope(q, capacity float64) (float64, error) {
	if capacity <= 0 {
		return 0, errZeroCapacity
	}
	x := q / capacity
	if x >= 1 {
		return 0, nil // saturated: slope saturates; reported separately
	}
	d1, err := MarginalDelay(x)
	if err != nil {
		return 0, err
	}
	return d1, nil
}

// SpillbackProbability is a crude indicator (0..1) that a queue spills back
// into the upstream link, growing with the queue length relative to a storage
// length Lstore (vehicles).
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
