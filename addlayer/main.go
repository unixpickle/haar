// Command addlayer adds a layer to a Haar cascade.
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

const MaxFeatures = 1000

func main() {
	if len(os.Args) != 6 {
		fmt.Fprintf(os.Stderr, "Usage: %s pos_dir neg_dir cascade_file retention exclusion\n",
			os.Args[0])
		os.Exit(1)
	}

	retention, err := strconv.ParseFloat(os.Args[4], 64)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Invalid retention:", os.Args[4])
		os.Exit(1)
	}
	exclusion, err := strconv.ParseFloat(os.Args[5], 64)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Invalid exclusion:", os.Args[4])
		os.Exit(1)
	}

	cascadeData, err := ioutil.ReadFile(os.Args[3])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to read cascade:", err)
		os.Exit(1)
	}
	var cascade haar.Cascade
	if err := json.Unmarshal(cascadeData, &cascade); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to parse cascade:", err)
		os.Exit(1)
	}

	log.Println("Loading samples ...")

	posDir, negDir := os.Args[1], os.Args[2]
	samples, err := haar.LoadSampleSource(posDir, negDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to load samples:", err)
		os.Exit(1)
	}

	log.Println("Adding layer ...")

	if len(samples.Positives()) > 0 {
		pos := samples.Positives()[0]
		if pos.Width() != cascade.WindowWidth || pos.Height() != cascade.WindowHeight {
			fmt.Fprintf(os.Stderr, "Positive size should be %dx%d but got %dx%d\n",
				cascade.WindowWidth, cascade.WindowHeight,
				pos.Width(), pos.Height())
			os.Exit(1)
		}
	}

	reqs := []*haar.Requirements{{
		PositiveRetention: retention,
		NegativeExclusion: exclusion,
		MaxFeatures:       MaxFeatures,
	}}
	haar.TrainMore(&cascade, reqs, samples, haar.ConsoleLogger{})

	data, err := json.Marshal(&cascade)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to marshal data:", err)
		os.Exit(1)
	}
	if err := ioutil.WriteFile(os.Args[3], data, 0755); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to write file:", err)
		os.Exit(1)
	}
}
