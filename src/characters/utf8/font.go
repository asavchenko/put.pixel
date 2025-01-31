package utf8

import (
	"math"
	"time"

	"assa.com/put.pixel/lib/freetype"
)

var font freetype.Font
var shapes map[rune][]byte
var availableCharCodes []uint32
var ctrlShapeCh chan map[string]interface{}
var emptyShape = GetEmptyShape()

const WIDTH = 18
const HEIGHT = 27

func init() {
	if f, err := freetype.LoadFont("src/characters/utf8/LiberationMono-Regular.ttf"); err != nil {
		panic(err)
	} else {
		font = f
		shapes = make(map[rune][]byte, 0)
		availableCharCodes = f.GetAvailableCharCodes()
	}
	ctrlShapeCh = make(chan map[string]interface{}, 1024)
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

			shape = make([]byte, WIDTH*HEIGHT)
			for y := dy; y < h; y++ {
				for x := dx; x < w; x++ {
					//_x := int(math.Round((float64(x)*float64(gd.XMax-gd.XMin) + float64(w)*float64(gd.XMin) - float64(dx)*float64(gd.XMax)) / float64(w-dx)))
					//_y := int(math.Round((float64(y)*float64(gd.YMax-gd.YMin) + float64(h)*float64(gd.YMin) - float64(dy)*float64(gd.YMax)) / float64(h-dy)))
					_x := int(math.Round((float64(x)*float64(font.GetMaxX()-font.GetMinY()) + float64(w)*float64(font.GetMinX()) - float64(dx)*float64(font.GetMaxX())) / float64(w-dx)))
					_y := int(math.Round((float64(y)*float64(font.GetMaxY()-font.GetMinY()) + float64(h)*float64(font.GetMinY()) - float64(dy)*float64(font.GetMaxY())) / float64(h-dy)))
					if gd.ContainsPoint(_x, _y) {
						shape[y*WIDTH+x] = 1
					} else {
						shape[y*WIDTH+x] = 0
					}
				}
			}
			shapes[ch] = shape
			sendResponseWithTimeout(out, shape, 20*time.Millisecond)
		}
	}()
}

func GetAvailableCharCodes() []uint32 {
	return font.GetAvailableCharCodes()
}

func GetShapeWidth() int {
	return WIDTH
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

func GetEmptyShape() []byte {
	return make([]byte, WIDTH*HEIGHT)
}

func GetShape(ch rune) []byte {
	//fmt.Println(string(ch), ch)
	out := make(chan interface{}, 1)
	select {
	case ctrlShapeCh <- map[string]interface{}{
		"ch":  ch,
		"out": out,
	}:
	case <-time.After(3 * time.Second):
		return GetEmptyShape()
	}
	select {
	case response := <-out:
		res := response.([]byte)
		if len(res) < 1 {
			return GetEmptyShape()
		}
		return res
	case <-time.After(3 * time.Second):
		return GetEmptyShape()
	}
}

func reverseShape(arr []byte) []byte {
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}

	return arr
}
