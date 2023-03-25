// edit from https://github.com/kabychow/go-cosinesimilarity

package math

import "math"

// Similarity len(x) must eq len(y)
func Similarity(x, y []uint8) float64 {
	var sum, s1, s2 uint64
	for i := 0; i < len(x); i++ {
		sum += uint64(x[i]) * uint64(y[i])
		s1 += uint64(x[i]) * uint64(x[i])
		s2 += uint64(y[i]) * uint64(y[i])
	}
	if s1 == 0 || s2 == 0 {
		return 0.0
	}
	return float64(sum) / (math.Sqrt(float64(s1)) * math.Sqrt(float64(s2)))
}
