// Package core implements the Greenshields traffic-flow fundamental diagram.
//
// Greenshields (1935) assumes a linear speed-density relationship. From it the
// flow-density and flow-speed diagrams follow by multiplication. This package
// is the single source of truth for the traffic math used by the HTTP API and
// the command-line entry point.
package core

import "errors"

// Domain and validation errors returned by the Greenshields core. Callers can
// wrap these with %w to add context while still being able to match them.
var (
	// ErrNonPositiveVf is returned when the free-flow speed is not positive.
	ErrNonPositiveVf = errors.New("greenshields: free-flow speed vf must be greater than 0")

	// ErrNonPositiveKj is returned when the jam density is not positive.
	ErrNonPositiveKj = errors.New("greenshields: jam density kj must be greater than 0")

	// ErrDensityTooLow is returned when a density is below zero.
	ErrDensityTooLow = errors.New("greenshields: density k must be greater than or equal to 0")

	// ErrDensityTooHigh is returned when a density exceeds the jam density.
	ErrDensityTooHigh = errors.New("greenshields: density k must be less than or equal to kj")

	// ErrUnreachableFlow is returned when a flow cannot be produced by any
	// density because it is above the capacity of the model.
	ErrUnreachableFlow = errors.New("greenshields: flow q exceeds capacity, no real density roots")
)
