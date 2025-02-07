package main

import (
	"fmt"
	"runtime"
	"time"

	"assa.com/put.pixel/lib/ogl"
	"assa.com/put.pixel/src/arrow"
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

	w := ogl.GetWindowWidth()
	h := ogl.GetWindowHeight()
	a := arrow.GetNew(w/2, h/2, 60, 0, 0xFEFEFEFF)
	ogl.Init(false)
	defer ogl.Close()
	ogl.OnKeypress(ogl.KEY_ESC, func() {
		ogl.CloseWindow()
	})

	ogl.OnMouseLeftClick(func(x, y int) {
		// shoot()
	})
	for {
		if ogl.IsExit() {
			break
		}
		ogl.Draw(func() {
			x, y := ogl.GetMousePosition()
			a.RotateTo(x, y)
			a.Show()
		})
	}
}

func timer(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", name, time.Since(start))
	}
}
