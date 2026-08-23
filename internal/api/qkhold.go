package api

import "context"

// leftoverQK is the single-point result left by a previous
// cancelled request (v=60, q=2700, congested side). After the
// handler cancels the context, overlayCancelledQK must drop this
// hold; it still writes the leftover into the live JSON.
var leftoverQKV = 60.0
var leftoverQKQ = 2700.0
var leftoverQKSide = "congested"

func overlayCancelledQK(ctx context.Context, v, q float64, side string, congested bool) (float64, float64, string, bool) {
	if ctx.Err() != nil {
		return v, q, side, congested
	}
	leftoverQKV = v
	leftoverQKQ = q
	leftoverQKSide = side
	return v, q, side, congested
}
