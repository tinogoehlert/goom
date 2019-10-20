package geometry

import "math"

// Vec2 is a coordinate in 2D space
type Vec2 struct {
	x      float32
	y      float32
	length float32
}

// V2 creates a new Vec2
func V2(x, y float32) Vec2 {
	return Vec2{
		x:      x,
		y:      y,
		length: float32(math.Sqrt(float64(x*x + y*y))),
	}
}

// X coordinate
func (a *Vec2) X() float32 {
	return a.x
}

// Y coordinate
func (a *Vec2) Y() float32 {
	return a.y
}

// XY coordinates
func (a *Vec2) XY() [2]float32 {
	return [2]float32{a.x, a.y}
}

// SetLength Set the length of the vector. A negative [value] will change the vectors
// orientation and a [value] of zero will set the vector to zero.
func (a *Vec2) SetLength(length float32) {
	if length == 0.0 {
		a.x = 0
		a.y = 0
		return
	}
	if a.length == 0.0 {
		return
	}
	l := length / a.length
	a.x *= l
	a.y *= l
}

// Scale this by factor.
func (a Vec2) Scale(factor float32) Vec2 {
	return V2(a.x*factor, a.y*factor)
}

// Mul Multiply two vectors.
func (a Vec2) Mul(b Vec2) Vec2 {
	return V2(a.x*b.x, a.y*b.y)
}

// Div Divide two vectors.
func (a Vec2) Div(b Vec2) Vec2 {
	return V2(a.x/b.x, a.y/b.y)
}

// DivScalar Divide by scalar.
func (a Vec2) DivScalar(s float32) Vec2 {
	return V2(a.x/s, a.y/s)
}

// Add two vectors.
func (a Vec2) Add(b Vec2) Vec2 {
	return V2(a.x+b.x, a.y+b.y)
}

// Sub substract two vectors.
func (a Vec2) Sub(b Vec2) Vec2 {
	return V2(a.x-b.x, a.y-b.y)
}

// Cross product.
func (a Vec2) Cross(b Vec2) float32 {
	return a.x*b.y - a.y*b.x
}

// CrossVec2 Cross product.
func (a Vec2) CrossVec2() Vec2 {
	return V2(-a.y, a.x)
}

// Dot Inner product.
func (a *Vec2) Dot(b Vec2) float32 {
	return a.x*b.x + a.y*b.y
}

// Normalize this.
func (a Vec2) Normalize() Vec2 {
	if a.length == 0.0 {
		return V2(0, 0)
	}
	d := 1.0 / a.length
	return V2(a.x*d, a.y*d)
}

// AngleTo Returns the angle between this vector and [other] in radians.
func (a Vec2) AngleTo(b Vec2) float32 {
	if a.x == b.x && a.y == b.y {
		return 0.0
	}
	d := a.Dot(b) / (a.length * b.length)
	return float32(math.Acos(float64(clamp(d, -1.0, 1.0))))
}

// DistanceTo Distance from this to b
func (a Vec2) DistanceTo(b Vec2) float32 {
	return float32(math.Sqrt(float64(a.DistanceToSqared(b))))
}

// DistanceToSqared Squared distance from this to b
func (a Vec2) DistanceToSqared(b Vec2) float32 {
	var (
		dx = a.x - b.x
		dy = a.y - b.y
	)
	return dx*dx + dy*dy
}
