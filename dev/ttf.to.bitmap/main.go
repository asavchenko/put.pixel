package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"runtime"
	"strings"
	"syscall"
	"time"

	"assa.com/put.pixel/lib/freetype"
	"assa.com/put.pixel/lib/mlib"
	"assa.com/put.pixel/lib/ogl"
)

func init() {
	fmt.Println("init")
	syscall.Setpriority(syscall.PRIO_PROCESS, os.Getpid(), -20)
	runtime.LockOSThread()
}

func main() {
	if f, err := freetype.LoadFont("VaryAlongQuads.ttf"); err != nil {
		//if f, err := freetype.LoadFont("Envy Code R.ttf"); err != nil {
		//if f, err := freetype.LoadFont("LiberationMono-Regular.ttf"); err != nil {
		logError(err)
	} else {
		//fmt.Println(f.GetGlyphIndex('Z'))
		//fmt.Println(f.GetGlyphOffset('A'))
		fmt.Println(f.GetGlyphOffset('6'))
		for _, c := range f.GetAvailableCharCodes() {
			b := make([]byte, 4)
			binary.BigEndian.PutUint32(b, c)
			fmt.Println(string(rune(c)), printBytes(b), "idx: ", f.GetGlyphIndex(uint(c)))
		}
		gd := f.GetGlyphData('c') // VaryAlongQuads.ttf
		//gd := f.GetGlyphData('A')
		//gd := f.GetGlyphData('9')
		if gd == nil {
			fmt.Println("something is not right...")
		} else {
			fmt.Println(gd.XMin, "<=x<=", gd.XMax, gd.YMin, "<=y<=", gd.YMax)
			for _, c := range gd.Points {
				for _, gp := range c {
					log("x:", gp.X, "y:", gp.Y, gp.OnCurve)
				}
				fmt.Println("--------")
			}
			showPoints(f, gd)
		}
		fmt.Println("")
	}
	fmt.Println("it works!")

}

func showPoints(f freetype.Font, gd *freetype.GlyphData) {
	ogl.Init(false)
	defer ogl.Close()
	ogl.OnKeypress(ogl.KEY_ESC, func() {
		ogl.CloseWindow()
	})

	// gd.XMax - gd.Xmin --------------- w - dx
	// x - gd.Xmin                       -?
	// gd.YMax - gd.Ymin ----------------h - dy
	debugged := false
	for {
		if ogl.IsExit() {
			break
		}
		ogl.Draw(func() {
			// draw points here!!!
			//drawPointsOnly(gd, debugged)
			//drawPointsWithLines(gd, debugged)
			//drawPointsWithBezier(gd, debugged)
			fillGlyph(f, gd, debugged)
			debugged = true
		})
	}
}

func fillGlyph(f freetype.Font, gd *freetype.GlyphData, debugged bool) {
	dx := 240
	dy := 110
	w := ogl.GetWindowWidth() - dx
	h := ogl.GetWindowHeight() - dy

	r, g, b := getColor(0, true)
	for y := dy; y < h; y++ {
		for x := dx; x < w; x++ {
			// we need to scale x, y to glyph space
			_x := int(math.Round((float64(x)*float64(gd.XMax-gd.XMin) + float64(w)*float64(gd.XMin) - float64(dx)*float64(gd.XMax)) / float64(w-dx)))
			_y := int(math.Round((float64(y)*float64(gd.YMax-gd.YMin) + float64(h)*float64(gd.YMin) - float64(dy)*float64(gd.YMax)) / float64(h-dy)))
			//if !debugged {
			//	log(_x, _y, x, y, scaleX, scaleY)
			//}
			if gd.ContainsPoint(_x, _y) {
				//if !debugged {
				//	fmt.Println(_x, _y, gd.XMax, gd.YMax)
				//}
				ogl.PutPixel(x, y, r, g, b)
			}
		}
	}
	if !debugged {
		log("done")
	}
}

func drawPointsOnly(gd *freetype.GlyphData, debugged bool) {
	dx := 240
	dy := 90
	w := ogl.GetWindowWidth() - dx
	h := ogl.GetWindowHeight() - dy
	for contourIdx, c := range gd.GetWithMiddlePoints() {
		for _, gp := range c {
			r, g, b := getColor(contourIdx, gp.OnCurve)

			x_ := (w-dx)*int(gp.X-gd.XMin)/int(gd.XMax-gd.XMin) + dx
			y_ := (h-dy)*int(gp.Y-gd.YMin)/int(gd.YMax-gd.YMin) + dy

			ogl.PutPixel(x_, y_, r, g, b)
			if !debugged && gp.OnCurve {
				log("x:", gp.X, "y:", gp.Y, "|", "x_:", x_, "y_:", y_)
			}
		}
		if !debugged {
			fmt.Println("-----")
		}

	}
}

func drawPointsWithBezier(gd *freetype.GlyphData, debugged bool) {
	dx := 240
	dy := 90
	w := ogl.GetWindowWidth() - dx
	h := ogl.GetWindowHeight() - dy
	for contourIdx, c := range gd.GetWithMiddlePoints() {
		i := 0
		for {
			if i+1 >= len(c) {
				break
			}
			//if !debugged {
			//	log("x:", c[i].x, "y:", c[i].y, "|", c[i].OnCurve)
			//}
			r, g, b := getColor(contourIdx, true)
			if c[i+1].OnCurve && c[i].OnCurve {
				//if !debugged {
				//	log("both points are on curve")
				//}
				ogl.Line((w-dx)*int(c[i].X-gd.XMin)/int(gd.XMax-gd.XMin)+dx, (h-dy)*int(c[i].Y-gd.YMin)/int(gd.YMax-gd.YMin)+dy,
					(w-dx)*int(c[i+1].X-gd.XMin)/int(gd.XMax-gd.XMin)+dx, (h-dy)*int(c[i+1].Y-gd.YMin)/int(gd.YMax-gd.YMin)+dy,
					r, g, b)
				i += 1
				continue
			}
			if i+2 < len(c) {
				if !c[i+1].OnCurve && c[i].OnCurve {
					//if !debugged {
					//	log("it looks like it's Bezier", "xa:", c[i].x, "ya:", c[i].y, "xb:", c[i+2].x, "yb:", c[i+2].y, "xc:", c[i+1].x, "yc:", c[i+1].y, c[i].OnCurve, c[i+1].OnCurve, c[i+2].OnCurve)
					//}
					points := mlib.GetBezierCoords2(mlib.Point2D{
						X: float64(c[i].X),
						Y: float64(c[i].Y),
					},
						mlib.Point2D{
							X: float64(c[i+2].X),
							Y: float64(c[i+2].Y),
						},
						mlib.Point2D{
							X: float64(c[i+1].X),
							Y: float64(c[i+1].Y),
						},
						5)
					j := 0
					for {
						if j+1 >= len(points) {
							break
						}
						//if !debugged {
						//	log("x1:", points[j].x, "y1:", points[j].y)
						//	log("x2:", points[j+1].x, "y2:", points[j+1].y)
						//}
						ogl.Line((w-dx)*(int(points[j].X)-int(gd.XMin))/int(gd.XMax-gd.XMin)+dx, (h-dy)*(int(points[j].Y)-int(gd.YMin))/int(gd.YMax-gd.YMin)+dy,
							(w-dx)*(int(points[j+1].X)-int(gd.XMin))/int(gd.XMax-gd.XMin)+dx, (h-dy)*(int(points[j+1].Y)-int(gd.YMin))/int(gd.YMax-gd.YMin)+dy,
							r, g, b)
						j += 1
					}
					i += 2
					continue
				}
			}
			//log("something is not right", c[i].OnCurve, c[i+1].OnCurve)
			i += 1
		}
		//if !debugged {
		//	log("----------")
		//}
	}
	//if !debugged {
	//	log("done")
	//}
}

func drawPointsWithLines(gd *freetype.GlyphData, debugged bool) {
	dx := 240
	dy := 90
	w := ogl.GetWindowWidth() - dx
	h := ogl.GetWindowHeight() - dy
	for contourIdx, c := range gd.GetWithMiddlePoints() {
		px := -1
		py := -1
		x0 := -1
		y0 := -1
		for _, gp := range c {
			r, g, b := getColor(contourIdx, gp.OnCurve)

			x_ := (w-dx)*int(gp.X-gd.XMin)/int(gd.XMax-gd.XMin) + dx
			y_ := (h-dy)*int(gp.Y-gd.YMin)/int(gd.YMax-gd.YMin) + dy

			//ogl.PutPixel(x_, y_, r, g, b)
			if !debugged && gp.OnCurve {
				log("x:", gp.X, "y:", gp.Y, "|", "x_:", x_, "y_:", y_)
			}
			if px < 0 || !gp.OnCurve {
				ogl.PutPixel(x_, y_, r, g, b)
			} else {
				ogl.Line(px, py, x_, y_, r, g, b)

				if !debugged {
					log("x0:", px, "y0:", py, "|", "x1:", x_, "y2:", y_)
				}
				if x0 < 0 {
					x0 = px
					y0 = py
				}
			}
			if gp.OnCurve {
				px = x_
				py = y_
			}
		}
		r, g, b := getColor(contourIdx, true)
		ogl.Line(px, py, x0, y0, r, g, b)
		if !debugged {
			log("x0:", px, "y0:", py, "|", "x1:", x0, "y2:", y0)
		}
		if !debugged {
			fmt.Println("-----")
		}

	}
}

func getColor(contourIdx int, onCurve bool) (byte, byte, byte) {
	r := byte(0x90)
	g := byte(0xee)
	b := byte(0x90)
	switch contourIdx {
	case 0:
		r = byte(0x90)
		g = byte(0xee)
		b = byte(0x90)
		if !onCurve {
			r = byte(0xf0)
			g = byte(0x80)
			b = byte(0x80)
		}
	case 1:
		r = byte(0x8f)
		g = byte(0xbc)
		b = byte(0x8f)
		if !onCurve {
			r = byte(0xfa)
			g = byte(0x80)
			b = byte(0x72)
		}
	case 2:
		r = byte(0xad)
		g = byte(0xff)
		b = byte(0x2f)
		if !onCurve {
			r = byte(0xe9)
			g = byte(0x96)
			b = byte(0x7a)
		}
	case 3:
		r = byte(0x00)
		g = byte(0xff)
		b = byte(0x00)
		if !onCurve {
			r = byte(0xff)
			g = byte(0x63)
			b = byte(0x47)
		}
	case 4:
		r = byte(0x00)
		g = byte(0xff)
		b = byte(0x7f)
		if !onCurve {
			r = byte(0xcd)
			g = byte(0x5c)
			b = byte(0x5c)
		}
	case 5:
		r = byte(0x7f)
		g = byte(0xff)
		b = byte(0x00)
		if !onCurve {
			r = byte(0xff)
			g = byte(0x45)
			b = byte(0x00)
		}
	case 6:
		r = byte(0x32)
		g = byte(0xcd)
		b = byte(0x32)
		if !onCurve {
			r = byte(0xdc)
			g = byte(0x14)
			b = byte(0x3c)
		}
	default:
		if !onCurve {
			r = byte(0xf0)
			g = byte(0x80)
			b = byte(0x80)
		}
	}
	return r, g, b
}

func printBytes(bytes []byte) string {
	msg := ""
	for _, b := range bytes {
		msg += fmt.Sprintf("0x%02X ", b)
	}
	return msg
}

func log(msgs ...interface{}) {
	_, f, line, _ := runtime.Caller(1)
	baseDir := ""
	pwd, err := os.Getwd()
	if err == nil {
		baseDir = pwd
	}

	relativePath := strings.Replace(f, baseDir+"/", "", -1)
	formatedMsg := fmt.Sprintln(time.Now().UTC().Format("15:04:05.999 02-01-2006"), fmt.Sprintf("%s:%d", relativePath, line), msgs)
	fmt.Print(formatedMsg)
}

func logError(msgs ...interface{}) {
	_, f, line, _ := runtime.Caller(1)
	baseDir := ""
	pwd, err := os.Getwd()
	if err == nil {
		baseDir = pwd
	}

	relativePath := strings.Replace(f, baseDir+"/", "", -1)
	formatedMsg := fmt.Sprintln("*********ERROR", time.Now().UTC().Format("15:04:05.999 02-01-2006"), fmt.Sprintf("%s:%d", relativePath, line), msgs)
	fmt.Print(formatedMsg)
}
