package snow

import (
	"assa.com/put.pixel/lib/mlib"
	"assa.com/put.pixel/lib/ogl"
	"assa.com/put.pixel/src/wind"
)

type Flake struct {
	Type    int
	x       int
	y       int
	speed   int
	dir     int
	windDir int
	color   byte
	shape   [][]int
	wH      int
	wW      int
}

func GetNew() *Flake {
	s := &Flake{}
	s.wH = ogl.GetWindowHeight()
	s.wW = ogl.GetWindowWidth()
	s.reset()
	s.windDir = wind.GetDirection()
	wind.Subscribe(func(windDir int) {
		if mlib.Rand(50) == 25 {
			s.windDir = windDir
		}
	})
	return s
}

func (s *Flake) reset() {
	s.x = mlib.Srand(s.wW)
	s.y = s.wH + mlib.Rand(500)
	s.Type = mlib.Rand(5)
	s.initType()
	s.speed = mlib.Rand(5)
	if s.windDir >= 0 {
		s.dir = 1
	} else {
		s.dir = 1
	}
	s.color = byte(mlib.Rand(255))
}

func (s *Flake) Move() {
	s.hide()
	if s.color > 0 {
		if mlib.Rand(5) >= 3 {
			s.color = s.color - 1
		}
	} else {
		s.y = s.wH + mlib.Rand(500)
		s.color = byte(mlib.Rand(256))
	}

	if s.dir > 0 {
		s.dir -= 1
	}
	if s.dir < 0 {
		s.dir += 1
	}
	if s.dir == 0 {
		if s.windDir >= 0 {
			s.dir = s.windDir + mlib.Rand(1+s.windDir)
		} else {
			s.dir = s.windDir + mlib.Rand(1-s.windDir)
		}
	}
	s.x += mlib.Sign(float64(s.dir))
	if s.x > s.wW {
		s.x = -mlib.Rand(200)
	}
	if s.x < 0 {
		s.x = mlib.Rand(200)
	}

	if s.y > 0 {
		s.y -= s.speed
	} else {
		s.color = byte(mlib.Rand(256))
		s.x = mlib.Rand(s.wW)
		s.y = s.wH + 100 + mlib.Rand(400)
		s.speed = mlib.Rand(5)
	}
	s.show()
}

func (s *Flake) initType() {
	switch s.Type {
	case 1:
		s.shape = [][]int{{1}}
	case 2:
		s.shape = [][]int{{1, 0}, {0, 1}}
	case 3:
		s.shape = [][]int{{1, 0, 1},
			{0, 1, 0},
			{1, 0, 1}}
	case 4:
		s.shape = [][]int{{1, 0, 0, 1},
			{0, 1, 1, 0},
			{0, 1, 1, 0},
			{1, 0, 0, 1}}
	case 5:
		s.shape = [][]int{{1, 0, 1, 0, 1},
			{0, 1, 0, 1, 0},
			{0, 0, 1, 0, 0},
			{0, 1, 0, 1, 0},
			{1, 0, 1, 0, 1}}
	default:
		s.shape = [][]int{{1}}
	}
}

func (s *Flake) hide() {
	s.draw(s.shape, s.x, s.y, 0)
}

func (s *Flake) show() {
	s.draw(s.shape, s.x, s.y, s.color)
}

func (s *Flake) draw(shape [][]int, x, y int, color byte) {
	var i, j int
	for i = len(shape) - 1; i > 0; i-- {
		for j = len(shape) - 1; j > 0; j-- {
			if shape[i][j] > 0 {
				ogl.PutPixel(x+i, y+j, color)
			}
		}
	}
}
