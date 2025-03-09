package main

import (
	"fmt"
	"runtime"
	"time"

	"assa.com/put.pixel/lib/ogl"
	"assa.com/put.pixel/lib/pngreader"
)

func init() {
	fmt.Println("init")
	//syscall.Setpriority(syscall.PRIO_PROCESS, os.Getpid(), -20)
	runtime.LockOSThread()
}

func main() {
	pathToPng := "sample.png"
	//pathToPng := "goose.png"
	//pathToPng := "ChunkDataTooLarge.png"
	//pathToPng := "face.png" // multiple IDAT chunks
	//pathToPng := "download.png" // it's a palette image with 256 colors
	//pathToPng := "logo5.png" // it's a RGB with alpha image with 256 shades

	pngReader, err := pngreader.GetNew(pathToPng)
	if err != nil {
		panic(err)
	}
	defer pngReader.Close()

	// Conceptually, a PNG image is a rectangular pixel array, with pixels appearing left-to-right within each scanline, and scanlines appearing top-to-bottom.
	//(For progressive display purposes, the data may actually be transmitted in a different order;
	// see Interlaced data order.)
	// The size of each pixel is determined by the bit depth, which is the number of bits per sample in the image data.
	// there is an extra byte at the start of each row that signals the filter type
	// The filter type byte is not considered part of the image data, but it is included in the datastream sent to the compression step.
	// Scanlines always begin on byte boundaries.
	// When pixels have fewer than 8 bits and the scanline width is not evenly divisible by the number of pixels per byte,
	// the low-order bits in the last byte of each scanline are wasted. The contents of these wasted bits are unspecified.

	// Filter Types are
	//  0       None
	//  1       Sub
	//  2       Up
	//  3       Average
	//  4       Paeth
	ogl.Init(false)
	defer ogl.Close()
	ogl.OnKeypress(ogl.KEY_ESC, func() {
		ogl.CloseWindow()
	})

	imgData, err := pngReader.GetImageData()
	if err != nil {
		logError(err)
		//return
	}
	for {
		if ogl.IsExit() {
			break
		}
		ogl.Draw(func() {
			for i := 0; i < pngReader.GetImageHeight(); i++ {
				for j := 0; j < pngReader.GetImageWidth()*pngReader.GetBytesPerPixel(); j += pngReader.GetBytesPerPixel() {
					index := i*pngReader.GetImageWidth()*pngReader.GetBytesPerPixel() + j
					if index >= len(imgData) {
						break
					}
					switch pngReader.GetColorType() {
					case pngreader.ColorTypeGrayscale:
						ogl.PutPixel(j, pngReader.GetImageHeight()-i, imgData[index], imgData[index], imgData[index])
					case pngreader.ColorTypeTruecolor:
						ogl.PutPixel(j, pngReader.GetImageHeight()-i, imgData[index], imgData[index+1], imgData[index+2])
					case pngreader.ColorTypeIndexedColor:
						// each pixel is a palette index
						ogl.PutPixel(j, pngReader.GetImageHeight()-i, pngReader.GetPalette()[imgData[index]]...)
					case pngreader.ColorTypeGrayscaleAlpha:
						if imgData[index+1] != 0 {
							ogl.PutPixel(j, pngReader.GetImageHeight()-i, imgData[index], imgData[index], imgData[index])
						}
					case pngreader.ColorTypeTruecolorAlpha:
						if imgData[index+3] != 0 {
							ogl.PutPixel(j, pngReader.GetImageHeight()-i, imgData[index], imgData[index+1], imgData[index+2])
						}
					}
				}
			}
		})
	}

}

func log(msgs ...interface{}) {
	_, f, line, _ := runtime.Caller(1)

	relativePath := f
	formatedMsg := fmt.Sprintln(time.Now().UTC().Format("15:04:05.999 02-01-2006"), fmt.Sprintf("%s:%d", relativePath, line), msgs)
	fmt.Print(formatedMsg)
}

func logError(msgs ...interface{}) {
	_, f, line, _ := runtime.Caller(1)

	relativePath := f
	formatedMsg := fmt.Sprintln("*********ERROR", time.Now().UTC().Format("15:04:05.999 02-01-2006"), fmt.Sprintf("%s:%d", relativePath, line), msgs)
	fmt.Print(formatedMsg)
}
