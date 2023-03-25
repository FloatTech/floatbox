package math

import (
	"math"
	"testing"
)

func TestSimilarity(t *testing.T) {
	r := Similarity([]uint8{1, 2, 3}, []uint8{1, 3, 4})
	t.Log(r)
	if math.Abs(r-0.9958705948858224) > 1e-6 {
		t.Fail()
	}
	r = Similarity([]uint8{3, 2, 1}, []uint8{1, 3, 4})
	t.Log(r)
	if math.Abs(r-0.6813851438692469) > 1e-6 {
		t.Fail()
	}
}
