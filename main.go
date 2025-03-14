package main

import (
	"fmt"
	"math"
	"runtime"
	"time"

	"assa.com/put.pixel/lib/mlib"
	"assa.com/put.pixel/lib/ogl"
	"assa.com/put.pixel/src/arrow"
	"assa.com/put.pixel/src/marble"
)

func init() {
	fmt.Println("init")
	//syscall.Setpriority(syscall.PRIO_PROCESS, os.Getpid(), -20)
	runtime.LockOSThread()
}

func main() {
	//f, err := os.Create("myprogram.prof") // then go tool pprof -http=:8080 myprogram.prof
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//pprof.StartCPUProfile(f)
	//defer pprof.StopCPUProfile()

	marbleSpeed := 6
	w := ogl.GetWindowWidth()
	h := ogl.GetWindowHeight()
	a := arrow.GetNew(w/2, h/9, 30, 0, 0xFEFEFEFF)
	size := 32
	marbles := marble.GetNewGroup(w, h, size)
	nextMarble := marble.GetNew(w/2-size*2, h/9, size, marble.GetRandomColor())
	curMarble := marble.GetNew(w/2, h/9, size, marble.GetRandomColor())
	curMarble.SetSpeed(float64(marbleSpeed))
	ogl.Init(false)
	defer ogl.Close()
	ogl.OnKeypress(ogl.KEY_ESC, func() {
		ogl.CloseWindow()
	})

	isMoving := false
	ogl.OnMouseLeftClick(func(x, y int) {
		// shoot()
		if isMoving {
			return
		}
		//fmt.Println(180 / math.Pi * a.A())
		curMarble.SetAngle(a.A())
		isMoving = true
	})
	for {
		if ogl.IsExit() {
			break
		}
		ogl.Draw(func() {
			x, y := ogl.GetMousePosition()
			a.RotateTo(x, y)
			nextMarble.Show()
			if isMoving {
				curMarble.Move()
				if isIntersection, ix, iy := marbles.GetIntersection(curMarble); isIntersection {
					curMarble.MoveTo(ix, iy)
					marbles.Add(curMarble)
					marbles.RemoveMatched(curMarble)
					curMarble = marble.GetNew(w/2, h/9, size, nextMarble.Color())
					curMarble.SetSpeed(float64(marbleSpeed))
					nextMarble = marble.GetNew(w/2-size*2, h/9, size, marble.GetRandomColor())
					isMoving = false
				} else {
					// checking Edges
					if mlib.AbsInt(curMarble.CX()) < curMarble.CR() {
						curMarble.Rotate(-math.Pi / 2)
						curMarble.Move()
						//fmt.Println("left edge was hit", minLeft, curMarble.X(), 180/math.Pi*curMarble.A())
					}
					if mlib.AbsInt(curMarble.CX()-w) < curMarble.CR() {
						curMarble.Rotate(math.Pi / 2)
						curMarble.Move()
						//fmt.Println("right edge was hit", 180/math.Pi*curMarble.A())
					}
				}
			}
			curMarble.Show()
			a.Show()
			marbles.Show()
		})
	}
}

func removeFromSliceByIdx(a []marble.Object, i []int) []marble.Object {
	res := make([]marble.Object, 0)
	for j, e := range a {
		found := false
		for _, k := range i {
			if j == k {
				found = true
				break
			}
		}
		if found {
			continue
		}
		res = append(res, e)
	}

	return res
}

func removeFromIntSliceByVal(a []int, e int) []int {
	res := make([]int, 0)
	for _, el := range a {
		if el == e {
			continue
		}
		res = append(res, el)
	}

	return res
}

func getNeighbors(mi int, marbles []marble.Object, arr *[]int) []int {
	log(mi, *arr)
	cm := marbles[mi]
	for i, m := range marbles {
		if m.Color() != cm.Color() || mi == i {
			continue
		}
		found := false
		for j := 0; j < len(*arr); j++ {
			if i == (*arr)[j] {
				found = true
				break
			}
		}
		if found || mi == i {
			continue
		}
		if cm.GetDistanceTo(m) > cm.D()+3 {
			continue
		}
		found = false
		for _, v := range *arr {
			if v == i {
				found = true
				break
			}
		}
		if !found {
			*arr = append(*arr, i)
		}
		getNeighbors(i, marbles, arr)
	}
	return *arr
}

func isOccupied(x float64, y float64, marbles []marble.Object) bool {
	for _, m := range marbles {
		if m.IsInside(x, y) {
			return true
		}
	}

	return false
}

func timer(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", name, time.Since(start))
	}
}

func log(msgs ...interface{}) {
	_, f, line, _ := runtime.Caller(1)

	relativePath := f
	formatedMsg := fmt.Sprintln(time.Now().UTC().Format("15:04:05.999 02-01-2006"), fmt.Sprintf("%s:%d", relativePath, line), msgs)
	fmt.Print(formatedMsg)
}

func logError(msgs ...interface{}) {
	_, f, line, _ := runtime.Caller(1)

	relativePath := f
	formatedMsg := fmt.Sprintln("*********ERROR", time.Now().UTC().Format("15:04:05.999 02-01-2006"), fmt.Sprintf("%s:%d", relativePath, line), msgs)
	fmt.Print(formatedMsg)
}
