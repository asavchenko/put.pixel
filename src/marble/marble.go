package marble

import (
	"assa.com/put.pixel/lib/mlib"
	"assa.com/put.pixel/lib/ogl"
	"assa.com/put.pixel/lib/pngreader"
	"math"
)

const RED = byte(0)
const BLUE = byte(1)
const GREEN = byte(2)
const YELLOW = byte(3)
const PURPLE = byte(4)

type marble struct {
	x, y    float64
	r       float64 // radius
	width   int
	height  int
	speed   float64
	a       float64 // angle
	color   byte
	imgData []byte
}

var colorPathToAssetMap = map[byte]string{
	RED:    "src/marble/red.png",
	BLUE:   "src/marble/blue.png",
	GREEN:  "src/marble/green.png",
	YELLOW: "src/marble/yellow.png",
	PURPLE: "src/marble/purple.png",
}

var colorImageMap map[byte][]byte
var imgWidth int
var imgHeight int

func init() {
	colorImageMap = make(map[byte][]byte)
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
		imgWidth = pngReader.GetImageWidth()
		imgHeight = pngReader.GetImageHeight()
		if pngReader.GetColorType() != pngreader.ColorTypeTruecolorAlpha {
			panic("unsupported color scheme for marble asset")
		}
		pngReader.Close()
		colorImageMap[color] = imgData
	}
}

func GetNew(x0, y0 int, size int, color byte) Object {
	data, exists := colorImageMap[color]
	if !exists {
		return nil
	}
	_copyData := make([]byte, len(data))
	for k, v := range data {
		_copyData[k] = v
	}
	_m := &marble{
		speed:   0,
		x:       float64(x0),
		y:       float64(y0),
		color:   color,
		width:   imgWidth,
		height:  imgHeight,
		r:       float64(max(imgWidth, imgHeight)) / 2,
		imgData: _copyData,
	}

	_m.resize(size)

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
	x := int(m.x - m.r)
	y := int(m.y + m.r)
	for i := 0; i < m.height; i++ {
		for j := 0; j < m.width; j += 1 {
			index := i*m.width*4 + j*4
			if index >= len(m.imgData) {
				break
			}
			if m.imgData[index+3] != 0 {
				ogl.PutPixel(x+j, y-i, m.imgData[index], m.imgData[index+1], m.imgData[index+2])
			}
		}
	}
}

func (m *marble) Move() {
	m.x = math.Cos(m.a)*m.speed + m.x
	m.y = math.Sin(m.a)*m.speed + m.y
}

func (m *marble) RotateTo(x1, y1 int) {
	if x1 == m.CX() {
		return
	}

	if x1 > m.CX() {
		m.a = math.Atan(float64(y1-m.CY()) / float64(x1-m.CX()))
	} else {
		m.a = math.Pi + math.Atan(float64(y1-m.CY())/float64(x1-m.CX()))
	}

}

func (m *marble) IsInside(x, y float64) bool {
	return math.Sqrt(math.Pow(m.X()-x, 2)+math.Pow(m.Y()-y, 2)) < m.R()
}

func (m *marble) IntersectsWith(_m Object) bool {
	return m.GetDistanceTo(_m) < m.D()
}

func (m *marble) WillIntersectsWith(cm Object) (bool, float64, float64) {
	i, x, y := m.willCenterIntersectsWith(cm)
	if !i {
		i, x, y = m.willBottomIntersectsWith(cm)
		if !i {
			i, x, y = m.willTopIntersectsWith(cm)
			if !i {
				return false, 0, 0
			} else {
				return true, x, y
			}
		} else {
			return true, x, y
		}
	} else {
		return true, x, y
	}
}

func (m *marble) willTopIntersectsWith(cm Object) (bool, float64, float64) {
	x := math.Cos(m.a+math.Pi/2)*m.R() + m.x
	y := math.Sin(m.a+math.Pi/2)*m.R() + m.y

	return m.willPointIntersectsWith(x, y, cm)
}

func (m *marble) willBottomIntersectsWith(cm Object) (bool, float64, float64) {
	x := math.Cos(m.a-math.Pi/2)*m.R() + m.x
	y := math.Sin(m.a-math.Pi/2)*m.R() + m.y

	return m.willPointIntersectsWith(x, y, cm)
}

func (m *marble) willCenterIntersectsWith(cm Object) (bool, float64, float64) {
	return m.willPointIntersectsWith(m.X(), m.Y(), cm)
}

func (m *marble) willPointIntersectsWith(x, y float64, cm Object) (bool, float64, float64) {
	x = math.Cos(m.a)*m.speed + x
	y = math.Sin(m.a)*m.speed + y
	ta := math.Tan(m.A())
	a := ta*ta + 1
	if a == 0 {
		//log(ta*ta+1, "=", 0)
		return false, 0, 0
	}
	c1 := y - cm.Y() - ta*x
	b := 2*ta*c1 - 2*cm.X()
	c := c1*c1 + cm.X()*cm.X() - cm.R()*cm.R()
	ds := b*b - 4*a*c
	if ds < 0 {
		//log(ds, "<", 0, "b=", b, "4*a*c=", 4*a*c)
		return false, 0, 0
	}
	ds = math.Sqrt(ds)

	x1 := (-b + ds) / 2 / a
	x2 := (-b - ds) / 2 / a
	if x1 < x || x1 > x+m.GetSpeed() {
		if x2 < x || x2 > x+cm.GetSpeed() {
			//log(x2, "<", x, "||", x2, ">", x+cm.GetSpeed())
			return false, 0, 0
		} else {
			//log("INTERSECTS! at", x2, ta*(x2-x)+y)
			return true, x2, ta*(x2-x) + y
		}
	} else {
		//log("INTERSECTS! at", x1, ta*(x1-x)+y)
		return true, x1, ta*(x1-x) + y
	}
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
	m.a += da
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
	return m.r * 2
}

func (m *marble) CR() int {
	return int(math.Round(m.r))
}

func (m *marble) CD() int {
	return m.CR() * 2
}

func (m *marble) MoveTo(x, y float64) {
	m.x = x
	m.y = y
}

func GetRandomColor() byte {
	return byte(mlib.GetRandomBtw(0, 4))
}

func (m *marble) Contains(x, y float64) bool {
	d := m.GetDistanceTo(x, y)
	return d >= 0 && d <= m.R()
}
