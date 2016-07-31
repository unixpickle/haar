package main

import (
	"encoding/json"
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/unixpickle/haar"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s cascade_file image_dir\n", os.Args[0])
		os.Exit(1)
	}

	cascadeData, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to read cascade:", err)
		os.Exit(1)
	}
	var cascade haar.Cascade
	if err := json.Unmarshal(cascadeData, &cascade); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to parse cascade:", err)
		os.Exit(1)
	}

	listing, err := ioutil.ReadDir(os.Args[2])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to read directory:", err)
		os.Exit(1)
	}

	log.Printf("Running %d-layer cascade on images...", len(cascade.Layers))

	var totalFeatures int
	var totalRuns int
	for _, item := range listing {
		if strings.HasPrefix(item.Name(), ".") {
			continue
		}
		path := filepath.Join(os.Args[2], item.Name())
		f, err := os.Open(path)
		if err != nil {
			log.Printf("Failed to read: %s: %s", path, err)
			continue
		}
		img, _, err := image.Decode(f)
		f.Close()
		if err != nil {
			log.Printf("Failed to decode: %s: %s", path, err)
			continue
		}
		intImg := haar.ImageIntegralImage(img)
		dualImg := haar.NewDualImage(intImg)

		for y := 0; y <= intImg.Height()-cascade.WindowHeight; y++ {
			for x := 0; x <= intImg.Width()-cascade.WindowWidth; x++ {
				cropped := dualImg.Window(x, y, cascade.WindowWidth,
					cascade.WindowHeight)
				totalFeatures += runCascade(cascade, cropped)
				totalRuns++
			}
		}
	}

	log.Println("Averaged", float64(totalFeatures)/float64(totalRuns), "features/run.")
}

func runCascade(c haar.Cascade, img haar.IntegralImage) int {
	var count int
	for _, layer := range c.Layers {
		count += len(layer.Features)
		if !layer.Classify(img) {
			break
		}
	}
	return count
}
