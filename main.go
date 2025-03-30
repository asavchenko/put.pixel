package main

import (
	"fmt"
	"runtime"
	"time"

	"assa.com/put.pixel/lib/ogl"
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
	log(w, h)
	ogl.Init(false)
	defer ogl.Close()
	ogl.OnKeypress(ogl.KEY_ESC, func() {
		ogl.CloseWindow()
	})

	start := time.Now()
	for {
		if ogl.IsExit() {
			break
		}
		ogl.Draw(func() {

		})
		// 59.95 Hz
		// 1 * time.Millisecond * x = 1/59.95
		// x = 1/59.95 / time.Millisecond
		d := 16*time.Millisecond - time.Since(start)
		if d > 0 {
			time.Sleep(d)
		}
		ogl.SwapBuffers()
		start = time.Now()
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
