package main

import (
	"fmt"
	"runtime"
	"time"

	"assa.com/put.pixel/lib/ogl"
	"assa.com/put.pixel/src/characters"
)

const numChrs = 2000

var chrs []*characters.Chr

func init() {
	fmt.Println("init")
	runtime.LockOSThread()
	chrs = make([]*characters.Chr, numChrs)
}

func main() {
	ogl.Init()
	defer ogl.Close()
	for i := 0; i < numChrs; i++ {
		chrs[i] = characters.GetNew()
	}
	for {
		if ogl.IsExit() {
			break
		}

		ogl.Draw(func() {
			for i := 0; i < numChrs; i++ {
				chrs[i].Move()
			}
		})
		time.Sleep(17 * time.Millisecond)
	}
}
