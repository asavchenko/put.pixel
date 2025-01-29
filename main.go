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
	//f, err := os.Create("myprogram.prof") // then go tool pprof -http=:8080 myprogram.prof
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//pprof.StartCPUProfile(f)
	//defer pprof.StopCPUProfile()

	ogl.Init(false)
	defer ogl.Close()
	ogl.OnKeypress(ogl.KEY_ESC, func() {
		ogl.CloseWindow()
	})
	dy := -1

	text := `What is Lorem Ipsum?
	Lorem Ipsum is simply dummy text of the printing and typesetting industry.
	Lorem Ipsum has been the industry's standard dummy text ever since the 1500s,
	when an unknown printer took a galley of type and scrambled it to make a type specimen book. 
	It has survived not only five centuries, but also the leap into electronic typesetting, 
	remaining essentially unchanged. It was popularised in the 1960s with the release of 
	Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker 
	including versions of Lorem Ipsum.

	Why do we use it?
	It is a long established fact that a reader will be distracted by the readable content of a page when looking at its layout.
	The point of using Lorem Ipsum is that it has a more-or-less normal distribution of letters, 
	as opposed to using 'Content here, content here', making it look like readable English. 
	Many desktop publishing packages and web page editors now use Lorem Ipsum as their default model text, 
	and a search for 'lorem ipsum' will uncover many web sites still in their infancy. 
	Various versions have evolved over the years, sometimes by accident, 
	sometimes on purpose (injected humour and the like).

	Where does it come from?
	Contrary to popular belief, Lorem Ipsum is not simply random text. 
	It has roots in a piece of classical Latin literature from 45 BC,
	making it over 2000 years old. Richard McClintock, a Latin professor at Hampden-Sydney College in Virginia, 
	looked up one of the more obscure Latin words, consectetur, from a Lorem Ipsum passage, 
	and going through the cites of the word in classical literature, discovered the undoubtable source. 
	Lorem Ipsum comes from sections 1.10.32 and 1.10.33 of \"de Finibus Bonorum et Malorum\" (The Extremes of Good and Evil) by Cicero, written in 45 BC. This book is a treatise on the theory of ethics, very popular during the Renaissance. The first line of Lorem Ipsum, \"Lorem ipsum dolor sit amet..\", comes from a line in section 1.10.32.

	The standard chunk of Lorem Ipsum used since the 1500s is reproduced below for those interested. Sections 1.10.32 and 1.10.33 from \"de Finibus Bonorum et Malorum\" by Cicero are also reproduced in their exact original form, accompanied by English versions from the 1914 translation by H. Rackham.

	Where can I get some?
	There are many variations of passages of Lorem Ipsum available, but the majority have suffered alteration in some form, by injected humour, or randomised words which don't look even slightly believable. If you are going to use a passage of Lorem Ipsum, you need to be sure there isn't anything embarrassing hidden in the middle of text. All the Lorem Ipsum generators on the Internet tend to repeat predefined chunks as necessary, making this the first true generator on the Internet. It uses a dictionary of over 200 Latin words, combined with a handful of model sentence structures, to generate Lorem Ipsum which looks reasonable. The generated Lorem Ipsum is therefore always free from repetition, injected humour, or non-characteristic words etc.`
	//text = "A"
	//text = "It works!"
	//text := "Я родился!"
	w := ogl.GetWindowWidth()
	h := ogl.GetWindowHeight()
	color := byte(0)
	fontSize := 12
	textHeight := characters.GetCharacterHeight(fontSize) + characters.GetLineSpaceSize(fontSize)
	//numChInRow := w/chWidth + 1
	//numColumns := h/textHeight + 1
	y := h
	x := 0
	i := 0
	isNewLine := true
	for {
		if i > len([]rune(text))-1 {
			break
		}
		r := []rune(text)[i]
		if int(r) == 10 {
			y -= textHeight
			if y < 0 {
				break
			}
			// new line
			isNewLine = true
			if int([]rune(text)[i+1]) == 9 {
				i += 2
			} else {
				i += 1
			}
			continue
		}
		i += 1
		if isNewLine {
			x = 0
		}
		if x > w {
			continue
		}
		ch := characters.GetNew(r, x, y, color).SetCharacterSize(fontSize)
		chrs = append(chrs, ch)
		isNewLine = false
		x += ch.GetWidth()

	}
	numChrs := len(chrs)
	fmt.Println(numChrs)
	ogl.OnKeypress(ogl.KEY_SPACE, func() {
		if dy == 0 {
			dy -= 1
		} else {
			dy = 0
		}
	})

	isInit := true
	for {
		if ogl.IsExit() {
			break
		}
		//start := time.Now()
		ogl.Draw(func() {
			if isInit {
				ogl.FillScreen(200)
				fmt.Println("filled")
				isInit = false
				return
			}
			//defer timer("inside draw")()
			for i := 0; i < numChrs; i++ {
				y := chrs[i].Y
				if y+chrs[i].GetMaxCharacterHeight() < 0 {
					chrs[i].MoveUnsafe(0, ogl.GetWindowHeight()+2*chrs[i].GetHeight())
				} else {
					chrs[i].MoveUnsafe(0, dy)
				}
				chrs[i].ShowUnsafe()
				chrs[i].HideUnsafe()
			}

			ogl.SwapBuffers()
			//defer timer("inside draw")()
			//ogl.ClearScreen()
			//for i := 0; i < numChrs; i++ {
			//	_, y := chrs[i].GetPosition()
			//	if y+chrs[i].GetMaxCharacterHeight() < 0 {
			//		chrs[i].Move(0, ogl.GetWindowHeight()+2*chrs[i].GetHeight())
			//	} else {
			//		chrs[i].Move(0, dy)
			//	}
			//}
		})
		//d := 29*time.Millisecond - time.Since(start)
		//if d > 0 {
		//	time.Sleep(d)
		//}
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
