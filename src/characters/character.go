package characters

import (
	"assa.com/put.pixel/lib/ogl"
	"assa.com/put.pixel/src/characters/utf8"
	"math"
	"time"
)

const INVISIBLE_COLOR = 0xffffffff

type Chr struct {
	Ch              rune
	Size            int
	X               int
	Y               int
	PX              int
	PY              int
	Color           byte
	shape           []byte
	prev            []uint32
	cur             []uint32
	commandsCh      chan map[string]interface{}
	wH              int
	wW              int
	width           int
	height          int
	shapeWidth      int
	shapeHeight     int
	shapeSize       int
	isVeryFirstShow bool
}

func (ch *Chr) GetCharacterSize() int {
	return ch.Size
}

func (ch *Chr) SetCharacterSize(size int) *Chr {
	ch.hideIgnoreVisible()
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
		return utf8.GetShapeHeight() + (14 - ch.Size)
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
	ch.shapeHeight = utf8.HEIGHT
	ch.shapeWidth = utf8.WIDTH
	ch.X = x
	ch.Y = y
	ch.PX = x
	ch.PY = y
	ch.Color = color
	ch.Size = 14
	ch.width = ch.GetCharacterWidth()
	ch.height = ch.GetCharacterHeight()
	ch.shapeSize = ch.shapeHeight * ch.shapeWidth
	ch.prev = make([]uint32, ch.shapeHeight*ch.shapeWidth)
	ch.cur = make([]uint32, ch.shapeHeight*ch.shapeWidth)
	di := 0
	for i := 0; i < ch.shapeHeight; i++ {
		for j := 0; j < ch.shapeWidth; j++ {
			dij := di + j
			ch.cur[dij] = INVISIBLE_COLOR
			ch.prev[dij] = INVISIBLE_COLOR
		}
		di += ch.shapeWidth
	}

	ch.commandsCh = make(chan map[string]interface{}, 0)
	go func() {
		for command := range ch.commandsCh {
			switch command["action"] {
			case "get_position":
				select {
				case command["output_to"].(chan []int) <- []int{ch.X, ch.Y}:
				case <-time.After(3 * time.Second):
				}
			case "move":
				data := command["data"].(map[string]interface{})
				dx := data["dx"].(int)
				dy := data["dy"].(int)
				ch.move(dx, dy)
				select {
				case command["done"].(chan bool) <- true:
				case <-time.After(3 * time.Second):
				}
			case "scale":
				data := command["data"].(map[string]interface{})
				size := data["size"].(int)
				ch.Scale(size)
				select {
				case command["done"].(chan bool) <- true:
				case <-time.After(3 * time.Second):
				}
			}
		}
	}()
	return ch
}

func (ch *Chr) GetWidth() int {
	return ch.width
}

func (ch *Chr) GetHeight() int {
	return ch.height
}

func (ch *Chr) IsInvisible() bool {
	return ch.X+ch.shapeWidth <= 0 || ch.X >= ch.wW || ch.Y+ch.shapeHeight <= 0 || ch.Y >= ch.wH
}

func (ch *Chr) isPixelVisible(x, y int) bool {
	return x >= 0 && y >= 0 && x < ch.wW && y < ch.wH
}

func (ch *Chr) Move(dx, dy int) chan bool {
	doneCh := make(chan bool, 1)
	select {
	case ch.commandsCh <- map[string]interface{}{
		"action": "move",
		"data": map[string]interface{}{
			"dx": dx,
			"dy": dy,
		},
		"done": doneCh,
	}:
	case <-time.After(3 * time.Second):
	}

	return doneCh
}

func (ch *Chr) Scale(size int) chan bool {
	doneCh := make(chan bool, 1)
	select {
	case ch.commandsCh <- map[string]interface{}{
		"action": "scale",
		"data": map[string]interface{}{
			"size": size,
		},
		"done": doneCh,
	}:
	case <-time.After(3 * time.Second):
	}

	return doneCh
}

func (ch *Chr) GetPosition() (int, int) {
	readFrom := make(chan []int, 1)
	select {
	case ch.commandsCh <- map[string]interface{}{
		"action":    "get_position",
		"output_to": readFrom,
	}:
	case <-time.After(3 * time.Second):
		return 0, 0
	}
	select {
	case res := <-readFrom:
		return res[0], res[1]
	case <-time.After(3 * time.Second):
		return 0, 0
	}
}

func (ch *Chr) move(dx, dy int) {
	ch.PX = ch.X
	ch.PY = ch.Y
	ch.X += dx
	ch.Y += dy
	ch.show()
	ch.hide()
}

func (ch *Chr) MoveUnsafe(dx, dy int) {
	ch.PX = ch.X
	ch.PY = ch.Y
	ch.X += dx
	ch.Y += dy
}

func (ch *Chr) HideUnsafe() {
	if ch.PX == ch.X && ch.PY == ch.Y {
		return
	}
	x := ch.PX
	y := ch.PY
	di := 0
	for i := 0; i < ch.shapeHeight; i++ {
		for j := 0; j < ch.shapeWidth; j++ {
			dij := di + j
			if ch.prev[dij] != INVISIBLE_COLOR && ch.curNotContainsOrInvisible(x, y) {
				ogl.UnsafePutPixelRGB(x, y, byte(ch.prev[dij]&0x000000FF), byte((ch.prev[dij]&0x0000FF00)>>8), byte((ch.prev[dij]&0x00FF0000)>>16))
			}
			x++
		}
		x = ch.PX
		y++
		di += ch.shapeWidth
	}
}

func (ch *Chr) hide() {
	if ch.PX == ch.X && ch.PY == ch.Y {
		return
	}

	x := ch.PX
	y := ch.PY
	di := 0
	for i := 0; i < ch.shapeHeight; i++ {
		for j := 0; j < ch.shapeWidth; j++ {
			dij := di + j
			if ch.prev[dij] != INVISIBLE_COLOR && ch.curNotContainsOrInvisible(x, y) {

				ogl.UnsafePutPixelRGB(x, y, byte(ch.prev[dij]&0x000000FF), byte((ch.prev[dij]&0x0000FF00)>>8), byte((ch.prev[dij]&0x00FF0000)>>16))
			}
			x++
		}
		x = ch.PX
		y++
		di += ch.shapeWidth
	}
}

func (ch *Chr) hideIgnoreVisible() {
	x := ch.PX
	y := ch.PY
	di := 0
	for i := 0; i < ch.shapeHeight; i++ {
		for j := 0; j < ch.shapeWidth; j++ {
			dij := di + j
			if ch.prev[dij] != INVISIBLE_COLOR {
				ogl.UnsafePutPixelRGB(x, y, byte(ch.prev[dij]&0x000000FF), byte((ch.prev[dij]&0x0000FF00)>>8), byte((ch.prev[dij]&0x00FF0000)>>16))
			}
			x++
		}
		x = ch.PX
		y++
		di += ch.shapeWidth
	}
}

func (ch *Chr) contains(x, y int) bool {
	return x < ch.X+ch.shapeWidth && x >= ch.X && y < ch.Y+ch.shapeHeight && y >= ch.Y
}

func (ch *Chr) show() {
	for i, e := range ch.cur {
		ch.prev[i] = e
	}
	if ch.IsInvisible() {
		di := 0
		for i := 0; i < ch.shapeHeight; i++ {
			for j := 0; j < ch.shapeWidth; j++ {
				ch.cur[di+j] = INVISIBLE_COLOR
			}
			di += ch.shapeWidth
		}

		return
	}

	x := ch.X
	y := ch.Y
	di := 0
	for i := 0; i < ch.shapeHeight; i++ {
		for j := 0; j < ch.shapeWidth; j++ {
			dij := di + j
			if ch.isPixelVisible(x, y) {
				rgb, contains := ch.prevContains(x, y)
				if contains {
					ch.cur[dij] = rgb
				} else {
					ch.cur[dij] = ogl.GetPixelUnsafe(x, y)
				}
			} else {
				ch.cur[dij] = INVISIBLE_COLOR
			}
			x++
		}
		x = ch.X
		y++
		di += ch.shapeWidth
	}

	ch.draw()
}

func (ch *Chr) ShowUnsafe() {
	for i, e := range ch.cur {
		ch.prev[i] = e
	}
	if ch.IsInvisible() {
		di := 0
		for i := 0; i < ch.shapeHeight; i++ {
			for j := 0; j < ch.shapeWidth; j++ {
				ch.cur[di+j] = INVISIBLE_COLOR
			}
			di += ch.shapeWidth
		}

		return
	}

	x := ch.X
	y := ch.Y
	di := 0
	for i := 0; i < ch.shapeHeight; i++ {
		for j := 0; j < ch.shapeWidth; j++ {
			dij := di + j
			if ch.isPixelVisible(x, y) {
				rgb, contains := ch.prevContains(x, y)
				if contains {
					ch.cur[dij] = rgb
				} else {
					ch.cur[dij] = ogl.GetPixelUnsafe(x, y)
				}
			} else {
				ch.cur[dij] = INVISIBLE_COLOR
			}
			x++
		}
		x = ch.X
		y++
		di += ch.shapeWidth
	}

	ch.draw()
}

func (ch *Chr) draw() {
	x := ch.X
	y := ch.Y
	di := 0
	for i := 0; i < ch.shapeHeight; i++ {
		for j := 0; j < ch.shapeWidth; j++ {
			if ch.shape[di+j] > 0 {
				ogl.PutPixelRGB(x, y, ch.Color, ch.Color, ch.Color)
			}
			x++
		}
		x = ch.X
		y++
		di += ch.shapeWidth
	}
}

func (ch *Chr) scale(size int) {
	defer func() {
		ch.prev = make([]uint32, ch.shapeHeight*ch.shapeWidth)
		ch.cur = make([]uint32, ch.shapeHeight*ch.shapeWidth)
		di := 0
		for i := 0; i < ch.shapeHeight; i++ {
			for j := 0; j < ch.shapeWidth; j++ {
				dij := di + j
				ch.cur[dij] = INVISIBLE_COLOR
				ch.prev[dij] = INVISIBLE_COLOR
			}
			di += ch.shapeWidth
		}
		ch.shapeSize = ch.shapeHeight * ch.shapeWidth
	}()
	original := utf8.GetShape(ch.Ch)
	ow := ch.GetCharacterWidth()
	oh := ch.GetCharacterHeight()
	ch.Size = size
	nw := ch.GetCharacterWidth()
	nh := ch.GetCharacterHeight()
	kw := float64(nw) / float64(ow)
	kh := float64(nh) / float64(oh)
	ch.shapeWidth = nw
	ch.shapeHeight = nh
	resized := make([]byte, nh*nw)
	if kw > 1 && kh > 1 {
		dj := 0
		for j := 0; j < nh; j++ {
			for i := 0; i < nw; i++ {
				y := int(math.Ceil(float64(j) / kh))
				x := int(math.Ceil(float64(i) / kw))
				if x >= ow {
					x = ow - 1
				}
				if y >= oh {
					y = oh - 1
				}

				resized[dj+i] = original[y*ow+x]
			}
			dj += nw
		}
		ch.shape = resized
		return
	}
	dj := 0
	for j := 0; j < oh; j++ {
		for i := 0; i < ow; i++ {
			if original[dj+i] < 1 {
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

			resized[y*nw+x] = original[dj+i]
		}
		dj += ow
	}
	ch.shape = resized
}

func (ch *Chr) prevContains(x, y int) (uint32, bool) {
	i := x - ch.PX
	if i < 0 || i >= ch.shapeWidth {
		return INVISIBLE_COLOR, false
	}
	j := y - ch.PY
	if j < 0 || j >= ch.shapeHeight {
		return INVISIBLE_COLOR, false
	}

	e := ch.prev[i+j*ch.shapeWidth]
	if e == INVISIBLE_COLOR {
		return INVISIBLE_COLOR, false
	}
	return e, true
}

func (ch *Chr) curNotContainsOrInvisible(x, y int) bool {
	i := x - ch.X
	if i < 0 || i >= ch.shapeWidth {
		return false
	}
	j := y - ch.Y
	if j < 0 || j >= ch.shapeHeight {
		return false
	}

	return ch.shape[i+j*ch.shapeWidth] == 0
}
