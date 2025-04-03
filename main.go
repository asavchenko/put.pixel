package main

import (
	"assa.com/put.pixel/src/characters"
	"assa.com/put.pixel/src/sprite"
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
	//h := ogl.GetWindowHeight()
	//log(w, h)
	x := float64(w) / 4
	y := float64(w) / 4
	mainCh := characters.GetSprite(x, y)
	ogl.Init(false)
	defer ogl.Close()
	ogl.OnKeypress(ogl.KEY_ESC, func() {
		ogl.CloseWindow()
	})
	currentKey := ogl.KEY_UNKNOWN
	mainCh.Stop()
	ogl.OnKeypress(ogl.KEY_RIGHT, func() {
		if currentKey == ogl.KEY_RIGHT {
			return
		}
		mainCh.Stop()
		currentKey = ogl.KEY_RIGHT
		mainCh.SetDirection(sprite.DIRECTION_RIGHT)
		mainCh.Start("walk")
	})
	ogl.OnKeypress(ogl.KEY_LEFT, func() {
		if currentKey == ogl.KEY_LEFT {
			return
		}
		mainCh.Stop()
		currentKey = ogl.KEY_LEFT
		mainCh.SetDirection(sprite.DIRECTION_LEFT)
		mainCh.Start("walk")
	})

	ogl.OnKeydown(ogl.KEY_RIGHT, func() {
		if currentKey == ogl.KEY_RIGHT {
			return
		}
		mainCh.Stop()
		currentKey = ogl.KEY_RIGHT
		mainCh.SetDirection(sprite.DIRECTION_RIGHT)
		mainCh.Start("walk")
	})
	ogl.OnKeydown(ogl.KEY_LEFT, func() {
		if currentKey == ogl.KEY_LEFT {
			return
		}
		mainCh.Stop()
		currentKey = ogl.KEY_LEFT
		mainCh.SetDirection(sprite.DIRECTION_LEFT)
		mainCh.Start("walk")
	})
	ogl.OnKeyup(ogl.KEY_RIGHT, func() {
		currentKey = ogl.KEY_UNKNOWN
		mainCh.Stop()
	})
	ogl.OnKeyup(ogl.KEY_LEFT, func() {
		currentKey = ogl.KEY_UNKNOWN
		mainCh.Stop()
	})

	start := time.Now()
	for {
		if ogl.IsExit() {
			break
		}
		ogl.Draw(func() {
			mainCh.Show()
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
