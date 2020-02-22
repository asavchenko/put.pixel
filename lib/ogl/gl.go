package ogl

// #include <string.h>
import "C"

import (
	"fmt"
	"log"
	"unsafe"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	width  = 2304
	height = 648
	Xm     = width - 1
	Ym     = height - 1
)

var screen []byte
var pixelArr1 []byte
var pixelArr2 []byte
var yTable []int
var window *glfw.Window
var buffers [2]uint32
var buffer uint32
var index int
var doubleBuffering bool
var keyCallbacks map[string]map[glfw.Key][]func()
var keyCombinationCallbacks map[string]map[string]interface{}
var lastX, lastY int
var lastWidth, lastHeight int
var pressedKeys []glfw.Key

func init() {
	keyCallbacks = make(map[string]map[glfw.Key][]func(), 0)
	keyCombinationCallbacks = make(map[string]map[string]interface{}, 0)
	pressedKeys = make([]glfw.Key, 0)
	doubleBuffering = true
	pixelArr1 = make([]byte, width*height*3)
	pixelArr2 = make([]byte, width*height*3)
	yTable = make([]int, height)
	for i := 0; i < height; i++ {
		yTable[i] = i * width
	}
}

func Init(fullScreen bool) {
	if err := glfw.Init(); err != nil {
		log.Fatal("failed to initialize glfw:", err)
	}
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.DoubleBuffer, 1)
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

	gl.BindBuffer(gl.PIXEL_UNPACK_BUFFER, buffers[0])
	gl.BufferData(gl.PIXEL_UNPACK_BUFFER, width*height*3, gl.Ptr(pixelArr1), gl.DYNAMIC_DRAW)

	gl.BindBuffer(gl.PIXEL_UNPACK_BUFFER, buffers[1])
	gl.BufferData(gl.PIXEL_UNPACK_BUFFER, width*height*3, gl.Ptr(pixelArr2), gl.DYNAMIC_DRAW)

	lastX, lastY = window.GetPos()
	lastWidth, lastHeight = window.GetSize()
	window.SetKeyCallback(onKeyPress)
	screen = pixelArr1
}

func Close() {
	glfw.Terminate()
}

func Draw(run func()) {
	if !window.ShouldClose() {
		draw(window, run)
	}
}

func IsExit() bool {
	return window.ShouldClose()
}

func DisableDoubleBuffering() {
	doubleBuffering = false
	gl.GenBuffers(1, &buffer)
	gl.BindBuffer(gl.PIXEL_UNPACK_BUFFER, buffer)
	gl.BufferData(gl.PIXEL_UNPACK_BUFFER, width*height*3, nil, gl.DYNAMIC_DRAW)
}

func draw(window *glfw.Window, run func()) {
	if doubleBuffering {
		showPrevScreen(window)
		screen = pixelArr2
		if index == 1 {
			screen = pixelArr1
		}
		run()
		glfw.PollEvents()
		return
	}
	pboPtr := gl.MapBuffer(gl.PIXEL_UNPACK_BUFFER, gl.WRITE_ONLY)
	if pboPtr == nil {
		return
	}
	if !gl.UnmapBuffer(gl.PIXEL_UNPACK_BUFFER) {
		return
	}

	screen = (*[width * height * 3]byte)(pboPtr)[:width*height*3]
	run()
	gl.BindBuffer(gl.PIXEL_UNPACK_BUFFER, buffer)
	gl.DrawPixels(width, height, gl.RGB, gl.UNSIGNED_BYTE, nil)
	gl.Flush()
	window.SwapBuffers()
	glfw.PollEvents()
	processInput(window)
}

func CloseWindow() {
	window.SetShouldClose(true)
}

func showPrevScreen(window *glfw.Window) {
	if index == 0 {
		gl.BindBuffer(gl.PIXEL_UNPACK_BUFFER, buffers[1])
	} else {
		gl.BindBuffer(gl.PIXEL_UNPACK_BUFFER, buffers[0])
	}
	gl.DrawPixels(width, height, gl.RGB, gl.UNSIGNED_BYTE, nil)
	gl.Flush()
	window.SwapBuffers()
}

func SwapBuffers() {
	if !doubleBuffering {
		return
	}
	index++
	if index > 1 {
		index = 0
	}
	gl.BindBuffer(gl.PIXEL_UNPACK_BUFFER, buffers[index])
	if index == 0 {
		gl.BufferData(gl.PIXEL_UNPACK_BUFFER, width*height*3, gl.Ptr(pixelArr1), gl.DYNAMIC_DRAW)
	} else {
		gl.BufferData(gl.PIXEL_UNPACK_BUFFER, width*height*3, gl.Ptr(pixelArr2), gl.DYNAMIC_DRAW)
	}
}

func ClearScreen() {
	if !doubleBuffering {
		C.memset(unsafe.Pointer(&screen[0]), 0, width*height*3)
		return
	}
	if index == 0 {
		C.memset(unsafe.Pointer(&pixelArr1[0]), 0, width*height*3)
		return
	}

	C.memset(unsafe.Pointer(&pixelArr2[0]), 0, width*height*3)
}
