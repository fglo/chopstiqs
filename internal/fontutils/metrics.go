package fontutils

import (
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type Metrics struct {
	Ascent  int
	Descent int
	Height  int
}

func NewMetrics(fontMetrics font.Metrics) Metrics {
	metrics := Metrics{
		Ascent:  fontMetrics.Ascent.Floor(),
		Descent: fontMetrics.Descent.Floor(),
		Height:  fontMetrics.Height.Floor(),
	}

	return metrics
}

func FixedInt26_6ToFloat64(i fixed.Int26_6) float64 {
	return float64(i) / (1 << 6)
}
