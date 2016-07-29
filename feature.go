package haar

import "fmt"

type FeatureType int

const (
	HorizontalPair FeatureType = iota
	VerticalPair
	HorizontalTriple
	VerticalTriple
	Diagonal
)

// A Feature is a Haar-like feature.
// Features are computed by subtracting the pixel sums
// in some rectangles from those in others.
type Feature struct {
	Type FeatureType

	// These coordinates define the bounding box of the
	// feature inside the window.
	X      int
	Y      int
	Width  int
	Height int
}

// AllFeatures builds a list of every haar-like feature
// that fits in the given window size.
func AllFeatures(width, height int) []*Feature {
	var res []*Feature
	for w := 1; w <= width; w++ {
		for h := 1; h <= height; h++ {
			if h == 1 && w == 1 {
				continue
			}
			for y := 0; y <= height-h; y++ {
				for x := 0; x <= width-w; x++ {
					f := Feature{X: x, Y: y, Width: w, Height: h}
					if w%2 == 0 {
						hf := f
						hf.Type = HorizontalPair
						res = append(res, &hf)
					}
					if h%2 == 0 {
						vf := f
						vf.Type = VerticalPair
						res = append(res, &vf)
					}
					if w%2 == 0 && h%2 == 0 {
						df := f
						df.Type = Diagonal
						res = append(res, &df)
					}
					if w%3 == 0 {
						hf := f
						hf.Type = HorizontalTriple
						res = append(res, &hf)
					}
					if h%3 == 0 {
						vf := f
						vf.Type = VerticalTriple
						res = append(res, &vf)
					}
				}
			}
		}
	}
	return res
}

// Value evaluates the feature on the given window.
func (f *Feature) Value(img IntegralImage) float64 {
	switch f.Type {
	case HorizontalPair, VerticalPair:
		return f.pair(img, f.Type == HorizontalPair)
	case HorizontalTriple, VerticalTriple:
		return f.triple(img, f.Type == HorizontalTriple)
	case Diagonal:
		return f.diagonal(img)
	default:
		panic(fmt.Sprintf("unknown feature type: %d", f.Type))
	}
}

func (f *Feature) pair(img IntegralImage, horizontal bool) float64 {
	var sum1, sum2 float64
	if horizontal {
		midTop := img.IntegralAt(f.X+f.Width/2, f.Y)
		midBottom := img.IntegralAt(f.X+f.Width/2, f.Y+f.Height)
		sum1 = midBottom + img.IntegralAt(f.X, f.Y) -
			(midTop + img.IntegralAt(f.X, f.Y+f.Height))
		sum2 = midTop + img.IntegralAt(f.X+f.Width, f.Y+f.Height) -
			(midBottom + img.IntegralAt(f.X+f.Width, f.Y))
	} else {
		midLeft := img.IntegralAt(f.X, f.Y+f.Height/2)
		midRight := img.IntegralAt(f.X+f.Width, f.Y+f.Height/2)
		sum1 = midRight + img.IntegralAt(f.X, f.Y) -
			(midLeft + img.IntegralAt(f.X+f.Width, f.Y))
		sum2 = midLeft + img.IntegralAt(f.X+f.Width, f.Y+f.Height) -
			(midRight + img.IntegralAt(f.X, f.Y+f.Height))
	}
	return sum1 - sum2
}

func (f *Feature) triple(img IntegralImage, horizontal bool) float64 {
	var integralValues [2][4]float64
	if horizontal {
		integralValues = [2][4]float64{
			{img.IntegralAt(f.X, f.Y), img.IntegralAt(f.X+f.Width/3, f.Y),
				img.IntegralAt(f.X+2*f.Width/3, f.Y), img.IntegralAt(f.X+f.Width, f.Y)},
			{img.IntegralAt(f.X, f.Y+f.Height), img.IntegralAt(f.X+f.Width/3, f.Y+f.Height),
				img.IntegralAt(f.X+2*f.Width/3, f.Y+f.Height), img.IntegralAt(f.X+f.Width, f.Y+f.Height)},
		}
	} else {
		integralValues = [2][4]float64{
			{img.IntegralAt(f.X, f.Y), img.IntegralAt(f.X, f.Y+f.Height/3),
				img.IntegralAt(f.X, f.Y+2*f.Height/3), img.IntegralAt(f.X, f.Y+f.Height)},
			{img.IntegralAt(f.X+f.Width, f.Y), img.IntegralAt(f.X+f.Width, f.Y+f.Height/3),
				img.IntegralAt(f.X+f.Width, f.Y+2*f.Height/3), img.IntegralAt(f.X+f.Width, f.Y+f.Height)},
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

func (f *Feature) diagonal(img IntegralImage) float64 {
	integralValues := [3][3]float64{
		{img.IntegralAt(f.X, f.Y), img.IntegralAt(f.X+f.Width/2, f.Y),
			img.IntegralAt(f.X+f.Width, f.Y)},
		{img.IntegralAt(f.X, f.Y+f.Height/2), img.IntegralAt(f.X+f.Width/2, f.Y+f.Height/2),
			img.IntegralAt(f.X+f.Width, f.Y+f.Height/2)},
		{img.IntegralAt(f.X, f.Y+f.Height), img.IntegralAt(f.X+f.Width/2, f.Y+f.Height),
			img.IntegralAt(f.X+f.Width, f.Y+f.Height)},
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
