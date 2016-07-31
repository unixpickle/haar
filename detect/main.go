// Command detect runs a cascade on an image.
package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"io/ioutil"
	"math"
	"os"

	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/unixpickle/haar"
)

const OverlapThreshold = 0.7

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintf(os.Stderr, "Usage: %s cascade.json input.png output.png\n", os.Args[0])
		os.Exit(1)
	}

	cascadeData, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to read cascade:", err)
		os.Exit(1)
	}
	var cascade haar.Cascade
	if err := json.Unmarshal(cascadeData, &cascade); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to unmarshal cascade:", err)
		os.Exit(1)
	}

	f, err := os.Open(os.Args[2])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to open image file:", err)
		os.Exit(1)
	}
	img, _, err := image.Decode(f)
	f.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to decode image:", err)
		os.Exit(1)
	}

	intImg := haar.ImageIntegralImage(img)
	matches := cascade.Scan(haar.NewDualImage(intImg), 0, 0)
	matches = matches.JoinOverlaps(OverlapThreshold)

	output, err := os.Create(os.Args[3])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to create output file:", err)
		os.Exit(1)
	}

	err = png.Encode(output, AnnotateImage(img, matches))
	output.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to save output file:", err)
		os.Exit(1)
	}
}

func AnnotateImage(img image.Image, matches haar.Matches) image.Image {
	dest := image.NewRGBA(image.Rect(0, 0, img.Bounds().Dx(), img.Bounds().Dy()))
	ctx := draw2dimg.NewGraphicContext(dest)

	ctx.DrawImage(img)
	ctx.SetStrokeColor(color.RGBA{R: 0xff, A: 0xff})
	ctx.SetLineWidth(math.Max(1, float64(img.Bounds().Dx())/500))
	for _, match := range matches {
		ctx.BeginPath()
		ctx.MoveTo(float64(match.X), float64(match.Y))
		ctx.LineTo(float64(match.X+match.Width), float64(match.Y))
		ctx.LineTo(float64(match.X+match.Width), float64(match.Y+match.Height))
		ctx.LineTo(float64(match.X), float64(match.Y+match.Height))
		ctx.Close()
		ctx.Stroke()
	}
	return dest
}
