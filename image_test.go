package haar

import (
	"math"
	"testing"
)

var imageTestBitmap = []float64{
	0.286735, 0.661007, 0.666167, 0.156079, 0.800531, 0.290782,
	0.505438, 0.417589, 0.419909, 0.044689, 0.425667, 0.026518,
	0.871354, 0.876751, 0.280542, 0.063563, 0.323771, 0.080933,
	0.105769, 0.749114, 0.074676, 0.136882, 0.602447, 0.778182,
	0.032917, 0.304980, 0.902359, 0.214043, 0.232776, 0.945587,
	0.256449, 0.533534, 0.147162, 0.271161, 0.899616, 0.232532,
}

var (
	imageTestBitmapWidth  = 6
	imageTestBitmapHeight = 6
)

func TestBitmapIntegralImage(t *testing.T) {
	bmp := BitmapIntegralImage(imageTestBitmap, imageTestBitmapWidth,
		imageTestBitmapHeight)
	for y := 0; y <= imageTestBitmapHeight; y++ {
		for x := 0; x <= imageTestBitmapWidth; x++ {
			actual := bmp.IntegralAt(x, y)
			expected := imageTestIntegral(x, y)
			if math.Abs(actual-expected) > 1e-5 {
				t.Errorf("at %d,%d expected %f but got %f", x, y, expected, actual)
			}
		}
	}
}

func imageTestIntegral(x, y int) float64 {
	var sum float64
	for stepX := 0; stepX < x; stepX++ {
		for stepY := 0; stepY < y; stepY++ {
			sum += imageTestBitmap[stepX+stepY*imageTestBitmapWidth]
		}
	}
	return sum
}
