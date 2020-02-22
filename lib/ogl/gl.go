package ogl

// #include <string.h>
import "C"

import (
	"fmt"
	"log"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
	width  = 640
	height = 480
)

var pixelArr []byte
var window *glfw.Window
var buffer uint32

func init() {
	pixelArr = make([]byte, width*height*3)
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

	gl.GenBuffers(1, &buffer)
	gl.BindBuffer(gl.PIXEL_UNPACK_BUFFER, buffer)
	gl.BufferData(gl.PIXEL_UNPACK_BUFFER, width*height*3, nil, gl.DYNAMIC_DRAW)
}

func Close() {
	glfw.Terminate()
}

func PutPixel(x, y int, color byte) {
	index := (x + y*width) * 3
	if index < 0 {
		return
	}
	if index+2 > len(pixelArr)-1 {
		return
	}
	pixelArr[index] = color
	pixelArr[index+1] = color
	pixelArr[index+2] = color
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
func Draw() {
	if !window.ShouldClose() {
		draw(buffer, window)
	}
}

func draw(buffer uint32, window *glfw.Window) {
	gl.Clear(gl.COLOR_BUFFER_BIT)
	pboPtr := gl.MapBuffer(gl.PIXEL_UNPACK_BUFFER, gl.WRITE_ONLY)
	if pboPtr == nil {
		return
	}
	if !gl.UnmapBuffer(gl.PIXEL_UNPACK_BUFFER) {
		return
	}

	pixelArr = (*[width * height * 3]byte)(pboPtr)[:width*height*3]

	gl.BindBuffer(gl.PIXEL_UNPACK_BUFFER, buffer)
	gl.DrawPixels(width, height, gl.RGB, gl.UNSIGNED_BYTE, nil)
	glfw.PollEvents()
	window.SwapBuffers()
}
