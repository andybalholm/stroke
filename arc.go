package stroke

import "math"

// heading returns the angle of the line from the origin to p.
func heading(p Point) float64 {
	return math.Atan2(float64(p.Y), float64(p.X))
}

// ArcSegment returns a Segment that approximates an arc of a circle,
// starting at start, centered at center, and extending angle radians
// counterclockwise. For a clockwise arc, use a negative angle.
// The accuracy of the approximation drops off as the angle increases;
// using a single segment for an arc longer than half a circle (angle > π)
// is not recommended (use AppendArc instead).
func ArcSegment(start, center Point, angle float32) Segment {
	startAngle := heading(start.Sub(center))
	endAngle := startAngle + float64(angle)
	radius := distance(start, center)
	sn, cs := math.Sincos(endAngle)
	end := Pt(center.X+radius*float32(cs), center.Y+radius*float32(sn))

	k := float32(math.Tan(float64(angle)/4)) * 4 / 3
	cp1 := start.Add(rot90CW(center.Sub(start)).Mul(k))
	cp2 := end.Add(rot90CW(end.Sub(center)).Mul(k))
	return Segment{start, cp1, cp2, end}
}

// AppendArc appends one or more segments to dst, forming an arc of a circle
// starting at start, centered at center, and extending angle radians
// counterclockwise. For a clockwise arc, use a negative angle.
func AppendArc(dst []Segment, start, center Point, angle float32) []Segment {
	// Break the arc into segments of less than a radian.
	n := int(math.Abs(float64(angle)) + 1)
	segmentAngle := angle / float32(n)

	pos := start
	for i := 0; i < n; i++ {
		seg := ArcSegment(pos, center, segmentAngle)
		dst = append(dst, seg)
		pos = seg.End
	}

	return dst
}

// AppendEllipticalArc appends one or more segments to dst, forming an arc of
// an ellipse starting at start, with foci at f1 and f2, and extending angle
// radians counterclockwise. For a clockwise arc, use a negative angle.
func AppendEllipticalArc(dst []Segment, start, f1, f2 Point, angle float32) []Segment {
	if f1 == f2 {
		return AppendArc(dst, start, f1, angle)
	}

	semiMajorAxis := float64(distance(start, f1)+distance(start, f2)) * 0.5
	center := f1.Add(f2.Sub(f1).Mul(0.5))
	focalDistance := float64(distance(center, f1))
	semiMinorAxis := math.Sqrt(semiMajorAxis*semiMajorAxis - focalDistance*focalDistance)
	majorAxisAngle := heading(f2.Sub(f1))

	var toUnitCircle affine2D
	toUnitCircle = toUnitCircle.Offset(center.Mul(-1))
	toUnitCircle = toUnitCircle.Rotate(Pt(0, 0), float32(-majorAxisAngle))
	toUnitCircle = toUnitCircle.Scale(Pt(0, 0), Pt(float32(1/semiMajorAxis), float32(1/semiMinorAxis)))

	toEllipse := toUnitCircle.Invert()

	// To avoid a heap allocation, we allocate enough space on the stack to store
	// a full circle's worth of segments (2π = 6.28, so we need 7 segments if
	// one radian is the maximum length).
	var storage [7]Segment

	unitCircleArc := AppendArc(storage[:0], toUnitCircle.Transform(start), Pt(0, 0), angle)
	for _, s := range unitCircleArc {
		dst = append(dst, Segment{
			Start: toEllipse.Transform(s.Start),
			CP1:   toEllipse.Transform(s.CP1),
			CP2:   toEllipse.Transform(s.CP2),
			End:   toEllipse.Transform(s.End),
		})
	}

	return dst
}
