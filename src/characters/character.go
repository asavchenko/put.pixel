package characters

import (
	"fmt"
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
	// we need to apply scaling
	// for now we are going to implement only sizes >= 14
	switch size {
	case 14:
		//
		ch.shape = utf8.GetShape(ch.Ch)
	default:
		//if size < 14 {
		//	ch.shape = utf8.GetShape(ch.Ch) // downscaling is not yet implemented
		//	break
		//}
		// magic 09.01.2024 we will use nearest neighbour
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

				fmt.Println("i:", i, "j:", j, "ow:", ow, "oh:", oh, "nw:", nw, "nh:", nh, "kw:", kw, "kh:", kh, "x:", x, "y:", y)
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

			fmt.Println("i:", i, "j:", j, "ow:", ow, "oh:", oh, "nw:", nw, "nh:", nh, "kw:", kw, "kh:", kh, "x:", x, "y:", y)
			resized[y][x] = original[j][i]
		}
	}
	ch.shape = resized
}
