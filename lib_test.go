package fire

import (
	"testing"

	"github.com/gonum/matrix/mat64"
)

var (
	b1 = []float64{0.0, 0.0, 0.0, 1.0}
	b2 = []float64{0.25, 0.25, 0.25, 0.25}
)

func TestCos(t *testing.T) {
	v1 := mat64.NewVector(4, b1)
	v2 := mat64.NewVector(4, b2)
	p := cos(v1, v2)
	t.Log(p)
}
