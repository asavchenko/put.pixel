package marble

import (
	"math"

	"assa.com/put.pixel/lib/mlib"
	"assa.com/put.pixel/lib/ogl"
	"assa.com/put.pixel/lib/pngreader"
)

const RED = byte(0)
const BLUE = byte(1)
const GREEN = byte(2)
const YELLOW = byte(3)
const PURPLE = byte(4)

type marble struct {
	x, y         float64
	bx, by       int
	gridSize     int
	r            float64 // radius
	width        int
	height       int
	windowWidth  int
	windowHeight int
	speed        float64
	a            float64 // angle
	isMoving     bool
	color        byte
	isRotated    bool
	imgData      []uint32
}

var colorPathToAssetMap = map[byte]string{
	RED:    "src/marble/red.png",
	BLUE:   "src/marble/blue.png",
	GREEN:  "src/marble/green.png",
	YELLOW: "src/marble/yellow.png",
	PURPLE: "src/marble/purple.png",
}

var colorImageMap map[byte][]uint32
var imgWidth int
var imgHeight int

func init() {
	colorImageMap = make(map[byte][]uint32)
	imgWidth = 0
	imgHeight = 0
	for color, pathToAsset := range colorPathToAssetMap {
		pngReader, err := pngreader.GetNew(pathToAsset)
		if err != nil {
			panic(err)
		}
		imgData, err := pngReader.GetImageData()
		if err != nil {
			panic(err)
		}
		imgDataRGBA := make([]uint32, len(imgData)>>2)
		for i := 0; i < len(imgData); i += 4 {
			imgDataRGBA[i/4] = uint32(imgData[i+3])<<24 | uint32(imgData[i+2])<<16 | uint32(imgData[i+1])<<8 | uint32(imgData[i])
		}
		imgWidth = pngReader.GetImageWidth()
		imgHeight = pngReader.GetImageHeight()
		flippedImgData := make([]uint32, 0)
		for j := 0; j < imgHeight; j++ {
			line := make([]uint32, imgWidth)
			startIdx := j * imgWidth
			copy(line[0:], imgDataRGBA[startIdx:startIdx+imgWidth])
			flippedImgData = append(line, flippedImgData...)
		}

		if pngReader.GetColorType() != pngreader.ColorTypeTruecolorAlpha {
			panic("unsupported color scheme for marble asset")
		}
		pngReader.Close()
		colorImageMap[color] = flippedImgData
	}
}

func GetNew(x0, y0 int, size int, color byte, windowWidth, windowHeight int, gridSize int) Object {
	data, exists := colorImageMap[color]
	if !exists {
		return nil
	}
	_copyData := make([]uint32, len(data))
	for k, v := range data {
		_copyData[k] = v
	}
	_m := &marble{
		speed:        0,
		x:            float64(x0),
		y:            float64(y0),
		color:        color,
		width:        imgWidth,
		windowWidth:  windowWidth,
		windowHeight: windowHeight,
		height:       imgHeight,
		r:            float64(max(imgWidth, imgHeight)) / 2,
		imgData:      _copyData,
		gridSize:     gridSize,
	}

	_m.resize(size)
	bx, by := _m.bXY()
	_m.bx = bx
	_m.by = by

	return _m
}

func (m *marble) SetSpeed(speed float64) {
	m.speed = speed
}

func (m *marble) GetSpeed() float64 {
	return m.speed
}

func (m *marble) resize(size int) {
	// todo
}

func (m *marble) Show() {
	x := m.CX() - m.CR()
	y := m.CY() - m.CR()
	ogl.PutBitmap(x, y, m.width, m.height, m.imgData)
}

func (m *marble) Move() {
	m.x += math.Cos(m.a) * m.speed
	m.y += math.Sin(m.a) * m.speed
	m.bx, m.by = m.bXY()
}

func (m *marble) RotateTo(x1, y1 int) {
	if x1 == m.CX() {
		return
	}
	dx := float64(x1) - m.X()
	dy := float64(y1) - m.Y()

	if x1 > m.CX() {
		m.a = math.Atan(dy / dx)
	} else {
		m.a = math.Pi + math.Atan(dy/dx)
	}

}

func (m *marble) IsInside(x, y float64) bool {
	dx := m.X() - x
	dy := m.Y() - y

	return math.Sqrt(dx*dx+dy*dy) < m.R()
}

func (m *marble) IntersectsWith(_m Object) bool {
	return m.GetDistanceTo(_m) < m.D()
}

func (m *marble) GetDistanceTo(_m ...interface{}) float64 {
	if len(_m) == 1 {
		switch t := _m[0].(type) {
		case Object:
			dx := m.X() - t.X()
			dy := m.Y() - t.Y()
			return math.Sqrt(dx*dx + dy*dy)
		}
	}
	if len(_m) == 2 {
		x := _m[0].(float64)
		y := _m[1].(float64)
		dx := m.X() - x
		dy := m.Y() - y

		return math.Sqrt(dx*dx + dy*dy)
	}

	return -1
}

func (m *marble) Rotate(da float64) {
	if m.isRotated {
		return
	}
	m.a += da
	m.isRotated = true
}

func (m *marble) SetAngle(a float64) {
	m.a = a
}

func (m *marble) X() float64 {
	return m.x
}

func (m *marble) Y() float64 {
	return m.y
}

func (m *marble) CX() int {
	return int(math.Round(m.x))
}

func (m *marble) CY() int {
	return int(math.Round(m.y))
}

func (m *marble) Color() byte {
	return m.color
}

func (m *marble) A() float64 {
	return m.a
}

func (m *marble) R() float64 {
	return m.r
}

func (m *marble) D() float64 {
	return m.R() * 2
}

func (m *marble) CR() int {
	return int(math.Round(m.R()))
}

func (m *marble) CD() int {
	return m.CR() * 2
}

func (m *marble) MoveTo(x, y float64) {
	m.x = x
	m.y = y
	m.bx, m.by = m.bXY()
}

func GetRandomColor() byte {
	return byte(mlib.GetRandomBtw(0, 4))
}

func (m *marble) Contains(x, y float64) bool {
	d := m.GetDistanceTo(x, y)

	return d >= 0 && d <= m.R()
}

func (m *marble) IsMoving() bool {
	return m.isMoving
}

func (m *marble) SetIsMoving(val bool) {
	m.isMoving = val
}
func (m *marble) bXY() (int, int) {
	// Bresenham dimension
	return m.b(m.X(), m.Y())

}

func (m *marble) BXY() (int, int) {
	// Bresenham dimension
	return m.bx, m.by
}

func (m *marble) IBXY(gridSize int, x, y int) (float64, float64) {
	gd := gridSize * 2
	gr := gridSize
	height := m.windowHeight
	width := m.windowWidth
	for j := 0; j < height/gd+1; j += 1 {
		for i := 0; i < width/gd+1; i += 1 {
			cx := i * gd
			if j%2 == 0 {
				cx += gr
			}
			if cx == 0 || cx+gr >= width {
				continue
			}
			cy := height - gr - j*gd
			if cx/gr == x && cy/gr == y {
				return float64(cx), float64(cy)
			}
		}
	}

	return 0, 0
}

func (m *marble) b(x, y float64) (int, int) {
	gd := m.gridSize * 2
	gr := m.gridSize
	height := m.windowHeight
	width := m.windowWidth
	ix := 0
	iy := 0
	deltas := make(map[int]map[int][]float64, 0)
	for j := 0; j < height/gd+1; j += 1 {
		for i := 0; i < width/gd+1; i += 1 {
			cx := i * gd
			if j%2 == 0 {
				cx += gr
			}
			if cx == 0 || cx+gr >= width {
				continue
			}
			cy := height - gr - j*gd
			if _, exists := deltas[cy/gr]; !exists {
				deltas[cy/gr] = make(map[int][]float64, 0)
			}
			deltas[cy/gr][cx/gr] = []float64{math.Abs(x - float64(cx)), math.Abs(y - float64(cy))}
		}
	}
	minDelta := float64(height+width) * 2
	// find the min delta
	for cy, row := range deltas {
		for cx, d := range row {
			if minDelta > d[0]+d[1] {
				ix = cx
				iy = cy
				minDelta = d[0] + d[1]
			}
		}
	}

	return ix, iy
}
