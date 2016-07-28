package haar

// A Classifier is a weighted sum of single-feature
// classifiers.
type Classifier struct {
	// Features is all the features.
	Features []Feature

	// Thresholds contains one threshold per feature.
	// If a feature returns a value greater than its
	// threshold, the feature's output is 1; otherwise,
	// the output is -1.
	Thresholds []float64

	// Weights contains one weight per feature.
	// The weight is applied after a feature has been
	// discretized into a 1 or -1.
	Weights []float64

	// If the weighted sum of the feature outputs
	// is greater than Threshold, a sample is
	// considered positive.
	Threshold float64
}

// Classify runs the classifier on an image.
// It returns true if the sample is positive.
func (c *Classifier) Classify(img IntegralImage) bool {
	var sum float64
	for i, feature := range c.Features {
		var output float64
		if feature.FeatureValue(img) > c.Thresholds[i] {
			output = 1
		} else {
			output = -1
		}
		sum += c.Weights[i] * output
	}
	return sum > c.Threshold
}

// A Cascade classifies images by running them through
// a series of classifiers, returning positive matches
// only if every classifier in the series returns
// positive.
type Cascade struct {
	Layers []*Classifier
}

// Classify classifies the given image by running it
// through the cascade.
// If the result is positive, this returns true.
func (c *Cascade) Classify(img IntegralImage) bool {
	for _, layer := range c.Layers {
		if !layer.Classify(img) {
			return false
		}
	}
	return true
}
