package main

import (
	"fmt"
	"runtime"
	"time"

	"assa.com/put.pixel/lib/ogl"
	"assa.com/put.pixel/src/characters"
)

var alphabet = []rune{' ', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-', '+', '.', '!', '<', '>', '?', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', ':'}

func init() {
	fmt.Println("init")
	runtime.LockOSThread()
}

func main() {
	ogl.Init(false)
	defer ogl.Close()
	ogl.OnKeypress(ogl.KEY_ESC, func() {
		ogl.Close()
	})
	ogl.OnCombination([]interface{}{ogl.KEY_LEFT_ALT, ogl.KEY_ENTER}, func() {
		if ogl.IsInFullScreen() {
			ogl.RestoreFromFullScreen()
			return
		}
		ogl.GoFullScreen()
	})
	balance := 0
	delta := 1
	go func() {
		for {
			if balance > 99 {
				delta = -1
			}

			if balance < 0 {
				delta = 1
			}
			time.Sleep(1 * time.Second)
			balance += delta
		}
	}()

	for {
		if ogl.IsExit() {
			break
		}

		ogl.Draw(func() {
			ogl.ClearScreen()
			if isOutOfService() {
				printStr("out of service", 0, 100, 0x77)
			} else {
				printStr("igt emulator", 0, 100, 0x77)
				printStr("balance: "+fmt.Sprint(getBalance()), 0, 200, 0x77)
			}
		})
		time.Sleep(17 * time.Millisecond)
	}
}

func printStr(str string, x, y int, color byte) {
	chrs := make([]*characters.Chr, len([]rune(str)))
	for i, ch := range []rune(str) {
		for _, a := range alphabet {
			if a == ch {
				chrs[i] = characters.GetNew(ch, color)
			}
		}
	}
	for _, ch := range chrs {
		ch.GoTo(x, y)
		x += 21
	}
}
