// Package wave computes kinematic-wave (shock) speeds between two traffic
// states. It is independent of the Greenshields model so that any pair of
// measured (density, flow) states can be compared, but it is designed to be
// fed states produced by the core package.
package wave

// State is one traffic state described by its density and flow.
type State struct {
	K float64
	Q float64
}
