package ogl

// #include <string.h>
import "C"

import (
	"fmt"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"os"
	"runtime"
	"strings"
	"time"
)

const (
	width  = 640
	height = 480
	Xm     = width - 1
	Ym     = height - 1
)

var pixelArr *[width * height * 4]byte
var pixelArrRGBA *[width * height]uint32

// var pixelArr []byte
var window *glfw.Window
var buffers [1]uint32

// var index int
var lastX, lastY int
var lastWidth, lastHeight int
var keyCallbacks map[string]map[glfw.Key][]func()
var mouseLeftCallbacks []func(x, y int)
var keyCombinationCallbacks map[string]map[string]interface{}
var pressedKeys []glfw.Key
var lookupTable []int
var lookupTableRGBA []int

func init() {
	keyCallbacks = make(map[string]map[glfw.Key][]func(), 0)
	mouseLeftCallbacks = make([]func(x, y int), 0)
	keyCombinationCallbacks = make(map[string]map[string]interface{}, 0)
	pressedKeys = make([]glfw.Key, 0)
	lookupTable = make([]int, height)
	for i := 0; i < height; i++ {
		lookupTable[i] = i * width * 4
	}
	lookupTableRGBA = make([]int, height)
	for i := 0; i < height; i++ {
		lookupTableRGBA[i] = i * width
	}
	//pixelArr = make([]byte, width*height*4)
}

func Init(fullScreen bool) {
	if err := glfw.Init(); err != nil {
		panic(fmt.Sprint("failed to initialize glfw:", err))
	}
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.DoubleBuffer, glfw.False)
	{
		var err error
		window, err = glfw.CreateWindow(width, height, "Title", nil, nil)
		if err != nil {
			panic(err)
		}
	}
	window.MakeContextCurrent()
	if fullScreen {
		GoFullScreen()
	}
	if err := gl.Init(); err != nil {
		panic(fmt.Sprint("failed to initialize gl bindings:", err))
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.ClearColor(0.0, 0.0, 0.0, 0.0)
	gl.GenBuffers(1, &buffers[0])
	glfw.SwapInterval(1)

	gl.BindBuffer(gl.PIXEL_UNPACK_BUFFER, buffers[0])
	gl.BufferData(gl.PIXEL_UNPACK_BUFFER, width*height*4, nil, gl.STREAM_DRAW)

	lastX, lastY = window.GetPos()
	lastWidth, lastHeight = window.GetSize()
	window.SetKeyCallback(onKeyPress)
	window.SetMouseButtonCallback(onMouseButtonEvent)
	window.SetCursorPosCallback(onMouseMoveEvent)
}

func CloseWindow() {
	window.SetShouldClose(true)
}

func Close() {
	gl.DeleteBuffers(1, &buffers[0])
	if window != nil {
		window.Destroy()
	}
	glfw.Terminate()
}

func PutPixel(x, y int, color ...byte) {
	r := byte(0)
	g := byte(0)
	b := byte(0)
	a := byte(0)
	if len(color) == 1 {
		r = color[0]
		g = color[0]
		b = color[0]
		a = color[0]
	} else if len(color) == 2 {
		r = color[0]
		g = color[1]
		b = color[1]
		a = color[1]
	} else if len(color) == 3 {
		r = color[0]
		g = color[1]
		b = color[2]
		a = color[2]
	} else if len(color) == 4 {
		r = color[0]
		g = color[1]
		b = color[2]
		a = color[3]
	}
	if x < 0 || x >= width || y < 0 || y >= height {
		return
	}

	copy(pixelArr[lookupTable[y]+x<<2:], []byte{r, g, b, a})
}

func UnsafePutPixelRGBA(x, y int, color uint32) {
	//                                                       r                   g                  b           a
	pixelArrRGBA[lookupTableRGBA[y]+x] = color
}

func PutPixelRGB(x, y int, color uint32) {
	if x < 0 || x >= width || y < 0 || y >= height {
		return
	}
	UnsafePutPixelRGBA(x, y, color)
}

func GetPixel(x, y int) uint32 {
	if x < 0 || x >= width || y < 0 || y >= height {
		return 0
	}

	return GetPixelUnsafe(x, y)
}

func GetPixelUnsafe(x, y int) uint32 {
	i := lookupTableRGBA[y] + x
	//             a                         b                          g                            r
	return pixelArrRGBA[i]
}

func GetPixelUnsafeRGBA(x, y int) uint32 {
	i := lookupTableRGBA[y] + x
	//             a                         b                          g                            r
	return pixelArrRGBA[i]
}

func GetWindowWidth() int {
	return width
}

func GetWindowHeight() int {
	return height
}

func IsExit() bool {
	return window.ShouldClose()
}

func Draw(run func()) {
	if !window.ShouldClose() {
		draw(window, run)
	}
}

func draw(window *glfw.Window, run func()) {
	pboPtr := gl.MapBuffer(gl.PIXEL_UNPACK_BUFFER, gl.WRITE_ONLY)
	if pboPtr == nil {
		return
	}
	if !gl.UnmapBuffer(gl.PIXEL_UNPACK_BUFFER) {
		return
	}
	//pixelArr = (*[width * height * 4]byte)(pboPtr)
	pixelArrRGBA = (*[width * height]uint32)(pboPtr)
	//pixelArr = unsafe.Slice((*byte)(pboPtr), width*height*3)
	//pixelArr = (*[width * height * 3]byte)(pboPtr)[:width*height*3]

	//copy((*[width * height * 3]byte)(pboPtr)[:width*height*3], pixelArr)
	ClearScreenRGBA()
	run()

	gl.DrawPixels(width, height, gl.RGBA, gl.UNSIGNED_BYTE, nil)
	//SwapBuffers()

	glfw.PollEvents()
	processInput(window)
}

func PutBitmap(x, y int, w, h int, bitmap []uint32) {
	di := 0
	for i := 0; i < h; i++ {
		if y >= len(lookupTable) || y < 0 {
			di += w
			y += 1
			continue
		}
		idx := lookupTableRGBA[y] + x
		//idx := lookupTable[y] + x

		if idx >= len(pixelArrRGBA) || idx < 0 {
			di += w
			y += 1
			continue
		}
		copy(pixelArrRGBA[idx:], bitmap[di:di+w])
		di += w
		y += 1
	}
}

func PutByteBitmap(x, y int, w, h int, bitmap []byte) {
	di := 0
	for i := 0; i < h; i++ {
		if y >= len(lookupTable) || y < 0 {
			di += w
			y += 1
			continue
		}
		idx := lookupTable[y] + x<<2

		if idx >= len(pixelArr) || idx < 0 {
			di += w
			y += 1
			continue
		}
		copy(pixelArr[idx:], bitmap[di:di+w])
		di += w
		y += 1
	}
}

func PutByteBitmapBordered(x, y int, il, ih, jl, jh int, w int, h int, bitmap []byte) {
	x4 := x << 2
	jl4 := jl << 2
	jh4 := jh << 2
	di := il * w
	width4 := width << 2
	idx := lookupTable[y] + x4
	for i := il; i <= ih; i++ {
		copy(pixelArr[idx:], bitmap[di+jl4:di+jh4])
		di += w
		idx += width4
	}
}

func toBytes(arr []uint32) []byte {
	res := make([]byte, 4*len(arr))
	i := 0
	for _, c := range arr {
		res[i] = byte(c >> 24)
		res[i+1] = byte(c >> 16)
		res[i+2] = byte(c >> 8)
		res[i+3] = byte(c)
		i += 4
	}

	return res
}

func SwapBuffers() {
	gl.Flush()
	//gl.Finish()
	//window.SwapBuffers()
	//gl.Flush()
}

func ClearScreen() {
	for i := range pixelArr {
		pixelArr[i] = 0
	}
}
func ClearScreenRGBA() {
	for i := range pixelArrRGBA {
		pixelArrRGBA[i] = 0
	}
}

func FillScreen(color byte) {
	for i := range pixelArr {
		pixelArr[i] = color
	}
}

func log(msgs ...interface{}) {
	_, f, line, _ := runtime.Caller(1)
	baseDir := getBaseDir()

	relativePath := strings.Replace(f, baseDir+"/", "", -1)
	formatedMsg := fmt.Sprintln(time.Now().UTC().Format("15:04:05.999 02-01-2006"), fmt.Sprintf("%s:%d", relativePath, line), msgs)
	fmt.Print(formatedMsg)
}

func logError(msgs ...interface{}) {
	_, f, line, _ := runtime.Caller(1)
	baseDir := getBaseDir()

	relativePath := strings.Replace(f, baseDir+"/", "", -1)
	formatedMsg := fmt.Sprintln("*********ERROR", time.Now().UTC().Format("15:04:05.999 02-01-2006"), fmt.Sprintf("%s:%d", relativePath, line), msgs)
	fmt.Print(formatedMsg)
}

func getBaseDir() string {
	pwd, err := os.Getwd()
	if err != nil {
		logError(err)

		return ""
	}

	return pwd
}
