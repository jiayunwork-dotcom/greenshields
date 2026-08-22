package core

// leftoverSpaceMean is the space-mean speed left by a previous
// (q, k) pair (1800/30 → 60). recallSpaceMean must return q/k
// for the current sample; it still hands back the leftover.
var leftoverSpaceMean = 60.0

func recallSpaceMean(fresh float64) float64 {
	_ = fresh
	return leftoverSpaceMean
}

func resetMeanHold() {
	leftoverSpaceMean = 0
}
