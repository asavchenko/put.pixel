package main

import (
	"runtime"
	"time"

	"assa.com/put.pixel/lib/mlib"
	"assa.com/put.pixel/lib/ogl"
	"assa.com/put.pixel/src/snow"
	"assa.com/put.pixel/src/wind"
)

const numFlakes = 2000

func init() {
	runtime.LockOSThread()
	wind.SetDirection(5)
}

func main() {
	ogl.Init()
	defer ogl.Close()
	for i := 0; i < numFlakes; i++ {
		go func(i int) {
			fl := snow.GetNew()
			time.Sleep(10 * time.Millisecond * time.Duration(mlib.Rand(i)))
			for {
				fl.Move()
				time.Sleep(16 * time.Millisecond)
			}
		}(i)
	}
	go func() {
		if mlib.Rand(150) > 25 {
			wind.SetDirection(mlib.Srand(5))
		}
		time.Sleep(1 * time.Second)
	}()
	for {
		if ogl.IsExit() {
			break
		}

		ogl.Draw()
		time.Sleep(17 * time.Millisecond)
	}
}
