package bitreader

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

type sliceBitReader struct {
	b   byte
	idx int // global index in the arr
	bi  int // bit index in the current byte
	arr []byte
}

func GetNewSliceBitReader(data []byte) *sliceBitReader {
	dat := make([]byte, len(data))
	for i, b := range data {
		dat[i] = b
	}
	return &sliceBitReader{0, 0, 0, dat}
}

func (br *sliceBitReader) GetPosition() int {
	return br.idx
}

func (br *sliceBitReader) ResetBitIndex() {
	br.bi = 0
}

func (br *sliceBitReader) GetRemainingData() ([]byte, error) {
	if br.idx == 0 {
		return br.arr, nil
	}

	return br.arr[br.idx:], nil
}

func (br *sliceBitReader) GoToNextByte() error {
	if br.bi == 0 {
		return nil
	}
	log("current byte is", br.b, "byte idx:", br.idx, "bit idx:", br.bi, "going to the next byte")
	br.idx++
	br.bi = 0
	if len(br.arr) <= br.idx {
		logError("the byte reader index is", br.idx, "but the length of the byte reader array is", len(br.arr))
		return fmt.Errorf("unexpected result")
	}
	br.b = br.arr[br.idx]

	log("after the move it's byte idx:", br.idx, "bit idx:", br.bi, "next byte:", br.b)

	return nil
}

func (br *sliceBitReader) GetRawData() []byte {
	data := make([]byte, len(br.arr))
	for i, b := range br.arr {
		data[i] = b
	}

	return data
}

func (br *sliceBitReader) GetBit() (byte, error) {
	if br.idx == 0 && br.bi == 0 {
		br.b = br.arr[0]
	}
	switch br.bi {
	case 0:
		br.bi += 1
		return br.b & 0b00000001, nil
	case 1:
		br.bi += 1
		return (br.b & 0b00000010) >> 1, nil
	case 2:
		br.bi += 1
		return (br.b & 0b00000100) >> 2, nil
	case 3:
		br.bi += 1
		return (br.b & 0b00001000) >> 3, nil
	case 4:
		br.bi += 1
		return (br.b & 0b00010000) >> 4, nil
	case 5:
		br.bi += 1
		return (br.b & 0b00100000) >> 5, nil
	case 6:
		br.bi += 1
		return (br.b & 0b01000000) >> 6, nil
	case 7:
		res := (br.b & 0b10000000) >> 7
		br.bi = 0
		if len(br.arr)-1 < br.idx {
			logError("the byte reader index is", br.idx, "but the length of the byte reader array is", len(br.arr))
			return 0, fmt.Errorf("unexpected result")
		}
		br.idx += 1
		br.b = br.arr[br.idx]

		return res, nil
	}

	return 0, fmt.Errorf("index is out of range")
}

func (br *sliceBitReader) GetBits(n int) ([]byte, error) {
	if n < 1 {
		return make([]byte, 0), fmt.Errorf("index is out of range")
	}
	var err error
	result := make([]byte, n)
	for i := 0; i < n; i++ {
		result[i], err = br.GetBit()
		if err != nil {
			logError(err, "on getting", n, "bits from the stream")
			return result, err
		}
	}

	return result, nil
}

func (br *sliceBitReader) GetByte() (byte, error) {
	bits, err := br.GetBits(8)
	if err != nil {
		logError(err)
		return 0, err
	}

	res := byte(0)
	i := byte(0b00000001)
	for _, b := range bits {
		res += b * i
		i = i << 1
	}

	return res, nil
}

func (br *sliceBitReader) GetBytes(n int) ([]byte, error) {
	if n < 1 {
		return br.GetRemainingData()
	}
	if n > len(br.arr) {
		n = len(br.arr)
	}
	res := make([]byte, n)
	var err error
	for i := 0; i < n; i++ {
		res[i], err = br.GetByte()
		if err != nil {
			logError(err)
			return res, err
		}
	}

	return res, nil
}

func (br *sliceBitReader) GetNthBitInByte(b byte, position int) byte {
	return getNthBitInByte(b, position)
}

func (br *sliceBitReader) BitsToInt(s []byte) int {
	return toInt(s)
}

func (br *sliceBitReader) BitsNot(s []byte) []byte {
	return not(s)
}

func (br *sliceBitReader) HasMoreData() bool {
	return br.idx < len(br.arr)
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
