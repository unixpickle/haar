package main

import (
	"fmt"
	"log"
	"os"

	"github.com/unixpickle/haar"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintf(os.Stderr, "Usage: %s pos_dir neg_dir output_file\n",
			os.Args[0])
		os.Exit(1)
	}

	log.Println("Loading samples ...")

	posDir, negDir := os.Args[1], os.Args[2]
	samples, err := haar.LoadSampleSource(posDir, negDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to load samples:", err)
		os.Exit(1)
	}

	reqs := []*haar.Requirements{
		{
			PositiveRetention: 0.99,
			NegativeExclusion: 0.6,
			MaxFeatures:       3,
		},
		{
			PositiveRetention: 0.99,
			NegativeExclusion: 0.5,
			MaxFeatures:       10,
		},
		{
			PositiveRetention: 0.99,
			NegativeExclusion: 0.6,
			MaxFeatures:       25,
		},
		{
			PositiveRetention: 0.99,
			NegativeExclusion: 0.6,
			MaxFeatures:       50,
		},
	}
	for i := 0; i < 3; i++ {
		reqs = append(reqs, &haar.Requirements{
			PositiveRetention: 1,
			NegativeExclusion: 0.6,
			MaxFeatures:       100,
		})
	}

	haar.Train(reqs, samples, haar.ConsoleLogger{})

	// TODO: save results here.
}
