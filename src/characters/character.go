package characters

import (
	"math"

	"assa.com/put.pixel/lib/ogl"
	"assa.com/put.pixel/src/characters/lang/eng"
	"assa.com/put.pixel/src/characters/utf8"
)

type Chr struct {
	Ch     rune
	Size   int
	X      int
	Y      int
	PX     int
	PY     int
	Color  byte
	shape  [][]int
	wH     int
	wW     int
	width  int
	height int
}

func (ch *Chr) GetCharacterSize() int {
	return ch.Size
}

func (ch *Chr) SetCharacterSize(size int) *Chr {
	switch size {
	case 14:
		ch.shape = utf8.GetShape(ch.Ch)
	default:
		ch.Scale(size)
	}
	ch.Size = size

	return ch
}

func (ch *Chr) GetCharacterWidth() int {
	switch ch.Size {
	case 14:
		return 17
	default:
		if ch.Size > 14 {
			return 17 + ch.Size - 14
		}
		if ch.Size < 0 {
			return 3
		}
		return 17 - (14 - ch.Size)
	}
}

func (ch *Chr) GetCharacterHeight() int {
	switch ch.Size {
	case 14:
		return 20
	default:
		if ch.Size > 14 {
			return 20 + ch.Size - 14
		}
		if ch.Size < 0 {
			return 5
		}
		return 20 - (14 - ch.Size)
	}
}

func (ch *Chr) GetSpaceSizeBtwCharacters() int {
	return ch.GetCharacterWidth() / 9
}

func (ch *Chr) GetLineSpaceSize() int {
	return ch.GetCharacterWidth() / 6
}

func GetNew(chRune rune, x, y int, color byte) *Chr {
	ch := &Chr{}
	ch.Ch = chRune
	ch.wH = ogl.GetWindowHeight()
	ch.wW = ogl.GetWindowWidth()
	ch.shape = utf8.GetShape(chRune)
	ch.X = x
	ch.Y = y
	ch.PX = x
	ch.PY = y
	ch.Color = color
	ch.Size = 14
	ch.width = ch.GetCharacterWidth()
	ch.height = ch.GetCharacterHeight()

	return ch
}

func (ch *Chr) GetWidth() int {
	return ch.width
}

func (ch *Chr) GetHeight() int {
	return ch.height
}

func (ch *Chr) IsVisible() bool {
	if ch.X < 0 && int(ch.X)+ch.width < 0 {
		return false
	}

	if ch.Y > 0 && int(ch.Y)+ch.height < 0 {
		return false
	}

	if int(ch.X) > ch.wW && int(ch.X)+ch.width > ch.wW {
		return false
	}

	if int(ch.Y) > ch.wH && int(ch.Y)+ch.height > ch.wH {
		return false
	}

	return true
}

func (ch *Chr) Move(dx, dy int) {
	ch.Hide()
	ch.PX = int(ch.X)
	ch.PY = int(ch.Y)
	ch.X += dx
	ch.Y += dy

	ch.Show()
}

func (ch *Chr) Hide() {
	var color = ch.Color
	ch.Color = 0
	ch.draw(ch.shape, ch.PX, ch.PY, 0)
	ch.Color = color
}

func abs(i int) int {
	if i >= 0 {
		return i
	}

	return -i
}

func (ch *Chr) Show() {
	ch.draw(ch.shape, int(ch.X), int(ch.Y), 255)
}

func (ch *Chr) draw(shape [][]int, x, y int, a byte) {
	var i, j int
	for i = len(shape) - 1; i > 0; i-- {
		for j = len(shape[i]) - 1; j > 0; j-- {
			if shape[i][j] > 0 {
				ogl.PutPixel(x+j, y-i, ch.Color, a)
			}
		}
	}
}

func (ch *Chr) Scale(size int) {
	original := eng.GetShape(ch.Ch)
	resized := make([][]int, 0)
	ow := ch.GetCharacterWidth()
	oh := ch.GetCharacterHeight()
	ch.Size = size
	nw := ch.GetCharacterWidth()
	nh := ch.GetCharacterHeight()
	kw := float64(nw) / float64(ow)
	kh := float64(nh) / float64(oh)
	resized = make([][]int, nh)
	if kw > 1 && kh > 1 {
		for j := 0; j < nh; j++ {
			resized[j] = make([]int, nw)
			for i := 0; i < nw; i++ {
				y := int(math.Ceil(float64(j) / kh))
				x := int(math.Ceil(float64(i) / kw))
				if x >= ow {
					x = ow - 1
				}
				if y >= oh {
					y = oh - 1
				}

				resized[j][i] = original[y][x]
			}
		}
		ch.shape = resized
		return
	}
	for i := 0; i < nh; i++ {
		resized[i] = make([]int, nw)
	}
	for j := 0; j < oh; j++ {
		for i := 0; i < ow; i++ {
			if original[j][i] < 1 {
				continue
			}
			y := int(math.Ceil(float64(j) * kh))
			x := int(math.Ceil(float64(i) * kw))
			if y >= nh {
				y = nh - 1
			}
			if x >= nw {
				x = nw - 1
			}

			resized[y][x] = original[j][i]
		}
	}
	ch.shape = resized
}
