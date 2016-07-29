package haar

import "log"

// A Logger logs information about the training process.
type Logger interface {
	// LogStartingLayer logs that a new layer of the
	// cascade is being trained.
	LogStartingLayer(index int)

	// LogCreatedNegatives logs that the negative samples
	// for the current layer have been created.
	// The count argument specifies the number of negative
	// samples.
	LogCreatedNegatives(count int)

	// LogFeature logs that a feature has been added to
	// the current layer.
	// The numFeatures argument specifies how many
	// features are currently in this layer.
	// The retention and exclusion arguments indicate the
	// positive retention rate and the negative exclusion
	// rate, respectively.
	LogFeature(numFeatures int, retention, exclusion float64, f *Feature)
}

// A ConsoleLogger logs output using the log package.
type ConsoleLogger struct{}

func (_ ConsoleLogger) LogStartingLayer(index int) {
	log.Printf("Starting layer %d ...", index)
}

func (_ ConsoleLogger) LogCreatedNegatives(count int) {
	log.Printf("Created %d negatives.", count)
}

func (_ ConsoleLogger) LogFeature(numFeatures int, retention, exclusion float64, f *Feature) {
	log.Printf("Feature %d: retention=%f exclusion=%f type=%d", numFeatures,
		retention, exclusion, f.Type)
}
