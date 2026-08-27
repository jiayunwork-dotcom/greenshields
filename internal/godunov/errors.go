package godunov

import "errors"

var (
	ErrBadCFL = errors.New("godunov: dx and dt must be positive and CFL must lie in (0, 1]")

	ErrZeroWave = errors.New("godunov: maximum wave speed is zero, time step undefined")

	ErrFluxMismatch = errors.New("godunov: demand-supply flux disagrees with entropy flux on the same Riemann data")

	ErrMassLeak = errors.New("godunov: two-cell periodic update did not conserve mass")

	ErrNotStationary = errors.New("godunov: equal-flow double root must stay put under Godunov flux")
)
