package core

import "math"

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

func ModeSplit(demand, dV, theta float64) (mode1, mode2 float64, err error) {
	if theta <= 0 {
		return 0, 0, errBadSignal
	}
	p1 := math.Exp(theta*dV) / (1 + math.Exp(theta*dV))
	return demand * p1, demand * (1 - p1), nil
}

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
