package pngreader

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"assa.com/put.pixel/lib/bitreader"
)

type pngReader struct {
	f                 *os.File
	width             int
	height            int
	bitDepth          byte
	colorType         byte
	compressionMethod byte
	filterMethod      byte
	interlaceMethod   byte
	imgData           []byte
	palette           [][]byte
	bpp               int
}

func (reader *pngReader) Close() error {
	return reader.f.Close()
}

func (reader *pngReader) GetImageWidth() int {
	return reader.width
}

func (reader *pngReader) GetImageHeight() int {
	return reader.height
}

func (reader *pngReader) GetBitDepth() byte {
	return reader.bitDepth
}

func (reader *pngReader) GetColorType() byte {
	return reader.colorType
}

func (reader *pngReader) GetCompressionMethod() byte {
	return reader.compressionMethod
}

func (reader *pngReader) GetFilterMethod() byte {
	return reader.filterMethod
}

func (reader *pngReader) GetInterlaceMethod() byte {
	return reader.interlaceMethod
}

func (reader *pngReader) GetBytesPerPixel() int {
	if reader.bpp != -1 {
		return reader.bpp
	}
	reader.bpp = 1
	switch reader.GetColorType() {
	case 0: // Each pixel is a grayscale sample.
		switch reader.GetBitDepth() {
		case 1:
			reader.bpp = 1
		case 2:
			reader.bpp = 1
		case 4:
			reader.bpp = 1
		case 8:
		case 16:
			reader.bpp = 2
			// *row++ = (png_byte)(color >> 8);
			//*row++ = (png_byte)(color & 0xFF);
		}
	case 2: // Each pixel is an R,G,B triple.
		switch reader.GetBitDepth() {
		case 8:
			reader.bpp = 3
		case 16:
			reader.bpp = 6
		}
	case 3: // Each pixel is a palette index; a PLTE chunk must appear.
		return reader.bpp
	case 4: // Each pixel is a grayscale sample, followed by an alpha sample..
		switch reader.GetBitDepth() {
		case 8:
			reader.bpp = 2
		case 16:
			reader.bpp = 4
		}
	case 6: // Each pixel is an R,G,B triple, followed by an alpha sample.
		switch reader.GetBitDepth() {
		case 8:
			reader.bpp = 4
		case 16:
			reader.bpp = 8
		}
	}

	return reader.bpp
}

func (reader *pngReader) GetImageData() ([]byte, error) {
	width := reader.GetImageWidth()
	//height := reader.GetImageHeight()
	//log("width:", width, "height:", height)
	// Bit depth is a single-byte integer giving the number of bits per sample or per palette index (not per pixel).
	// Valid values are 1, 2, 4, 8, and 16, although not all values are allowed for all color types.
	// Color type is a single-byte integer that describes the interpretation of the image data.
	// Color type codes represent sums of the following values:
	// 1 (palette used),
	// 2 (color used),
	// and 4 (alpha channel used)
	// Valid values are 0, 2, 3, 4, and 6.
	rowLen := width
	bpp := 1
	switch reader.GetColorType() {
	case 0: // Each pixel is a grayscale sample.
		switch reader.GetBitDepth() {
		case 1:
			log("it's a black and white image")
			bpp = 1
			rowLen = width / 8
		case 2:
			log("it's a grayscale image with 4 shades")
			bpp = 1
			rowLen = width / 4
		case 4:
			log("it's a grayscale image with 16 shades")
			bpp = 1
			rowLen = width / 2
		case 8:
			log("it's a grayscale image with 256 shades")
		case 16:
			log("it's a grayscale image with 65536 shades")
			log("which means we will need to cut it under 256 shades")
			bpp = 2
			rowLen = 2 * width
			// *row++ = (png_byte)(color >> 8);
			//*row++ = (png_byte)(color & 0xFF);
		}
	case 2: // Each pixel is an R,G,B triple.
		switch reader.GetBitDepth() {
		case 8:
			log("it's a RGB image with 256 shades")
			bpp = 3
			rowLen = 3 * width
		case 16:
			log("it's a RGB image with 65536 shades")
			bpp = 6
			rowLen = 6 * width
		}
	case 3: // Each pixel is a palette index; a PLTE chunk must appear.
		switch reader.GetBitDepth() {
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
		switch reader.GetBitDepth() {
		case 8:
			log("it's a grayscale with alpha image with 256 shades")
			rowLen = 2 * width
			bpp = 2
		case 16:
			log("it's a grayscale with alpha image with 65536 shades")
			rowLen = 4 * width
			bpp = 4
		}
	case 6: // Each pixel is an R,G,B triple, followed by an alpha sample.
		switch reader.GetBitDepth() {
		case 8:
			//log("it's a RGB with alpha image with 256 shades")
			rowLen = 4 * width
			bpp = 4
		case 16:
			//log("it's a RGB with alpha image with 65536 shades")
			rowLen = 8 * width
			bpp = 8
		}
	}
	return reader.filterImage(rowLen, bpp)
}

func (reader *pngReader) GetPalette() [][]byte {
	return reader.palette
}

func (reader *pngReader) filterImage(rowLen int, bpp int) ([]byte, error) {
	height := reader.height
	if reader.interlaceMethod == 1 {
		return reader.filterInterlacedImage(rowLen, bpp)
	}
	filteredData := make([]byte, rowLen*height)
	//log(len(filteredData), rowLen, height, bpp)
	data := reader.GetRawImageData()
	//log("decoded", len(data))
	br := bitreader.GetNewSliceBitReader(data)
	//log("row length is", rowLen)
	for i := 0; i < height; i++ {
		if err := br.GoToNextByte(); err != nil {
			logError(err)
			break
		}
		filterType, err := br.GetBits(8)
		if err != nil {
			return filteredData, err
		}
		scanLine, err := br.GetBytes(rowLen)
		if err != nil {
			return filteredData, err
		}
		scanLine = scanLine

		prevScanLine := make([]byte, rowLen)
		if i > 0 {
			prevScanLine = filteredData[(i-1)*rowLen : i*rowLen]
		}
		x := i * rowLen
		switch int(filterType[0] + filterType[1]*2 + filterType[2]*4) {
		case 0:
			for k, v := range scanLine {
				filteredData[x+k] = v
			}
		case 1:
			// The Sub() filter transmits the difference between each byte and the value of the corresponding byte of the prior pixel.

			// To compute the Sub() filter, apply the following formula to each byte of the scanline:

			// Sub(x) = Raw(x) - Raw(x-bpp)
			// where x ranges from zero to the number of bytes representing the scanline minus one,
			// Raw() refers to the raw data byte at that byte position in the scanline,
			// and bpp is defined as the number of bytes per complete pixel, rounding up to one.
			//	For example, for color type 2 with a bit depth of 16, bpp is equal to 6 (three samples, two bytes per sample);
			//	for color type 0 with a bit depth of 2, bpp is equal to 1 (rounding up);
			//	for color type 4 with a bit depth of 16, bpp is equal to 4 (two-byte grayscale sample, plus two-byte alpha sample).

			// Note this computation is done for each byte, regardless of bit depth.
			// In a 16-bit image, each MSB is predicted from the preceding MSB and each LSB from the preceding LSB,
			// because of the way that bpp is defined.

			// Unsigned arithmetic modulo 256 is used, so that both the inputs and outputs fit into bytes. The sequence of Sub values is transmitted as the filtered scanline.

			// For all x < 0, assume Raw(x) = 0.

			// To reverse the effect of the Sub() filter after decompression, output the following value:

			// Sub(x) + Raw(x-bpp) (computed mod 256),
			// where Raw() refers to the bytes already decoded.

			for j := 0; j < len(scanLine); j++ {
				el := scanLine[j]
				if j >= bpp {
					el += filteredData[x+j-bpp]
				}
				filteredData[x+j] = el
			}
		case 2:
			// The Up() filter is just like the Sub() filter except that the pixel immediately above the current pixel, rather than just to its left, is used as the predictor.

			// To compute the Up() filter, apply the following formula to each byte of the scanline:

			// Up(x) = Raw(x) - Prior(x)
			// where x ranges from zero to the number of bytes representing the scanline minus one, Raw() refers to the raw data byte at that byte position in the scanline,
			// and Prior(x) refers to the unfiltered bytes of the prior scanline.

			// Note this is done for each byte, regardless of bit depth. Unsigned arithmetic modulo 256 is used,
			// so that both the inputs and outputs fit into bytes. The sequence of Up values is transmitted as the filtered scanline.

			// On the first scanline of an image (or of a pass of an interlaced image), assume Prior(x) = 0 for all x.

			// To reverse the effect of the Up() filter after decompression, output the following value:

			// Up(x) + Prior(x) (computed mod 256),
			// where Prior() refers to the decoded bytes of the prior scanline.
			for j := 0; j < len(scanLine); j++ {
				filteredData[x+j] = scanLine[j] + prevScanLine[j]
			}
		case 3:
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
				a := int(prevScanLine[j])
				b := 0
				if j >= bpp {
					b = int(filteredData[x+j-bpp])
				}

				filteredData[x+j] = scanLine[j] + byte((a+b)>>1)
			}
		case 4:
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
			//	Paeth(x) + PaethPredictor(Raw(x-bpp), Prior(x), Prior(x-bpp)) (computed mod 256),
			//	where Raw() and Prior() refer to bytes already decoded.
			//	Exactly the same PaethPredictor() function is used by both encoder and decoder.
			for j := 0; j < len(scanLine); j++ {
				a := byte(0)
				if j >= bpp {
					a = filteredData[x+j-bpp]
				}

				b := prevScanLine[j]

				c := byte(0)
				if j >= bpp {
					c = prevScanLine[j-bpp]
				}

				filteredData[x+j] = scanLine[j] + byte(paeth(int(a), int(b), int(c)))
			}
		default:
			logError(i, filterType, "unexpected filter type")
			return filteredData, fmt.Errorf("unexpected filter type %d", i)
			//return fmt.Errorf("unexpected filter type %d %s %s", i, printBits(filterType), printBytes(scanLine))
		}
	}

	return filteredData, nil
}

func paeth(a, b, c int) int {
	p := a + b - c // extend the gradient
	// return whatever input is closest to p
	pa := abs(p - a)
	pb := abs(p - b)
	pc := abs(p - c)
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

func (reader *pngReader) GetRawImageData() []byte {
	return reader.imgData
}

func GetNew(pathToImage string) (PNGReader, error) {
	reader := &pngReader{}
	reader.bpp = -1
	reader.imgData = make([]byte, 0)
	pngFile, err := os.Open(pathToImage)
	if err != nil {
		logError(err)

		return nil, err
	}

	reader.f = pngFile

	r := bitreader.GetNewFileBitReader(pngFile)
	//header, err := r.GetBytes(8)
	_, err = r.GetBytes(8)
	if err != nil {
		logError(err)
		return nil, err
	}
	reader.width = 0                   //  4 bytes
	reader.height = 0                  //  4 bytes
	reader.bitDepth = byte(0)          //  1 byte
	reader.colorType = byte(0)         //  1 byte
	reader.compressionMethod = byte(0) //  1 byte
	reader.filterMethod = byte(0)      //  1 byte
	reader.interlaceMethod = byte(0)   //  1 byte
	//log(printBytes(header))
	idatChunks := make([]byte, 0)
	for {
		//log("==================================================================================================")
		chunkLenBytes, err := r.GetBytes(4)
		if err != nil {
			logError(err)

			return nil, err
		}
		// it's in big endian
		// so

		chunkLen := int(chunkLenBytes[3]) + int(uint32(chunkLenBytes[2])<<8) + int(uint32(chunkLenBytes[1])<<16) + int(uint32(chunkLenBytes[0])<<24)
		//log("chunk length:", chunkLen, printBytes(chunkLenBytes))
		chunkTypeBytes, err := r.GetBytes(4)
		if err != nil {
			logError(err)

			return nil, err
		}
		/**
		Ancillary bit: bit 5 of first byte
		0 (uppercase) = critical, 1 (lowercase) = ancillary.
		Chunks that are not strictly necessary in order to meaningfully display the contents of the file are known as "ancillary" chunks. A decoder encountering an unknown chunk in which the ancillary bit is 1 can safely ignore the chunk and proceed to display the image. The time chunk (tIME) is an example of an ancillary chunk.

		Chunks that are necessary for successful display of the file's contents are called "critical" chunks. A decoder encountering an unknown chunk in which the ancillary bit is 0 must indicate to the user that the image contains information it cannot safely interpret. The image header chunk (IHDR) is an example of a critical chunk.
		*/
		//log("Ancillary bit:", r.GetNthBitInByte(chunkTypeBytes[0], 5))
		/**
		Private bit: bit 5 of second byte
		0 (uppercase) = public, 1 (lowercase) = private.
		A public chunk is one that is part of the PNG specification or is registered in the list of PNG special-purpose public chunk types. Applications can also define private (unregistered) chunks for their own purposes. The names of private chunks must have a lowercase second letter, while public chunks will always be assigned names with uppercase second letters. Note that decoders do not need to test the private-chunk property bit, since it has no functional significance; it is simply an administrative convenience to ensure that public and private chunk names will not conflict. See Additional chunk types, and Recommendations for Encoders: Use of private chunks.
		*/
		//log("Private bit:", r.GetNthBitInByte(chunkTypeBytes[1], 5))
		/**
		Reserved bit: bit 5 of third byte
		Must be 0 (uppercase) in files conforming to this version of PNG.
		The significance of the case of the third letter of the chunk name is reserved for possible future expansion. At the present time all chunk names must have uppercase third letters. (Decoders should not complain about a lowercase third letter, however, as some future version of the PNG specification could define a meaning for this bit. It is sufficient to treat a chunk with a lowercase third letter in the same way as any other unknown chunk type.)
		*/
		//log("Reserved bit:", r.GetNthBitInByte(chunkTypeBytes[2], 5))

		/*
			Safe-to-copy bit: bit 5 of fourth byte
			0 (uppercase) = unsafe to copy, 1 (lowercase) = safe to copy.
				This property bit is not of interest to pure decoders, but it is needed by PNG editors (programs that modify PNG files).
				This bit defines the proper handling of unrecognized chunks in a file that is being modified.
				If a chunk's safe-to-copy bit is 1, the chunk may be copied to a modified PNG file whether
					or not the software recognizes the chunk type, and regardless of the extent of the file modifications.
				If a chunk's safe-to-copy bit is 0, it indicates that the chunk depends on the image data.
					If the program has made any changes to critical chunks, including addition, modification, deletion,
					or reordering of critical chunks, then unrecognized unsafe chunks must not be copied to the output PNG file.
					(Of course, if the program does recognize the chunk, it can choose to output an appropriately modified version.)

				A PNG editor is always allowed to copy all unrecognized chunks if it has only added, deleted, modified,
					or reordered ancillary chunks. This implies that it is not permissible for ancillary chunks to depend on other ancillary chunks.

				PNG editors that do not recognize a critical chunk must report an error and refuse to process that PNG file at all.
				The safe/unsafe mechanism is intended for use with ancillary chunks. The safe-to-copy bit will always be 0 for critical chunks.

				Rules for PNG editors are discussed further in Chunk Ordering Rules.
		*/
		//log("Safe-to-copy bit:", r.GetNthBitInByte(chunkTypeBytes[3], 5))

		switch string(rune(chunkTypeBytes[0])) + string(rune(chunkTypeBytes[1])) + string(rune(chunkTypeBytes[2])) + string(rune(chunkTypeBytes[3])) {
		case "IHDR":
			if err := reader.handleIHDRChunk(r, chunkLen); err != nil {
				logError(err)
				return nil, err
			}
			continue
		case "IEND":
			//log("it's IEND")
			if err := reader.handleIDATChunks(idatChunks); err != nil {
				logError(err)
				return nil, err
			}
			return reader, nil
		case "PLTE":
			log("it's PLTE")
			if data, err := r.GetBytes(chunkLen + 4); err != nil {
				logError(err)
				return nil, err
			} else {
				// The PLTE chunk contains from 1 to 256 palette entries, each a three-byte series of the form:

				//Red:   1 byte (0 = black, 255 = red)
				//Green: 1 byte (0 = black, 255 = green)
				//Blue:  1 byte (0 = black, 255 = blue)
				br := bitreader.GetNewSliceBitReader(data)
				reader.palette = make([][]byte, 0)
				for i := 0; i < chunkLen/3; i++ {
					if len(data) == 0 {
						break
					}
					r, err := br.GetByte()
					if err != nil {
						logError(err)
						return nil, err
					}
					g, err := br.GetByte()
					if err != nil {
						logError(err)
						return nil, err
					}
					b, err := br.GetByte()
					if err != nil {
						logError(err)
						return nil, err
					}
					reader.palette = append(reader.palette, []byte{r, g, b})
				}
			}
			continue
		case "IDAT":
			//log("it's IDAT")
			if data, err := r.GetBytes(chunkLen + 4); err != nil {
				return nil, err
			} else {
				idatChunks = append(idatChunks, data[:chunkLen]...)
				continue
			}
		default:
			//log("it's", string(rune(chunkTypeBytes[0]))+string(rune(chunkTypeBytes[1]))+string(rune(chunkTypeBytes[2]))+string(rune(chunkTypeBytes[3])))
		}
		if r.GetNthBitInByte(chunkTypeBytes[0], 5) == 1 {
			//log("skipping it's an optional chunk")
		}
		if _, err := r.GetBytes(chunkLen + 4); err != nil {
			return nil, err
		}
	}
}

func (reader *pngReader) handleIDATChunks(idatChunks []byte) error {
	cr := bitreader.GetNewSliceBitReader(idatChunks)
	// There can be multiple IDAT chunks;
	// if so, they must appear consecutively with no other intervening chunks.
	// The compressed datastream is then the concatenation of the contents of all the IDAT chunks.
	// The encoder can divide the compressed datastream into IDAT chunks however it wishes.
	// (Multiple IDAT chunks are allowed so that encoders can work in a fixed amount of memory;
	// typically the chunk size will correspond to the encoder's buffer size.)
	// It is important to emphasize that IDAT chunk boundaries have no semantic significance
	// and can occur at any point in the compressed datastream.
	// A PNG file in which each IDAT chunk contains only one data byte is valid,
	// though remarkably wasteful of space. (For that matter, zero-length IDAT chunks are valid, though even more wasteful.)

	// 1 byte: zlib compression method (named cmf).
	// 1 byte: zlib extra flags.
	//  	n byte: actual compressed data.
	// 2 bytes: check value, algorithm used is called Adler32.
	// CMF (Compression Method and flags)
	// This byte is divided into a 4-bit compression method and a 4- bit information field depending on the compression method.
	//
	// bits 0 to 3  CM     Compression method
	// bits 4 to 7  CINFO  Compression info
	//
	// CM (Compression method)
	// This identifies the compression method used in the file. CM = 8
	// denotes the "deflate" compression method with a window size up
	// to 32K.  This is the method used by gzip and PNG (see
	// references [1] and [2] in Chapter 3, below, for the reference
	// documents).  CM = 15 is reserved.  It might be used in a future
	// version of this specification to indicate the presence of an
	// extra field before the compressed data.
	//
	//	CINFO (Compression info)
	// For CM = 8, CINFO is the base-2 logarithm of the LZ77 window
	// size, minus eight (CINFO=7 indicates a 32K window size). Values
	// of CINFO above 7 are not allowed in this version of the
	// specification.  CINFO is not defined in this specification for
	//	CM not equal to 8.
	//cmf, err := cr.GetBits(8)
	_, err := cr.GetBits(8)
	if err != nil {
		logError(err)
		return err
	}
	//log("CM:", cmf[7], cmf[6], cmf[5], cmf[4], "=", cmf[4]+cmf[5]*2+cmf[6]*2*2+cmf[7]*2*2*2)
	//cinfo := []byte{cmf[3], cmf[2], cmf[1], cmf[0]}

	//log("CINFO:", cinfo[0], cinfo[1], cinfo[2], cinfo[3], "=", cinfo[3]+cinfo[2]*2+cinfo[1]*2*2+cinfo[0]*2*2*2)
	// FLG (FLaGs)
	// This flag byte is divided as follows:
	//
	// bits 0 to 4  FCHECK  (check bits for CMF and FLG)
	// bit  5       FDICT   (preset dictionary)
	// bits 6 to 7  FLEVEL  (compression level)
	//
	// The FCHECK value must be such that CMF and FLG, when viewed as a 16-bit unsigned integer stored in MSB order (CMF*256 + FLG),
	//
	// FDICT (Preset dictionary)
	// If FDICT is set, a DICT dictionary identifier is present
	// immediately after the FLG byte. The dictionary is a sequence of
	// bytes which are initially fed to the compressor without
	// producing any compressed output. DICT is the Adler-32 checksum
	// of this sequence of bytes (see the definition of ADLER32
	// below).  The decompressor can use this identifier to determine
	// which dictionary has been used by the compressor.
	//
	// FLEVEL (Compression level)
	// These flags are available for use by specific compression
	// methods.  The "deflate" method (CM = 8) sets these flags as follows:
	//
	// 0 - compressor used fastest algorithm
	// 1 - compressor used fast algorithm
	// 2 - compressor used default algorithm
	// 3 - compressor used maximum compression, slowest algorithm
	//
	// The information in FLEVEL is not needed for decompression; it
	// is there to indicate if recompression might be worthwhile.
	_, err = cr.GetBits(5)
	//fcheck, err := cr.GetBits(5)
	if err != nil {
		logError(err)
		return err
	}
	fdict, err := cr.GetBits(1)
	if err != nil {
		logError(err)
		return err
	}
	if fdict[0] == 1 {
		if _, err := cr.GetBytes(4); err != nil {
			logError(err)
			return err
		}
	}
	//flevel, err := cr.GetBits(2)
	_, err = cr.GetBits(2)
	if err != nil {
		logError(err)
		return err
	}
	//log("FCHECK", fcheck[0], fcheck[1], fcheck[2], fcheck[3], fcheck[4])
	//log("FDICT", fdict[0])
	//log("FLEVEL", flevel[0], flevel[1])
	for {
		if !cr.HasMoreData() {
			break
		}
		headerBits, err := cr.GetBits(3)
		if err != nil {
			logError(err)
			return err
		}
		//Each block of compressed data begins with 3 header bits containing the following data:
		//first bit
		//next 2 bits
		// BFINAL - BFINAL is set if and only if this is the last block of the data set.
		// BTYPE - specifies how the data are compressed, as follows:
		//				00 - no compression
		//				01 - compressed with fixed Huffman codes
		//				10 - compressed with dynamic Huffman codes
		//				11 - reserved (error)
		//log("HEADER", headerBits[0], headerBits[1], headerBits[2])
		switch fmt.Sprint(headerBits[2]) + fmt.Sprint(headerBits[1]) {
		case "00": // - no compression
			if err := reader.handleNoCompression(cr); err != nil {
				logError(err)
				return err
			}
		case "01": // - compressed with fixed Huffman codes
			if err := reader.handleFixedHuffman(cr); err != nil {
				logError(err)
				return err
			}
		case "10": // - compressed with dynamic Huffman codes
			if err := reader.handleDynamicHuffman(cr); err != nil {
				logError(err)
				return err
			}
		case "11": // - reserved (error)
			return fmt.Errorf("unexpected input")
		}
		if headerBits[0] == 1 {
			break
		}
	}

	return nil
}

func (reader *pngReader) handleIHDRChunk(r bitreader.BitReader, chunkLen int) error {
	//log("it's IHDR")
	/*
		The IHDR chunk must appear FIRST. It contains:

		   Width:              4 bytes
		   Height:             4 bytes
		   Bit depth:          1 byte
		   Color type:         1 byte
		   Compression method: 1 byte
		   Filter method:      1 byte
		   Interlace method:   1 byte
	*/
	if data, err := r.GetBytes(chunkLen + 4); err != nil {
		return err
	} else {
		reader.width = int(data[3]) + int(uint32(data[2])<<8) + int(uint32(data[1])<<16) + int(uint32(data[0])<<24)
		reader.height = int(data[7]) + int(uint32(data[6])<<8) + int(uint32(data[5])<<16) + int(uint32(data[4])<<24)
		reader.bitDepth = data[8]
		reader.colorType = data[9]
		reader.compressionMethod = data[10]
		reader.filterMethod = data[11]
		reader.interlaceMethod = data[12]
		//log("width:", reader.width)
		//log("height:", reader.height)
		// Bit depth is a single-byte integer giving the number of bits per sample or per palette index (not per pixel).
		// Valid values are 1, 2, 4, 8, and 16, although not all values are allowed for all color types.
		//log("bit depth:", reader.bitDepth)
		// Color type is a single-byte integer that describes the interpretation of the image data.
		// Color type codes represent sums of the following values:
		// 1 (palette used),
		// 2 (color used),
		// and 4 (alpha channel used).
		// Valid values are 0, 2, 3, 4, and 6.
		//log("color type:", reader.colorType)
		//  Color    Allowed    Interpretation
		//   Type    Bit Depths
		//
		//   0       1,2,4,8,16  Each pixel is a grayscale sample.
		//
		//   2       8,16        Each pixel is an R,G,B triple.
		//
		//   3       1,2,4,8     Each pixel is a palette index;
		//                       a PLTE chunk must appear.
		//
		//   4       8,16        Each pixel is a grayscale sample,
		//                       followed by an alpha sample.
		//
		//   6       8,16        Each pixel is an R,G,B triple,
		//                       followed by an alpha sample.
		//log("compression method:", reader.compressionMethod)
		// Compression method is a single-byte integer that indicates the method used to compress the image data.
		// At present, only compression method 0 (deflate/inflate compression with a sliding window of at most 32768 bytes) is defined.
		// All standard PNG images must be compressed with this scheme.
		// The compression method field is provided for possible future expansion or proprietary variants.
		// Decoders must check this byte and report an error if it holds an unrecognized code.
		// See Deflate/Inflate Compression for details.

		//log("filter method:", reader.filterMethod)
		// Filter method is a single-byte integer that indicates the preprocessing method
		// applied to the image data before compression.
		// At present, only filter method 0 (adaptive filtering with five basic filter types) is defined.
		// As with the compression method field, decoders must check this byte and report an error
		// if it holds an unrecognized code. See Filter Algorithms for details.
		//log("interlace method:", reader.interlaceMethod)
		// Interlace method is a single-byte integer that indicates the transmission order of the image data.
		// Two values are currently defined: 0 (no interlace) or 1 (Adam7 interlace).
		// See Interlaced data order for details.
	}
	return nil
}

func (reader *pngReader) handleDynamicHuffman(cr bitreader.BitReader) error {
	// finally compressed data follows
	// The first section of a GZIP-compressed block is three integers indicating the number of length codes, the number of literal codes, and the number of distance codes.
	nLit, err := cr.GetBits(5)
	if err != nil {
		logError(err)
		return err
	}
	hlit := int(nLit[0]+nLit[1]*2+nLit[2]*2*2+nLit[3]*2*2*2+nLit[4]*2*2*2*2) + 257
	//log("the number of literal codes:", nLit[0], nLit[1], nLit[2], nLit[3], nLit[4], "=", hlit)

	nDist, err := cr.GetBits(5)
	if err != nil {
		logError(err)
		return err
	}
	hdist := int(nDist[0]+nDist[1]*2+nDist[2]*2*2+nDist[3]*2*2*2+nDist[4]*2*2*2*2) + 1
	//log("the number of distance codes:", nDist[0], nDist[1], nDist[2], nDist[3], nDist[4], "=", hdist)
	nLength, err := cr.GetBits(4)
	if err != nil {
		logError(err)
		return err
	}
	nLen := int(nLength[0] + nLength[1]*2 + nLength[2]*2*2 + nLength[3]*2*2*2)
	//log("the number of length codes:", nLength[0], nLength[1], nLength[2], nLength[3], "=", nLen)
	lengthCodes, err := cr.GetBits((nLen + 4) * 3)
	if err != nil {
		logError(err)
		return err
	}
	lenDictionary := []int{16, 17, 18, 0, 8, 7, 9, 6, 10, 5, 11, 4, 12, 3, 13, 2, 14, 1, 15}
	lenDictionaryMap := make(map[int]int, 0)
	j := 0
	lenHCodes := make(map[int]int, 0)
	for i := 0; i < len(lengthCodes); i += 3 {
		l := int(lengthCodes[i] + lengthCodes[i+1]*2 + lengthCodes[i+2]*2*2)
		lenDictionaryMap[lenDictionary[j]] = l
		if _, exists := lenHCodes[l]; !exists {
			lenHCodes[l] = 0
		}
		if l > 0 {
			lenHCodes[l]++
		}
		j++
	}

	canonicalHuffmanCodingMapForLengthsTree := getCanonicalHuffmanCodingMapForLengthsTree(lenHCodes, lenDictionaryMap)
	//log("canonicalHuffmanCodingMapForLengthsTree:", canonicalHuffmanCodingMapForLengthsTree)
	// we have our tree that has compressed the bit lengths for the other two trees,
	// as you see this tree doesn't have a stop code, so how many symbols must we decode?
	// well we have a total of ( hlit + hdist ) codes we will decode using this tree,
	// after that we can build the two huffman trees and we can decompress the original data.
	// they come one after another, the first one is hlit elements, and the other one is hdist elements

	// the first 255 characters (symbols/codes) actually still
	// map to the 255 literals of the ASCII table,
	// the 256 symbol/code means "stop token",
	// 257 - 285 indicate "length codes/duplication length",
	// immediately following one of these length codes comes a distance code,
	// see that is the clever part, if you hit a 257-285 code,
	// you know anything after that indicates "distance" code,
	// but wait there is more Cleverness and shenanigans with this algorithm,
	// you see how 285 - 257 = 28 length codes?
	// Well instead of 257 being equal of a length of (1) and 285 being a length of (28),
	// they made those numbers index into a table that indicates whether or not to read more bits
	// for a more flexible number of length codes, i.e: if you get a code of 257 this actually means
	// the length to duplicate is 3, 258 means 4, until 264 means 10 bytes of length to duplicate,
	// but 265 means read one more bit from the stream,
	// if the bit is 0 you have 11 bytes of length,
	// but if its 1 then you have 12 bytes of length,
	// and so on for the other codes after 265, here is the complete table from the standard:
	//
	//     Extra                 Extra               Extra
	//Code Bits Length(s)   Code Bits Lengths   Code Bits Length(s)
	//---- ---- ------      ---- ---- -------    ---- ---- -------
	//257   0      3        267   1   15,16     277   4   67-82
	//258   0      4        268   1   17,18     278   4   83-98
	//259   0      5        269   2   19-22     279   4   99-114
	//260   0      6        270   2   23-26     280   4  115-130
	//261   0      7        271   2   27-30     281   5  131-162
	//262   0      8        272   2   31-34     282   5  163-194
	//263   0      9        273   3   35-42     283   5  195-226
	//264   0     10        274   3   43-50     284   5  227-257
	//265   1   11,12       275   3   51-58     285   0    258
	//266   1   13,14       276   3   59-66
	// we have Huffman for 0-18 OK, so
	// codes 0 to 15: they represent literal code bit lengths,
	// i.e: if you see number 15 it means the code bit length is 15 bits long.
	// 15 is the maximum because that is the longest a code bit length can be in when making the "two trees"
	//code 16:
	//	means repeat the previous code 3 to 6 times depending on the next 2 bits.
	//	So if you get number 16, you read 2 more bits and interpret them as an integer (2 bits only mean decimal 0 to 3)
	//	then add them to the number 3.
	//code 17:
	//	means repeat 0 for 3 - 10 times depending on the next 3 bits.
	//	if you see 17, read the next 3 bits and add the integer to the number 3 and you get repeat count.
	//code 18:
	//	means repeat 0 for 11 - 138 times depending on the next 7 bits.
	//	if you see 18, read the next 7 bits and add the integer to the number 11 and you get repeat count.

	uncompressedLitTreeData, err := decodeTree(hlit, cr, canonicalHuffmanCodingMapForLengthsTree)
	if err != nil {
		logError(err)
		return err
	}
	litDictionary := make(map[int]int, 0)
	litFreq := make(map[int]int, 0)
	for i := 0; i < hlit; i++ {
		if _, exists := litDictionary[int(uncompressedLitTreeData[i])]; !exists {
			litDictionary[int(uncompressedLitTreeData[i])] = 0
		}
		litDictionary[int(uncompressedLitTreeData[i])] += 1
		litFreq[i] = int(uncompressedLitTreeData[i])
	}
	uncompressedDistTreeData, err := decodeTree(hdist, cr, canonicalHuffmanCodingMapForLengthsTree)
	if err != nil {
		logError(err)
		return err
	}
	distDictionary := make(map[int]int, 0)
	distFreq := make(map[int]int, 0)
	for i := 0; i < hdist; i++ {
		if _, exists := distDictionary[int(uncompressedDistTreeData[i])]; !exists {
			distDictionary[int(uncompressedDistTreeData[i])] = 0
		}
		distDictionary[int(uncompressedDistTreeData[i])] += 1
		distFreq[i] = int(uncompressedDistTreeData[i])
	}
	canonicalHuffmanCodingMapForLit := getCanonicalHuffmanCodingMap(litDictionary)
	//log("canonicalHuffmanCodingMapForLit:", canonicalHuffmanCodingMapForLit)
	canonicalHuffmanCodingMapForDist := getCanonicalHuffmanCodingMap(distDictionary)
	litTree := buildNaturalOrderTree(canonicalHuffmanCodingMapForLit, litFreq)
	//log("canonicalHuffmanCodingMapForDist:", canonicalHuffmanCodingMapForDist)
	distTree := buildNaturalOrderTree(canonicalHuffmanCodingMapForDist, distFreq)
	//log("-------------------")
	//log("lit tree:", litTree)
	//log("-------------------")
	//log("dist tree:", distTree)
	//log("-------------------")
	//log("starting to decompress")
	_, err = reader.decodeLZ77(cr, litTree, distTree)
	//d, err := reader.decodeLZ77(cr, litTree, distTree)
	//log("decoded LZ77", len(d))
	if err != nil {
		logError(err)
		return err
	}

	return nil
}

func (reader *pngReader) handleFixedHuffman(cr bitreader.BitReader) error {
	log("FIXED HUFFMAN")
	// The Huffman codes for the two alphabets are fixed, and are not represented explicitly in the data. The Huff-
	// man code lengths for the literal/length alphabet are:

	// Lit Value             Bits            Codes
	// 0-143                 8               00110000  - 10111111
	// 144 - 255             9               110010000 - 111111111
	// 256 - 279             7               0000000   - 0010111
	// 280 - 287             8               11000000  - 11000111

	// Distance codes 0-31 are represented by (fixed-length) 5-bit codes, with possible additional bits as shown in
	//the table shown in Paragraph 3.2.5, above. Note that distance codes 30-31 will never actually occur in the
	//compressed data.

	litTree := make(map[string]int, 0)
	// 00000
	code := 0b00110000
	for i := 0; i < 144; i++ {
		litTree[fmt.Sprintf("%08b", code)] = i
		code += 1
	}
	code = 0b110010000
	for i := 144; i < 256; i++ {
		litTree[fmt.Sprintf("%09b", code)] = i
		code += 1
	}
	code = 0b0000000
	for i := 256; i < 280; i++ {
		litTree[fmt.Sprintf("%07b", code)] = i
		code += 1
	}
	code = 0b11000000
	for i := 280; i < 287; i++ {
		litTree[fmt.Sprintf("%08b", code)] = i
		code += 1
	}
	code = 0b11000000
	for i := 280; i < 287; i++ {
		litTree[fmt.Sprintf("%08b", code)] = i
		code += 1
	}
	distTree := make(map[string]int, 0)
	code = 0b00000
	for i := 0; i < 32; i++ {
		distTree[fmt.Sprintf("%05b", code)] = i
		code += 1
	}
	d, err := reader.decodeLZ77(cr, litTree, distTree)
	if err != nil {
		logError(err)
		return err
	}
	log("decoded LZ77", len(d))

	reader.imgData = append(reader.imgData, d...)

	return nil
}

func (reader *pngReader) handleNoCompression(cr bitreader.BitReader) error {
	// Any bits of input up to the next byte boundary are ignored. The rest of the block consists of the following
	// information:
	// 0  1   2   3   4...
	//+---+---+---+---+================================+
	//| LEN   | NLEN  |... LEN bytes of literal data...|
	//+---+---+---+---+================================+
	// LEN is the number of data bytes in the block. NLEN is the one’s complement of LEN.

	// so we need to skip any remaining bits in current partially
	// processed byte
	// read LEN and NLEN
	// copy LEN bytes of data to output
	//
	// 00101110 10110111
	// 00110001 11110011
	if err := cr.GoToNextByte(); err != nil {
		logError(err)
		return err
	}

	lbytes, err := cr.GetBytes(4)
	if err != nil {
		logError(err)
		return err
	}

	l := int(lbytes[0]) + int(lbytes[1])<<8
	cl := int(lbytes[2]) + int(lbytes[3])<<8

	if uint16(cl) != uint16(^l) {
		return fmt.Errorf("not complement")
	}
	log("l:", l)
	log("cl:", cl)
	if l == 0 {
		return nil
	}
	d, err := cr.GetBytes(l)
	if err != nil {
		logError(err)
	}
	log("read", len(d), "bytes of uncompressed data")
	reader.imgData = append(reader.imgData, d...)

	return nil
}

func buildNaturalOrderTree(codingMap map[int][]string, freq map[int]int) map[string]int {
	result := make(map[string]int, 0)
	for l := 0; l < len(freq); l++ {
		n := freq[l]
		if _, exists := codingMap[n]; !exists {
			continue
		}
		var bitSeq string
		bitSeq, codingMap[n] = codingMap[n][0], codingMap[n][1:]
		result[bitSeq] = l
	}

	return result
}

func (reader *pngReader) decodeLZ77(cr bitreader.BitReader, litTree, distTree map[string]int) ([]byte, error) {
	uncompressedTreeData := make([]byte, 0)
	if len(reader.imgData) != 0 {
		for _, b := range reader.imgData {
			uncompressedTreeData = append(uncompressedTreeData, b)
		}
	}
	key := ""
	for {
		b, err := cr.GetBit()
		if err != nil {
			logError(err)
			return uncompressedTreeData, err
		}
		key += fmt.Sprint(b)
		if val, exists := litTree[key]; !exists {
			continue
		} else {
			//log("the val is", val, key)
			key = ""
			if val < 256 {
				uncompressedTreeData = append(uncompressedTreeData, byte(val))
				//uncompressedTreeData = append([]byte{byte(val)}, uncompressedTreeData...)
				continue
			}
			if val == 256 { // STOP
				//log("STOP", cr.HasMoreData())
				break
			}
			d, l, err := getDistanceLength(cr, val, distTree)
			if err != nil {
				logError(err, val)
				return uncompressedTreeData, err
			}
			idx := len(uncompressedTreeData) - d
			startIdx := idx
			for i := 0; i < l; i++ {
				if idx >= len(uncompressedTreeData) {
					idx = startIdx
				}
				uncompressedTreeData = append(uncompressedTreeData, uncompressedTreeData[idx])
				idx++
			}
		}
	}

	reader.imgData = uncompressedTreeData

	return uncompressedTreeData, nil
}

func (reader *pngReader) filterInterlacedImage(rowLen int, bpp int) ([]byte, error) {
	startingRow := []int{0, 0, 4, 0, 2, 0, 1}
	startingCol := []int{0, 4, 0, 2, 0, 1, 0}
	rowIncrement := []int{8, 8, 8, 4, 4, 2, 2}
	colIncrement := []int{8, 8, 4, 4, 2, 2, 1}
	width := reader.width
	height := reader.height
	filteredData := make([]byte, rowLen*height*bpp)
	data := reader.GetRawImageData()
	log("decoded", len(data))
	br := bitreader.GetNewSliceBitReader(data)
	log("row length is", rowLen)
	for pass := 0; pass < 7; pass++ {
		rowsNum := (height - startingRow[pass] + rowIncrement[pass] - 1) / rowIncrement[pass]
		if rowsNum < 0 {
			rowsNum = 0
		}

		colsNum := (width - startingCol[pass] + colIncrement[pass] - 1) / colIncrement[pass]
		if colsNum < 0 {
			colsNum = 0
		}

		//log("number of rows:", rowsNum, "number of columns:", colsNum)
		//FILTER ... colsNum
		//.
		//. rowsNum
		//.
		//FILTER ...
		passFilteredData := make([]byte, rowsNum*colsNum*bpp)
		for j := 0; j < rowsNum; j++ {
			if err := br.GoToNextByte(); err != nil {
				logError(err)
				break
			}
			filterType, err := br.GetBits(8)
			if err != nil {
				return filteredData, err
			}
			scanLine, err := br.GetBytes(colsNum * bpp)
			if err != nil {
				return filteredData, err
			}
			prevScanLine := make([]byte, colsNum*bpp)
			if j > 0 {
				prevScanLine = passFilteredData[(j-1)*colsNum*bpp : j*colsNum*bpp]
			}
			x := j * colsNum * bpp
			switch int(filterType[0] + filterType[1]*2 + filterType[2]*4) {
			case 0:
				for k, v := range scanLine {
					passFilteredData[x+k] = v
				}
			case 1:
				for i := 0; i < len(scanLine); i++ {
					el := scanLine[i]
					if i >= bpp {
						el += passFilteredData[x+i-bpp]
					}
					passFilteredData[x+i] = el
				}
			case 2:
				for i := 0; i < len(scanLine); i++ {
					passFilteredData[x+i] = scanLine[i] + prevScanLine[i]
				}
			case 3:
				for i := 0; i < len(scanLine); i++ {
					a := int(prevScanLine[i])
					b := 0
					if i >= bpp {
						b = int(passFilteredData[x+i-bpp])
					}

					passFilteredData[x+i] = scanLine[i] + byte((a+b)>>1)
				}
			case 4:
				for i := 0; i < len(scanLine); i++ {
					a := byte(0)
					if i >= bpp {
						a = passFilteredData[x+i-bpp]
					}

					b := prevScanLine[i]

					c := byte(0)
					if i >= bpp {
						c = prevScanLine[i-bpp]
					}

					passFilteredData[x+i] = scanLine[i] + byte(paeth(int(a), int(b), int(c)))
				}
			default:
				logError(j, filterType, "unexpected filter type")
				return filteredData, fmt.Errorf("unexpected filter type %d", j)
				//return fmt.Errorf("unexpected filter type %d %s %s", i, printBits(filterType), printBytes(scanLine))
			}
		}
		row := startingRow[pass]
		j := 0
		for {
			if row >= height {
				break
			}

			i := 0
			col := startingCol[pass]
			for {
				if col >= width {
					break
				}
				for k := 0; k < bpp; k++ {
					filteredData[row*rowLen+col*bpp+k] = passFilteredData[j*colsNum*bpp+i*bpp+k]
				}
				col += colIncrement[pass]
				i++
			}
			row += rowIncrement[pass]
			j++
		}
	}

	return filteredData, nil
}

func getDistanceLength(cr bitreader.BitReader, val int, distTree map[string]int) (int, int, error) {
	startVal := 0
	numExtraBits := 0
	l := 0
	var err error
	switch val {
	case 257:
		l = 3
	case 258:
		l = 4
	case 259:
		l = 5
	case 260:
		l = 6
	case 261:
		l = 7
	case 262:
		l = 8
	case 263:
		l = 9
	case 264:
		l = 10
	case 265:
		startVal = 11
		numExtraBits = 1
	case 266:
		startVal = 13
		numExtraBits = 1
	case 267:
		startVal = 15
		numExtraBits = 1
	case 268:
		startVal = 17
		numExtraBits = 1
	case 269:
		startVal = 19
		numExtraBits = 2
	case 270:
		startVal = 23
		numExtraBits = 2
	case 271:
		startVal = 27
		numExtraBits = 2
	case 272:
		startVal = 31
		numExtraBits = 2
	case 273:
		startVal = 35
		numExtraBits = 3
	case 274:
		startVal = 43
		numExtraBits = 3
	case 275:
		startVal = 51
		numExtraBits = 3
	case 276:
		startVal = 59
		numExtraBits = 3
	case 277:
		startVal = 67
		numExtraBits = 4
	case 278:
		startVal = 83
		numExtraBits = 4
	case 279:
		startVal = 99
		numExtraBits = 4
	case 280:
		startVal = 115
		numExtraBits = 4
	case 281:
		startVal = 131
		numExtraBits = 5
	case 282:
		startVal = 163
		numExtraBits = 5
	case 283:
		startVal = 195
		numExtraBits = 5
	case 284:
		startVal = 227
		numExtraBits = 5
	case 285:
		l = 258
	default:
		return 0, 0, fmt.Errorf("bad input")
	}

	if l < 1 {
		l, err = __getValWithExtraBits(cr, startVal, numExtraBits)
		if err != nil {
			logError(err, val, startVal, numExtraBits)
			return 0, 0, err
		}
	}
	d, err := getDistance(cr, distTree)
	if err != nil {
		logError(err, val, startVal, numExtraBits)
		return 0, 0, err
	}

	return d, l, nil
}

func getDistance(cr bitreader.BitReader, tree map[string]int) (int, error) {
	maxLen := 0
	for s, _ := range tree {
		if len(s) > maxLen {
			maxLen = len(s)
		}
	}
	key := ""
	for {
		if len(key) > maxLen {
			log(key, len(key), maxLen)
			break
		}
		b, err := cr.GetBit()
		if err != nil {
			logError(err)

			logError(err)
			return 0, err
		}
		key += fmt.Sprint(b)
		startVal := 0
		numExtraBits := 0
		if val, exists := tree[key]; !exists {
			continue
		} else {
			switch val {
			case 0:
				return 1, nil
			case 1:
				return 2, nil
			case 2:
				return 3, nil
			case 3:
				return 4, nil
			case 4:
				startVal = 5
				numExtraBits = 1
			case 5:
				startVal = 7
				numExtraBits = 1
			case 6:
				startVal = 9
				numExtraBits = 2
			case 7:
				startVal = 13
				numExtraBits = 2
			case 8:
				startVal = 17
				numExtraBits = 3
			case 9:
				startVal = 25
				numExtraBits = 3
			case 10:
				startVal = 33
				numExtraBits = 4
			case 11:
				startVal = 49
				numExtraBits = 4
			case 12:
				startVal = 65
				numExtraBits = 5
			case 13:
				startVal = 97
				numExtraBits = 5
			case 14:
				startVal = 129
				numExtraBits = 6
			case 15:
				startVal = 193
				numExtraBits = 6
			case 16:
				startVal = 257
				numExtraBits = 7
			case 17:
				startVal = 385
				numExtraBits = 7
			case 18:
				startVal = 513
				numExtraBits = 8
			case 19:
				startVal = 769
				numExtraBits = 8
			case 20:
				startVal = 1025
				numExtraBits = 9
			case 21:
				startVal = 1537
				numExtraBits = 9
			case 22:
				startVal = 2049
				numExtraBits = 10
			case 23:
				startVal = 3073
				numExtraBits = 10
			case 24:
				startVal = 4097
				numExtraBits = 11
			case 25:
				startVal = 6145
				numExtraBits = 11
			case 26:
				startVal = 8193
				numExtraBits = 12
			case 27:
				startVal = 12289
				numExtraBits = 12
			case 28:
				startVal = 16385
				numExtraBits = 13
			case 29:
				startVal = 24577
				numExtraBits = 13
			}
			res, err := __getValWithExtraBits(cr, startVal, numExtraBits)
			if err != nil {
				logError(err, val, startVal, numExtraBits)
			}
			return res, err
		}

	}

	return 0, fmt.Errorf("unexpected input")
}

func __getValWithExtraBits(cr bitreader.BitReader, startVal int, numExtraBits int) (int, error) {
	bitsE, err := cr.GetBits(numExtraBits)
	if err != nil {
		logError(err, startVal, numExtraBits)
		return 0, err
	}

	return startVal + cr.BitsToInt(bitsE), nil
}

func decodeTree(dist int, cr bitreader.BitReader, canonicalHuffmanCodingMapForLengthsTree map[string]int) ([]byte, error) {
	maxLen := 0
	for s, _ := range canonicalHuffmanCodingMapForLengthsTree {
		if len(s) > maxLen {
			maxLen = len(s)
		}
	}
	uncompressedTreeData := make([]byte, 0)
	key := ""
	k := ""
	n := 0
	for {
		if len(key) > maxLen {
			log(key, len(key), maxLen)
			return nil, fmt.Errorf("bad input")
		}
		if n >= dist {
			break
		}
		b, err := cr.GetBit()
		if err != nil {
			logError(err)

			logError(err)
			return nil, err
		}
		key += fmt.Sprint(b)
		k += fmt.Sprint(b)
		if val, exists := canonicalHuffmanCodingMapForLengthsTree[key]; !exists {
			continue
		} else {
			key = ""
			if byte(val) < 16 {
				uncompressedTreeData = append(uncompressedTreeData, byte(val))
			}
			if byte(val) == 16 {
				d, err := __getValWithExtraBits(cr, 3, 2)
				if err != nil {
					logError(err)

					return nil, err
				}
				prevVal := byte(0)
				if n > 0 {
					prevVal = uncompressedTreeData[n-1]
				}

				for i := 0; i < d; i++ {
					uncompressedTreeData = append(uncompressedTreeData, prevVal)
					n++
				}
				continue
			}
			if byte(val) == 17 {
				d, err := __getValWithExtraBits(cr, 3, 3)
				if err != nil {
					logError(err)

					return nil, err
				}

				for i := 0; i < d; i++ {
					uncompressedTreeData = append(uncompressedTreeData, 0)
					n++
				}
				continue
			}
			if byte(val) == 18 {
				d, err := __getValWithExtraBits(cr, 11, 7)
				if err != nil {
					logError(err)

					return nil, err
				}

				for i := 0; i < d; i++ {
					uncompressedTreeData = append(uncompressedTreeData, 0)
					n++
				}
				continue
			}
			n++
		}
	}

	return uncompressedTreeData, nil
}

func getCanonicalHuffmanCodingMapForLengthsTree(lenHCodes map[int]int, lenDictionaryMap map[int]int) map[string]int {
	canonicalHuffmanCodingMap := getCanonicalHuffmanCodingMap(lenHCodes)
	//log(canonicalHuffmanCodingMap)
	result := make(map[string]int, 0)
	for l := 0; l <= 18; l++ {
		n := lenDictionaryMap[l]
		if _, exists := canonicalHuffmanCodingMap[n]; !exists {
			continue
		}
		var bitSeq string
		bitSeq, canonicalHuffmanCodingMap[n] = canonicalHuffmanCodingMap[n][0], canonicalHuffmanCodingMap[n][1:]
		result[bitSeq] = l
	}

	return result
}

func getCanonicalHuffmanCodingMap(dictionary map[int]int) map[int][]string {
	//log(dictionary)
	keysSorted := make([]int, 0)
	for k, _ := range dictionary {
		keysSorted = append(keysSorted, k)
	}
	sort.Ints(keysSorted)
	canonicalHuffmanCodingMap := make(map[int][]string, 0)
	sequence := "0"
	isStart := true
	i := 0
	for _, k := range keysSorted {
		if k == 0 {
			continue
		}
		// k represents the sequence
		// lenHCodes[k] represents the number of variations for that k of 0 and 1 sequences
		// ok
		number, _ := strconv.ParseUint(sequence, 2, 32)
		n := int(number) //Convert uint64 To int
		if !isStart {
			n += 1
			for i := 0; i < k-len(sequence); i++ {
				if getSequenceLength(n) < k {
					n = n << 1
				} else {
					break
				}
			}
		}
		sequence = fmt.Sprintf("%0"+fmt.Sprint(k)+"b", n)

		canonicalHuffmanCodingMap[k] = []string{sequence}
		for j := 1; j < dictionary[k]; j++ {
			nmbr, _ := strconv.ParseUint(canonicalHuffmanCodingMap[k][j-1], 2, 32)
			n := int(nmbr)
			sequence = fmt.Sprintf("%0"+fmt.Sprint(k)+"b", n+1)
			canonicalHuffmanCodingMap[k] = append(canonicalHuffmanCodingMap[k], sequence)
		}
		i++
		isStart = false
	}

	return canonicalHuffmanCodingMap
}

func getSequenceLength(sequence int) int {
	// we go thru bits of the sequence and increment, so we count only significant positions
	return len(fmt.Sprintf("%b", sequence))
}

func printBytes(bytes []byte) string {
	msg := ""
	for _, b := range bytes {
		msg += fmt.Sprintf("0x%02X ", b)
	}

	return msg
}

func bytesToInt(bytes []byte) int {
	result := 0
	for i := 0; i < len(bytes); i++ {
		if i > 3 {
			break
		}
		result = result << 8
		result += int(bytes[i])

	}

	return result
}

func reverseSliceByte(a []byte) []byte {
	s := make([]byte, len(a))
	copy(s, a)
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}

	return s
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
		logError(err)

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

func getDecompressedDataUsingZlib(data []byte) ([]byte, error) {
	b := bytes.NewReader(data)
	z, err := zlib.NewReader(b)
	if err != nil {
		return nil, err
	}
	defer z.Close()
	p, err := io.ReadAll(z)
	if err != nil {
		return nil, err
	}

	return p, nil
}
