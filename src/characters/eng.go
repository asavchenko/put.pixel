package characters

import (
	"assa.com/put.pixel/lib/mlib"
	"assa.com/put.pixel/lib/ogl"
)

type Chr struct {
	Type  int
	X     int
	Y     int
	Color byte
	shape [][]int
	frame int
	wH    int
	wW    int
}

func GetNew() *Chr {
	ch := &Chr{}
	ch.wH = ogl.GetWindowHeight()
	ch.wW = ogl.GetWindowWidth()
	ch.reset()
	return ch
}

func (ch *Chr) reset() {
	ch.X = mlib.Rand(ch.wW)
	ch.Y = ch.wH + mlib.Rand(500)
	switch mlib.Rand(9) {
	case 1:
		ch.Type = 18
	//case 2:
	//	ch.Type = 7
	//case 3:
	//	ch.Type = 9
	//case 4:
	//	ch.Type = 25
	case 5:
		ch.Type = 18
	//case 6:
	//	ch.Type = 11
	default:
		ch.Type = 18
	}
	ch.Type = 18
	ch.initType()
	ch.Color = byte(mlib.Rand(255))
}

func (ch *Chr) Move() {
	ch.hide()
	ch.show()
	ch.frame++
}

func (ch *Chr) colorLogic() {
	if ch.Color < 0 {
		ch.Y = ch.wH + mlib.Rand(50)
		ch.Color = byte(mlib.Rand(255))
		return
	}

	if mlib.Rand(5) == 2 {
		ch.Color--
	}
}

func (ch *Chr) directionLogic() {
}

func (ch *Chr) fallLogic() {
	if ch.Y-1 > 0 {
		ch.Y--
		return
	}
	ch.reset()
}

func (ch *Chr) blowLogic() {
}

func (ch *Chr) initType() {
	switch ch.Type {
	case 18:
		ch.shape = [][]int{
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		}
	default:
		ch.shape = [][]int{{1}}
	}
}

func (ch *Chr) hide() {
	ch.draw(ch.shape, ch.X, ch.Y, 0)
}

func (ch *Chr) show() {
	ch.draw(ch.shape, ch.X, ch.Y, 255)
}

func (ch *Chr) draw(shape [][]int, x, y int, a byte) {
	var i, j int
	for i = len(shape) - 1; i > 0; i-- {
		for j = len(shape) - 1; j > 0; j-- {
			if shape[i][j] > 0 {
				ogl.PutPixel(x+i, y+j, ch.Color, a)
			}
		}
	}
}
