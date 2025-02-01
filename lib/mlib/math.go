package mlib

import (
	"math/rand"
	"time"
)

type Point2DInt struct {
	X int
	Y int
}

type Point2DFloat32 struct {
	X float32
	Y float32
}

type Point2DFloat64 struct {
	X float64
	Y float64
}

type Point2D struct {
	X float64
	Y float64
}

var r1 *rand.Rand

func init() {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 = rand.New(s1)
}

func Rand(n int) int {
	if n < 0 {
		return -r1.Intn(-n)
	}
	if n == 0 {
		return 0
	}
	return r1.Intn(n + 1)
}

func GetRandomBtw(n1, n2 int) int {
	if n1 > n2 {
		n2 = -n1 + n2
		n1 = n2 + n1
		n2 = -n2 + n1
	}

	res := Rand(n2-n1) + n1
	if res < n1 {
		return n1
	}
	if res > n2 {
		return n2
	}

	return res
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
func Srand(n int) int {
	if Rand(2) == 1 {
		return -Rand(n)
	}

	return Rand(n)
}

func AbsInt(n int) int {
	return Sign(float64(n)) * n
}

func SignInt(n int) int {
	return Sign(float64(n))
}

/*
*

	.A_____________________________.B

	           .C
*/
func GetBezierCoords2(a, b Point2D, c Point2D, numberOfSteps int) []Point2D {
	result := make([]Point2D, 0)
	// x = a.x * (1 - t)^2 + 2 * c.x * (1 - t) * t + b.x * t^2
	// y = a.y * (1 - t)^2 + 2 * c.y * (1 - t) * t + b.y * t^2,
	d := float64(1) / float64(numberOfSteps)
	for i := 0; i < numberOfSteps; i++ {
		t := float64(i) * d
		//result = append(result, Point2D{
		//	x: a.x*(1-t)*(1-t) + 2*c.x*(1-t)*t + b.x*t*t,
		//	y: a.y*(1-t)*(1-t) + 2*c.y*(1-t)*t + b.y*t*t,
		//})
		t1 := (1 - t) * (1 - t)
		t2 := t * t
		t3 := 2 * (1 - t) * t
		p := Point2D{
			X: a.X*t1 + c.X*t3 + b.X*t2,
			Y: a.Y*t1 + c.Y*t3 + b.Y*t2,
		}
		//fmt.Println(p.x, p.y, t)
		result = append(result, p)
	}
	result = append(result, b)

	return result
}

/*
*

	.A_____________________________.B

	           .C      .D
*/
func GetBezierCoords3(a, b Point2D, c, d Point2D, numberOfSteps int) []Point2D {
	result := make([]Point2D, 0)
	// x = a.x * (1 - t)^3 + 3 * c.x * (1 - t)^2 * t + d.x * 3 * (1-t) * t^2 + b.x * t^3
	// y = a.y * (1 - t)^3 + 3 * c.y * (1 - t)^2 * t + d.y * 3 * (1-t) * t^2 + b.y * t^3
	e := float64(1 / numberOfSteps)
	for i := 0; i < numberOfSteps; i++ {
		t := float64(i) * e
		//result = append(result, Point2D{
		//	x: a.x*(1-t)*(1-t)*(1-t) + 3*c.x*(1-t)*(1-t)*t + d.x*3*(1-t)*t*t + b.x*t*t*t,
		//	y: a.y*(1-t)*(1-t)*(1-t) + 3*c.y*(1-t)*(1-t)*t + d.y*3*(1-t)*t*t + b.y*t*t*t,
		//})

		//t1 := (1 - t) * (1 - t) * (1 - t)
		//t2 := 3 * (1 - t) * (1 - t) * t
		//t3 := 3 * (1 - t) * t * t
		//t4 := t * t * t

		_t := (1 - t) * (1 - t)
		t1 := _t * (1 - t)
		t2 := 3 * _t * t
		__t := t * t
		t3 := 3 * (1 - t) * __t
		t4 := __t * t
		result = append(result, Point2D{
			X: a.X*t1 + c.X*t2 + d.X*t3 + b.X*t4,
			Y: a.Y*t1 + c.Y*t2 + d.Y*t3 + b.Y*t4,
		})
	}

	return result
}
