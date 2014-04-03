// Package newmath is a trivial example package.
package newmath

// Sqrt returns an approximation to the square root of x.
func Sqrt(x float64) float64 {
	z := 1.0
	for i := 0; i < 1000; i++ {
		// Can see both exported and non-exported symbols from mul.go, because
		// it's all in the same package.
		z -= minus(z*z, x) / Mul(2, z)
	}
	return z
}
