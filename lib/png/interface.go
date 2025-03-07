package png

type (
	PNGReader interface {
		GetImageWidth() int
		GetImageHeight() int
		GetBitDepth() byte
		GetColorType() byte
		GetCompressionMethod() byte
		GetFilterMethod() byte
		GetInterlaceMethod() byte
		GetImageData() []byte
		Close() error
	}
)
