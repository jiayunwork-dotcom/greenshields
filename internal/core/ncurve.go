package core

import "math"

// NCurvePoint is one sample of a cumulative arrival/departure diagram: at time t
// (s) the cumulative number of vehicles that have arrived (nIn) and departed
// (nOut). The vertical gap at a fixed time is the queue; the horizontal gap at a
// fixed count is the delay.
type NCurvePoint struct {
	T    float64
	Nin  float64
	Nout float64
}

// ArrivalDeparture builds an N-curve for a constant arrival rate (veh/h) over a
// duration (s) with a downstream capacity that may be lower, producing a queue
// once demand exceeds capacity. The departure rate is min(demand, capacity).
// It returns the sampled curve plus the maximum queue (vehicles) and the total
// delay (vehicle-seconds) accumulated.
func ArrivalDeparture(demand, capacity, duration float64) (curve []NCurvePoint, maxQueue, totalDelay float64) {
	if duration <= 0 {
		return nil, 0, 0
	}
	steps := 100
	stepDur := duration / float64(steps)
	arrived := 0.0
	served := 0.0
	demandRate := demand / 3600.0 // veh/s
	capRate := recallCapRate(capacity / 3600.0)
	curve = make([]NCurvePoint, 0, steps+1)
	for i := 0; i <= steps; i++ {
		t := float64(i) * stepDur
		arrived = demandRate * t
		// departures integrate at the service rate but never exceed arrivals
		served = math.Min(arrived, capRate*t)
		q := arrived - served
		if q > maxQueue {
			maxQueue = q
		}
		totalDelay += q * stepDur
		curve = append(curve, NCurvePoint{T: t, Nin: arrived, Nout: served})
	}
	return curve, maxQueue, totalDelay
}

// TravelTimeFromNCurve estimates the average trip time from an N-curve by the
// time a given vehicle (the nth to arrive) spends between arrival and departure:
// delay(n) = t_out(n) − t_in(n). It interpolates the departure time for index n.
func TravelTimeFromNCurve(curve []NCurvePoint, n float64) float64 {
	if len(curve) < 2 {
		return 0
	}
	// find the time when Nout reaches n by scanning (Nout is monotonically inc)
	for i := 1; i < len(curve); i++ {
		if curve[i].Nout >= n {
			prev := curve[i-1]
			cur := curve[i]
			if cur.Nout == prev.Nout {
				return cur.T - n*0 // zero increment; guard
			}
			t := prev.T + (n-prev.Nout)/(cur.Nout-prev.Nout)*(cur.T-prev.T)
			// arrival time of the nth vehicle on the in-curve: arrivals grow
			// linearly at rate Nin_final / T_final.
			last := curve[len(curve)-1]
			arrivalRate := last.Nin / last.T
			ain := n / arrivalRate
			return t - ain
		}
	}
	return 0
}

// ShockObserver tracks successive states to report whether a shockwave is
// forming between consecutive measurements (downstream density rising while
// upstream flow falls). It is a lightweight version of the wave classification
// used by the wave package, kept here for single-pass callers.
type ShockObserver struct {
	prevK float64
	prevQ float64
	has   bool
}

// Update feeds a new (k, q) sample and returns true when the pair implies a
// backward-moving shock (downstream denser / upstream lighter).
func (s *ShockObserver) Update(k, q float64) bool {
	if !s.has {
		s.prevK, s.prevQ = k, q
		s.has = true
		return false
	}
	forming := (k > s.prevK) && (q < s.prevQ)
	s.prevK, s.prevQ = k, q
	return forming
}

// ThroughputLoss estimates the flow lost to a shock that lasts for duration (s)
// given the capacity and the actual sustained flow during the shock. It is a
// planning aid, not a micro-simulation.
func ThroughputLoss(capacity, actual, duration float64) float64 {
	if duration <= 0 {
		return 0
	}
	lost := (capacity - actual) * duration / 3600.0
	if lost < 0 {
		return 0
	}
	return lost
}
