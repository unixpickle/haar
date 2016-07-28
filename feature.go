package haar

// A Feature is anything capable of computing a value
// for itself given an image window.
type Feature interface {
	// Bounds returns the rectangle in each window that
	// this feature takes into account.
	Bounds() (x, y, width, height int)

	// FeatureValue evaluates the feature in the given
	// window image.
	FeatureValue(img IntegralImage) float64
}

// AllFeatures builds a list of every haar-like feature
// that fits in the given window size.
func AllFeatures(width, height int) []Feature {
	var res []Feature
	for w := 1; w <= width; w++ {
		for h := 1; h <= height; h++ {
			if h == 1 && w == 1 {
				continue
			}
			for y := 0; y <= height-h; y++ {
				for x := 0; x <= width-w; x++ {
					rect := featureRect{x, y, w, h}
					if w%2 == 0 {
						res = append(res, &rectPair{rect, true})
					}
					if h%2 == 0 {
						res = append(res, &rectPair{rect, false})
					}
					if w%2 == 0 && h%2 == 0 {
						res = append(res, &diagonalRects{rect})
					}
					if w%3 == 0 {
						res = append(res, &tripleRects{rect, true})
					}
					if h%3 == 0 {
						res = append(res, &tripleRects{rect, false})
					}
				}
			}
		}
	}
	return res
}

type featureRect struct {
	X int
	Y int
	W int
	H int
}

func (f *featureRect) Bounds() (x, y, width, height int) {
	return f.X, f.Y, f.W, f.H
}

type rectPair struct {
	featureRect

	// If Horizontal is true, the two adjacent rects are
	// next to each other; otherwise, they are on top of
	// each other.
	Horizontal bool
}

func (r *rectPair) FeatureValue(img IntegralImage) float64 {
	var sum1, sum2 float64
	if r.Horizontal {
		midTop := img.IntegralAt(r.X+r.W/2, r.Y)
		midBottom := img.IntegralAt(r.X+r.W/2, r.Y+r.H)
		sum1 = midBottom + img.IntegralAt(r.X, r.Y) -
			(midTop + img.IntegralAt(r.X, r.Y+r.H))
		sum2 = midTop + img.IntegralAt(r.X+r.W, r.Y+r.H) -
			(midBottom + img.IntegralAt(r.X+r.W, r.Y))
	} else {
		midLeft := img.IntegralAt(r.X, r.Y+r.H/2)
		midRight := img.IntegralAt(r.X+r.W, r.Y+r.H/2)
		sum1 = midRight + img.IntegralAt(r.X, r.Y) -
			(midLeft + img.IntegralAt(r.X+r.W, r.Y))
		sum2 = midLeft + img.IntegralAt(r.X+r.W, r.Y+r.H) -
			(midRight + img.IntegralAt(r.X, r.Y+r.H))
	}
	return sum1 - sum2
}

type diagonalRects struct {
	featureRect
}

func (d *diagonalRects) FeatureValue(img IntegralImage) float64 {
	integralValues := [3][3]float64{
		{img.IntegralAt(d.X, d.Y), img.IntegralAt(d.X+d.W/2, d.Y),
			img.IntegralAt(d.X+d.W, d.Y)},
		{img.IntegralAt(d.X, d.Y+d.H/2), img.IntegralAt(d.X+d.W/2, d.Y+d.H/2),
			img.IntegralAt(d.X+d.W, d.Y+d.H/2)},
		{img.IntegralAt(d.X, d.Y+d.H), img.IntegralAt(d.X+d.W/2, d.Y+d.H),
			img.IntegralAt(d.X+d.W, d.Y+d.H)},
	}
	topLeft := integralValues[0][0] + integralValues[1][1] -
		(integralValues[1][0] + integralValues[0][1])
	topRight := integralValues[0][1] + integralValues[1][2] -
		(integralValues[1][1] + integralValues[0][2])
	bottomLeft := integralValues[1][0] + integralValues[2][1] -
		(integralValues[2][0] + integralValues[1][1])
	bottomRight := integralValues[1][1] + integralValues[2][2] -
		(integralValues[2][1] + integralValues[1][2])
	return topLeft + bottomRight - (topRight + bottomLeft)
}

type tripleRects struct {
	featureRect

	// If Horizontal is true, the three rectangles go
	// from left to right; otherwise, they go from top
	// to bottom.
	Horizontal bool
}

func (t *tripleRects) FeatureValue(img IntegralImage) float64 {
	var integralValues [2][4]float64
	if t.Horizontal {
		integralValues = [2][4]float64{
			{img.IntegralAt(t.X, t.Y), img.IntegralAt(t.X+t.W/3, t.Y),
				img.IntegralAt(t.X+2*t.W/3, t.Y), img.IntegralAt(t.X+t.W, t.Y)},
			{img.IntegralAt(t.X, t.Y+t.H), img.IntegralAt(t.X+t.W/3, t.Y+t.H),
				img.IntegralAt(t.X+2*t.W/3, t.Y+t.H), img.IntegralAt(t.X+t.W, t.Y+t.H)},
		}
	} else {
		integralValues = [2][4]float64{
			{img.IntegralAt(t.X, t.Y), img.IntegralAt(t.X, t.Y+t.H/3),
				img.IntegralAt(t.X, t.Y+2*t.H/3), img.IntegralAt(t.X, t.Y+t.H)},
			{img.IntegralAt(t.X+t.W, t.Y), img.IntegralAt(t.X+t.W, t.Y+t.H/3),
				img.IntegralAt(t.X+t.W, t.Y+2*t.H/3), img.IntegralAt(t.X+t.W, t.Y+t.H)},
		}
	}
	firstRect := integralValues[0][0] + integralValues[1][1] -
		(integralValues[0][1] + integralValues[1][0])
	secondRect := integralValues[0][1] + integralValues[1][2] -
		(integralValues[0][2] + integralValues[1][1])
	thirdRect := integralValues[0][2] + integralValues[1][3] -
		(integralValues[0][3] + integralValues[1][2])
	return firstRect + thirdRect - secondRect
}
