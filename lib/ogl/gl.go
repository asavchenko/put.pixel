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
var buffers [2]uint32
var index int
var lastX, lastY int
var lastWidth, lastHeight int
var keyCallbacks map[string]map[glfw.Key][]func()
var keyCombinationCallbacks map[string]map[string]interface{}
var pressedKeys []glfw.Key

func init() {
	keyCallbacks = make(map[string]map[glfw.Key][]func(), 0)
	keyCombinationCallbacks = make(map[string]map[string]interface{}, 0)
	pressedKeys = make([]glfw.Key, 0)
	pixelArr = make([]byte, width*height*4)
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

	gl.GenBuffers(2, &buffers[0])
	glfw.SwapInterval(1)

	gl.BindBuffer(gl.PIXEL_UNPACK_BUFFER, buffers[0])
	gl.BufferData(gl.PIXEL_UNPACK_BUFFER, width*height*4, nil, gl.DYNAMIC_DRAW)

	gl.BindBuffer(gl.PIXEL_UNPACK_BUFFER, buffers[1])
	gl.BufferData(gl.PIXEL_UNPACK_BUFFER, width*height*4, nil, gl.DYNAMIC_DRAW)
	lastX, lastY = window.GetPos()
	lastWidth, lastHeight = window.GetSize()
	window.SetKeyCallback(onKeyPress)
}

func CloseWindow() {
	window.SetShouldClose(true)
}

func Close() {
	glfw.Terminate()
}

func PutPixel(x, y int, color byte, alpha byte) {
	index := (x + y*width) * 4
	if index < 0 {
		return
	}
	if index+2 > len(pixelArr)-1 {
		return
	}
	pixelArr[index] = color
	pixelArr[index+1] = color
	pixelArr[index+2] = color
	pixelArr[index+3] = alpha
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

// CPU = draw | math (performance algorythm)
// CPU (math) + GPU (draw) - performance
func draw(window *glfw.Window, run func()) {
	pboPtr := gl.MapBuffer(gl.PIXEL_UNPACK_BUFFER, gl.WRITE_ONLY)
	if pboPtr == nil {
		return
	}
	if !gl.UnmapBuffer(gl.PIXEL_UNPACK_BUFFER) {
		return
	}

	pixelArr = (*[width * height * 4]byte)(pboPtr)[:width*height*4]
	run()
	gl.DrawPixels(width, height, gl.RGBA, gl.UNSIGNED_BYTE, nil)
	glfw.PollEvents()
	processInput(window)
	SwapBuffers()
	ClearScreen()
}

func SwapBuffers() {
	index++
	if index > 1 {
		index = 0
	}
	gl.BindBuffer(gl.PIXEL_UNPACK_BUFFER, buffers[index])
	if index == 0 {
		gl.BufferData(gl.PIXEL_UNPACK_BUFFER, width*height*4, nil, gl.DYNAMIC_DRAW)
	} else {
		gl.BufferData(gl.PIXEL_UNPACK_BUFFER, width*height*4, nil, gl.DYNAMIC_DRAW)
	}
	window.SwapBuffers()
}

func ClearScreen() {
	for i := range pixelArr {
		pixelArr[i] = 0
	}
	//C.memset(unsafe.Pointer(&pixelArr[0]), 0, width*height*4)
}
