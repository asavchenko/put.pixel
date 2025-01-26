package characters

import (
	"math"

	"assa.com/put.pixel/lib/ogl"
	"assa.com/put.pixel/src/characters/utf8"
)

type Chr struct {
	Ch              rune
	Size            int
	X               int
	Y               int
	PX              int
	PY              int
	Color           byte
	shape           [][]byte
	prev            [][][]byte
	cur             [][][]byte
	wH              int
	wW              int
	width           int
	height          int
	shapeWidth      int
	shapeHeight     int
	isVeryFirstShow bool
}

func (ch *Chr) GetCharacterSize() int {
	return ch.Size
}

func (ch *Chr) SetCharacterSize(size int) *Chr {
	ch.HideIgnoreVisible()
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
		return utf8.GetShapeWidth()
	default:
		if ch.Size > 14 {
			return utf8.GetShapeWidth() + ch.Size - 14
		}
		if ch.Size < 0 {
			return 3
		}
		return utf8.GetShapeWidth() - (14 - ch.Size)
	}
}

func (ch *Chr) GetCharacterHeight() int {
	switch ch.Size {
	case 14:
		return utf8.GetShapeHeight()
	default:
		if ch.Size > 14 {
			return utf8.GetShapeHeight() + ch.Size - 14
		}
		if ch.Size < 0 {
			return 5
		}
		return utf8.GetShapeHeight() - (14 - ch.Size)
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
	ch.prev = make([][][]byte, 0)
	ch.cur = make([][][]byte, 0)
	ch.X = x
	ch.Y = y
	ch.PX = x
	ch.PY = y
	ch.Color = color
	ch.Size = 14
	ch.width = ch.GetCharacterWidth()
	ch.height = ch.GetCharacterHeight()
	ch.shapeHeight = len(ch.shape)
	for _, l := range ch.shape {
		ch.shapeWidth = len(l)
		break
	}

	return ch
}

func (ch *Chr) GetWidth() int {
	return ch.width
}

func (ch *Chr) GetHeight() int {
	return ch.height
}

func (ch *Chr) IsVisible() bool {
	if ch.X < 0 && ch.X+ch.shapeWidth < 0 {
		return false
	}

	if ch.X >= ch.wW && ch.X+ch.shapeWidth >= ch.wW {
		return false
	}

	if ch.Y < 0 && ch.Y-ch.shapeHeight < 0 {
		return false
	}

	if ch.Y >= ch.wH && ch.Y-ch.shapeHeight >= ch.wH {
		return false
	}

	return true
}

func (ch *Chr) IsPixelVisible(x, y int) bool {
	return x >= 0 && y >= 0 && x < ch.wW && y < ch.wH
}

func (ch *Chr) Move(dx, dy int) {
	ch.PX = ch.X
	ch.PY = ch.Y
	ch.X += dx
	ch.Y += dy
	ch.Show()
	ch.Hide()
}

func (ch *Chr) Hide() {
	if ch.PX == ch.X && ch.PY == ch.Y {
		return
	}
	var i, j int
	for i = ch.shapeHeight - 1; i >= 0; i-- {
		for j = ch.shapeWidth - 1; j >= 0; j-- {
			if ch.IsPixelVisible(ch.PX+j, ch.PY-i) && ch.curNotContainsOrInvisible(ch.PX+j, ch.PY-i) && len(ch.prev) > 0 && len(ch.prev[i]) > 0 && len(ch.prev[i][j]) > 0 {
				ogl.PutPixel(ch.PX+j, ch.PY-i, ch.prev[i][j]...)
			}
		}
	}
}

func (ch *Chr) HideIgnoreVisible() {
	var i, j int
	for i = ch.shapeHeight - 1; i >= 0; i-- {
		for j = ch.shapeWidth - 1; j >= 0; j-- {
			if ch.IsPixelVisible(ch.PX+j, ch.PY-i) && len(ch.prev) > 0 && len(ch.prev[i]) > 0 && len(ch.prev[i][j]) > 0 {
				ogl.PutPixel(ch.PX+j, ch.PY-i, ch.prev[i][j]...)
			}
		}
	}
}

func (ch *Chr) contains(x, y int) bool {
	return x < ch.X+ch.shapeWidth && x >= ch.X && y > ch.Y-ch.shapeHeight && y <= ch.Y
}

func (ch *Chr) Show() {
	ch.prev = make([][][]byte, ch.shapeHeight)
	copy(ch.prev, ch.cur)
	ch.cur = make([][][]byte, ch.shapeHeight)
	if !ch.IsVisible() {
		return
	}
	var i, j int
	for i = ch.shapeHeight - 1; i >= 0; i-- {
		line := make([][]byte, ch.shapeWidth)
		for j = ch.shapeWidth - 1; j >= 0; j-- {
			if ch.IsPixelVisible(ch.X+j, ch.Y-i) {
				rgbArr, contains := ch.prevContains(ch.X+j, ch.Y-i)
				if contains && len(rgbArr) > 0 {
					line[j] = rgbArr
				} else {
					line[j] = ogl.GetPixel(ch.X+j, ch.Y-i)
				}
			} else {
				line[j] = []byte{}
			}
		}
		ch.cur[i] = line
	}
	ch.draw(ch.shape, ch.X, ch.Y)
}

func (ch *Chr) draw(shape [][]byte, x, y int) {
	if !ch.IsVisible() {
		return
	}
	var i, j int
	for i = len(shape) - 1; i >= 0; i-- {
		for j = len(shape[i]) - 1; j >= 0; j-- {
			if shape[i][j] > 0 && ch.IsPixelVisible(x+j, y-i) {
				ogl.PutPixel(x+j, y-i, ch.Color)
			}
		}
	}
}

func (ch *Chr) Scale(size int) {
	defer func() {
		ch.shapeHeight = len(ch.shape)
		for _, l := range ch.shape {
			ch.shapeWidth = len(l)
			break
		}
	}()
	original := utf8.GetShape(ch.Ch)
	resized := make([][]byte, 0)
	ow := ch.GetCharacterWidth()
	oh := ch.GetCharacterHeight()
	ch.Size = size
	nw := ch.GetCharacterWidth()
	nh := ch.GetCharacterHeight()
	kw := float64(nw) / float64(ow)
	kh := float64(nh) / float64(oh)
	resized = make([][]byte, nh)
	if kw > 1 && kh > 1 {
		for j := 0; j < nh; j++ {
			resized[j] = make([]byte, nw)
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
		resized[i] = make([]byte, nw)
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

func (ch *Chr) prevContains(x int, y int) ([]byte, bool) {
	if len(ch.prev) < 1 {
		return nil, false
	}
	i := ch.PY - y
	if i < 0 || i >= ch.shapeHeight {
		return nil, false
	}
	j := x - ch.PX
	if j < 0 || j >= ch.shapeWidth {
		return nil, false
	}
	if len(ch.prev[i]) < 1 {
		return nil, false
	}
	if len(ch.prev[i][j]) < 1 {
		return nil, false
	}

	return ch.prev[i][j], true
}

func (ch *Chr) curNotContainsOrInvisible(x int, y int) bool {
	if len(ch.cur) < 1 {
		return true
	}
	if !ch.IsPixelVisible(x, y) {
		return true
	}
	i := ch.Y - y
	if i < 0 || i >= ch.shapeHeight {
		return true
	}
	j := x - ch.X
	if j < 0 || j >= ch.shapeWidth {
		return true
	}

	return ch.shape[i][j] == 0
}
