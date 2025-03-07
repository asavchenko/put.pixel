package main

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

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
		switch pngReader.GetBitDepth() {
		case 1:
		case 2:
		case 4:
		case 8:
		case 16:
		}
	case 2: // Each pixel is an R,G,B triple.
		switch pngReader.GetBitDepth() {
		case 8:
		case 16:
		}
	case 3: // Each pixel is a palette index; a PLTE chunk must appear.
		switch pngReader.GetBitDepth() {
		case 1:
		case 2:
		case 4:
		case 8:
		}
	case 4: // Each pixel is a grayscale sample, followed by an alpha sample..
		switch pngReader.GetBitDepth() {
		case 8:
		case 16:
		}
	case 6: // Each pixel is an R,G,B triple, followed by an alpha sample.
		switch pngReader.GetBitDepth() {
		case 8:
		case 16:
		}
	}
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
