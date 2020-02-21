package mlib

import "math/rand"

/**
 *
 */
func Rand(n int) int {
	if n <= 0 {
		n = 1
	}

	return rand.Intn(SignInt(n)*n) + 1
}

/**
 *
 */
func Sign(x float64) int {
	if x >= 0 {
		return 1
	}

	return -1
}

/**
 *
 */
func SignInt(x int) int {
	if x >= 0 {
		return 1
	}

	return -1
}

/**
 *
 */
func Srand(n int) int {
	if Rand(2) == 1 {
		return -Rand(n)
	}

	return Rand(n)
}
