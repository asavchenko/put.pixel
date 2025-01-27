package ogl

// #include <string.h>
import "C"

import (
	"fmt"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"

	"log"
)

const (
	width  = 640
	height = 480
	Xm     = width - 1
	Ym     = height - 1
)

var pixelArr []byte
var window *glfw.Window
var buffers [1]uint32
var index int
var lastX, lastY int
var lastWidth, lastHeight int
var keyCallbacks map[string]map[glfw.Key][]func()
var keyCombinationCallbacks map[string]map[string]interface{}
var pressedKeys []glfw.Key
var lookupTable [][]int

func init() {
	keyCallbacks = make(map[string]map[glfw.Key][]func(), 0)
	keyCombinationCallbacks = make(map[string]map[string]interface{}, 0)
	pressedKeys = make([]glfw.Key, 0)
	pixelArr = make([]byte, width*height*3)
	lookupTable = make([][]int, height)
	for i := 0; i < height; i++ {
		line := make([]int, width)
		for j := 0; j < width; j++ {
			line[j] = (j + i*width) * 3
		}
		lookupTable[i] = line
	}
}

func Init(fullScreen bool) {
	if err := glfw.Init(); err != nil {
		log.Fatal("failed to initialize glfw:", err)
	}
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
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
		log.Fatal("failed to initialize gl bindings:", err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	gl.GenBuffers(1, &buffers[0])
	glfw.SwapInterval(1)

	gl.BindBuffer(gl.PIXEL_UNPACK_BUFFER, buffers[0])
	gl.BufferData(gl.PIXEL_UNPACK_BUFFER, width*height*3, nil, gl.DYNAMIC_DRAW)

	lastX, lastY = window.GetPos()
	lastWidth, lastHeight = window.GetSize()
	window.SetKeyCallback(onKeyPress)
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
	if x < 0 || x >= width || y < 0 && y >= height {
		return
	}
	index := lookupTable[y][x]
	r := byte(0)
	g := byte(0)
	b := byte(0)
	if len(color) == 1 {
		r = color[0]
		g = color[0]
		b = color[0]
	} else if len(color) == 2 {
		r = color[0]
		g = color[1]
	} else if len(color) == 3 {
		r = color[0]
		g = color[1]
		b = color[2]
	}

	pixelArr[index] = r
	pixelArr[index+1] = g
	pixelArr[index+2] = b
}

func UnsafePutPixel(x, y int, color ...byte) {
	index := lookupTable[y][x]
	r := byte(0)
	g := byte(0)
	b := byte(0)
	if len(color) == 1 {
		r = color[0]
		g = color[0]
		b = color[0]
	} else if len(color) == 2 {
		r = color[0]
		g = color[1]
	} else if len(color) == 3 {
		r = color[0]
		g = color[1]
		b = color[2]
	}

	pixelArr[index] = r
	pixelArr[index+1] = g
	pixelArr[index+2] = b
}

func GetPixel(x, y int) []byte {
	if x < 0 || x >= width || y < 0 && y >= height {
		return []byte{}
	}
	index := lookupTable[y][x]

	return []byte{pixelArr[index], pixelArr[index+1], pixelArr[index+2]}
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
	copy((*[width * height * 3]byte)(pboPtr)[:width*height*3], pixelArr)
	run()

	gl.DrawPixels(width, height, gl.RGB, gl.UNSIGNED_BYTE, nil)

	glfw.PollEvents()
	processInput(window)
	SwapBuffers()
	//ClearScreen()
}

func GetCurrentIndex() int {
	return index
}

func SwapBuffers() {
	window.SwapBuffers()
}

func ClearScreen() {
	for i := range pixelArr {
		pixelArr[i] = 0
	}
}
