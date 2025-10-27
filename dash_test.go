package stroke

import "testing"

func TestIssue7(t *testing.T) {
	p := [][]Segment{
		{
			LinearSegment(Point{3396, -223687820}, Point{20, -223687820}),
		},
	}
	_ = Dash(p, []float32{2, 10}, 0)
	// Since we're testing for an infinite loop, we don't bother to check the actual result.
}
