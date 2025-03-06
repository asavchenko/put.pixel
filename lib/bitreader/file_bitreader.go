package bitreader

import (
	"fmt"
	"os"
)

type fileBitReader struct {
	b   byte
	idx int // global index in the file
	bi  int // bit index in the current byte
	f   *os.File
}

func GetNewFileBitReader(r *os.File) *fileBitReader {
	return &fileBitReader{0, 0, 0, r}
}

func (br *fileBitReader) GetBit() (byte, error) {
	if br.idx == 0 {
		br.idx += 1
		data := make([]byte, 1)

		n, err := br.f.Read(data)
		if err != nil {
			return 0, err
		}
		if n != 1 {
			return 0, fmt.Errorf("unexpected result")
		}
		br.b = data[0]
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
		br.idx += 1
		data := make([]byte, 1)

		n, err := br.f.Read(data)
		if err != nil {
			return 0, err
		}
		if n != 1 {
			return 0, fmt.Errorf("unexpected result")
		}
		br.b = data[0]
		return res, nil
	}

	return 0, fmt.Errorf("index is out of range")
}

func (br *fileBitReader) GetBits(n int) ([]byte, error) {
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

func (br *fileBitReader) GetByte() (byte, error) {
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

func (br *fileBitReader) GetBytes(n int) ([]byte, error) {
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

func (br *fileBitReader) GetNthBitInByte(b byte, position int) byte {
	return getNthBitInByte(b, position)
}

func (br *fileBitReader) ToInt(s []byte) int {
	return toInt(s)
}
