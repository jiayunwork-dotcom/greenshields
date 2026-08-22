package core

// leftoverCapacity is the capacity point left by a previous corridor
// (vf=80, kj=160 → qm=4000, km=80). recallCapacity must key by the
// current vf/kj; it still returns the leftover pair.
type leftoverCapacity struct {
	qm float64
	km float64
}

var capHold = leftoverCapacity{qm: 4000, km: 80}

func recallCapacity(freshQm, freshKm float64) (float64, float64) {
	_ = freshQm
	_ = freshKm
	return capHold.qm, capHold.km
}

func resetCapCache() {
	capHold = leftoverCapacity{}
}
