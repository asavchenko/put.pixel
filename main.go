package main

import (
	"fmt"
	"math"
	"runtime"
	"time"

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

	marbleSpeed := float64(9)
	w := ogl.GetWindowWidth()
	h := ogl.GetWindowHeight()
	a := arrow.GetNew(w/2, h/9, 30, 0, 0xFEFEFEFF)
	size := 32
	gridSize := size/2 + 2
	marbles := marble.GetNewGroup(w, h, size, gridSize)
	nextMarble := marble.GetNew(w/2-size*2, h/9, size, marble.GetRandomColor(), w, h, gridSize)
	curMarble := marble.GetNew(w/2, h/9, size, marble.GetRandomColor(), w, h, gridSize)
	curMarble.SetSpeed(marbleSpeed)
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
		curMarble.SetIsMoving(true)
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
			curMarble.Show()
			a.Show()
			marbles.Show()
			if !isMoving {
				return
			}
			if isIntersection, cx, cy := marbles.GetBIntersection(curMarble); isIntersection {
				isMoving = false
				curMarble.SetIsMoving(false)
				curMarble.MoveTo(cx, cy)
				marbles.Add(curMarble)
				marbles.RemoveMatched(curMarble)
				curMarble = marble.GetNew(w/2, h/9, size, nextMarble.Color(), w, h, gridSize)
				curMarble.SetSpeed(marbleSpeed)
				nextMarble = marble.GetNew(w/2-size*2, h/9, size, marble.GetRandomColor(), w, h, gridSize)
				return

			}
			// checking Edges
			if curMarble.X()-curMarble.R()-curMarble.GetSpeed() <= 0 {
				curMarble.Rotate(-math.Pi / 2)
				//fmt.Println("left edge was hit", curMarble.X(), 180/math.Pi*curMarble.A())
			} else if curMarble.X()+curMarble.R()+curMarble.GetSpeed() >= float64(w) {
				//log("rotate!")
				curMarble.Rotate(math.Pi / 2)
				//fmt.Println("right edge was hit", curMarble.X(), 180/math.Pi*curMarble.A())
			}
		})
	}
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
