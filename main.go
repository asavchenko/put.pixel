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

var snowFlakes []*snow.Flake

func init() {
	runtime.LockOSThread()
	wind.SetDirection(5)
	snowFlakes = make([]*snow.Flake, numFlakes)
}

func main() {
	ogl.Init()
	defer ogl.Close()
	for i := 0; i < numFlakes; i++ {
		snowFlakes[i] = snow.GetNew()
	}
	go func() {
		for {
			if mlib.Rand(150) > 25 {
				wind.SetDirection(mlib.Srand(5))
			}
			time.Sleep(1 * time.Second)
		}
	}()
	for {
		if ogl.IsExit() {
			break
		}
		for i := 0; i < numFlakes; i++ {
			snowFlakes[i].Move()
		}
		ogl.Draw()
		time.Sleep(17 * time.Millisecond)
	}
}
