package bitreader

type (
	BitReader interface {
		GetBit() (byte, error)
		GetBits(n int) ([]byte, error)
		GetByte() (byte, error)
		GetBytes(n int) ([]byte, error)
		GetRemainingData() ([]byte, error)
		GetNthBitInByte(b byte, position int) byte
		ToInt([]byte) int
		GetRawData() []byte
		GoToNextByte() error
		HasMoreData() bool
		GetPosition() int
	}
)

func toInt(s []byte) int {
	val := 1
	result := 0
	for i := 0; i < len(s); i++ {
		result += int(s[i]) * val
		val = val << 1
	}

	return result
}

func getNthBitInByte(b byte, position int) byte {
	switch position {
	case 0:
		return b & 0b00000001
	case 1:
		return (b & 0b00000010) >> 1
	case 2:
		return (b & 0b00000100) >> 2
	case 3:
		return (b & 0b00001000) >> 3
	case 4:
		return (b & 0b00010000) >> 4
	case 5:
		return (b & 0b00100000) >> 5
	case 6:
		return (b & 0b01000000) >> 6
	case 7:
		return (b & 0b10000000) >> 7
	}

	return 0
}

func reverseSliceByte(a []byte) []byte {
	s := make([]byte, len(a))
	copy(s, a)
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}

	return s
}
