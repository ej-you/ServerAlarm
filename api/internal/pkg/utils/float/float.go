package float

import "math"

const eps = 1e-4

func Equal(a, b float32) bool {
	return math.Abs(float64(a-b)) < eps
}
