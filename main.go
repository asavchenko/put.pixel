package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"syscall"
	"time"

	"assa.com/put.pixel/lib/ogl"
	"assa.com/put.pixel/src/characters"
)

var chrs []*characters.Chr

func init() {
	fmt.Println("init")
	syscall.Setpriority(syscall.PRIO_PROCESS, os.Getpid(), -20)
	runtime.LockOSThread()
	chrs = make([]*characters.Chr, 0)
}

func main() {
	ogl.Init(false)
	defer ogl.Close()
	ogl.OnKeypress(ogl.KEY_ESC, func() {
		ogl.CloseWindow()
	})
	dx := 0
	ogl.OnKeypress(ogl.KEY_SPACE, func() {
		if dx == 0 {
			dx -= 1
		} else {
			dx = 0
		}
	})
	text := "Very long string goes here without question"
	//text := "Я родился!"
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

	numChrs := len([]rune(text))
	for {
		if ogl.IsExit() {
			break
		}
		ogl.Draw(func() {
			for i := 0; i < numChrs; i++ {
				chrs[i].Move(0, dx)
				if chrs[i].Y < 0 {
					chrs[i].Move(0, ogl.GetWindowHeight()+2*chrs[i].GetHeight())
				}
			}
		})
	}
}

func log(msgs ...interface{}) {
	_, f, line, _ := runtime.Caller(1)
	baseDir := ""
	pwd, err := os.Getwd()
	if err == nil {
		baseDir = pwd
	}

	relativePath := strings.Replace(f, baseDir+"/", "", -1)
	formatedMsg := fmt.Sprintln(time.Now().UTC().Format("15:04:05.999 02-01-2006"), fmt.Sprintf("%s:%d", relativePath, line), msgs)
	fmt.Print(formatedMsg)
}

func logError(msgs ...interface{}) {
	_, f, line, _ := runtime.Caller(1)
	baseDir := ""
	pwd, err := os.Getwd()
	if err == nil {
		baseDir = pwd
	}

	relativePath := strings.Replace(f, baseDir+"/", "", -1)
	formatedMsg := fmt.Sprintln("*********ERROR", time.Now().UTC().Format("15:04:05.999 02-01-2006"), fmt.Sprintf("%s:%d", relativePath, line), msgs)
	fmt.Print(formatedMsg)
}
