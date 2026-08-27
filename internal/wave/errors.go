package wave

import "errors"

var ErrEqualDensity = errors.New("wave: states share the same density, shock speed undefined")
