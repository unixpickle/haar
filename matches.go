package haar

import "fmt"

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

// Overlaps returns true if m overlaps m1.
func (m *Match) Overlaps(m1 *Match) bool {
	return !(m.X >= m1.X+m1.Width || m.X+m.Width <= m1.X ||
		m.Y >= m1.Y+m1.Height || m.Y+m.Height <= m1.Y)
}

// Matches is a slice of (possibly overlapping) matches.
type Matches []*Match

// Overlaps returns true if m1 overlaps any matches in m.
func (m Matches) Overlaps(m1 *Match) bool {
	for _, match := range m {
		if match.Overlaps(m1) {
			return true
		}
	}
	return false
}

// Unique produces a list of matches in which overlapping
// matches have been averaged together into one match.
func (m Matches) JoinOverlaps() Matches {
	var clusters []Matches

	for _, match := range m {
		var overlaps []int
		for i, cluster := range clusters {
			if cluster.Overlaps(match) {
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
