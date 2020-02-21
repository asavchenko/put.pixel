package main

import (
	"fmt"
	"log"
	"math/rand"
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
	width  = 640
	height = 480
)

var pixelArr []byte

func init() {
	runtime.LockOSThread()
}

func main() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	window, err := glfw.CreateWindow(640, 480, "Title", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		log.Fatalln("failed to initialize gl bindings:", err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)
	var buffer uint32
	gl.GenBuffers(1, &buffer)
	gl.BindBuffer(gl.PIXEL_UNPACK_BUFFER, buffer)
	gl.BufferData(gl.PIXEL_UNPACK_BUFFER, width*height*3, nil, gl.STATIC_DRAW)
	for !window.ShouldClose() {
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
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			putPixel(x, y, getRandValue(0, 255))
		}
	}
	gl.BindBuffer(gl.PIXEL_UNPACK_BUFFER, buffer)
	gl.DrawPixels(width, height, gl.RGB, gl.UNSIGNED_BYTE, nil)
	glfw.PollEvents()
	window.SwapBuffers()
}

func getRandValue(min, max int) byte {
	return byte(rand.Intn(max-min) + min)
}

func putPixel(x, y int, color byte) {
	index := x + y*width
	pixelArr[index*3] = color
	pixelArr[index*3+1] = color
	pixelArr[index*3+2] = color
}
