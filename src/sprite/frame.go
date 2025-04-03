package sprite

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"

	"assa.com/put.pixel/lib/log"
	"assa.com/put.pixel/lib/ogl"
	"assa.com/put.pixel/lib/pngreader"
)

type frame struct {
	imgData           []uint32
	width             int
	height            int
	direction         int
	allowedDirections []int
	imgName           string
}

func GetFrames(pathToFolder string, allowedDirections []int, direction int) []Frame {
	frames := make([]Frame, 0)
	if err := filepath.Walk(pathToFolder, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && pathToFolder != path {
			return filepath.SkipDir
		}
		if !strings.Contains(info.Name(), ".png") {
			return nil
		}
		f := GetNewFrame(path)
		if f != nil {
			f.SetImgName(info.Name())
			f.SetAllowedDirections(allowedDirections)
			f.SetDirection(direction)
			frames = append(frames, f)
		}
		return nil
	}); err != nil {
		log.Info(err)
	}

	return frames
}

func GetNewFrame(pathToImg string) Frame {
	f := &frame{}
	if err := f.SetImage(pathToImg); err != nil {
		log.Info(err)
		return nil
	}

	return f
}

func (f *frame) SetAllowedDirections(directions []int) Frame {
	f.allowedDirections = directions

	return f
}

func (f *frame) GetImgName() string {
	return f.imgName
}

func (f *frame) SetImgName(name string) Frame {
	f.imgName = name

	return f
}

func (f *frame) GetDirection() int {
	return f.direction
}

func (f *frame) GetWidth() int {
	return f.width
}

func (f *frame) GetHeight() int {
	return f.height
}

func (f *frame) SetImage(pathToImg string) error {
	reader, err := pngreader.GetNew(pathToImg)
	if err != nil {
		return err
	}

	f.width = reader.GetImageWidth()
	f.height = reader.GetImageHeight()

	imgData, err := reader.GetImageData()
	if err != nil {
		return err
	}

	imgDataRGBA := make([]uint32, len(imgData)>>2)
	for i := 0; i < len(imgData); i += 4 {
		imgDataRGBA[i/4] = uint32(imgData[i+3])<<24 | uint32(imgData[i+2])<<16 | uint32(imgData[i+1])<<8 | uint32(imgData[i])
	}
	flippedImgData := make([]uint32, 0)
	for j := 0; j < f.height; j++ {
		line := make([]uint32, f.width)
		startIdx := j * f.width
		copy(line[0:], imgDataRGBA[startIdx:startIdx+f.width])
		flippedImgData = append(line, flippedImgData...)
	}

	if reader.GetColorType() != pngreader.ColorTypeTruecolorAlpha {
		return fmt.Errorf("unsupported color scheme for marble asset")
	}

	f.imgData = flippedImgData

	return reader.Close()
}

func (f *frame) SetSize(width, height int) Frame {
	f.resize(width, height)
	f.width = width
	f.height = height

	return f
}

func (f *frame) SetDirection(direction int) Frame {
	found := false
	for _, ad := range f.allowedDirections {
		if ad != direction {
			continue
		}
		found = true
		break
	}

	if !found {
		log.Error("direction", direction, "is not allowed", "allowed directions are:", f.allowedDirections)
		return f
	}

	if direction == DIRECTION_LEFT && f.direction == DIRECTION_RIGHT {
		f.flip()
	} else if direction == DIRECTION_RIGHT && f.direction == DIRECTION_LEFT {
		f.flip()
	}

	f.direction = direction

	return f
}

func (f *frame) resize(width, height int) Frame {
	original := f.imgData
	ow := f.width
	oh := f.height
	nw := width
	nh := height
	kw := float64(nw) / float64(ow)
	kh := float64(nh) / float64(oh)
	resized := make([]uint32, nh*nw)
	if kw > 1 && kh > 1 {
		dj := 0
		for j := 0; j < nh; j++ {
			for i := 0; i < nw; i++ {
				y := int(math.Ceil(float64(j) / kh))
				x := int(math.Ceil(float64(i) / kw))
				if x >= ow {
					x = ow - 1
				}
				if y >= oh {
					y = oh - 1
				}

				resized[dj+i] = original[y*ow+x]
			}
			dj += nw
		}
		f.imgData = resized
		return f
	}
	dj := 0
	for j := 0; j < oh; j++ {
		for i := 0; i < ow; i++ {
			if original[dj+i] < 1 {
				continue
			}
			y := int(math.Ceil(float64(j) * kh))
			x := int(math.Ceil(float64(i) * kw))
			if y >= nh {
				y = nh - 1
			}
			if x >= nw {
				x = nw - 1
			}
			resized[y*nw+x] = original[dj+i]
		}
		dj += ow
	}

	f.imgData = resized

	return f
}

func (f *frame) Show(x, y int) Frame {
	w2 := f.width / 2
	h2 := f.height / 2
	ogl.PutBitmap(x-w2, y+h2, f.width, f.height, f.imgData)

	return f
}

func (f *frame) flip() {
	idx := 0
	for y := 0; y < f.height; y++ {
		for i, j := 0, f.width-1; i < j; i, j = i+1, j-1 {
			f.imgData[idx+i], f.imgData[idx+j] = f.imgData[idx+j], f.imgData[idx+i]
		}

		idx += f.width
	}
}
