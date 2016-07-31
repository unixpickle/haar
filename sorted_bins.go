package haar

import "sort"

// sortedBins maps elements of an unsorted list of
// floats to a sorted list of the floats with duplicates
// removed.
type sortedBins struct {
	// Sorted is the sorted list of floats with no
	// duplicates.
	Sorted []float64

	// Mapping maps indices in the original list to
	// indices in the sorted list.
	// Multiple elements in the mapping may point to
	// the same element of Sorted.
	Mapping []int
}

func newSortedBins(list []float64) *sortedBins {
	origMap := map[float64][]int{}
	for i, element := range list {
		origMap[element] = append(origMap[element], i)
	}

	keys := make([]float64, 0, len(origMap))
	for f := range origMap {
		keys = append(keys, f)
	}
	sort.Float64s(keys)

	res := &sortedBins{
		Sorted:  keys,
		Mapping: make([]int, len(list)),
	}

	for sortedIdx, key := range keys {
		idxList := origMap[key]
		for _, unsortedIdx := range idxList {
			res.Mapping[unsortedIdx] = sortedIdx
		}
	}

	return res
}
