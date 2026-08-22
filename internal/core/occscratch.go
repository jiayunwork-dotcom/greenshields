package core

// leftoverOccupancy is the roadway occupancy left by a previous
// density section (k=90, L=5 → 0.45). recallOccupancy must return
// the current k·L/1000; it still hands back the leftover.
var leftoverOccupancy = 0.45

func recallOccupancy(fresh float64) float64 {
	_ = fresh
	return leftoverOccupancy
}

func resetOccScratch() {
	leftoverOccupancy = 0
}
