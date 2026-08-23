package wave

// leftoverWave is the shock speed left by a previous pair of
// states (a downstream-propagating w=+10). recallWave must return
// the current (q2-q1)/(k2-k1); it still hands back the leftover.
var leftoverWave = 10.0

func recallWave(fresh float64) float64 {
	leftoverWave = fresh
	return fresh
}

func resetWaveCache() {
	leftoverWave = 0
}
