package snow

import (
	"assa.com/put.pixel/lib/mlib"
	"assa.com/put.pixel/lib/ogl"
	"assa.com/put.pixel/src/wind"
)

type Flake struct {
	Type    int
	X       int
	Y       int
	Speed   int
	Dir     int
	WindDir int
	Color   byte
	shape   [][]int
	frame   int
	wH      int
	wW      int
}

func GetNew() *Flake {
	s := &Flake{}
	s.wH = ogl.GetWindowHeight()
	s.wW = ogl.GetWindowWidth()
	s.reset()
	s.WindDir = wind.GetDirection()
	wind.Subscribe(func(windDir int) {
		s.WindDir = windDir
	})
	return s
}

func (s *Flake) reset() {
	s.X = mlib.Rand(s.wW)
	s.Y = s.wH + mlib.Rand(500)
	s.Type = mlib.Rand(5)
	s.initType()
	s.Speed = mlib.Rand(5)
	s.Dir = mlib.Srand(s.WindDir)
	s.Color = byte(mlib.Rand(255))
}

func (s *Flake) Move() {
	s.hide()
	if s.frame%s.Speed == 0 {
		s.fallLogic()
	}
	if s.frame%s.Dir == 0 {
		s.blowLogic()
	}
	s.show()
	s.frame++
}

func (s *Flake) colorLogic() {
	if s.Color < 0 {
		s.Y = s.wH + mlib.Rand(50)
		s.Color = byte(mlib.Rand(255))
		return
	}

	if mlib.Rand(5) == 2 {
		s.Color--
	}
}

func (s *Flake) directionLogic() {
	if s.Dir > 0 && mlib.Rand(150) == 2 {
		s.Dir--
	}

	if s.Dir < 0 && mlib.Rand(150) == 2 {
		s.Dir++
	}

	if s.Dir != 0 {
		return
	}

	windDir := s.WindDir
	if windDir > 0 {
		s.Dir = windDir + mlib.Rand(windDir)
	} else {
		s.Dir = windDir - mlib.Rand(1-windDir)
	}
}

func (s *Flake) fallLogic() {
	if s.Y-1 > 0 {
		s.Y--
		return
	}
	s.reset()
}

func (s *Flake) blowLogic() {
	s.X += mlib.Sign(float64(s.Dir))
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
	s.draw(s.shape, s.X, s.Y, s.Color/8)
}

func (s *Flake) show() {
	s.draw(s.shape, s.X, s.Y, s.Color)
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
