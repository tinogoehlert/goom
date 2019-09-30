package geometry

// Coord is a coordinate in 2D space
type Coord [2]float32

// X coordinate
func (c *Coord) X() float32 { return c[0] }

// Y coordinate
func (c *Coord) Y() float32 { return c[1] }
