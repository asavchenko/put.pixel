package main

import (
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"time"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
	width  = 640
	height = 480
)

var pixelArr []byte
var pixelArr1 []byte
var pixelArr2 []byte

func init() {
	runtime.LockOSThread()
}

func main() {
	pixelArr = make([]byte, width*height*3)
	pixelArr1 = make([]byte, width*height*3)
	pixelArr2 = make([]byte, width*height*3)
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
	window.MakeContextCurrent() // Now that the context has been made current, you can call gl.Init().

	if err := gl.Init(); err != nil {
		log.Fatalln("failed to initialize gl bindings:", err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)
	var buffers [2]uint32
	gl.GenBuffers(2, &buffers[0])

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			color := getRandValue(0, 255)
			index := x + y*width
			pixelArr1[index*3] = color
			pixelArr1[index*3+1] = color
			pixelArr1[index*3+2] = color
		}
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			color := getRandValue(0, 255)
			index := x + y*width
			pixelArr2[index*3] = color
			pixelArr2[index*3+1] = color
			pixelArr2[index*3+2] = color
		}
	}
	gl.BindBuffer(gl.PIXEL_UNPACK_BUFFER, buffers[0])
	gl.BufferData(gl.PIXEL_UNPACK_BUFFER, width*height*3, gl.Ptr(pixelArr1), gl.DYNAMIC_DRAW)

	gl.BindBuffer(gl.PIXEL_UNPACK_BUFFER, buffers[0])
	gl.BufferData(gl.PIXEL_UNPACK_BUFFER, width*height*3, gl.Ptr(pixelArr2), gl.DYNAMIC_DRAW)
	index := 0
	t1 := time.NewTicker(200 * time.Millisecond)
	t2 := time.NewTicker(100 * time.Millisecond)

	for !window.ShouldClose() {
		select {
		case <-t1.C:
			draw(index, buffers, window)
		case <-t2.C:
			index++
			if index > 1 {
				index = 0
			}
			for y := 0; y < height; y++ {
				for x := 0; x < width; x++ {
					putPixel(x, y, getRandValue(0, 255), index)
				}
			}
			gl.BindBuffer(gl.PIXEL_UNPACK_BUFFER, buffers[index])
			if index == 0 {
				gl.BufferData(gl.PIXEL_UNPACK_BUFFER, width*height*3, gl.Ptr(pixelArr1), gl.DYNAMIC_DRAW)
			} else {
				gl.BufferData(gl.PIXEL_UNPACK_BUFFER, width*height*3, gl.Ptr(pixelArr2), gl.DYNAMIC_DRAW)
			}
		}

	}
}

func draw(index int, buffers [2]uint32, window *glfw.Window) {
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.BindBuffer(gl.PIXEL_UNPACK_BUFFER, buffers[index])
	gl.DrawPixels(width, height, gl.RGB, gl.UNSIGNED_BYTE, nil)
	glfw.PollEvents()
	window.SwapBuffers()
}

func getRandValue(min, max int) byte {
	return byte(rand.Intn(max-min) + min)
}

func putPixel(x, y int, color byte, i int) {
	index := x + y*width
	if i == 0 {
		pixelArr1[index*3] = color
		pixelArr1[index*3+1] = color
		pixelArr1[index*3+2] = color
		return
	}
	pixelArr2[index*3] = color
	pixelArr2[index*3+1] = color
	pixelArr2[index*3+2] = color
}
