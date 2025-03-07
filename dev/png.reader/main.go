package main

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"assa.com/put.pixel/lib/bitreader"
	"assa.com/put.pixel/lib/png"
)

func main() {
	pathToPng := "sample.png"
	pngReader, err := png.GetNew(pathToPng)
	if err != nil {
		panic(err)
	}
	defer pngReader.Close()

	width := pngReader.GetImageWidth()
	height := pngReader.GetImageHeight()
	log("width:", width, "height:", height)
	// Bit depth is a single-byte integer giving the number of bits per sample or per palette index (not per pixel).
	// Valid values are 1, 2, 4, 8, and 16, although not all values are allowed for all color types.
	// Color type is a single-byte integer that describes the interpretation of the image data.
	// Color type codes represent sums of the following values:
	// 1 (palette used),
	// 2 (color used),
	// and 4 (alpha channel used).
	// Valid values are 0, 2, 3, 4, and 6.

	switch pngReader.GetColorType() {
	case 0: // Each pixel is a grayscale sample.
		log("it's a grayscale image")
		switch pngReader.GetBitDepth() {
		case 1:
			log("it's a black and white image")
		case 2:
			log("it's a grayscale image with 4 shades")
		case 4:
			log("it's a grayscale image with 16 shades")
		case 8:
			log("it's a grayscale image with 256 shades")
			data := pngReader.GetImageData()
			log("size of the data is", len(data))
			br := bitreader.GetNewSliceBitReader(data)
			i := 0
			for {
				if i >= height-1 {
					break
				}
				sline, err := br.GetBytes(width + 1)
				if err != nil {
					panic(err)
				}
				i++
				log("filter type for scanline:", i, "is", sline[0], printBits(sline[0]))
			}
		case 16:
			log("it's a grayscale image with 65536 shades")
			log("which means we will need to cut it under 256 shades")
			// *row++ = (png_byte)(color >> 8);
			//*row++ = (png_byte)(color & 0xFF);
		}
	case 2: // Each pixel is an R,G,B triple.
		log("it's a RGB image")
		switch pngReader.GetBitDepth() {
		case 8:
			log("it's a RGB image with 256 shades")
		case 16:
			log("it's a RGB image with 65536 shades")
		}
	case 3: // Each pixel is a palette index; a PLTE chunk must appear.
		log("it's a palette image")
		switch pngReader.GetBitDepth() {
		case 1:
			log("it's a palette image with 2 colors")
		case 2:
			log("it's a palette image with 4 colors")
		case 4:
			log("it's a palette image with 16 colors")
		case 8:
			log("it's a palette image with 256 colors")
		}
	case 4: // Each pixel is a grayscale sample, followed by an alpha sample..
		log("it's a grayscale with alpha image")
		switch pngReader.GetBitDepth() {
		case 8:
			log("it's a grayscale with alpha image with 256 shades")
		case 16:
			log("it's a grayscale with alpha image with 65536 shades")
		}
	case 6: // Each pixel is an R,G,B triple, followed by an alpha sample.
		log("it's a RGB with alpha image")
		switch pngReader.GetBitDepth() {
		case 8:
			log("it's a RGB with alpha image with 256 shades")
		case 16:
			log("it's a RGB with alpha image with 65536 shades")
		}
	}
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
	//   1       Sub
	//   2       Up
	//   3       Average
	//   4       Paeth

}

// ToInt converts a unknown value to int
func toIntRaw(value interface{}) (int, error) {
	switch v := value.(type) {
	case bool:
		if v {
			return 1, nil
		} else {
			return 0, nil
		}
	}
	result, err := toFloat64(value)
	if err != nil {
		return 0, fmt.Errorf("on conversion to int got unexpected value %#v", value)
	}

	return int(result), nil
}

func toInt(value interface{}) int {
	v, _ := toIntRaw(value)

	return v
}

// ToFloat64 converts a unknown value to float64
func toFloat64(i interface{}) (float64, error) {
	switch s := i.(type) {
	case float64:
		return s, nil
	case float32:
		return float64(s), nil
	case int64:
		return float64(s), nil
	case int32:
		return float64(s), nil
	case int16:
		return float64(s), nil
	case int8:
		return float64(s), nil
	case int:
		return float64(s), nil
	case uint:
		return float64(s), nil
	case uint8:
		return float64(s), nil
	case uint16:
		return float64(s), nil
	case uint32:
		return float64(s), nil
	case uint64:
		return float64(s), nil
	case uintptr:
		return float64(s), nil
	case time.Duration:
		return float64(s), nil
	case time.Month:
		return float64(s), nil
	case time.Weekday:
		return float64(s), nil
	case nil:
		return float64(0), nil
	case time.Time:
		return float64(s.Unix()), nil
	case bool:
		if s {
			return float64(1), nil
		}
		return float64(0), nil
	case json.Number:
		if result, err := strconv.ParseFloat(string(s), 64); err != nil {
			return float64(0), fmt.Errorf("on conversion to float64 got unexpected value %#v", i)

		} else {
			return result, nil
		}
	case string:
		if v, err := strconv.ParseFloat(s, 64); err == nil {
			return v, nil
		}
		return float64(0), fmt.Errorf("on conversion to float64 got unexpected value %#v", i)
	default:
		return float64(0), fmt.Errorf("on conversion to float64 got unexpected value %#v, %T", i, i)
	}
}

func log(msgs ...interface{}) {
	_, f, line, _ := runtime.Caller(1)
	baseDir := getBaseDir()

	relativePath := strings.Replace(f, baseDir+"/", "", -1)
	formatedMsg := fmt.Sprintln(time.Now().UTC().Format("15:04:05.999 02-01-2006"), fmt.Sprintf("%s:%d", relativePath, line), msgs)
	fmt.Print(formatedMsg)
}

func logError(msgs ...interface{}) {
	_, f, line, _ := runtime.Caller(1)
	baseDir := getBaseDir()

	relativePath := strings.Replace(f, baseDir+"/", "", -1)
	formatedMsg := fmt.Sprintln("*********ERROR", time.Now().UTC().Format("15:04:05.999 02-01-2006"), fmt.Sprintf("%s:%d", relativePath, line), msgs)
	fmt.Print(formatedMsg)
}

func getBaseDir() string {
	pwd, err := os.Getwd()
	if err != nil {
		return ""
	}

	return pwd
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

func printBytes(bytes []byte) string {
	msg := ""
	for _, b := range bytes {
		msg += fmt.Sprintf("0x%02X ", b)
	}
	return msg
}
