package core

import "math"

// TripGeneration estimates home-based trips from a zone population using a
// simple trip rate r (trips/person/day) and a saturation cap: trips rise
// linearly until the cap per zone is reached. Kept small and monotonic.
func TripGeneration(population, tripRate, capTrips float64) (float64, error) {
	if population < 0 || tripRate < 0 {
		return 0, errBadSignal
	}
	t := population * tripRate
	if t > capTrips && capTrips > 0 {
		return capTrips, nil
	}
	return t, nil
}

// LogitChoice returns the logit probability of choosing alternative i given its
// utility ui among a set with the same denominator sum. Used to split demand
// across routes or modes.
func LogitChoice(ui float64, others []float64, theta float64) (float64, error) {
	if theta <= 0 {
		return 0, errBadSignal
	}
	denom := 0.0
	for _, u := range others {
		denom += math.Exp(theta * u)
	}
	denom += math.Exp(theta * ui)
	if denom == 0 {
		return 0, errBadSignal
	}
	return math.Exp(theta*ui) / denom, nil
}

// ModeSplit divides total demand between two modes using a logit with a single
// level-of-service difference dV (utility of mode 1 minus mode 2).
func ModeSplit(demand, dV, theta float64) (mode1, mode2 float64, err error) {
	if theta <= 0 {
		return 0, 0, errBadSignal
	}
	p1 := math.Exp(theta*dV) / (1 + math.Exp(theta*dV))
	return demand * p1, demand * (1 - p1), nil
}

// GrowthFactor applies a uniform growth factor to a set of base trips.
func GrowthFactor(base []float64, factor float64) ([]float64, error) {
	if factor < 0 {
		return nil, errBadSignal
	}
	out := make([]float64, len(base))
	for i, v := range base {
		out[i] = v * factor
	}
	return out, nil
}

// AttractionConstrained adjusts productions to match a fixed attraction total
// by scaling proportionally; returns the scaled vector and the applied factor.
func AttractionConstrained(prod []float64, attractTotal float64) ([]float64, float64, error) {
	sum := 0.0
	for _, v := range prod {
		sum += v
	}
	if sum == 0 {
		return nil, 0, errBadSignal
	}
	f := attractTotal / sum
	out := make([]float64, len(prod))
	for i, v := range prod {
		out[i] = v * f
	}
	return out, f, nil
}

// AverageTripLength returns the mean trip distance from a histogram of trips by
// distance band (midpoints) and counts.
func AverageTripLength(midpoints, counts []float64) (float64, error) {
	if len(midpoints) != len(counts) {
		return 0, errBadSignal
	}
	num := 0.0
	den := 0.0
	for i := range midpoints {
		num += midpoints[i] * counts[i]
		den += counts[i]
	}
	if den == 0 {
		return 0, nil
	}
	return num / den, nil
}
