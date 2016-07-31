package haar

import (
	"math"
	"runtime"
	"sort"
	"sync"

	"github.com/unixpickle/num-analysis/linalg"
	"github.com/unixpickle/weakai/boosting"
)

const (
	trainingPositiveBias = 2
	trainingBatchSize    = 10
)

// Requirements stores minimum requirements for training
// a layer in a cascade.
type Requirements struct {
	// PositiveRetention is the minimum fraction of
	// positive samples this layer should return a
	// positive classification for.
	//
	// A good value of PositiveRetention is somewhere
	// in the high 0.9's, for example 0.99.
	PositiveRetention float64

	// NegativeExclusion is the minimum fraction of
	// negative samples for this layer to correctly
	// classify as negative.
	// An example value would be 0.5, which would be
	// reasonable for the first classifier in a
	// cascade.
	NegativeExclusion float64

	// MaxFeatures specifies the maximum number of
	// features to be used in this layer.
	// If MaxFeatures is exceeded before NegativeExclusion,
	// the layer should be used despite its sub-par
	// exclusion capability.
	MaxFeatures int
}

// Train trains a cascade classifier given the
// requirements for its layers.
//
// If the given Logger is nil, nothing will be logged.
//
// This may return fewer than the requested number of
// layers if all negative samples are dealt with.
func Train(layerReqs []*Requirements, s SampleSource, l Logger) *Cascade {
	var res Cascade

	TrainMore(&res, layerReqs, s, l)

	return &res
}

// TrainMore is like Train, but it adds layers onto an
// existing cascade instead of building a cascade from
// scratch.
//
// The list of requirements is for the layers added,
// so the number of requirements specifies how many
// layers to add, not the total number of layers in
// the final cascade.
func TrainMore(c *Cascade, addReqs []*Requirements, s SampleSource, l Logger) {
	positives := s.Positives()
	if len(c.Layers) > 0 {
		positives = acceptedPositives(positives, c)
	}
	if len(positives) == 0 {
		return
	}

	c.WindowWidth = positives[0].Width()
	c.WindowHeight = positives[0].Height()

	features := AllFeatures(positives[0].Width(), positives[0].Height())

	for _, reqs := range addReqs {
		if l != nil {
			l.LogStartingLayer(len(c.Layers))
		}
		var negs []IntegralImage
		if len(c.Layers) == 0 {
			negs = s.InitialNegatives()
		} else {
			negs = s.AdversarialNegatives(c)
		}
		if l != nil {
			l.LogCreatedNegatives(len(negs))
		}
		if len(negs) == 0 {
			break
		}
		layer := trainLayer(reqs, positives, negs, features, l)
		c.Layers = append(c.Layers, layer)
		positives = acceptedPositives(positives, layer)
	}
}

func trainLayer(reqs *Requirements, pos, neg []IntegralImage, features []*Feature,
	l Logger) *Layer {
	allSamples := make([]IntegralImage, len(pos)+len(neg))
	copy(allSamples, pos)
	copy(allSamples[len(pos):], neg)
	desired := make(linalg.Vector, len(allSamples))
	for i := range pos {
		desired[i] = 1
	}
	for i := range neg {
		desired[i+len(pos)] = -1
	}

	gradient := boosting.Gradient{
		Loss: &boosting.WeightedExpLoss{
			PosWeight: trainingPositiveBias * float64(len(neg)) / float64(len(pos)),
		},
		Desired: desired,
		List:    boostingSamples(allSamples),
		Pool:    &boostingPool{Features: features},
	}

	var threshold float64
	for i := 0; i < reqs.MaxFeatures; i++ {
		gradient.Step()
		threshold = necessaryThreshold(gradient.OutCache, desired, reqs.PositiveRetention)
		ret, exc := boostingScores(gradient.OutCache, desired, threshold)
		if l != nil {
			latestFeature := gradient.Sum.Classifiers[i].(*boostingClassifier).Feature
			if exc > 0 {
				l.LogFeature(i+1, ret, exc, latestFeature)
			} else {
				rawRet, rawExc := boostingScores(gradient.OutCache, desired, 0)
				l.LogFeature(i+1, rawRet, rawExc, latestFeature)
			}
		}
		if ret >= reqs.PositiveRetention && exc >= reqs.NegativeExclusion {
			break
		}
	}

	layer := &Layer{}
	for i, feature := range gradient.Sum.Classifiers {
		c := feature.(*boostingClassifier)
		weight := gradient.Sum.Weights[i]
		layer.Features = append(layer.Features, c.Feature)
		layer.Thresholds = append(layer.Thresholds, c.Threshold)
		layer.Weights = append(layer.Weights, weight)
	}
	layer.Threshold = threshold

	return layer
}

type boostingSamples []IntegralImage

func (b boostingSamples) Len() int {
	return len(b)
}

type boostingClassifier struct {
	Feature   *Feature
	Threshold float64
}

func (c *boostingClassifier) Classify(s boosting.SampleList) linalg.Vector {
	res := make(linalg.Vector, s.Len())
	for i, sample := range s.(boostingSamples) {
		output := c.Feature.Value(sample)
		if output > c.Threshold {
			res[i] = 1
		} else {
			res[i] = -1
		}
	}
	return res
}

type boostingOption struct {
	Classifier *boostingClassifier
	WeightDot  float64
}

type boostingPool struct {
	Features []*Feature

	bins []*sortedBins
}

func (b *boostingPool) BestClassifier(s boosting.SampleList, w linalg.Vector) boosting.Classifier {
	needBins := b.bins == nil
	if needBins {
		b.bins = make([]*sortedBins, len(b.Features))
	}

	featureChan := make(chan []int, len(b.Features)/trainingBatchSize+1)
	for i := 0; i < len(b.Features); i += trainingBatchSize {
		idxList := make([]int, 0, trainingBatchSize)
		for j := i; j < i+trainingBatchSize && j < len(b.Features); j++ {
			idxList = append(idxList, j)
		}
		featureChan <- idxList
	}
	close(featureChan)

	optionChan := make(chan boostingOption)

	var wg sync.WaitGroup
	for i := 0; i < runtime.GOMAXPROCS(0); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for featureIdxs := range featureChan {
				for _, featureIdx := range featureIdxs {
					feature := b.Features[featureIdx]
					if needBins {
						b.bins[featureIdx] = buildFeatureBins(feature,
							s.(boostingSamples))
					}
					optionChan <- bestFeatureSplit(feature, b.bins[featureIdx], w)
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(optionChan)
	}()

	var bestOption boostingOption
	for option := range optionChan {
		if math.Abs(option.WeightDot) >= math.Abs(bestOption.WeightDot) {
			bestOption = option
		}
	}

	return bestOption.Classifier
}

func buildFeatureBins(feature *Feature, s boostingSamples) *sortedBins {
	values := make([]float64, s.Len())
	for i, sample := range s {
		values[i] = feature.Value(sample)
	}
	return newSortedBins(values)
}

func bestFeatureSplit(feature *Feature, s *sortedBins, w linalg.Vector) boostingOption {
	var bestOption boostingOption

	weightSums := make([]float64, len(s.Sorted))
	for unsorted, sorted := range s.Mapping {
		weightSums[sorted] += w[unsorted]
	}

	bestOption.Classifier = &boostingClassifier{
		Feature:   feature,
		Threshold: s.Sorted[len(s.Sorted)-1],
	}

	// Start with the dot product where all outputs are
	// -1 because the threshold is high.
	for _, x := range w {
		bestOption.WeightDot -= x
	}

	// Compute a rolling dot product as the -1's in the
	// weak learner's output change to 1's.
	weightDot := bestOption.WeightDot
	for i := len(s.Sorted) - 1; i > 0; i-- {
		weightDot += 2 * weightSums[i]
		if math.Abs(weightDot) > math.Abs(bestOption.WeightDot) {
			bestOption.Classifier.Threshold = (s.Sorted[i-1] + s.Sorted[i]) / 2
			bestOption.WeightDot = weightDot
		}
	}

	return bestOption
}

func necessaryThreshold(boostOut, desired linalg.Vector, retention float64) float64 {
	var positiveOuts []float64
	var positiveCount int
	positiveCounts := map[float64]int{}

	for i, des := range desired {
		if des < 0 {
			continue
		}
		positiveCount++
		out := boostOut[i]
		if _, ok := positiveCounts[out]; !ok {
			positiveOuts = append(positiveOuts, out)
		}
		positiveCounts[out]++
	}

	sort.Float64s(positiveOuts)

	var res float64

	neededPositives := int(math.Ceil(retention * float64(positiveCount)))
	for i := len(positiveOuts) - 1; i > 0; i-- {
		neededPositives -= positiveCounts[positiveOuts[i]]
		if neededPositives <= 0 {
			res = (positiveOuts[i-1] + positiveOuts[i]) / 2
			break
		}
	}
	if neededPositives > 0 {
		res = math.Nextafter(positiveOuts[0], math.Inf(-1))
	}

	return math.Min(0, res)
}

func boostingScores(boostOut, desired linalg.Vector, thresh float64) (retention,
	exclusion float64) {
	var retained, positive int
	var excluded, negative int

	for i, des := range desired {
		if des > 0 {
			positive++
			if boostOut[i] > thresh {
				retained++
			}
		} else {
			negative++
			if boostOut[i] <= thresh {
				excluded++
			}
		}
	}

	retention = float64(retained) / float64(positive)
	exclusion = float64(excluded) / float64(negative)
	return
}

func acceptedPositives(pos []IntegralImage, c Classifier) []IntegralImage {
	res := make([]IntegralImage, 0, len(pos))
	for _, x := range pos {
		if c.Classify(x) {
			res = append(res, x)
		}
	}
	return res
}
