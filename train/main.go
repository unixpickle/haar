// Command train trains a new Haar cascade.
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/unixpickle/haar"
)

const defaultInitialRetention = 0.99

func main() {
	if len(os.Args) != 4 && len(os.Args) != 5 {
		fmt.Fprintf(os.Stderr, "Usage: %s pos_dir neg_dir output_file [initial_retention]\n",
			os.Args[0])
		os.Exit(1)
	}

	initialRetention := defaultInitialRetention
	if len(os.Args) == 5 {
		var err error
		initialRetention, err = strconv.ParseFloat(os.Args[4], 64)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Invalid initial retention:", os.Args[4])
			os.Exit(1)
		}
	}

	log.Println("Loading samples ...")

	posDir, negDir := os.Args[1], os.Args[2]
	samples, err := haar.LoadSampleSource(posDir, negDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to load samples:", err)
		os.Exit(1)
	}

	var reqs []*haar.Requirements
	for i := 0; i < 4; i++ {
		reqs = append(reqs, &haar.Requirements{
			PositiveRetention: 0.995,
			NegativeExclusion: 0.6,
			MaxFeatures:       100,
		})
	}
	for i := 0; i < 7; i++ {
		reqs = append(reqs, &haar.Requirements{
			PositiveRetention: 1,
			NegativeExclusion: 0.8,
			MaxFeatures:       100,
		})
	}
	reqs[0].PositiveRetention = initialRetention

	cascade := haar.Train(reqs, samples, haar.ConsoleLogger{})

	data, err := json.Marshal(cascade)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to marshal data:", err)
		os.Exit(1)
	}
	if err := ioutil.WriteFile(os.Args[3], data, 0755); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to write file:", err)
		os.Exit(1)
	}
}
