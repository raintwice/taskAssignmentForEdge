package common

import "math"

//for float64 value
const eps = 0.0000001
func IsValueEqual(f1, f2 float64) bool {
    return math.Abs(f1-f2) < eps
}