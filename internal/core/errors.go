package core

import "errors"

var (
	ErrNonPositiveVf = errors.New("greenshields: free-flow speed vf must be greater than 0")

	ErrNonPositiveKj = errors.New("greenshields: jam density kj must be greater than 0")

	ErrDensityTooLow = errors.New("greenshields: density k must be greater than or equal to 0")

	ErrDensityTooHigh = errors.New("greenshields: density k must be less than or equal to kj")

	ErrUnreachableFlow = errors.New("greenshields: flow q exceeds capacity, no real density roots")
)
