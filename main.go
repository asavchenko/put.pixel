package main

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"time"

	"assa.com/put.pixel/lib/log"
	"assa.com/put.pixel/lib/ogl"
	"assa.com/put.pixel/src/characters"
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
	text := getText()
	a := characters.GetNewArea()
	w := ogl.GetWindowWidth()
	h := ogl.GetWindowHeight()
	fontSize := 12
	a.SetFontSize(fontSize)
	a.SetFontColor(0xFEFEFEFF)
	a.SetViewPortHeight(h)
	a.SetViewPortWidth(w)
	a.SetText(text)
	ogl.Init(false)
	defer ogl.Close()
	ogl.OnKeypress(ogl.KEY_ESC, func() {
		ogl.CloseWindow()
	})

	ogl.OnKeypress(ogl.KEY_RIGHT, func() {
		a.MoveView(10, 0)
	})

	ogl.OnKeypress(ogl.KEY_LEFT, func() {
		a.MoveView(-10, 0)
	})

	ogl.OnKeypress(ogl.KEY_UP, func() {
		a.MoveView(0, 10)
	})

	ogl.OnKeypress(ogl.KEY_DOWN, func() {
		a.MoveView(0, -10)
	})

	ogl.OnKeydown(ogl.KEY_RIGHT, func() {
		a.MoveView(10, 0)
	})

	ogl.OnKeydown(ogl.KEY_LEFT, func() {
		a.MoveView(-10, 0)
	})

	ogl.OnKeydown(ogl.KEY_UP, func() {
		a.MoveView(0, 10)
	})

	ogl.OnKeydown(ogl.KEY_DOWN, func() {
		a.MoveView(0, -10)
	})

	for {
		if ogl.IsExit() {
			break
		}
		ogl.Draw(func() {
			vpx, vpy := a.GetViewPortPosition()
			vpw, vph := a.GetViewPortWidthHeight()
			//defer timer("inside draw")()
			for _, ch := range a.GetCharacters() {
				ch.ShowInViewPort(vpx, vpy, vpw, vph)
			}
		})
	}
}

func getText() string {
	file, err := os.Open("goroutine.txt")
	if err != nil {
		return ""
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Println(err)
		}
	}()

	b, err := io.ReadAll(file)
	if err != nil {
		return ""
	}

	return string(b)
}

func timer(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", name, time.Since(start))
	}
}
