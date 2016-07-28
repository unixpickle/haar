package haar

// A SampleSource provides images for use while training
// cascade classifiers.
type SampleSource interface {
	// Positives returns positive training samples.
	Positives() []IntegralImage

	// InitialNegatives returns negative training samples
	// to use for training the first layer of a cascade.
	InitialNegatives() []IntegralImage

	// AdversarialNegatives returns negative training
	// samples which fool the existing cascade.
	AdversarialNegatives(c *Cascade) []IntegralImage
}
