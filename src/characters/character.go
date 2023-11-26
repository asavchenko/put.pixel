package characters

import (
	"assa.com/put.pixel/lib/ogl"
	"assa.com/put.pixel/src/characters/lang/eng"
)

type Chr struct {
	X      int
	Y      int
	Color  byte
	shape  [][]int
	wH     int
	wW     int
	width  int
	height int
}

var characterSize = 14

func GetCharacterSize() int {
	return characterSize
}

func SetCharacterSize(size int) {
	characterSize = size
}

func GetCharacterWidth() int {
	switch characterSize {
	case 14:
		return 17
	default:
		if characterSize > 14 {
			return 17 + characterSize - 14
		}
		if characterSize < 0 {
			return 3
		}
		return 17 - (14 - characterSize)
	}
}

func GetCharacterHeight() int {
	switch characterSize {
	case 14:
		return 20
	default:
		if characterSize > 14 {
			return 20 + characterSize - 14
		}
		if characterSize < 0 {
			return 5
		}
		return 20 - (14 - characterSize)
	}
}

func GetSpaceSizeBtwCharacters() int {
	return 0
	return GetCharacterWidth() / 9
}

func GetLineSpaceSize() int {
	return 0
	return GetCharacterWidth() / 6
}

func GetNew(chRune rune, x, y int, color byte) *Chr {
	ch := &Chr{}
	ch.wH = ogl.GetWindowHeight()
	ch.wW = ogl.GetWindowWidth()
	ch.shape = eng.GetShape(chRune)
	ch.X = x
	ch.Y = y
	ch.Color = color
	ch.width = GetCharacterWidth()
	ch.height = GetCharacterHeight()

	return ch
}

func (ch *Chr) GetWidth() int {
	return ch.width
}

func (ch *Chr) GetHeight() int {
	return ch.height
}

func (ch *Chr) IsVisible() bool {
	if ch.X < 0 && ch.X+ch.width < 0 {
		return false
	}

	if ch.Y > 0 && ch.Y+ch.height < 0 {
		return false
	}

	if ch.X > ch.wW && ch.X+ch.width > ch.wW {
		return false
	}

	if ch.Y > ch.wH && ch.Y+ch.height > ch.wH {
		return false
	}

	return true
}

func (ch *Chr) Move(dx, dy int) {
	ch.Hide()
	ch.X += dx
	ch.Y += dy
	ch.Show()
}

func (ch *Chr) Hide() {
	ch.draw(ch.shape, ch.X, ch.Y, 0)
}

func (ch *Chr) Show() {
	ch.draw(ch.shape, ch.X, ch.Y, 255)
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
