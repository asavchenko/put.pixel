package arrow

import (
	"math"

	"assa.com/put.pixel/lib/ogl"
)

type arrow struct {
	x0, y0 int     // base of the arrow
	a      float64 // angle
	l      int     // length
	color  uint32
}

func GetNew(x0, y0 int, l int, a float64, color uint32) *arrow {
	_a := arrow{
		x0:    x0,
		y0:    y0,
		a:     a,
		color: color,
		l:     l,
	}

	return &_a
}

func (a *arrow) Rotate(da float64) {
	a.a += da
}

func (a *arrow) A() float64 {
	return a.a
}

func (a *arrow) RotateTo(x1, y1 int) {
	if x1 == a.x0 {
		return
	}

	if x1 > a.x0 {
		a.a = math.Atan(float64(y1-a.y0) / float64(x1-a.x0))
	} else {
		a.a = math.Pi + math.Atan(float64(y1-a.y0)/float64(x1-a.x0))
	}
}

func (a *arrow) Show() {
	x := int(math.Round(math.Cos(a.a)*float64(a.l) + float64(a.x0)))
	y := int(math.Round(math.Sin(a.a)*float64(a.l) + float64(a.y0)))

	ogl.LineRGBA(a.x0, a.y0, x, y, a.color)
	dl := a.l - a.l/6
	ogl.LineRGBA(x, y,
		int(math.Round(math.Cos(a.a-math.Pi/18)*float64(dl)+float64(a.x0))),
		int(math.Round(math.Sin(a.a-math.Pi/18)*float64(dl)+float64(a.y0))),
		a.color)
	ogl.LineRGBA(x, y,
		int(math.Round(math.Cos(a.a+math.Pi/18)*float64(dl)+float64(a.x0))),
		int(math.Round(math.Sin(a.a+math.Pi/18)*float64(dl)+float64(a.y0))),
		a.color)
}
