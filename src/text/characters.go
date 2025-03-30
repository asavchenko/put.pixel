package text

import (
	"assa.com/put.pixel/lib/log"
	"assa.com/put.pixel/src/text/utf8"
)

type area struct {
	ch        int // height btw two text rows
	color     uint32
	fontSize  int
	nw, nh    int
	minRow    int
	maxRow    int
	minColumn int
	maxColumn int
	curRow    int
	curColumn int
	r         []int
	c         []int
	chrs      []*Chr
	text      string
	vp        *viewPort
}

type viewPort struct {
	x, y int
	w, h int
}

func GetNewArea() *area {
	return &area{vp: &viewPort{}}
}

func (a *area) MoveView(dx, dy int) {
	a.vp.x += dx
	a.vp.y += dy

	if a.vp.x < 0 {
		a.vp.x = 0 // there is no text, it's the left limit, a starting position
	}
	if a.vp.y > 0 {
		a.vp.y = 0 // there is no text to go further, we are at the beginning of it
	}
	// detect if new load is needed
	// 		load
	i := a.curColumn - 9
	if i < 0 {
		i = 0
	}

	for {
		if i >= len(a.c) {
			a.curColumn = i
			break
		}
		if a.c[i] >= a.vp.x {
			a.curColumn = i
			break
		}
		i++
	}

	j := a.curRow - 9
	if j < 0 {
		j = 0
	}

	for {
		if j >= len(a.r) {
			a.curRow = j
			break
		}
		if a.r[j] < a.vp.y {
			a.curRow = j
			break
		}
		j++
	}
	if a.curRow < a.minRow+a.nh && a.minRow > 0 {
		a.load()
		return
	}
	if a.curRow > a.maxRow-a.nh {
		a.load()
		return
	}
	if a.curColumn > a.maxColumn-a.nw {
		a.load()
		return
	}
	if a.curColumn < a.minColumn+a.nw && a.minColumn > 0 {
		a.load()
		return
	}
}

func (a *area) SetViewPortWidth(w int) {
	a.vp.w = w
	a.nw = a.vp.w / GetCharacterWidth(a.fontSize)
	log.Println("nw:", a.nw)
}

func (a *area) SetViewPortHeight(h int) {
	a.vp.h = h
	a.nh = a.vp.h / a.ch
	log.Println("nh:", a.nh)
}

func (a *area) SetFontSize(fontSize int) {
	a.fontSize = fontSize
	a.ch = GetLineSpaceSize(fontSize) + GetCharacterHeight(fontSize)
	a.nw = a.vp.w / GetCharacterWidth(fontSize)
	a.nh = a.vp.h / a.ch
	log.Println("fontSize:", a.fontSize, "a.ch:", a.ch, "a.nw:", a.nw, "a.nh:", a.nh)
	if a.text != "" {
		log.Println("calling init")
		a.init()
	}
}

func (a *area) SetFontColor(color uint32) {
	a.color = color
	for _, ch := range a.chrs {
		ch.Color = color
	}
}

func (a *area) SetText(text string) {
	log.Println("loading text")
	a.text = text
	a.init()
}

func (a *area) init() {
	a.curRow = 0
	a.curColumn = 0
	a.r = make([]int, 0)
	a.c = make([]int, 0)
	a.load()
}

func (a *area) GetCharacters() []*Chr {
	result := make([]*Chr, 0)

	for _, ch := range a.chrs {
		if a.vp.inViewPortCh(ch) {
			result = append(result, ch)
		}
	}

	return result
}

func (a *area) load() {
	a.minRow = a.curRow - a.nh*5
	if a.minRow < 0 {
		a.minRow = 0
	}
	a.maxRow = a.minRow + a.nh*10
	a.minColumn = a.curColumn - a.nw*5
	if a.minColumn < 0 {
		a.minColumn = 0
	}
	a.maxColumn = a.minColumn + a.nw*10
	// recalculate trigger next load area:
	a.chrs = make([]*Chr, 0, a.nw*a.nh*100)
	y := a.vp.h - a.ch
	x := 0
	i := 0
	isNewLine := true
	curRow := 0
	curColumn := 0
	runes := []rune(a.text)
	for {
		if i > len(runes)-1 {
			break
		}
		r := runes[i]
		if int(r) == 10 {
			y -= a.ch
			// new line
			isNewLine = true
			if int(runes[i+1]) == 9 {
				i += 2
			} else {
				i += 1
			}
			if curRow >= len(a.r) {
				a.r = append(a.r, y)
			}
			curRow++

			if curRow > a.maxRow {
				break
			}
			continue
		}
		i += 1
		if isNewLine {
			curColumn = 0
			x = 0
		}
		if curRow < a.minRow {
			continue
		}
		if curColumn < a.minColumn || curColumn > a.maxColumn {
			// we still need to advance x properly
			ch := GetNew(r, x, y, a.color).SetCharacterSize(a.fontSize)
			if curColumn >= len(a.c) {
				a.c = append(a.c, x)
			} else if x < a.c[curColumn] {
				a.c[curColumn] = x
			}
			curColumn++
			x += ch.GetWidth()
			continue
		}

		ch := GetNew(r, x, y, a.color).SetCharacterSize(a.fontSize)
		a.chrs = append(a.chrs, ch)
		isNewLine = false
		if curColumn >= len(a.c) {
			a.c = append(a.c, x)
		} else if x < a.c[curColumn] {
			a.c[curColumn] = x
		}
		x += ch.GetWidth()
		curColumn++
	}
	log.Println("num characters:", len(a.chrs))
}

func (a *area) GetViewPortPosition() (int, int) {
	return a.vp.x, a.vp.y
}

func (a *area) GetViewPortWidthHeight() (int, int) {
	return a.vp.w, a.vp.h
}

func (vp *viewPort) inViewPort(x, y int) bool {
	return x >= vp.x && x < vp.x+vp.w && y >= vp.y && y < vp.y+vp.h
}

func (vp *viewPort) inViewPortCh(ch *Chr) bool {
	return !(ch.x+ch.GetWidth() <= vp.x || ch.x >= vp.x+vp.w || ch.y+ch.GetHeight() <= vp.y || ch.y >= vp.y+vp.h)
}

func GetCharacterWidth(size int) int {
	switch size {
	case 14:
		return utf8.GetShapeWidth()
	default:
		return utf8.GetShapeWidth() * size / 14
	}
}

func GetCharacterHeight(size int) int {
	switch size {
	case 14:
		return utf8.GetShapeHeight()
	default:
		return utf8.GetShapeHeight() * size / 14
	}
}

func GetSpaceSizeBtwCharacters(size int) int {
	return 1
}

func GetLineSpaceSize(size int) int {
	return 2
}

func GetAvailableCharCodes() []uint32 {
	return utf8.GetAvailableCharCodes()
}
