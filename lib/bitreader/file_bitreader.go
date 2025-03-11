package bitreader

import (
	"bytes"
	"fmt"
	"os"
)

type fileBitReader struct {
	b    byte
	idx  int // global index in the file
	bi   int // bit index in the current byte
	file *os.File
}

func GetNewFileBitReader(r *os.File) *fileBitReader {
	return &fileBitReader{0, 0, 0, r}
}

func (br *fileBitReader) GetPosition() int {
	return br.idx
}

func (br *fileBitReader) ResetBitIndex() {
	br.bi = 0
}

func (br *fileBitReader) GetRemainingData() ([]byte, error) {
	fi, err := br.file.Stat()
	if err != nil {
		return nil, err
	}

	data := make([]byte, fi.Size()-int64(br.idx))
	n, err := br.file.Read(data)
	if err != nil {
		return nil, err
	}
	if n != len(data) {
		return nil, fmt.Errorf("unexpected result")
	}

	return data, nil
}

func (br *fileBitReader) HasMoreData() bool {
	fi, err := br.file.Stat()
	if err != nil {
		return false
	}

	return br.idx < int(fi.Size())
}

func (br *fileBitReader) GetRawData() []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(br.file)
	br.file.Seek(0, 0)

	return buf.Bytes()
}

func (br *fileBitReader) GoToNextByte() error {
	if br.bi == 0 {
		return nil
	}
	data := make([]byte, 1)
	br.bi = 0
	n, err := br.file.Read(data)
	if err != nil {
		return err
	}
	if n != 1 {
		return fmt.Errorf("unexpected result")
	}
	br.b = data[0]
	br.idx += 1

	return nil
}

func (br *fileBitReader) GetBit() (byte, error) {
	if br.idx == 0 && br.bi == 0 {
		b, err := br.getNextByte()
		if err != nil {
			return 0, err
		}
		br.b = b
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
		b, err := br.getNextByte()
		if err != nil {
			return 0, err
		}
		br.b = b

		return res, nil
	}

	return 0, fmt.Errorf("index is out of range")
}

func (br *fileBitReader) getNextByte() (byte, error) {
	data := make([]byte, 1)

	n, err := br.file.Read(data)
	if err != nil {
		return 0, err
	}
	if n != 1 {
		return 0, fmt.Errorf("unexpected result")
	}

	return data[0], nil
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

func (br *fileBitReader) BitsToInt(s []byte) int {
	return toInt(s)
}

func (br *fileBitReader) BitsNot(s []byte) []byte {
	return not(s)
}
