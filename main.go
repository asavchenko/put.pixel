package main

import (
	"runtime"
	"time"

	"assa.com/put.pixel/lib/mlib"
	"assa.com/put.pixel/lib/ogl"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	ogl.Init(false)
	defer ogl.Close()
	ogl.DisableDoubleBuffering()
	width := ogl.GetWindowWidth()
	height := ogl.GetWindowHeight()
	Wm := width * 2
	Hm := height * 2
	ogl.OnKeypress(ogl.KEY_ESC, func() {
		ogl.CloseWindow()
	})
	ogl.OnCombination([]interface{}{ogl.KEY_LEFT_ALT, ogl.KEY_ENTER}, func() {
		if ogl.IsInFullScreen() {
			ogl.RestoreFromFullScreen()
			return
		}
		ogl.GoFullScreen()
	})
	for {
		if ogl.IsExit() {
			break
		}

		ogl.Draw(func() {
			for i := 0; i < 200; i++ {
				ogl.Line(mlib.Srand(Wm), mlib.Srand(Hm), mlib.Srand(Wm), mlib.Srand(Hm), byte(mlib.Rand(256)))
			}
			time.Sleep(16 * time.Millisecond)
		})
	}
}
