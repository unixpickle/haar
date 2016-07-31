package haar

import (
	"fmt"
	"math"
)

// A Match is a region inside an image in which an object
// was detected.
type Match struct {
	X      int
	Y      int
	Width  int
	Height int
}

// String returns a human readable version of the
// match's bounding rectangle.
func (m *Match) String() string {
	return fmt.Sprintf("(x=%d y=%d width=%d height=%d)", m.X, m.Y, m.Width, m.Height)
}

// Overlap returns the amount of overlap between two
// match rectangles.
// The overlap is the fraction of the smaller match
// that is covered by the other match.
func (m *Match) Overlap(m1 *Match) float64 {
	if m.X >= m1.X+m1.Width || m.X+m.Width <= m1.X ||
		m.Y >= m1.Y+m1.Height || m.Y+m.Height <= m1.Y {
		return 0
	}

	var minX, maxX int
	if m.X < m1.X {
		minX = m1.X
	} else {
		minX = m.X
	}
	if m.X+m.Width < m1.X+m1.Width {
		maxX = m.X + m.Width
	} else {
		maxX = m1.X + m1.Width
	}

	var minY, maxY int
	if m.Y < m1.Y {
		minY = m1.Y
	} else {
		minY = m.Y
	}
	if m.Y+m.Height < m1.Y+m1.Height {
		maxY = m.Y + m.Height
	} else {
		maxY = m1.Y + m1.Height
	}

	area := float64((maxX - minX) * (maxY - minY))
	if m.Width*m.Height > m1.Width*m1.Height {
		area /= float64(m1.Width * m1.Height)
	} else {
		area /= float64(m.Width * m.Height)
	}

	return area
}

// Matches is a slice of (possibly overlapping) matches.
type Matches []*Match

// Overlaps returns maximum overlap between m1 and any
// of the matches in m.
func (m Matches) MaxOverlap(m1 *Match) float64 {
	var max float64
	for _, match := range m {
		max = math.Max(max, match.Overlap(m1))
	}
	return max
}

// Unique produces a list of matches in which overlapping
// matches have been averaged together into one match.
//
// The threshold argument specifies how much overlap two
// matches must have between they are merged.
// Any overlap greater than threshold is considered enough
// to merge two images.
// An overlap of 0 joins any regions which overlap at all.
func (m Matches) JoinOverlaps(threshold float64) Matches {
	var clusters []Matches

	for _, match := range m {
		var overlaps []int
		for i, cluster := range clusters {
			if cluster.MaxOverlap(match) > threshold {
				overlaps = append(overlaps, i)
			}
		}
		if len(overlaps) == 0 {
			clusters = append(clusters, Matches{match})
		} else {
			first := overlaps[0]
			clusters[first] = append(clusters[first], match)
			for i := len(overlaps) - 1; i > 0; i-- {
				k := overlaps[i]
				clusters[first] = append(clusters[first], clusters[k]...)
				clusters[k] = clusters[len(clusters)-1]
				clusters = clusters[:len(clusters)-1]
			}
		}
	}

	res := make(Matches, len(clusters))
	for i, cluster := range clusters {
		res[i] = cluster.average()
	}
	return res
}

func (m Matches) average() *Match {
	var sum Match
	for _, match := range m {
		sum.X += match.X
		sum.Y += match.Y
		sum.Width += match.Width
		sum.Height += match.Height
	}
	sum.X /= len(m)
	sum.Y /= len(m)
	sum.Width /= len(m)
	sum.Height /= len(m)
	return &sum
}
