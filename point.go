package stroke

import "strconv"

// A Point is a two dimensional point.
type Point struct {
	X, Y float32
}

// String return a string representation of p.
func (p Point) String() string {
	return "(" + strconv.FormatFloat(float64(p.X), 'f', -1, 32) +
		"," + strconv.FormatFloat(float64(p.Y), 'f', -1, 32) + ")"
}

// Pt is shorthand for Point{X: x, Y: y}.
func Pt(x, y float32) Point {
	return Point{X: x, Y: y}
}

// Add return the point p+p2.
func (p Point) Add(p2 Point) Point {
	return Point{X: p.X + p2.X, Y: p.Y + p2.Y}
}

// Sub returns the vector p-p2.
func (p Point) Sub(p2 Point) Point {
	return Point{X: p.X - p2.X, Y: p.Y - p2.Y}
}

// Mul returns p scaled by s.
func (p Point) Mul(s float32) Point {
	return Point{X: p.X * s, Y: p.Y * s}
}

// Div returns the vector p/s.
func (p Point) Div(s float32) Point {
	return Point{X: p.X / s, Y: p.Y / s}
}
