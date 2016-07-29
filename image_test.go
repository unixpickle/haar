package haar

import (
	"math"
	"testing"
)

var imageTestBitmap = []float64{
	0.862977, 0.575527, 0.108108, 0.613100, 0.139519, 0.669601, 0.191301,
	0.555602, 0.951677, 0.578089, 0.615650, 0.790867, 0.131685, 0.213610,
	0.803027, 0.242205, 0.248390, 0.117146, 0.457930, 0.832474, 0.379080,
	0.858395, 0.835004, 0.124126, 0.732274, 0.718383, 0.074130, 0.116138,
	0.608587, 0.653907, 0.322213, 0.247214, 0.559763, 0.465253, 0.275907,
	0.453643, 0.824009, 0.360547, 0.371745, 0.914475, 0.476705, 0.984499,
	0.687608, 0.192326, 0.620609, 0.846832, 0.137359, 0.231368, 0.725900,
}

var (
	imageTestBitmapWidth  = 7
	imageTestBitmapHeight = 7
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

func TestDualImage(t *testing.T) {
	bmp := BitmapIntegralImage(imageTestBitmap, imageTestBitmapWidth,
		imageTestBitmapHeight)
	cropping := NewDualImage(bmp).Window(2, 1, 3, 4)

	var sum float64
	var squareSum float64
	for x := 2; x < 5; x++ {
		for y := 1; y < 5; y++ {
			pixel := imageTestBitmap[x+imageTestBitmapWidth*y]
			sum += pixel
			squareSum += pixel * pixel
		}
	}

	mean := sum / 12.0
	stddev := math.Sqrt(squareSum/12 - mean*mean)

	for x := 0; x <= 3; x++ {
		for y := 0; y <= 4; y++ {
			area := float64((x + 2) * (y + 1))
			expected := (bmp.IntegralAt(x+2, y+1) - mean*area) / stddev
			actual := cropping.IntegralAt(x, y)
			if math.Abs(actual-expected) > 1e-5 {
				t.Errorf("at %d,%d expected %f but got %f", x, y, expected, actual)
			}
		}
	}

	for x := 0; x < 3; x++ {
		for y := 0; y < 4; y++ {
			actual := cropping.IntegralAt(x+1, y+1) + cropping.IntegralAt(x, y) -
				(cropping.IntegralAt(x, y+1) + cropping.IntegralAt(x+1, y))
			expected := (imageTestBitmap[(x+2)+imageTestBitmapWidth*(y+1)] - mean) / stddev
			if math.Abs(actual-expected) > 1e-5 {
				t.Errorf("pixel at %d,%d should be %f but it's %f", x, y, expected, actual)
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
