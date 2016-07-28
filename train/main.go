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
			NegativeExclusion: 0.4,
			MaxFeatures:       10,
		},
		{
			PositiveRetention: 0.99,
			NegativeExclusion: 0.4,
			MaxFeatures:       20,
		},
		{
			PositiveRetention: 0.99,
			NegativeExclusion: 0.4,
			MaxFeatures:       100,
		},
	}

	haar.Train(reqs, samples, haar.ConsoleLogger{})

	// TODO: save results here.
}
