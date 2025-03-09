package pngreader

type (
	PNGReader interface {
		GetImageWidth() int
		GetImageHeight() int
		GetBitDepth() byte
		GetColorType() byte
		GetCompressionMethod() byte
		GetFilterMethod() byte
		GetInterlaceMethod() byte
		GetRawImageData() []byte
		GetImageData() ([]byte, error)
		GetPalette() [][]byte
		Close() error
		GetBytesPerPixel() int
	}
)

const (
	ColorTypeGrayscale      = 0
	ColorTypeTruecolor      = 2
	ColorTypeIndexedColor   = 3
	ColorTypeGrayscaleAlpha = 4
	ColorTypeTruecolorAlpha = 6
)
