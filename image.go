package haar

import "math"

// An IntegralImage is a grayscale image optimized for
// Haar-like feature computatiod.
type IntegralImage interface {
	// Width returns the width of the image in pixels.
	Width() int

	// Height returns the height of the image in pixels.
	Height() int

	// IntegralAt returns the integral of all the pixels
	// above and to the left of the given coordinate.
	//
	// Coordinates start at 0 and the point (0,0) refers
	// to the top-left pixel of the image.
	//
	// The integral around the top and left parts of the
	// image needn't be zero, since an image may be the
	// cropped version of another image.
	IntegralAt(x, y int) float64
}

// BitmapIntegralImage creates an IntegralImage from a
// grayscale bitmap image.
// The pixels in the bitmap should be packed going left
// to right, then top to bottom.
func BitmapIntegralImage(pixels []float64, width, height int) IntegralImage {
	if len(pixels) != width*height {
		panic("invalid bitmap size")
	}

	res := &sliceIntegralImage{
		integrals: make([]float64, width*height),
		width:     width,
		height:    height,
	}

	var idx int
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixel := pixels[idx]
			aboveLeft := res.IntegralAt(x, y)
			left := res.IntegralAt(x, y+1)
			above := res.IntegralAt(x+1, y)
			res.integrals[idx] = pixel + above + left - aboveLeft
			idx++
		}
	}

	return res
}

// A DualImage stores an image in such a way that it
// can be cropped and normalized efficiently.
type DualImage struct {
	// image is the underlying image, which has the
	// extra property that integrals around the top
	// and left edges are 0.
	image IntegralImage

	// squared is an IntegralImage for the image
	// computed by squaring the brightness values
	// in Image.
	squared IntegralImage
}

// NewDualImage creates a DualImage based on the data
// in an IntegralImage.
func NewDualImage(img IntegralImage) *DualImage {
	bitmap := make([]float64, img.Width()*img.Height())
	squaredBmp := make([]float64, img.Width()*img.Height())

	var idx int
	for y := 0; y < img.Height(); y++ {
		for x := 0; x < img.Width(); x++ {
			brightness := img.IntegralAt(x+1, y+1) + img.IntegralAt(x, y) -
				(img.IntegralAt(x, y+1) + img.IntegralAt(x+1, y))
			bitmap[idx] = brightness
			squaredBmp[idx] = brightness * brightness
			idx++
		}
	}

	return &DualImage{
		image:   BitmapIntegralImage(bitmap, img.Width(), img.Height()),
		squared: BitmapIntegralImage(squaredBmp, img.Width(), img.Height()),
	}
}

// Width returns the width of the underlying image.
func (d *DualImage) Width() int {
	return d.image.Width()
}

// Height returns the height of the underlying image.
func (d *DualImage) Height() int {
	return d.image.Height()
}

// Window returns a normalized, cropped version of the
// underlying image.
func (d *DualImage) Window(x, y, width, height int) IntegralImage {
	if x < 0 || y < 0 {
		panic("crop coordinates cannot be negative")
	}
	if x+width > d.image.Width() || y+height > d.image.Height() {
		panic("crop rectangle goes out of bounds")
	}
	area := float64(width * height)
	totalSum := d.image.IntegralAt(x+width, y+height) + d.image.IntegralAt(x, y) -
		(d.image.IntegralAt(x+width, y) + d.image.IntegralAt(x, y+height))
	squareSum := d.squared.IntegralAt(x+width, y+height) + d.squared.IntegralAt(x, y) -
		(d.squared.IntegralAt(x+width, y) + d.squared.IntegralAt(x, y+height))
	mean := totalSum / area
	return &croppedImage{
		img:    d.image,
		x:      x,
		y:      y,
		w:      width,
		h:      height,
		mean:   mean,
		stddev: math.Sqrt(squareSum/area - math.Pow(mean, 2)),
	}
}

type sliceIntegralImage struct {
	integrals []float64
	width     int
	height    int
}

func (s *sliceIntegralImage) Width() int {
	return s.width
}

func (s *sliceIntegralImage) Height() int {
	return s.height
}

func (s *sliceIntegralImage) IntegralAt(x, y int) float64 {
	if x <= 0 || y <= 0 {
		return 0
	}
	return s.integrals[(x-1)+s.width*(y-1)]
}

type croppedImage struct {
	img IntegralImage
	x   int
	y   int
	w   int
	h   int

	mean   float64
	stddev float64
}

func (c *croppedImage) Width() int {
	return c.w
}

func (c *croppedImage) Height() int {
	return c.h
}

func (c *croppedImage) IntegralAt(x, y int) float64 {
	area := float64((x + c.x) * (y + c.y))
	rawVal := c.img.IntegralAt(x+c.x, y+c.y)
	return (rawVal - area*c.mean) / c.stddev
}
