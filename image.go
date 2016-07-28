package haar

// An IntegralImage is a grayscale image optimized for
// Haar-like feature computation.
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
