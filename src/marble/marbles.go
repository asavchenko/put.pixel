package marble

import (
	"fmt"
	"math"
	"os"
	"runtime"
	"strings"
	"time"
)

type group struct {
	windowWidth  int
	windowHeight int
	marbleSize   int
	marbles      []Object
}

func GetNewGroup(width, height, size int) Group {
	items := make([]Object, 0)
	for j := 0; j < height/size-5; j += 1 {
		for i := 0; i < width/size; i += 1 {
			x := i*size + size/2
			y := height - size/2 - j*size
			item := GetNew(x, y, size, GetRandomColor())
			item.SetSpeed(item.R())
			item.SetAngle(math.Pi * 3 / 2)
			items = append(items, item)
		}
	}

	return &group{width, height, size, items}
}

func (g *group) IntersectsWith(cm Object) bool {
	d := cm.D() * 2
	for _, m := range g.marbles {
		if m.X() < cm.X()-d || m.X() > cm.X()+d {
			continue
		}
		if m.Y() > cm.Y()+d || m.Y() < cm.Y()-d {
			continue
		}
		if m.IntersectsWith(cm) {
			return true
		}
	}

	return false
}

func (g *group) GetLeftBorder() int {
	leftBorder := -1
	for _, m := range g.marbles {
		if leftBorder == -1 {
			leftBorder = m.CX() - m.CR()
			continue
		}
		if leftBorder > m.CX() {
			leftBorder = m.CX()
		}
	}

	return leftBorder
}

func (g *group) GetRightBorder() int {
	rightBorder := -1
	for _, m := range g.marbles {
		if rightBorder == -1 {
			rightBorder = m.CX() + m.CR()
			continue
		}
		if rightBorder < m.CX() {
			rightBorder = m.CX()
		}
	}

	return rightBorder
}

func (g *group) GetBottomBorder() int {
	bottomBorder := -1
	for _, m := range g.marbles {
		if bottomBorder == -1 {
			bottomBorder = m.CY() + m.CR()
			continue
		}
		if bottomBorder < m.CY() {
			bottomBorder = m.CY()
		}
	}

	return bottomBorder
}

func (g *group) MoveDown() {
	items := make([]Object, 0)
	for i := 0; i < g.windowWidth/g.marbleSize; i += 1 {
		x := i*g.marbleSize + g.marbleSize/2
		y := g.windowHeight - g.marbleSize/2
		item := GetNew(x, y, g.marbleSize, GetRandomColor())
		item.SetSpeed(item.R())
		item.SetAngle(math.Pi * 3 / 2)
		items = append(items, item)
	}

	for _, m := range g.marbles {
		m.Move()
	}
	g.marbles = append(items, g.marbles...)
}

func (g *group) GetIntersection(cm Object) (bool, float64, float64) {
	points := make([][]float64, 0)
	d := cm.D() * 2
	for _, m := range g.marbles {
		if m.X() < cm.X()-d || m.X() > cm.X()+d {
			continue
		}
		if m.Y() > cm.Y()+d || m.Y() < cm.Y()-d {
			continue
		}

		i, x, y := cm.WillIntersectsWith(m)
		if i {
			points = append(points, []float64{x, y})
		}
	}
	if len(points) < 1 {
		//log("no intersection points", cm.X(), cm.Y())
		return false, 0, 0
	}
	x := float64(0)
	y := float64(0)
	if len(points) > 1 {
		sumX := float64(0)
		sumY := float64(0)
		for _, p := range points {
			sumX += p[0]
			sumY += p[1]
		}
		x = sumX / float64(len(points))
		y = sumY / float64(len(points))
	} else {
		x = points[0][0]
		y = points[0][1]
	}
	log(x, y)
	r2 := float64(g.marbleSize * g.marbleSize)
	for i := 0; i < g.windowWidth/g.marbleSize; i++ {
		cx := float64(i*g.marbleSize + g.marbleSize/2)
		for j := 0; j < g.windowHeight/g.marbleSize; j++ {
			cy := float64(g.windowHeight - g.marbleSize/2 - j*g.marbleSize)
			if !g.Contains(cx, cy) && (x-cx)*(x-cx)+(y-cy)*(y-cy) < r2 {
				log("intersection point in grid is", x, y)
				return true, cx, cy
			}
		}
	}

	return false, 0, 0
}

func (g *group) Show() {
	for _, m := range g.marbles {
		m.Show()
	}
}

func (g *group) RemoveMatched(cm Object) {
	itemsToRemove := g.getMarblesToRemove(g.marbles, cm)
	log(len(itemsToRemove))
	if len(itemsToRemove) >= 3 {
		g.marbles = g.getMarblesWithout(g.marbles, itemsToRemove)
	}
}

func (g *group) getMarblesToRemove(marbles []Object, cm Object) []Object {
	d := cm.D() * 2

	toRemove := make([]Object, 0)
	for _, m := range g.getMarblesWithout(marbles, cm) {
		if m.Color() != cm.Color() {
			continue
		}
		if m.X() < cm.X()-d || m.X() > cm.X()+d {
			continue
		}
		if m.Y() > cm.Y()+d || m.Y() < cm.Y()-d {
			continue
		}

		if cm.GetDistanceTo(m) > cm.D() {
			continue
		}
		toRemove = append(toRemove, m)
	}

	for _, m := range toRemove {
		toRemove = append(toRemove, g.getMarblesToRemove(g.getMarblesWithout(marbles, toRemove), m)...)
	}

	return toRemove
}

func (g *group) Add(cm Object) {
	g.marbles = append(g.marbles, cm)
}

func (g *group) Contains(x, y float64) bool {
	for _, m := range g.marbles {
		if m.X() == x && m.Y() == y {
			return true
		}
	}

	return false
}

func (g *group) GetMarbles() []Object {
	res := make([]Object, 0)
	for _, m := range g.marbles {
		res = append(res, m)
	}

	return res
}

func (g *group) getMarblesWithout(marbles []Object, cmArr ...interface{}) []Object {
	res := make([]Object, 0)
	if len(cmArr) < 1 {
		return marbles
	}
	cmEl := cmArr[0]
	for _, m := range marbles {
		switch cm := cmEl.(type) {
		case Object:
			if m.CX() == cm.CX() && m.CY() == cm.CY() {
				continue
			}
			res = append(res, m)
		case []Object:
			found := false
			for _, _m := range cm {
				if m.CX() == _m.CX() && m.CY() == _m.CY() {
					found = true
					break
				}
			}
			if !found {
				res = append(res, m)
			}
		}
	}

	return res
}

func log(msgs ...interface{}) {
	_, f, line, _ := runtime.Caller(1)
	baseDir := getBaseDir()

	relativePath := strings.Replace(f, baseDir+"/", "", -1)
	formatedMsg := fmt.Sprintln(time.Now().UTC().Format("15:04:05.999 02-01-2006"), fmt.Sprintf("%s:%d", relativePath, line), msgs)
	fmt.Print(formatedMsg)
}

func logError(msgs ...interface{}) {
	_, f, line, _ := runtime.Caller(1)
	baseDir := getBaseDir()

	relativePath := strings.Replace(f, baseDir+"/", "", -1)
	formatedMsg := fmt.Sprintln("*********ERROR", time.Now().UTC().Format("15:04:05.999 02-01-2006"), fmt.Sprintf("%s:%d", relativePath, line), msgs)
	fmt.Print(formatedMsg)
}

func getBaseDir() string {
	pwd, err := os.Getwd()
	if err != nil {
		return ""
	}

	return pwd
}
