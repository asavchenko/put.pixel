package characters

import (
	"assa.com/put.pixel/lib/log"
	"assa.com/put.pixel/src/characters/utf8"
)

type triggerPoint struct {
	x, y int
	r, c int
}

type area struct {
	isFirstLoad bool
	ch          int // height btw two text rows
	color       uint32
	fontSize    int
	x, y        int
	w, h        int
	nw, nh      int
	tp          [2][2]triggerPoint
	chrs        []*Chr
	text        string
	vp          *viewPort
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
	//log.Println(a.vp.x, a.vp.y)
	//log.Println(a.tp[0][0].x, a.tp[0][0].y)
	//log.Println(a.tp[0][1].x, a.tp[0][1].y)
	//log.Println(a.tp[1][0].x, a.tp[1][0].y)
	//log.Println(a.tp[1][1].x, a.tp[1][1].y)

	if a.vp.x < 0 {
		a.vp.x = 0 // there is no text, it's the left limit, a starting position
	}
	if a.vp.y > 0 {
		a.vp.y = 0 // there is no text to go further, we are at the beginning of it
	}
}

func (a *area) SetPosition(x, y int) {
	a.x = x
	a.y = y
}

func (a *area) SetWidth(w int) {
	a.w = w
}

func (a *area) SetHeight(h int) {
	a.h = h
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

func (a *area) GetCharacters() []*Chr {
	vp := a.vp
	result := make([]*Chr, 0)
	if vp.inViewPort(a.tp[0][0].x, a.tp[0][0].y) {
		a.load(a.tp[0][0])
	} else if vp.inViewPort(a.tp[0][1].x, a.tp[0][1].y) {
		a.load(a.tp[0][1])
	} else if vp.inViewPort(a.tp[1][0].x, a.tp[1][0].y) {
		a.load(a.tp[1][0])
	} else if vp.inViewPort(a.tp[1][1].x, a.tp[1][1].y) {
		a.load(a.tp[1][1])
	}
	for _, ch := range a.chrs {
		if vp.inViewPortCh(ch) {
			result = append(result, ch)
		}
	}

	return result
}

func (a *area) load(tp triggerPoint) {
	log.Println("loading more, trigger point reached")
	minColumn := tp.c - a.nw*5/2
	if minColumn < 0 {
		minColumn = 0
	}
	maxColumn := minColumn + a.nw*5

	minRow := tp.r - a.nh*5/2
	if minRow < 0 {
		minRow = 0
	}
	maxRow := minRow + a.nh*5

	xmin := tp.x - a.w*5/2
	if xmin < 0 {
		xmin = 0
	}
	xmax := xmin + a.w*5

	ymax := tp.y + a.h*5/2
	if ymax > 0 {
		ymax = 0
	}
	ymin := ymax - a.w*5
	log.Println("nw:", a.nw, "nh:", a.nh)
	log.Println("xmin:", xmin, "ymin:", ymin, "xmax:", xmax, "ymax:", ymax, "minRow:", minRow, "maxRow:", maxRow, "minColumn:", minColumn, "maxColumn:", maxColumn)
	a.x = xmin
	a.y = ymin
	// recalculate trigger Points:
	a.tp[0][0].x = xmin + a.vp.w
	a.tp[0][0].y = ymin + a.vp.h
	a.tp[0][0].r = maxRow - a.nh
	a.tp[0][0].c = minColumn + a.nw

	a.tp[0][1].x = xmin + a.vp.w
	a.tp[0][1].y = ymax - a.vp.h
	a.tp[0][1].r = minRow + a.nh
	a.tp[0][1].c = minColumn + a.nw

	a.tp[1][0].x = xmax - a.vp.w
	a.tp[1][0].y = ymin + a.vp.h
	a.tp[1][0].r = maxRow - a.nh
	a.tp[1][0].c = maxColumn - a.nw

	a.tp[1][1].x = xmax - a.vp.w
	a.tp[1][1].y = ymax - a.vp.h
	a.tp[1][1].r = minRow + a.nh
	a.tp[1][1].c = maxColumn - a.nw

	a.chrs = make([]*Chr, 0)
	y := a.vp.h - a.ch
	x := xmin
	i := 0
	isNewLine := true
	curRow := minRow
	curColumn := minColumn
	for {
		if i > len([]rune(a.text))-1 {
			break
		}
		r := []rune(a.text)[i]
		if int(r) == 10 {
			y -= a.ch
			if y < 0 {
				break
			}
			// new line
			isNewLine = true
			if int([]rune(a.text)[i+1]) == 9 {
				i += 2
			} else {
				i += 1
			}
			curRow++
			if curRow > maxRow {
				break
			}
			continue
		}
		i += 1
		if isNewLine {
			curColumn = 0
			x = 0
		}
		if curRow < minRow {
			continue
		}
		if curColumn < minColumn || curColumn > maxColumn {
			continue
		}

		//log.Println("ch.X:", x, "ch.Y:", y, "ch:", r)
		ch := GetNew(r, x, y, a.color).SetCharacterSize(a.fontSize)
		a.chrs = append(a.chrs, ch)
		isNewLine = false
		x += ch.GetWidth()
		curColumn++
	}
}

func (a *area) init() {
	if a.isFirstLoad {
		a.load(triggerPoint{
			x: a.vp.x,
			y: a.vp.y,
			r: 0,
			c: 0,
		})
		return
	}
	a.load(triggerPoint{
		x: (a.tp[0][0].x + a.tp[0][1].x) / 2,
		y: (a.tp[0][0].y + a.tp[0][1].y) / 2,
		r: (a.tp[0][0].r + a.tp[0][1].r) / 2,
		c: (a.tp[0][0].c + a.tp[1][0].c) / 2,
	})
	return
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
	return !(ch.x+ch.GetWidth() <= vp.x || ch.x >= vp.x+vp.w || ch.y+ch.GetHeight() <= 0 || ch.y >= vp.y+vp.h)
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
