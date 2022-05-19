// SPDX-License-Identifier: Unlicense OR MIT

// The following code has been extracted from Gioui.org

package stroke

import (
	"math"
)

// arcTransform computes a transformation that can be used for generating quadratic bézier
// curve approximations for an arc.
//
// The math is extracted from the following paper:
//
//	"Drawing an elliptical arc using polylines, quadratic or
//	 cubic Bezier curves", L. Maisonobe
//
// An electronic version may be found at:
//
//	http://spaceroots.org/documents/ellipse/elliptical-arc.pdf
func arcTransform(p, f1, f2 Point, angle float32) (transform affine2D, segments int) {
	const segmentsPerCircle = 16
	const anglePerSegment = 2 * math.Pi / segmentsPerCircle

	s := angle / anglePerSegment
	if s < 0 {
		s = -s
	}
	segments = int(math.Ceil(float64(s)))
	if segments <= 0 {
		segments = 1
	}

	var rx, ry, alpha float64
	if f1 == f2 {
		// degenerate case of a circle.
		rx = dist(f1, p)
		ry = rx
	} else {
		// semi-major axis: 2a = |PF1| + |PF2|
		a := 0.5 * (dist(f1, p) + dist(f2, p))
		// semi-minor axis: c^2 = a^2 - b^2 (c: focal distance)
		c := dist(f1, f2) * 0.5
		b := math.Sqrt(a*a - c*c)
		switch {
		case a > b:
			rx = a
			ry = b
		default:
			rx = b
			ry = a
		}
		if f1.X == f2.X {
			// special case of a "vertical" ellipse.
			alpha = math.Pi / 2
			if f1.Y < f2.Y {
				alpha = -alpha
			}
		} else {
			x := float64(f1.X-f2.X) * 0.5
			if x < 0 {
				x = -x
			}
			alpha = math.Acos(x / c)
		}
	}

	var (
		θ   = angle / float32(segments)
		ref affine2D // transform from absolute frame to ellipse-based one
		rot affine2D // rotation matrix for each segment
		inv affine2D // transform from ellipse-based frame to absolute one
	)
	center := Point{
		X: 0.5 * (f1.X + f2.X),
		Y: 0.5 * (f1.Y + f2.Y),
	}
	ref = ref.Offset(Point{}.Sub(center))
	ref = ref.Rotate(Point{}, float32(-alpha))
	ref = ref.Scale(Point{}, Point{
		X: float32(1 / rx),
		Y: float32(1 / ry),
	})
	inv = ref.Invert()
	rot = rot.Rotate(Point{}, 0.5*θ)

	// Instead of invoking math.Sincos for every segment, compute a rotation
	// matrix once and apply for each segment.
	// Before applying the rotation matrix rot, transform the coordinates
	// to a frame centered to the ellipse (and warped into a unit circle), then rotate.
	// Finally, transform back into the original frame.
	return inv.Mul(rot).Mul(ref), segments
}

func dist(p1, p2 Point) float64 {
	var (
		x1 = float64(p1.X)
		y1 = float64(p1.Y)
		x2 = float64(p2.X)
		y2 = float64(p2.Y)
		dx = x2 - x1
		dy = y2 - y1
	)
	return math.Hypot(dx, dy)
}
