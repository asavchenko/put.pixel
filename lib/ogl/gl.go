package ogl

// #include <string.h>
import "C"

import (
	"fmt"
	"log"
	"unsafe"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
	width  = 640
	height = 480
)

var pixelArr1 []byte
var pixelArr2 []byte
var window *glfw.Window
var buffers [2]uint32
var index int

func init() {
	pixelArr1 = make([]byte, width*height*3)
	pixelArr2 = make([]byte, width*height*3)
}

func Init() {
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
}

func Close() {
	glfw.Terminate()
}

func PutPixel(x, y int, color byte) {
	i := (x + y*width) * 3
	if i < 0 {
		return
	}
	if index == 0 {
		if i+2 > len(pixelArr1)-1 {
			return
		}
		pixelArr1[i] = color
		pixelArr1[i+1] = color
		pixelArr1[i+2] = color
		return
	}
	if i+2 > len(pixelArr2)-1 {
		return
	}
	pixelArr2[i] = color
	pixelArr2[i+1] = color
	pixelArr2[i+2] = color
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
	showPrevScreen(window)
	ClearScreen()
	run()
	glfw.PollEvents()
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
	if index == 0 {
		C.memset(unsafe.Pointer(&pixelArr1[0]), 0, width*height*3)
		return
	}

	C.memset(unsafe.Pointer(&pixelArr2[0]), 0, width*height*3)
}
