package haar

// These are defaults recommended in the original
// Viola-Jones paper on face detection.
const (
	DefaultScanScale  = 1.25
	DefaultScanStride = 1
)

// A Classifier classifies image windows.
type Classifier interface {
	// Classify returns whether an image is positive (true)
	// or negative (false).
	Classify(img IntegralImage) bool
}

// A Layer is a weighted sum of single-feature
// classifiers.
type Layer struct {
	// Features is all the features.
	Features []*Feature

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

// Sum returns the weighted sum of all the feature
// outputs when run on an image window.
func (c *Layer) Sum(img IntegralImage) float64 {
	var sum float64
	for i, feature := range c.Features {
		var output float64
		if feature.Value(img) > c.Thresholds[i] {
			output = 1
		} else {
			output = -1
		}
		sum += c.Weights[i] * output
	}
	return sum
}

// Classify runs the layer on an image.
// It returns true if the sample is positive.
func (c *Layer) Classify(img IntegralImage) bool {
	return c.Sum(img) > c.Threshold
}

// A Cascade classifies images by running them through
// a series of classifiers, returning positive matches
// only if every classifier in the series returns
// positive.
type Cascade struct {
	Layers       []*Layer
	WindowWidth  int
	WindowHeight int
}

// Classify classifies the given window by running it
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

// Scan looks for instances of this cascade within an
// entire image.
//
// The cascade is scaled by various powers of scale to
// find objects of different sizes.
// If scale is 0, DefaultScanScale is used.
//
// The cascade is moved horizontally and vertically
// during scanning according to the stride argument.
// If the stride is 0, DefaultScanStride is used.
// The stride specifies how many pixels (relative to the
// size of the cascade) to move the cascade for each
// iteration of the scanning process.
// A value of 1 means that the unscaled cascade is moved
// one pixel at a time.
//
// The result may contain overlapping matches.
func (c *Cascade) Scan(img *DualImage, scale, stride float64) Matches {
	if scale == 0 {
		scale = DefaultScanScale
	}
	if stride == 0 {
		stride = DefaultScanStride
	}

	var res Matches
	curScale := 1.0
	for {
		cropWidth := int(float64(c.WindowWidth)*curScale + 0.5)
		cropHeight := int(float64(c.WindowHeight)*curScale + 0.5)
		if cropWidth > img.Width() || cropHeight > img.Height() {
			break
		}

		curStride := curScale * stride

		for y := 0.0; int(y) <= img.Height()-cropHeight; y += curStride {
			for x := 0.0; int(x) <= img.Width()-cropWidth; x += curStride {
				cropping := img.Window(int(x), int(y), cropWidth, cropHeight)
				if curScale != 1 {
					cropping = ScaleIntegralImage(cropping, c.WindowWidth, c.WindowHeight)
				}
				if c.Classify(cropping) {
					res = append(res, &Match{
						X:      int(x),
						Y:      int(y),
						Width:  cropWidth,
						Height: cropHeight,
					})
				}
			}
		}

		curScale *= scale
	}

	return res
}
