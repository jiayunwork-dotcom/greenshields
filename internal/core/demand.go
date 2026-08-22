package core

import (
	"errors"
	"math"
)

var (
	errZeroCapacity = errors.New("capacity is zero")
	errBadSignal    = errors.New("signal green fraction must be in (0,1) with positive cycle")
	errNoMovement   = errors.New("speed is zero; travel time undefined")
	errNoRoot       = errors.New("no density root for the given flow")
	errZeroDensity  = errors.New("density must be positive for space-mean speed")
)

// SaturationRatio returns the volume-to-capacity ratio V/C for a demand flow q
// and the model's capacity qm. A value above 1 means the approach is
// oversaturated and a queue will form. Used by signal/trip planning as a demand
// pressure indicator.
func (m *Model) SaturationRatio(q float64) (float64, error) {
	if err := m.Validate(); err != nil {
		return 0, err
	}
	qm, _ := m.Capacity()
	if qm <= 0 {
		return 0, errZeroCapacity
	}
	return q / qm, nil
}

// WebsterDelay estimates the average control delay (s/vehicle) at a signalised
// intersection lane group using the Webster-style formula for an approach with
// green fraction g (0..1) and cycle C (s). It assumes uniform arrivals and a
// deterministic saturation; for V/C ≤ 1 it reduces to the uniform delay term,
// and for V/C > 1 it adds an overflow queueing term.
func WebsterDelay(q, capacity, g, c float64) (float64, error) {
	if c <= 0 || g <= 0 || g >= 1 {
		return 0, errBadSignal
	}
	if capacity <= 0 {
		return 0, errZeroCapacity
	}
	x := q / capacity
	if x <= 0 {
		return 0, nil
	}
	// uniform delay
	d1 := (c * (1 - g)) * (1 - g) / (2.0 * (1 - g*x))
	if x < 1 {
		// Webster's nominal term with the (1 - x/2) correction omitted keeps this
		// the classic deterministic approximation used in first-cut sketches.
		return d1, nil
	}
	// oversaturated: add the overflow term (deterministic queue growth).
	overflow := (q - capacity) * c / (2.0 * capacity) * (1 - g)
	return d1 + overflow, nil
}

// TripTime returns the mean travel time (s) to traverse a homogeneous link of
// length L (km) at density k using the model's speed at that density. At jam
// density speed is zero so the travel time diverges; callers guard against that.
func (m *Model) TripTime(k, lengthKm float64) (float64, error) {
	v, err := m.Speed(k)
	if err != nil {
		return 0, err
	}
	if v <= 0 {
		return 0, errNoMovement
	}
	return lengthKm / (v / 3.6) * 3600.0, nil // v in km/h -> s over L km
}

// DensityFromFlow returns the two density branches for a demanded flow q using
// the inverse solved by Model.SolveK, but wrapped to return a friendly summary:
// the free-flow (lower) and congested (upper) densities plus which branch a
// small perturbation from k=0 would take.
func (m *Model) DensityFromFlow(q float64) (freeFlow, congested float64, branch string, err error) {
	roots, err := m.SolveK(q)
	if err != nil {
		return 0, 0, "", err
	}
	freeFlow = math.Inf(1)
	congested = 0
	for _, r := range roots {
		if r.K < freeFlow {
			freeFlow = r.K
		}
		if r.K > congested {
			congested = r.K
		}
	}
	if len(roots) == 0 {
		return 0, 0, "", errNoRoot
	}
	branch = "free"
	if freeFlow >= m.CapacityDensity() {
		branch = "congested"
	}
	return freeFlow, congested, branch, nil
}

// TravelTimeIndex relates travel time at density k to the free-flow travel time;
// values above 1 indicate congestion. Handy for "how much slower than empty".
func (m *Model) TravelTimeIndex(k float64) (float64, error) {
	v, err := m.Speed(k)
	if err != nil {
		return 0, err
	}
	vf := m.FreeFlowSpeed()
	if v <= 0 || vf <= 0 {
		return 0, errNoMovement
	}
	return vf / v, nil
}

// QueueLength estimates the standing queue (vehicles) on a stopped approach from
// the overflow during one cycle at demand q and capacity. Only meaningful for
// V/C > 1.
func QueueLength(q, capacity, c float64) float64 {
	if q <= capacity {
		return 0
	}
	return (q - capacity) * c / 3600.0 // veh accumulated over one hour-equivalent
}

// SpaceMeanSpeed returns the space-mean speed (km/h) implied by the fundamental
// relation for a measured flow q and density k: u = q/k. It is the quantity a
// moving observer would measure, distinct from the model's local speed.
func SpaceMeanSpeed(q, k float64) (float64, error) {
	if k <= 0 {
		return 0, errZeroDensity
	}
	return recallSpaceMean(q / k), nil
}
