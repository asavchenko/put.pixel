package marble

import (
	"assa.com/put.pixel/lib/mlib"
	"fmt"
	"math"
	"os"
	"runtime"
	"strings"
	"time"
)

const DIRECTION_TO_THE_RIGHT = 0
const DIRECTION_TO_THE_TOP = 1
const DIRECTION_TO_THE_LEFT = 2
const DIRECTION_TO_THE_BOTTOM = 3
const DIRECTION_STOP = 4

type group struct {
	windowWidth  int
	windowHeight int
	marbleSize   int
	gridCellSize int
	marbles      []Object
	bcells       [][]int
	bscreen      [][]int
}

func GetNewGroup(width, height, d int, gridSize int) Group {
	gd := gridSize * 2
	items := make([]Object, 0)
	bcells := make([][]int, 0)
	bscreen := make([][]int, height/gridSize)
	for j := 0; j < height/gd+1; j += 1 {
		line := make([]int, width/gridSize)
		y := height - gridSize - j*gd
		if y/gridSize < 1 {
			continue
		}
		for i := 0; i < width/gd+1; i += 1 {
			x := i * gd
			if j%2 == 0 {
				x += gridSize
			}
			if x == 0 || x+gridSize >= width {
				continue
			}
			line[x/gridSize] = 1
		}
		bscreen[y/gridSize] = line
	}
	for j := 0; j < height/gd-5; j += 1 {
		y := height - gridSize - j*gd
		for i := 0; i < width/gd+1; i += 1 {
			x := i * gd
			if j%2 == 0 {
				x += gridSize
			}
			if x == 0 || x+gridSize >= width {
				continue
			}
			item := GetNew(x, y, d, GetRandomColor(), width, height, gridSize)
			item.SetSpeed(item.R())
			item.SetAngle(math.Pi * 3 / 2)
			items = append(items, item)
			fmt.Print(y/gridSize, x/gridSize, " | ")
			bcells = append(bcells, []int{x / gridSize, y / gridSize})
		}
		if j%2 == 0 {
			fmt.Print("\n   ")
		} else {
			fmt.Println("")
		}
	}

	return &group{width, height, d, gridSize, items, bcells, bscreen}
}

func (g *group) MoveDown() {
	// todo
}

func (g *group) GetBIntersection(cm Object) (bool, float64, float64) {
	x0, y0 := cm.BXY()
	cm.Move()
	x1, y1 := cm.BXY()
	if x0 == x1 && y1 == y0 {
		return false, 0, 0
	}

	dx := mlib.AbsInt(x1 - x0)
	sx := 1
	if x0 >= x1 {
		sx = -1
	}
	dy := -mlib.AbsInt(y1 - y0)
	sy := 1
	if y0 >= y1 {
		sy = -1
	}
	e := dx + dy

	for {
		if g.inBcells(x0, y0) || y0 >= len(g.bscreen)-1 {
			cx, cy := g.GetFirstFreeBCell(x0, y0)
			x, y := cm.IBXY(g.gridCellSize, cx, cy)
			return true, x, y
		}
		e2 := e << 2
		if e2 >= dy {
			if x0 == x1 {
				break
			}
			e += dy
			x0 += sx
		}
		if e2 <= dx {
			if y0 == y1 {
				break
			}
			e += dx
			y0 += sy
		}
	}

	return false, 0, 0
}

func (g *group) isBXYValid(x, y int) bool {
	if y < 0 || x < 0 || y >= len(g.bscreen) {
		return false
	}
	if len(g.bscreen[y]) < 1 || x >= len(g.bscreen[y]) {
		return false
	}

	return g.bscreen[y][x] == 1
}

func (g *group) GetFirstFreeBCell(x, y int) (int, int) {
	log(x, y, x+2, y+2, x-2, y-2)
	if g.isBXYValid(x, y) && !g.inBcells(x, y) {
		log(x, y)
		return x, y
	}
	if g.isBXYValid(x, y+2) && !g.inBcells(x, y+2) {
		log(x, y+2)
		return x, y + 2
	}
	if g.isBXYValid(x, y-2) && !g.inBcells(x, y-2) {
		log(x, y-2)
		return x, y - 2
	}

	if g.isBXYValid(x-1, y) && !g.inBcells(x-1, y) {
		log(x-1, y)
		return x - 1, y
	}

	if g.isBXYValid(x-1, y-2) && !g.inBcells(x-1, y-2) {
		log(x-1, y-2)
		return x - 1, y - 2
	}
	if g.isBXYValid(x+1, y-2) && !g.inBcells(x+1, y-2) {
		log(x+1, y-2)
		return x + 1, y - 2
	}
	if g.isBXYValid(x+1, y) && !g.inBcells(x+1, y) {
		log(x+1, y)
		return x + 1, y
	}
	if g.isBXYValid(x+1, y-2) && !g.inBcells(x+1, y-2) {
		log(x+1, y-2)
		return x + 1, y - 2
	}
	if g.isBXYValid(x-1, y+2) && !g.inBcells(x-1, y+2) {
		log(x-1, y+2)
		return x - 1, y + 2
	}
	if g.isBXYValid(x+1, y+2) && !g.inBcells(x+1, y+2) {
		log(x+1, y+2)
		return x + 1, y + 2
	}

	if g.isBXYValid(x-2, y) && !g.inBcells(x-2, y) {
		log(x-2, y)
		return x - 2, y
	}

	if g.isBXYValid(x+2, y) && !g.inBcells(x+2, y) {
		log(x+2, y)
		return x + 2, y
	}

	if g.isBXYValid(x-2, y+2) && !g.inBcells(x-2, y+2) {
		log(x-2, y+2)
		return x - 2, y + 2
	}

	if g.isBXYValid(x+2, y+2) && !g.inBcells(x+2, y+2) {
		log(x+2, y+2)
		return x + 2, y + 2
	}

	if g.isBXYValid(x-2, y-2) && !g.inBcells(x-2, y-2) {
		log(x-2, y-2)
		return x - 2, y - 2
	}

	if g.isBXYValid(x+2, y-2) && !g.inBcells(x+2, y-2) {
		log(x+2, y-2)
		return x + 2, y - 2
	}

	return x, y
}

func (g *group) Show() {
	for _, m := range g.marbles {
		m.Show()
	}
}

func (g *group) RemoveMatched(cm Object) {
	itemsToRemove := g.getMarblesToRemove(g.marbles, cm)
	if len(itemsToRemove) < 3 {
		return
	}
	bcells := make([][]int, 0)
	bcellsToRemove := make([][]int, 0)
	g.marbles = g.getMarblesWithout(g.marbles, itemsToRemove)
	for _, m := range itemsToRemove {
		mbx, mby := m.BXY()
		for _, bc := range g.bcells {
			if bc[0] != mbx || mby != bc[1] {
				continue
			}
			bcellsToRemove = append(bcellsToRemove, bc)
		}
	}
	for _, bc := range g.bcells {
		found := false
		for _, bcr := range bcellsToRemove {
			if bcr[0] == bc[0] && bcr[1] == bc[1] {
				found = true
				break
			}
		}
		if !found {
			bcells = append(bcells, bc)
		}
	}

	g.bcells = bcells
	g.removeHanging()
}

func (g *group) removeHanging() {
	labels := make([][]int, len(g.bscreen))
	for y := 0; y < len(g.bscreen); y++ {
		labels[y] = make([]int, len(g.bscreen[y]))
		for x := 0; x < len(g.bscreen[y]); x++ {
			labels[y][x] = -1
		}
	}
	label := 0
	for y := 0; y < len(g.bscreen); y++ {
		for x := 0; x < len(g.bscreen[y]); x++ {
			l := g.getBNeighborsLabel(x, y, labels, label)
			if l < 0 {
				continue
			}
			if label < l {
				label = l
			}
			labels[y][x] = l
		}
	}

	labelSets := make(map[int][]int, 0)
	for y := 0; y < len(g.bscreen); y++ {
		for x := 0; x < len(g.bscreen[y]); x++ {
			l := labels[y][x]
			if l == -1 {
				continue
			}
			if _, exists := labelSets[l]; !exists {
				labelSets[l] = getUniqueIntSlice(g.getBNeighborLabels(x, y, labels))
			} else {
				labelSets[l] = getUniqueIntSlice(append(labelSets[l], g.getBNeighborLabels(x, y, labels)...))
			}
		}
	}
	//log(labelSets)
	for k, arr := range labelSets {
		for _, a := range arr {
			labelSets[k] = append(labelSets[k], labelSets[a]...)
		}
		labelSets[k] = getUniqueIntSlice(labelSets[k])
	}
	labelMapping := make(map[int]int, 0)
	for k, arr := range labelSets {
		if len(arr) == 0 {
			labelMapping[k] = k
			continue
		}
		labelMapping[k] = minIntSlice(arr)
	}
	//log(labelMapping)
	for y := 0; y < len(g.bscreen); y++ {
		for x := 0; x < len(g.bscreen[y]); x++ {
			l := labels[y][x]
			if l == -1 {
				labels[y][x] = 0
			}
			labels[y][x] = labelMapping[l]
		}
	}
	labelsToKeep := make([]int, 0)
	for {
		y := len(g.bscreen) - 1
		for x := 0; x < len(g.bscreen[y]); x++ {
			l := labels[y][x]
			if l > 0 {
				labelsToKeep = append(labelsToKeep, l)
			}
		}
		break
	}
	//log(getUniqueIntSlice(labelsToKeep))
	marblesToRemove := make([]Object, 0)
	bcellsToRemove := make([][]int, 0)
	for _, bc := range g.bcells {
		x := bc[0]
		y := bc[1]
		l := labels[y][x]
		toRemove := true
		for _, lk := range labelsToKeep {
			if lk == l {
				toRemove = false
				break
			}
		}
		if !toRemove {
			continue
		}
		bcellsToRemove = append(bcellsToRemove, []int{x, y})
		marblesToRemove = append(marblesToRemove, g.getMarbleByBxBy(x, y))
	}
	g.marbles = g.getMarblesWithout(g.marbles, marblesToRemove)
	//log(labels)
	bcells := make([][]int, 0)
	for _, bc := range g.bcells {
		found := false
		for _, bcr := range bcellsToRemove {
			if bcr[0] == bc[0] && bcr[1] == bc[1] {
				found = true
				break
			}
		}
		if !found {
			bcells = append(bcells, bc)
		}
	}

	g.bcells = bcells
}

func (g *group) getMarblesToRemove(marbles []Object, cm Object) []Object {
	toRemove := make([]Object, 0)
	cmbx, cmby := cm.BXY()
	for _, m := range g.getMarblesWithout(marbles, cm) {
		if m.Color() != cm.Color() {
			continue
		}
		i, j := m.BXY()
		dx := mlib.AbsInt(cmbx - i)
		dy := mlib.AbsInt(cmby - j)
		if (dx <= 2 && dy < 1) || (dy == 2 && dx <= 2) {
			log(cmbx, cmby, i, j)
			toRemove = append(toRemove, m)
		}
	}

	for _, m := range toRemove {
		toRemove = append(toRemove, g.getMarblesToRemove(g.getMarblesWithout(marbles, toRemove), m)...)
	}

	return toRemove
}

func (g *group) Add(cm Object) {
	g.marbles = append(g.marbles, cm)
	bx, by := cm.BXY()

	g.bcells = append(g.bcells, []int{bx, by})
}

func (g *group) Contains(x, y float64) bool {
	for _, m := range g.marbles {
		if m.Contains(x, y) {
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
		mbx, mby := m.BXY()
		switch cm := cmEl.(type) {
		case Object:
			cmx, cmy := cm.BXY()
			if mbx == cmx && mby == cmy {
				continue
			}
			res = append(res, m)
		case []Object:
			found := false
			for _, _m := range cm {
				cmx, cmy := _m.BXY()
				if mbx == cmx && mby == cmy {
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

func (g *group) GetNeighbor(m Object, neighborType int) Object {

	return nil
}

func (g *group) GetBresenhamCells() [][]int {
	res := make([][]int, 0)
	for _, c := range g.bcells {
		res = append(res, []int{c[0], c[1]})
	}

	return res
}

func (g *group) getMarbleByBxBy(x, y int) Object {
	for _, m := range g.marbles {
		bx, by := m.BXY()
		if bx == x && by == y {
			return m
		}
	}

	return nil
}

func (g *group) inBcells(x0 int, y0 int) bool {
	for _, c := range g.bcells {
		if c[0] == x0 && c[1] == y0 {
			return true
		}
	}

	return false
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

func (g *group) GetFIntersection(cm Object) (bool, float64, float64) {
	points := make([][]float64, 0)
	d := float64(g.marbleSize * 2)
	for _, m := range g.marbles {
		if m.X() < cm.X()-d || m.X() > cm.X()+d {
			continue
		}
		if m.Y() > cm.Y()+d || m.Y() < cm.Y()-d {
			continue
		}

		i, _, _ := cm.WillIntersectsWith(m)
		if i {
			points = append(points, []float64{m.X(), m.Y()})
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
		sumX += cm.X()
		sumY += cm.Y()
		x = sumX / (float64(len(points)) + 1)
		y = sumY / (float64(len(points)) + 1)
	} else {
		x = (points[0][0] + cm.X()) / 2
		y = (points[0][1] + cm.Y()) / 2
	}
	r := d / 4
	d = r * 2
	h := g.windowHeight
	w := g.windowWidth
	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			if g.Contains(float64(j)*d+r, float64(h)-float64(i)*d-r) {
				continue
			}
			dx := float64(j)*d + r - x
			dy := float64(h) - float64(i)*d - r - y
			if math.Sqrt(dx*dx+dy*dy) < r {
				return true, float64(j)*d + r, float64(h) - float64(i)*d - r
			}
		}
	}

	//i * r = x
	//j * r = y
	return true, x, y
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

func (g *group) getBNeighborsLabel(x int, y int, labels [][]int, label int) int {
	if !g.inBcells(x, y) {
		return -1
	}

	neighborLabels := g.getBNeighborLabels(x, y, labels)
	if len(neighborLabels) < 1 {
		return label + 1
	}

	return minIntSlice(neighborLabels)
}

func (g *group) getBNeighborLabels(x int, y int, labels [][]int) []int {
	neighborLabels := make([]int, 0)
	if !g.inBcells(x, y) {
		return neighborLabels
	}

	if len(labels[y]) < 1 {
		return neighborLabels
	}
	// left
	if x-2 >= 0 && labels[y][x-2] != -1 {
		neighborLabels = append(neighborLabels, labels[y][x-2])
	}
	// right
	if x+2 < len(labels[y]) && labels[y][x+2] != -1 {
		neighborLabels = append(neighborLabels, labels[y][x+2])
	}
	// top
	if y-2 >= 0 && len(labels[y-2]) > 0 && labels[y-2][x] != -1 {
		neighborLabels = append(neighborLabels, labels[y-2][x])
	}
	// bottom
	if y+2 < len(labels) && len(labels[y+2]) > 0 && labels[y+2][x] != -1 {
		neighborLabels = append(neighborLabels, labels[y+2][x])
	}
	// top left
	if y-2 >= 0 && x-1 >= 0 && len(labels[y-2]) > 0 && labels[y-2][x-1] != -1 {
		neighborLabels = append(neighborLabels, labels[y-2][x-1])
	}
	// top right
	if y-2 >= 0 && x+1 < len(labels[y-2]) && len(labels[y-2]) > 0 && labels[y-2][x+1] != -1 {
		neighborLabels = append(neighborLabels, labels[y-2][x+1])
	}
	// bottom left
	if y+2 < len(labels) && x-1 >= 0 && len(labels[y+2]) > 0 && labels[y+2][x-1] != -1 {
		neighborLabels = append(neighborLabels, labels[y+2][x-1])
	}
	// bottom right
	if y+2 < len(labels) && x+1 < len(labels[y+2]) && len(labels[y+2]) > 0 && labels[y+2][x+1] != -1 {
		neighborLabels = append(neighborLabels, labels[y+2][x+1])
	}

	return neighborLabels
}

func minIntSlice(a []int) int {
	m := a[0]
	for _, e := range a {
		if m > e {
			m = e
		}
	}

	return m
}

func getUniqueIntSlice(a []int) []int {
	b := make([]int, 0)
	for _, ea := range a {
		found := false
		for _, eb := range b {
			if eb == ea {
				found = true
				break
			}
		}
		if !found {
			b = append(b, ea)
		}
	}

	return b
}
