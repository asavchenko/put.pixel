package bitreader

import (
	"fmt"
)

type sliceBitReader struct {
	b   byte
	idx int // global index in the arr
	bi  int // bit index in the current byte
	arr []byte
}

func GetNewSliceBitReader(data []byte) *sliceBitReader {
	return &sliceBitReader{0, 0, 0, data}
}

func (br *sliceBitReader) GetBit() (byte, error) {
	//fmt.Println("byte idx:", br.idx, "bit idx:", br.bi)
	if br.idx == 0 {
		br.idx += 1
		if len(br.arr) < 1 {
			fmt.Println("the length of the byte reader array is 0")
			return 0, fmt.Errorf("unexpected result")
		}
		br.b = br.arr[0]
		br.bi += 1
		return br.b & 0b00000001, nil
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
			fmt.Println("the byte reader index is", br.idx, "but the length of the byte reader array is", len(br.arr))
			return 0, fmt.Errorf("unexpected result")
		}
		br.b = br.arr[br.idx]
		br.idx += 1

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
			return result, err
		}
	}

	return result, nil
}

func (br *sliceBitReader) GetByte() (byte, error) {
	bits, err := br.GetBits(8)
	if err != nil {
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
		return make([]byte, 0), fmt.Errorf("index is out of range")
	}
	res := make([]byte, n)
	var err error
	for i := 0; i < n; i++ {
		res[i], err = br.GetByte()
		if err != nil {
			return res, err
		}
	}

	return res, nil
}

func (br *sliceBitReader) GetNthBitInByte(b byte, position int) byte {
	return getNthBitInByte(b, position)
}

func (br *sliceBitReader) ToInt(s []byte) int {
	return toInt(s)
}
