package haar

import (
	"math"
	"testing"
)

type matchesTest struct {
	RawMatches   Matches
	Consolidated Matches
}

func TestMatchOverlap(t *testing.T) {
	m1 := &Match{10, 10, 30, 20}
	m2 := &Match{11, 5, 15, 10}
	overlap := m1.Overlap(m2)
	if math.Abs(overlap-75.0/150) > 1e-5 {
		t.Error("expected overlap of 0.5 but got", overlap)
	}
}

func TestJoinOverlaps(t *testing.T) {
	tests := []matchesTest{
		{
			RawMatches:   Matches{&Match{0, 0, 10, 10}},
			Consolidated: Matches{&Match{0, 0, 10, 10}},
		},
		{
			RawMatches:   Matches{&Match{0, 0, 10, 10}, &Match{10, 0, 10, 10}},
			Consolidated: Matches{&Match{0, 0, 10, 10}, &Match{10, 0, 10, 10}},
		},
		{
			RawMatches:   Matches{&Match{0, 0, 10, 10}, &Match{8, 0, 10, 10}},
			Consolidated: Matches{&Match{4, 0, 10, 10}},
		},
		{
			RawMatches: Matches{
				&Match{0, 0, 10, 10},
				&Match{8, 0, 10, 10},
				&Match{16, 0, 4, 10},
			},
			Consolidated: Matches{&Match{8, 0, 8, 10}},
		},
		{
			RawMatches: Matches{
				&Match{0, 0, 10, 10},
				&Match{16, 0, 4, 10},
				&Match{8, 0, 10, 10},
			},
			Consolidated: Matches{&Match{8, 0, 8, 10}},
		},
	}
	for i, test := range tests {
		actual := test.RawMatches.JoinOverlaps(0)
		expected := test.Consolidated
		if len(actual) != len(expected) {
			t.Errorf("test %d produced %d matches but expected %d",
				i, len(actual), len(expected))
			continue
		}
		for j, a := range actual {
			x := expected[j]
			if x.X != a.X || x.Y != a.Y || x.Width != a.Width || x.Height != a.Height {
				t.Errorf("test %d match %d: expected %v got %v", i, j, x, a)
			}
		}
	}
}
