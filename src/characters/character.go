package characters

import (
	"assa.com/put.pixel/lib/mlib"
	"math"
	"time"

	"assa.com/put.pixel/lib/ogl"
	"assa.com/put.pixel/src/characters/utf8"
)

const INVISIBLE_COLOR = 0xffffffff

type Chr struct {
	Ch            rune
	a             float64
	fallingSpeed  int
	rotationSpeed float64
	Size          int
	X             int
	Y             int
	PX            int
	PY            int
	Color         uint32
	shape         []byte
	bitmap        []byte
	commandsCh    chan map[string]interface{}
	wH            int
	wW            int
	shapeWidth    int
	shapeHeight   int
	shapeSize     int
}

func (ch *Chr) SetChar(chRune rune) {
	ch.Ch = chRune
	<-ch.Scale(ch.Size)
	//ch.generateBitmap()
}

func (ch *Chr) GetCharacterSize() int {
	return ch.Size
}

func (ch *Chr) SetCharacterSize(size int) *Chr {
	switch size {
	case 14:
		ch.shape = ch.trim(utf8.GetShape(ch.Ch), utf8.WIDTH, utf8.HEIGHT)
	default:
		<-ch.Scale(size)
	}
	ch.Size = size

	return ch
}

func (ch *Chr) GetMaxCharacterWidth() int {
	switch ch.Size {
	case 14:
		return utf8.GetShapeWidth()
	default:
		return utf8.GetShapeWidth() * ch.Size / 14
	}
}

func (ch *Chr) GetMaxCharacterHeight() int {
	switch ch.Size {
	case 14:
		return utf8.GetShapeHeight()
	default:
		return utf8.GetShapeHeight() * ch.Size / 14
	}
}

func (ch *Chr) GetSpaceSizeBtwCharacters() int {
	return 1
}

func (ch *Chr) GetLineSpaceSize() int {
	return ch.GetMaxCharacterHeight() + 1
}

func GetNew(chRune rune, x, y int, color uint32) *Chr {
	ch := &Chr{}
	ch.Ch = chRune
	ch.wH = ogl.GetWindowHeight()
	ch.wW = ogl.GetWindowWidth()
	ch.shape = ch.trim(utf8.GetShape(chRune), utf8.WIDTH, utf8.HEIGHT)
	ch.X = x
	ch.Y = y
	ch.PX = x
	ch.PY = y
	ch.Color = color
	ch.Size = 14
	ch.shapeSize = ch.shapeHeight * ch.shapeWidth
	ch.generateBitmap()
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
				ch.scale(size)
				select {
				case command["done"].(chan bool) <- true:
				case <-time.After(3 * time.Second):
				}
			}
		}
	}()
	return ch
}

func (ch *Chr) MoveUnsafe(dx, dy int) {
	ch.PX = ch.X
	ch.PY = ch.Y
	ch.X += dx
	ch.Y += dy
	ch.ShowUnsafe()
}

func (ch *Chr) contains(x, y int) bool {
	return x < ch.X+ch.shapeWidth && x >= ch.X && y < ch.Y+ch.shapeHeight && y >= ch.Y
}

func (ch *Chr) ShowUnsafe() {
	if ch.IsInvisible() {
		return
	}

	ch.draw()
}

func (ch *Chr) draw() {
	ogl.PutByteBitmap(ch.X, ch.Y, ch.shapeWidth*4, ch.shapeHeight, ch.bitmap)
	//x := ch.X
	//y := ch.Y
	//di := 0
	//for i := 0; i < ch.shapeHeight; i++ {
	//	for j := 0; j < ch.shapeWidth; j++ {
	//		if ch.shape[di+j] > 0 {
	//			ogl.PutPixel(x, y, 200, 200, 200)
	//		}
	//		x++
	//	}
	//	di += ch.shapeWidth
	//	x = ch.X
	//	y++
	//}
	//ogl.Line(ch.X, ch.Y, ch.X+ch.shapeWidth, ch.Y, ch.Color, ch.Color, ch.Color)
	//ogl.Line(ch.X+ch.shapeWidth, ch.Y, ch.X+ch.shapeWidth, ch.Y+ch.shapeHeight, ch.Color, ch.Color, ch.Color)
	//ogl.Line(ch.X+ch.shapeWidth, ch.Y+ch.shapeHeight, ch.X, ch.Y+ch.shapeHeight, ch.Color, ch.Color, ch.Color)
	//ogl.Line(ch.X, ch.Y+ch.shapeHeight, ch.X, ch.Y, ch.Color, ch.Color, ch.Color)
}

func (ch *Chr) trim(shape []byte, w, h int) []byte {
	li := 0     // left idx
	ri := w - 1 // right idx
	for i := 0; i < w; i++ {
		dj := 0
		isEmpty := true
		for j := 0; j < h; j++ {
			if shape[dj+i] != 0 {
				isEmpty = false
				break
			}
			dj += w
		}
		if isEmpty {
			li++
			continue
		}
		break
	}
	for i := w - 1; i > li; i-- {
		dj := 0
		isEmpty := true
		for j := 0; j < h; j++ {
			if shape[dj+i] != 0 {
				isEmpty = false
				break
			}
			dj += w
		}
		if isEmpty {
			ri--
			continue
		}
		break
	}
	nw := ri - li + 1
	if nw < 1 {
		li = 0
		ri = w/2 - 1
		nw = ri - li + 1
	}
	trimmed := make([]byte, nw*h)
	for i := li; i <= ri; i++ {
		dj := 0
		for j := 0; j < h; j++ {
			trimmed[i+dj-li] = shape[j*w+i]
			dj += nw
		}
	}
	ch.shapeHeight = h
	ch.shapeWidth = nw

	return trimmed
}

func (ch *Chr) scale(size int) {
	defer func() {
		ch.shapeSize = ch.shapeHeight * ch.shapeWidth
		ch.generateBitmap()
	}()
	original := utf8.GetShape(ch.Ch)
	ow := utf8.WIDTH
	oh := utf8.HEIGHT
	ch.Size = size
	nw := ow * size / 14
	nh := oh * size / 14
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
		ch.shape = ch.trim(resized, nw, nh)
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
	ch.shape = ch.trim(resized, nw, nh)
}

// /////////////////////////////////////////////////////
func (ch *Chr) show() {
	if ch.IsInvisible() {
		return
	}

	ch.draw()
}

func (ch *Chr) GetWidth() int {
	return ch.shapeWidth
}

func (ch *Chr) GetHeight() int {
	return ch.shapeHeight
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
}

func (ch *Chr) Rotate(a float64) {
	dj := 0
	y := ch.Y
	x0 := float64(ch.X + ch.shapeWidth/2)
	for j := 0; j < ch.shapeHeight; j++ {
		for i := 0; i < ch.shapeWidth; i++ {
			if ch.shape[dj+i] > 0 {
				x := int(math.Ceil((float64(i)-float64(ch.shapeWidth)/2)*math.Cos(a) + x0))
				ogl.PutPixelRGB(x, y, ch.Color)
			}
		}
		dj += ch.shapeWidth
		y++
	}
}

func (ch *Chr) generateBitmap() {
	ch.bitmap = make([]byte, ch.shapeHeight*ch.shapeWidth*4)
	di := 0
	idx := 0
	for i := 0; i < ch.shapeHeight; i++ {
		for j := 0; j < ch.shapeWidth; j++ {
			if ch.shape[di+j] > 0 {
				ch.bitmap[idx] = byte(ch.Color >> 24)
				ch.bitmap[idx+1] = byte(ch.Color >> 16)
				ch.bitmap[idx+2] = byte(ch.Color >> 8)
				ch.bitmap[idx+3] = byte(ch.Color)
			} else {
				ch.bitmap[idx] = 0
				ch.bitmap[idx+1] = 0
				ch.bitmap[idx+2] = 0
				ch.bitmap[idx+3] = 0
			}
			idx += 4
		}
		di += ch.shapeWidth
	}
}

func (ch *Chr) SetRotationSpeed(f float64) {
	ch.rotationSpeed = f
}

func (ch *Chr) SetFallingSpeed(i int) {
	ch.fallingSpeed = i
}

func (ch *Chr) Run() {
	ch.a += ch.rotationSpeed
	if ch.a > 2*math.Pi {
		for {
			if ch.a < 2*math.Pi {
				break
			}
			ch.a -= math.Pi
		}
	}
	if mlib.Rand(999) == 9 {
		ch.SetChar(rune(availableCharCodes[mlib.GetRandomBtw(0, len(availableCharCodes)-1)]))
	}
	//if ch.a-math.Pi <= ch.rotationSpeed {
	//	}
	ch.Y -= ch.fallingSpeed
	ch.Rotate(ch.a)
}

func (ch *Chr) GetRotationAngle() float64 {
	return ch.a
}
