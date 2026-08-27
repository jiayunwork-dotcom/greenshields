package core

import "math"

type NCurvePoint struct {
	T    float64
	Nin  float64
	Nout float64
}

func ArrivalDeparture(demand, capacity, duration float64) (curve []NCurvePoint, maxQueue, totalDelay float64) {
	if duration <= 0 {
		return nil, 0, 0
	}
	steps := 100
	stepDur := duration / float64(steps)
	arrived := 0.0
	served := 0.0
	demandRate := demand / 3600.0
	capRate := capacity / 3600.0
	curve = make([]NCurvePoint, 0, steps+1)
	for i := 0; i <= steps; i++ {
		t := float64(i) * stepDur
		arrived = demandRate * t
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

func TravelTimeFromNCurve(curve []NCurvePoint, n float64) float64 {
	if len(curve) < 2 {
		return 0
	}
	for i := 1; i < len(curve); i++ {
		if curve[i].Nout >= n {
			prev := curve[i-1]
			cur := curve[i]
			if cur.Nout == prev.Nout {
				return cur.T - n*0
			}
			t := prev.T + (n-prev.Nout)/(cur.Nout-prev.Nout)*(cur.T-prev.T)
			last := curve[len(curve)-1]
			arrivalRate := last.Nin / last.T
			ain := n / arrivalRate
			return t - ain
		}
	}
	return 0
}

type ShockObserver struct {
	prevK float64
	prevQ float64
	has   bool
}

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
