package characters

import (
	"assa.com/put.pixel/lib/ogl"
	"assa.com/put.pixel/src/characters/utf8"
	"math"
)

const INVISIBLE_COLOR = 0xffffffff

type Chr struct {
	Ch          rune
	Size        int
	x           int
	y           int
	Color       uint32
	shape       []byte
	bitmap      []byte
	wH          int
	wW          int
	shapeWidth  int
	shapeHeight int
	shapeSize   int
}

func (ch *Chr) SetChar(chRune rune) {
	ch.Ch = chRune
	ch.scale(ch.Size)
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
		ch.scale(size)
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
	ch.shape = nil
	ch.x = x
	ch.y = y
	ch.Color = color
	ch.Size = 14
	ch.shapeSize = ch.shapeHeight * ch.shapeWidth
	ch.generateBitmap()
	return ch
}

func (ch *Chr) MoveUnsafe(dx, dy int) {
	ch.x += dx
	ch.y += dy
	ch.ShowUnsafe()
}

func (ch *Chr) contains(x, y int) bool {
	return x < ch.x+ch.shapeWidth && x >= ch.x && y < ch.y+ch.shapeHeight && y >= ch.y
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
func (ch *Chr) ShowUnsafe() {
	if ch.IsInvisible() {
		return
	}

	ch.draw()
}

func (ch *Chr) draw() {
	ogl.PutByteBitmap(ch.x, ch.y, ch.shapeWidth*4, ch.shapeHeight, ch.bitmap)
}

func (ch *Chr) GetWidth() int {
	return ch.shapeWidth
}

func (ch *Chr) GetHeight() int {
	return ch.shapeHeight
}

func (ch *Chr) IsInvisible() bool {
	return ch.x+ch.shapeWidth <= 0 || ch.x >= ch.wW || ch.y+ch.shapeHeight <= 0 || ch.y >= ch.wH
}

func (ch *Chr) isPixelVisible(x, y int) bool {
	return x >= 0 && y >= 0 && x < ch.wW && y < ch.wH
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

func (ch *Chr) ShowInViewPort(vpx, vpy int, vpw, vph int) {
	il, ih, jl, jh := ch.getInViewPort(vpx, vpy, vpw, vph)
	ogl.PutByteBitmapBordered(ch.x-vpx+jl, ch.y-vpy+il, il, ih, jl, jh, ch.shapeWidth<<2, ch.shapeHeight, ch.bitmap)
}

func (ch *Chr) getInViewPort(vpx int, vpy int, vpw int, vph int) (int, int, int, int) {
	il := 0
	ih := ch.shapeHeight - 1
	for i := 0; i < ch.shapeHeight; i++ {
		if ch.y+i >= vpy+vph {
			ih = i - 1
			il = 0
			break
		}
		if ch.y+i < vpy {
			il = vpy - (ch.y + i)
			ih = ch.shapeHeight - 1
			break
		}
	}
	jl := 0
	jh := ch.shapeWidth - 1
	for j := 0; j < ch.shapeWidth; j++ {
		if ch.x+j >= vpx+vpw {
			jh = j - 1
			jl = 0
			break
		}
		if ch.x+j < vpx {
			jl = vpx - (ch.x + j)
			jh = ch.shapeWidth - 1
			break
		}
	}

	return il, ih, jl, jh
}
