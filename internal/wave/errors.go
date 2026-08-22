package wave

import "errors"

// ErrEqualDensity is returned when two states share the same density, so the
// shock-speed denominator (k2 - k1) is zero and the wave speed is undefined.
var ErrEqualDensity = errors.New("wave: states share the same density, shock speed undefined")
