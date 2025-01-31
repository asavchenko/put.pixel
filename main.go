package main

import (
	"assa.com/put.pixel/lib/mlib"
	"fmt"
	"math"
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

	w := ogl.GetWindowWidth()
	h := ogl.GetWindowHeight()
	fontSize := 18
	codes := characters.GetAvailableCharCodes()
	x := -9
	for k := 1; k < 3; k++ {
		for j := 0; j < 9; j++ {
			for i := 0; i < 54; i++ {
				r := codes[mlib.GetRandomBtw(0, len(codes)-1)]
				size := mlib.GetRandomBtw(9, fontSize)
				ch := characters.GetNew(rune(r), x, h+mlib.GetRandomBtw(0, h*3), getColor(size)).SetCharacterSize(size)
				ch.SetRotationSpeed(2 * math.Pi / float64(mlib.GetRandomBtw(30, 360)))
				ch.SetFallingSpeed(getFallingSpeed(size))
				chrs = append(chrs, ch)
			}
			x += mlib.GetRandomBtw(k*9, (k+j+1)*9) + mlib.GetRandomBtw(k, 9-j+k) + 9*k + 9*j
			if x > w {
				for {
					if x < w {
						break
					}
					x -= mlib.GetRandomBtw(0, w)
				}
			}
		}
	}
	numChrs := len(chrs)
	fmt.Println("num chars:", numChrs)
	for {
		if ogl.IsExit() {
			break
		}
		ogl.Draw(func() {
			//defer timer("inside draw")()
			for i := 0; i < numChrs; i++ {
				ch := chrs[i]
				y := ch.Y

				ch.Run()
				if y+chrs[i].GetMaxCharacterHeight() < 0 {
					ch.X += mlib.GetRandomBtw(0, 1)
					if ch.X > w {
						for {
							if ch.X < w {
								break
							}
							ch.X -= mlib.GetRandomBtw(9, 18)
						}
					}
					ch.Y = h + ch.GetHeight() + mlib.Rand(h)
				}
			}
		})
	}
}

func getFallingSpeed(size int) int {
	switch size - 9 {
	case 9:
		return 5
	case 8:
		return 4
	case 7:
		return 4
	case 6:
		return 3
	case 5:
		return 3
	case 4:
		return 2
	case 3:
		return 2
	case 2:
		return 2
	case 1:
		return 1
	case 0:
		return 1
	}

	fmt.Println("WTF>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	return 1
}

func getColor(n int) uint32 {
	switch n - 9 {
	case 9:
		return 0x42826cff
	case 8:
		return 0x38705dff
	case 7:
		return 0x2f5f4eff
	case 6:
		return 0x254e40ff
	case 5:
		return 0x1c3e32ff
	case 4:
		return 0x142f25ff
	case 3:
		return 0x0b2019ff
	case 2:
		return 0x05120dff
	case 1:
		return 0x010604ff
	case 0:
		return 0x000100ff
	}

	fmt.Println("WTF>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	return 0x000100ff
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
