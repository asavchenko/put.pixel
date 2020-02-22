package mlib

import (
	"math/rand"
	"time"
)

var r1 *rand.Rand

func init() {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 = rand.New(s1)
}

/**
 *
 */
func Rand(n int) int {
	if n <= 0 {
		n = 1
	}

	return r1.Intn(n) + 1
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
func AbsInt(x int) int {
	if x >= 0 {
		return x
	}

	return -x
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
