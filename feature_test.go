package haar

import (
	"math"
	"testing"
)

type featureTest struct {
	Desc     string
	Feature  Feature
	Expected float64
}

func TestBuiltinFeatures(t *testing.T) {
	img := featureTestImage()

	tests := []featureTest{
		{
			Desc:     "horizontal pair",
			Feature:  &rectPair{featureRect{1, 1, 4, 2}, true},
			Expected: rectangleSum(img, 1, 1, 2, 2) - rectangleSum(img, 3, 1, 2, 2),
		},
		{
			Desc:     "horizontal pair (short)",
			Feature:  &rectPair{featureRect{1, 3, 4, 1}, true},
			Expected: rectangleSum(img, 1, 3, 2, 1) - rectangleSum(img, 3, 3, 2, 1),
		},
		{
			Desc:     "vertical pair",
			Feature:  &rectPair{featureRect{1, 1, 2, 4}, false},
			Expected: rectangleSum(img, 1, 1, 2, 2) - rectangleSum(img, 1, 3, 2, 2),
		},
		{
			Desc:     "vertical pair (thin)",
			Feature:  &rectPair{featureRect{2, 0, 1, 4}, false},
			Expected: rectangleSum(img, 2, 0, 1, 2) - rectangleSum(img, 2, 2, 1, 2),
		},
		{
			Desc:    "diagonal",
			Feature: &diagonalRects{featureRect{1, 1, 4, 4}},
			Expected: rectangleSum(img, 1, 1, 2, 2) + rectangleSum(img, 3, 3, 2, 2) -
				(rectangleSum(img, 1, 3, 2, 2) + rectangleSum(img, 3, 1, 2, 2)),
		},
		{
			Desc:    "diagonal (short)",
			Feature: &diagonalRects{featureRect{1, 1, 4, 2}},
			Expected: rectangleSum(img, 1, 1, 2, 1) + rectangleSum(img, 3, 2, 2, 1) -
				(rectangleSum(img, 1, 2, 2, 1) + rectangleSum(img, 3, 1, 2, 1)),
		},
		{
			Desc:    "diagonal (thin)",
			Feature: &diagonalRects{featureRect{1, 1, 2, 4}},
			Expected: rectangleSum(img, 1, 1, 1, 2) + rectangleSum(img, 2, 3, 1, 2) -
				(rectangleSum(img, 1, 3, 1, 2) + rectangleSum(img, 2, 1, 1, 2)),
		},
		{
			Desc:    "horizontal triple",
			Feature: &tripleRects{featureRect{1, 1, 6, 3}, true},
			Expected: rectangleSum(img, 1, 1, 2, 3) + rectangleSum(img, 5, 1, 2, 3) -
				rectangleSum(img, 3, 1, 2, 3),
		},
		{
			Desc:    "vertical triple",
			Feature: &tripleRects{featureRect{1, 1, 3, 6}, false},
			Expected: rectangleSum(img, 1, 1, 3, 2) + rectangleSum(img, 1, 5, 3, 2) -
				rectangleSum(img, 1, 3, 3, 2),
		},
	}
	for _, test := range tests {
		actual := test.Feature.FeatureValue(img)
		if math.Abs(test.Expected-actual) > 1e-5 {
			t.Errorf("%s: expected %f got %f", test.Desc, test.Expected, actual)
		}
	}
}

func featureTestImage() IntegralImage {
	return BitmapIntegralImage(imageTestBitmap, imageTestBitmapWidth,
		imageTestBitmapHeight)
}

func rectangleSum(img IntegralImage, x, y, w, h int) float64 {
	return img.IntegralAt(x+w, y+h) + img.IntegralAt(x, y) -
		(img.IntegralAt(x+w, y) + img.IntegralAt(x, y+h))
}
