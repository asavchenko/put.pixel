package main

import (
	"assa.com/put.pixel/lib/ogl"
	"assa.com/put.pixel/src/characters"
	"fmt"
	"runtime"
	"time"
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
	fontSize := 14
	textWidth := len(text) * (characters.GetCharacterWidth(fontSize) + characters.GetSpaceSizeBtwCharacters(fontSize))
	textHeight := characters.GetCharacterHeight(fontSize) + characters.GetLineSpaceSize(fontSize)
	y := (h + 2*textHeight) / 2
	x := (w - textWidth) / 2
	for _, r := range text {
		chrs = append(chrs, characters.GetNew(r, x, y, color).SetCharacterSize(fontSize))
		x += characters.GetCharacterWidth(fontSize) + characters.GetSpaceSizeBtwCharacters(fontSize)
	}
	numChrs := len(text)
	for {
		if ogl.IsExit() {
			break
		}

		ogl.Draw(func() {
			for i := 0; i < numChrs; i++ {
				chrs[i].Move(0, -1)
				if int(chrs[i].Y) < 0 {
					chrs[i].Y = ogl.GetWindowHeight() + chrs[i].GetHeight()
				}
			}
		})
		time.Sleep(1 * time.Millisecond)
	}
}
