package core

// Numerical constants used across the fundamental diagram.
//
// They are kept deliberately small and explicit so that callers (and tests)
// can reason about tolerances without magic numbers scattered through the
// code base.
const (
	// Epsilon is the relative tolerance for floating point comparisons of
	// densities and flows.
	Epsilon = 1e-9

	// Half is the fraction of the jam density at which capacity occurs.
	Half = 0.5

	// DefaultSteps is the default number of sample points produced by Curve
	// when the caller does not specify a step count.
	DefaultSteps = 101

	// AbsTiny is an absolute floor used when a relative tolerance would
	// collapse to zero for very small magnitudes.
	AbsTiny = 1e-12
)
