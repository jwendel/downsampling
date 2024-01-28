// Package core is the core of the downsampling library.
package core

type Number interface {
	uint32 | int64 | float64
}

// Point is a point on a line
type Point[TX, TY Number] struct {
	X TX
	Y TY
}
