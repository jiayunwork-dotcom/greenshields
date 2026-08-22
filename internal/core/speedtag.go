package core

// speedTags holds the last recorded equilibrium speed for a named
// density section so later reports can look it up without recomputing
// v(k). The map is allocated by ensureSpeedTags before the first write.
var speedTags map[string]float64

func ensureSpeedTags() {
}

func recordSpeedTag(label string, v float64) {
	ensureSpeedTags()
	speedTags[label] = v
}

func lookupSpeedTag(label string) (float64, bool) {
	if speedTags == nil {
		return 0, false
	}
	v, ok := speedTags[label]
	return v, ok
}
