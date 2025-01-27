package ogl

// #include <string.h>
import "C"

import (
	"fmt"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"strconv"

	"log"
)

const (
	width  = 640
	height = 480
	Xm     = width - 1
	Ym     = height - 1
)

//var pixelArr *[width * height * 3]byte

var pixelArr []byte
var window *glfw.Window
var buffers [1]uint32
var index int
var lastX, lastY int
var lastWidth, lastHeight int
var keyCallbacks map[string]map[glfw.Key][]func()
var keyCombinationCallbacks map[string]map[string]interface{}
var pressedKeys []glfw.Key
var lookupTable []int

func init() {
	keyCallbacks = make(map[string]map[glfw.Key][]func(), 0)
	keyCombinationCallbacks = make(map[string]map[string]interface{}, 0)
	pressedKeys = make([]glfw.Key, 0)
	lookupTable = make([]int, height)
	for i := 0; i < height; i++ {
		lookupTable[i] = i * width * 3
	}
	pixelArr = make([]byte, width*height*3)
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
	if x < 0 || x >= width || y < 0 || y >= height {
		return
	}
	idx := lookupTable[y] + x<<1 + x
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

	pixelArr[idx] = r
	pixelArr[idx+1] = g
	pixelArr[idx+2] = b
}

func UnsafePutPixel(x, y int, color ...byte) {
	idx := lookupTable[y] + x<<1 + x
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

	pixelArr[idx] = r
	pixelArr[idx+1] = g
	pixelArr[idx+2] = b
}

func UnsafePutPixelRGB(x, y int, r, g, b byte) {
	i := lookupTable[y] + x<<1 + x

	pixelArr[i] = r
	pixelArr[i+1] = g
	pixelArr[i+2] = b
}

func PutPixelRGB(x, y int, r, g, b byte) {
	if x < 0 || x >= width || y < 0 || y >= height {
		return
	}
	idx := lookupTable[y] + x<<1 + x

	pixelArr[idx] = r
	pixelArr[idx+1] = g
	pixelArr[idx+2] = b
}

func GetPixel(x, y int) []byte {
	if x < 0 || x >= width || y < 0 && y >= height {
		return []byte{}
	}
	index := lookupTable[y] + x<<1 + x

	return []byte{pixelArr[index], pixelArr[index+1], pixelArr[index+2]}
}

func GetPixelUnsafe(x, y int) uint32 {
	i := lookupTable[y] + x<<1 + x

	return uint32(pixelArr[i]) + uint32(pixelArr[i+1])<<8 + uint32(pixelArr[i+2])<<16
}

func printBitsUint32(b uint32) string {
	// to binary representation
	str := strconv.FormatInt(int64(b), 2)
	// with leading zeros
	delta := 32 - len(str)
	for i := 0; i < delta; i++ {
		str = "0" + str
	}
	// from string to byte array
	chunk := make([]byte, 0)
	for _, ds := range str {
		d, _ := strconv.Atoi(string(ds))
		chunk = append(chunk, byte(d))
	}

	return fmt.Sprint(chunk)
}

func printBits(b byte) string {
	// to binary representation
	str := strconv.FormatInt(int64(b), 2)
	// with leading zeros
	delta := 8 - len(str)
	for i := 0; i < delta; i++ {
		str = "0" + str
	}
	// from string to byte array
	chunk := make([]byte, 0)
	for _, ds := range str {
		d, _ := strconv.Atoi(string(ds))
		chunk = append(chunk, byte(d))
	}

	return fmt.Sprint(chunk)
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
	//pixelArr = (*[width * height * 3]byte)(pboPtr)
	//pixelArr = unsafe.Slice((*byte)(pboPtr), width*height*3)
	//pixelArr = (*[width * height * 3]byte)(pboPtr)[:width*height*3]

	copy((*[width * height * 3]byte)(pboPtr)[:width*height*3], pixelArr)
	run()

	gl.DrawPixels(width, height, gl.RGB, gl.UNSIGNED_BYTE, nil)

	glfw.PollEvents()
	processInput(window)
	SwapBuffers()
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

func FillScreen(color byte) {
	for i := range pixelArr {
		pixelArr[i] = color
	}
}
