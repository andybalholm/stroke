package stroke

import "testing"

func TestZeroLengthPath(t *testing.T) {
	p := [][]Segment{[]Segment{LinearSegment(Pt(0, 0), Pt(0, 0))}}
	s := Stroke(p, Options{
		Width: 10,
		Cap:   RoundCap,
		Join:  RoundJoin,
	})

	if len(s) > 0 {
		t.Errorf("expected empty result, got %d subpaths", len(s))
	}
}
