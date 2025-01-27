package utf8

import (
	"math"
	"time"

	"assa.com/put.pixel/lib/freetype"
)

var font freetype.Font
var shapes map[rune][][]byte
var availableCharCodes []uint32
var ctrlShapeCh chan map[string]interface{}

const WIDTH = 24
const HEIGHT = 32

func init() {
	if f, err := freetype.LoadFont("src/characters/utf8/LiberationMono-Regular.ttf"); err != nil {
		panic(err)
	} else {
		font = f
		shapes = make(map[rune][][]byte, 0)
		availableCharCodes = f.GetAvailableCharCodes()
	}
	ctrlShapeCh = make(chan map[string]interface{}, 1024)
	emptyShape := getEmptyShape()
	go func() {
		for request := range ctrlShapeCh {
			ch := request["ch"].(rune)
			out := request["out"].(chan interface{})

			shape, exists := shapes[ch]
			if exists {
				sendResponseWithTimeout(out, shape, 20*time.Millisecond)
				continue
			}
			found := false
			for _, c := range availableCharCodes {
				if rune(c) == ch {
					found = true
					break
				}
			}
			if !found {
				sendResponseWithTimeout(out, emptyShape, 20*time.Millisecond)
				continue
			}
			if ch == ' ' {
				shapes[ch] = emptyShape

				sendResponseWithTimeout(out, shape, 20*time.Millisecond)
				continue
			}

			gd := font.GetGlyphData(uint(ch))
			if gd == nil {
				sendResponseWithTimeout(out, emptyShape, 20*time.Millisecond)
				continue
			}
			dx := 0
			dy := 0
			w := WIDTH
			h := HEIGHT

			shape = make([][]byte, 0)
			for y := dy; y < h; y++ {
				line := make([]byte, 0)
				for x := dx; x < w; x++ {
					//_x := int(math.Round((float64(x)*float64(gd.XMax-gd.XMin) + float64(w)*float64(gd.XMin) - float64(dx)*float64(gd.XMax)) / float64(w-dx)))
					//_y := int(math.Round((float64(y)*float64(gd.YMax-gd.YMin) + float64(h)*float64(gd.YMin) - float64(dy)*float64(gd.YMax)) / float64(h-dy)))
					_x := int(math.Round((float64(x)*float64(font.GetMaxX()-font.GetMinY()) + float64(w)*float64(font.GetMinX()) - float64(dx)*float64(font.GetMaxX())) / float64(w-dx)))
					_y := int(math.Round((float64(y)*float64(font.GetMaxY()-font.GetMinY()) + float64(h)*float64(font.GetMinY()) - float64(dy)*float64(font.GetMaxY())) / float64(h-dy)))
					if gd.ContainsPoint(_x, _y) {
						line = append(line, 1)
					} else {
						line = append(line, 0)
					}
				}
				shape = append(shape, line)
			}
			shapes[ch] = reverseShape(shape)
			sendResponseWithTimeout(out, shape, 20*time.Millisecond)
		}
	}()
}

func GetShapeWidth() int {
	w := float64(WIDTH)
	return int(w / 1.5)
}

func GetShapeHeight() int {
	return HEIGHT
}

func sendResponseWithTimeout(out chan interface{}, response interface{}, timeout time.Duration) {
	select {
	case out <- response:
	case <-time.After(timeout):
	}
}

func getEmptyShape() [][]byte {
	shape := make([][]byte, 0)
	for i := 0; i < HEIGHT; i++ {
		line := make([]byte, 0)
		for j := 0; j < WIDTH; j++ {
			line = append(line, 0)
		}
		shape = append(shape, line)
	}

	return shape
}

func GetShape(ch rune) [][]byte {
	//fmt.Println(string(ch), ch)
	out := make(chan interface{}, 1)
	select {
	case ctrlShapeCh <- map[string]interface{}{
		"ch":  ch,
		"out": out,
	}:
	case <-time.After(3 * time.Second):
		return getEmptyShape()
	}
	select {
	case response := <-out:
		return response.([][]byte)
	case <-time.After(3 * time.Second):
		return getEmptyShape()
	}
}

func reverseShape(arr [][]byte) [][]byte {
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}

	return arr
}
