package main

import (
	"fmt"
	"runtime"
	"time"

	"assa.com/put.pixel/lib/ogl"
	"assa.com/put.pixel/src/characters"
)

var chrs []*characters.Chr

func init() {
	fmt.Println("init")
	runtime.LockOSThread()
	chrs = make([]*characters.Chr, 0)
}

func main() {
	ogl.Init(false)
	defer ogl.Close()
	ogl.OnKeypress(ogl.KEY_ESC, func() {
		ogl.CloseWindow()
	})
	text := "It works!"
	w := ogl.GetWindowWidth()
	h := ogl.GetWindowHeight()
	color := byte(200)
	textWidth := len(text) * (characters.GetCharacterWidth() + characters.GetSpaceSizeBtwCharacters())
	textHeight := characters.GetCharacterHeight() + characters.GetLineSpaceSize()
	y := (h + 2*textHeight) / 2
	x := (w - textWidth) / 2
	for _, r := range text {
		chrs = append(chrs, characters.GetNew(r, x, y, color))
		x += characters.GetCharacterWidth() + characters.GetSpaceSizeBtwCharacters()
	}
	numChrs := len(text)
	for {
		if ogl.IsExit() {
			break
		}

		ogl.Draw(func() {
			for i := 0; i < numChrs; i++ {
				chrs[i].Show()
			}
		})
		time.Sleep(17 * time.Millisecond)
	}
}
