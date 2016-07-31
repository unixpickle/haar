package haar

import "testing"

func TestSortedBins(t *testing.T) {
	values := []float64{1, 3, 2, 3, 5, 0.5, 0.5, 4.5, 5}
	sorted := []float64{0.5, 1, 2, 3, 4.5, 5}
	idxs := []int{1, 3, 2, 3, 5, 0, 0, 4, 5}

	actual := newSortedBins(values)
	if len(actual.Sorted) != len(sorted) {
		t.Error("sorted had wrong length: got", len(actual.Sorted),
			"expected", len(sorted))
		return
	}
	if len(actual.Mapping) != len(idxs) {
		t.Error("mapping had wrong length: got", len(actual.Mapping),
			"expected", len(idxs))
		return
	}
	for i, x := range idxs {
		a := actual.Mapping[i]
		if a != x {
			t.Errorf("mapping %d: expected %d got %d", i, x, a)
		}
	}
	for i, x := range sorted {
		a := actual.Sorted[i]
		if a != x {
			t.Errorf("sorted %d: expected %f got %f", i, x, a)
		}
	}
}
