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
	"assa.com/put.pixel/lib/ogl"
	"assa.com/put.pixel/lib/pngreader"
)

func init() {
	fmt.Println("init")
	//syscall.Setpriority(syscall.PRIO_PROCESS, os.Getpid(), -20)
	runtime.LockOSThread()
}

func main() {
	//pathToPng := "sample.png"
	pathToPng := "face.png" // multiple IDAT chunks
	//pathToPng := "download.png" // it's a palette image with 256 colors
	//pathToPng := "logo5.png" // it's a RGB with alpha image with 256 shades

	pngReader, err := pngreader.GetNew(pathToPng)
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
	// and 4 (alpha channel used)
	// Valid values are 0, 2, 3, 4, and 6.
	filteredData := make([]byte, 0)
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

			// Compare the two byte slices
			log("size of the data is", len(data))
			br := bitreader.GetNewSliceBitReader(data)
			i := 0
			for {
				if i >= height-1 {
					break
				}
				scanLine, err := br.GetBytes(width + 1)
				if err != nil {
					panic(err)
				}

				filterType := byte(0)
				filterType, scanLine = scanLine[0], scanLine[1:]
				bpp := 1
				switch filterType {
				case 0:
					log("Filter Type None")
					for _, v := range scanLine {
						filteredData = append(filteredData, v)
					}
				case 1:
					log("Filter Type Sub")
					//The Sub() filter transmits the difference between each byte and the value of the corresponding byte of the prior pixel.

					//To compute the Sub() filter, apply the following formula to each byte of the scanline:

					//Sub(x) = Raw(x) - Raw(x-bpp)
					// where x ranges from zero to the number of bytes representing the scanline minus one,
					// Raw() refers to the raw data byte at that byte position in the scanline,
					// and bpp is defined as the number of bytes per complete pixel, rounding up to one.
					//	For example, for color type 2 with a bit depth of 16, bpp is equal to 6 (three samples, two bytes per sample);
					//	for color type 0 with a bit depth of 2, bpp is equal to 1 (rounding up);
					//	for color type 4 with a bit depth of 16, bpp is equal to 4 (two-byte grayscale sample, plus two-byte alpha sample).

					// Note this computation is done for each byte, regardless of bit depth.
					// In a 16-bit image, each MSB is predicted from the preceding MSB and each LSB from the preceding LSB,
					// because of the way that bpp is defined.

					//Unsigned arithmetic modulo 256 is used, so that both the inputs and outputs fit into bytes. The sequence of Sub values is transmitted as the filtered scanline.

					//For all x < 0, assume Raw(x) = 0.

					//To reverse the effect of the Sub() filter after decompression, output the following value:

					//Sub(x) + Raw(x-bpp)
					//(computed mod 256), where Raw() refers to the bytes already decoded.

					for j := 0; j < len(scanLine); j++ {
						el := scanLine[j]
						if j >= bpp {
							el += filteredData[j-bpp]
						}
						filteredData = append(filteredData, el)
					}
				case 2:
					log("Filter Type Up")
					// The Up() filter is just like the Sub() filter except that the pixel immediately above the current pixel, rather than just to its left, is used as the predictor.

					// To compute the Up() filter, apply the following formula to each byte of the scanline:

					// Up(x) = Raw(x) - Prior(x)
					// where x ranges from zero to the number of bytes representing the scanline minus one, Raw() refers to the raw data byte at that byte position in the scanline,
					// and Prior(x) refers to the unfiltered bytes of the prior scanline.

					// Note this is done for each byte, regardless of bit depth. Unsigned arithmetic modulo 256 is used,
					// so that both the inputs and outputs fit into bytes. The sequence of Up values is transmitted as the filtered scanline.

					// On the first scanline of an image (or of a pass of an interlaced image), assume Prior(x) = 0 for all x.

					// To reverse the effect of the Up() filter after decompression, output the following value:

					// Up(x) + Prior(x)
					// (computed mod 256), where Prior() refers to the decoded bytes of the prior scanline.
					for j := 0; j < len(scanLine); j++ {
						el := scanLine[j]
						if i > 0 {
							el += filteredData[(i-1)*width+j]
						}
						filteredData = append(filteredData, el)
					}
				case 3:
					log("Filter Type Average")
					// The Average() filter uses the average of the two neighboring pixels (left and above) to predict the value of a pixel.

					// To compute the Average() filter, apply the following formula to each byte of the scanline:

					// Average(x) = Raw(x) - floor((Raw(x-bpp)+Prior(x))/2)
					// where x ranges from zero to the number of bytes representing the scanline minus one,
					// Raw() refers to the raw data byte at that byte position in the scanline,
					// Prior() refers to the unfiltered bytes of the prior scanline, and bpp is defined as for the Sub() filter.

					// Note this is done for each byte, regardless of bit depth. The sequence of Average values is transmitted as the filtered scanline.

					// The subtraction of the predicted value from the raw byte must be done modulo 256,
					// so that both the inputs and outputs fit into bytes. However, the sum Raw(x-bpp)+Prior(x) must be formed without overflow
					// (using at least nine-bit arithmetic). floor() indicates that the result of the division is rounded to the next lower integer
					// if fractional; in other words, it is an integer division or right shift operation.

					// For all x < 0, assume Raw(x) = 0. On the first scanline of an image (or of a pass of an interlaced image), assume Prior(x) = 0 for all x.

					// To reverse the effect of the Average() filter after decompression, output the following value:

					// Average(x) + floor((Raw(x-bpp)+Prior(x))/2)
					// where the result is computed mod 256, but the prediction is calculated in the same way as for encoding.
					// Raw() refers to the bytes already decoded, and Prior() refers to the decoded bytes of the prior scanline.
					for j := 0; j < len(scanLine); j++ {
						a, b := 0, 0
						if i > 0 {
							a = int(filteredData[(i-1)*width+j])
						}
						if j >= bpp {
							b = int(filteredData[j-bpp])
						}

						filteredData = append(filteredData, scanLine[j]+byte((a+b)/2))
					}
				case 4:
					log("Filter Type Paeth")
					// The Paeth() filter computes a simple linear function of the three neighboring pixels
					// (left, above, upper left), then chooses as predictor the neighboring pixel closest
					//	to the computed value. This technique is due to Alan W. Paeth [PAETH].

					// To compute the Paeth() filter, apply the following formula to each byte of the scanline:
					//
					// Paeth(x) = Raw(x) - PaethPredictor(Raw(x-bpp), Prior(x), Prior(x-bpp))
					//	where x ranges from zero to the number of bytes representing the scanline minus one,
					//	Raw() refers to the raw data byte at that byte position in the scanline,
					//	Prior() refers to the unfiltered bytes of the prior scanline, and bpp is defined as for the Sub() filter.
					//
					// Note this is done for each byte, regardless of bit depth.
					// Unsigned arithmetic modulo 256 is used, so that both the inputs and outputs fit into bytes.
					// The sequence of Paeth values is transmitted as the filtered scanline.
					//
					// The PaethPredictor() function is defined by the following pseudocode:
					//
					// function PaethPredictor (a, b, c)
					//		begin
					//			; a = left, b = above, c = upper left
					//			p := a + b - c        ; initial estimate
					//			pa := abs(p - a)      ; distances to a, b, c
					//			pb := abs(p - b)
					//			pc := abs(p - c)
					//			; return nearest of a,b,c,
					//			; breaking ties in order a,b,c.
					//			if pa <= pb AND pa <= pc then
					//				return a
					//			else if pb <= pc then
					//				return b
					//			else
					//				return c
					//		end
					//	The calculations within the PaethPredictor() function must be performed exactly, without overflow.
					//	Arithmetic modulo 256 is to be used only for the final step of subtracting the function result from the target byte value.
					//
					// Note that the order in which ties are broken is critical and must not be altered.
					// The tie break order is: pixel to the left, pixel above, pixel to the upper left.
					// (This order differs from that given in Paeth's article.)
					//
					// For all x < 0, assume Raw(x) = 0 and Prior(x) = 0.
					//	On the first scanline of an image (or of a pass of an interlaced image), assume Prior(x) = 0 for all x.
					//
					//	To reverse the effect of the Paeth() filter after decompression, output the following value:
					//
					//	Paeth(x) + PaethPredictor(Raw(x-bpp), Prior(x), Prior(x-bpp))
					//	(computed mod 256), where Raw() and Prior() refer to bytes already decoded. Exactly the same PaethPredictor() function is used by both encoder and decoder.
					for j := 0; j < len(scanLine); j++ {
						a, b, c := byte(0), byte(0), byte(0)
						if i > 0 {
							b = filteredData[(i-1)*width+j]
						}
						if j >= bpp {
							a = filteredData[j-bpp]
						}
						if i > 0 && j >= bpp {
							c = filteredData[(i-1)*width+j-bpp]
						}

						filteredData = append(filteredData, scanLine[j]+paeth(a, b, c))
					}
				default:
					log(filterType)
					panic("unexpected filter type")
				}
				i++
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
	//  1       Sub
	//  2       Up
	//  3       Average
	//  4       Paeth
	ogl.Init(false)
	defer ogl.Close()
	ogl.OnKeypress(ogl.KEY_ESC, func() {
		ogl.CloseWindow()
	})

	filteredData = filteredData
	for {
		if ogl.IsExit() {
			break
		}
		ogl.Draw(func() {
			for i := 0; i < height; i++ {
				for j := 0; j < width; j++ {
					index := i*width + j
					if index >= len(filteredData) {
						break
					}
					color := filteredData[index]
					ogl.PutPixel(j, height-i, color)
				}
			}
		})
	}

}

func paeth(a, b, c byte) byte {
	p := int(a) + int(b) - int(c) // extend the gradient
	// return whatever input is closest to p
	pa := abs(p - int(a))
	pb := abs(p - int(b))
	pc := abs(p - int(c))
	if pa <= pb && pa <= pc {
		return a
	}
	if pb <= pc {
		return b
	}
	return c
}

func abs(a int) int {
	if a < 0 {
		return -a
	}

	return a
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

func reverseSliceByte(a []byte) []byte {
	s := make([]byte, len(a))
	copy(s, a)
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}

	return s
}
