package stroke

import "math"

// ArcSegment returns a Segment that approximates an arc of a circle,
// starting at start, centered at center, and extending angle radians
// counterclockwise. For a clockwise arc, use a negative angle.
// The accuracy of the approximation drops off as the angle increases;
// using a single segment for an arc longer than half a circle (angle > π)
// is not recommended (use AppendArc instead).
func ArcSegment(start, center Point, angle float32) Segment {
	startAngle := math.Atan2(float64(start.Y-center.Y), float64(start.X-center.X))
	endAngle := startAngle + float64(angle)
	radius := distance(start, center)
	end := Pt(center.X+radius*float32(math.Cos(endAngle)), center.Y+radius*float32(math.Sin(endAngle)))

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
