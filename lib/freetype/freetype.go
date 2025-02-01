package freetype

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type font struct {
	Xmin                int16
	Ymin                int16
	Xmax                int16
	Ymax                int16
	idRangeOffsetsStart uint
	glyfTableStart      uint
	glyphIdArrStart     uint
	segCount            uint
	indexToLocFormat    int
	loca                []byte
	glyf                []byte
	startCodeArr        []uint
	endCodeArr          []uint
	idDelta             []uint
	idRangeOffset       []uint
	glyphIdArr          []uint
	ascent              int16 // The distance from the baseline to the highest or upper grid coordinate used to place an outline point.
	//                           It is a positive value, due to the grid's orientation with the y axis upwards.
	//
	descent int16 //             The distance from the baseline to the lowest grid coordinate used to place an outline point.
	//                           In FreeType, this is a negative value, due to the grid's orientation.
	//                           Note that in some font formats this is a positive value.
	//
	lineGap int16 //             The distance that must be placed between two lines of text.
	//                           The baseline-to-baseline distance should be computed as
	//
}

func LoadFont(pathToFont string) (Font, error) {
	fontFile, err := os.Open(pathToFont)
	if err != nil {
		return nil, err
	}

	defer fontFile.Close()

	//The file structure starts with a table called the font directory
	//and it must be the first table in the file,
	//it consists of 2 parts; the offset subtable and table directory,
	//offset subtable contains info on the amount of tables
	//	Here is the structure for the Offset subtable:
	//
	//	typedef struct {
	//		u32	scaler_type;
	//		u16	numTables; r
	//		u16	searchRange;
	//		u16	entrySelector;
	//		u16	rangeShift;
	//	} offset_subtable;
	subTableOffset := make([]byte, 4+2+2+2+2)
	nBytesRead, err := fontFile.Read(subTableOffset)
	if err != nil {
		return nil, err
	}
	if nBytesRead != 4+2+2+2+2 {
		return nil, fmt.Errorf("unexpected input")
	}
	if !areEqual(subTableOffset, []byte{0x74, 0x72, 0x75, 0x65}, 4) && !areEqual(subTableOffset, []byte{0x00, 0x01, 0x00, 0x00}, 4) {
		return nil, fmt.Errorf("unsupported format %s | %s", printBytes(subTableOffset[:4]), printBytes(subTableOffset))
	}

	//typedef struct {
	//	union {
	//	char tag_c[4];
	//	u32	tag;
	//};
	//	u32	checkSum;
	//	u32	offset;
	//	u32	length;
	//} table_directory;

	numTables := bytesToInt([]byte{subTableOffset[4], subTableOffset[5]})
	//log("numTables:", numTables)
	offsetToCMAP := 0
	lenCMAP := 0
	offsetToHEAD := 0
	lenHEAD := 0
	offsetToLOCA := 0
	lenLOCA := 0
	offsetToGLYF := 0
	lenGLYF := 0
	offsetToHHEA := 0
	lenHHEA := 0
	for i := 0; i < numTables; i++ {
		tableDir := make([]byte, 4+4+4+4)
		nBytesRead, err := fontFile.Read(tableDir)
		if err != nil {
			return nil, err
		}
		if nBytesRead != 4+4+4+4 {
			return nil, fmt.Errorf("unexpected input")
		}
		fmt.Println("tag #", i, ":", bytesToStr(tableDir[:4]))
		fmt.Println("    offset:", bytesToInt([]byte{tableDir[8], tableDir[9], tableDir[10], tableDir[11]}))
		fmt.Println("    length:", bytesToInt([]byte{tableDir[12], tableDir[13], tableDir[14], tableDir[15]}))
		if bytesToStr(tableDir[:4]) == "cmap" {
			offsetToCMAP = bytesToInt([]byte{tableDir[8], tableDir[9], tableDir[10], tableDir[11]})
			lenCMAP = bytesToInt([]byte{tableDir[12], tableDir[13], tableDir[14], tableDir[15]})
		}
		if bytesToStr(tableDir[:4]) == "head" {
			offsetToHEAD = bytesToInt([]byte{tableDir[8], tableDir[9], tableDir[10], tableDir[11]})
			lenHEAD = bytesToInt([]byte{tableDir[12], tableDir[13], tableDir[14], tableDir[15]})
		}

		if bytesToStr(tableDir[:4]) == "loca" {
			offsetToLOCA = bytesToInt([]byte{tableDir[8], tableDir[9], tableDir[10], tableDir[11]})
			lenLOCA = bytesToInt([]byte{tableDir[12], tableDir[13], tableDir[14], tableDir[15]})
		}

		if bytesToStr(tableDir[:4]) == "glyf" {
			offsetToGLYF = bytesToInt([]byte{tableDir[8], tableDir[9], tableDir[10], tableDir[11]})
			lenGLYF = bytesToInt([]byte{tableDir[12], tableDir[13], tableDir[14], tableDir[15]})
		}

		if bytesToStr(tableDir[:4]) == "hhea" {
			offsetToHHEA = bytesToInt([]byte{tableDir[8], tableDir[9], tableDir[10], tableDir[11]})
			lenHHEA = bytesToInt([]byte{tableDir[12], tableDir[13], tableDir[14], tableDir[15]})
		}
	}

	f, err := loadCMAP(offsetToCMAP, lenCMAP, fontFile)
	if err != nil {
		return nil, err
	}
	if err := loadHEAD(f, fontFile, offsetToHEAD, lenHEAD); err != nil {
		return nil, err
	}
	if err := loadLOCA(f, fontFile, offsetToLOCA, lenLOCA); err != nil {
		return nil, err
	}

	if err := loadGLYF(f, fontFile, offsetToGLYF, lenGLYF); err != nil {
		return nil, err
	}

	if err := loadHHEA(f, fontFile, offsetToHHEA, lenHHEA); err != nil {
		return nil, err
	}

	return f, nil
}

func loadCMAP(offsetToCMAP int, lenCMAP int, fontFile *os.File) (*font, error) {
	if offsetToCMAP < 1 || lenCMAP < 1 {
		return nil, fmt.Errorf("unexpected input")
	}
	// cmap
	// typedef struct {
	//	u16 version;
	//	u16 numberSubTables;
	//  cmap_encoding_subTable* subTables;
	//} cmap;

	//typedef struct {
	//	u16 platformID;
	//	u16 platformSpecificID;
	//	u32 offset; // is an offset from the start of the cmap table not from the start of the file
	//} cmap_encoding_subTable;

	//the platformID and platformSpecificID tell you which platform this encoding subtable
	//was intended for.
	//
	//platformID can have these values:
	//0 is called Unicode Platform and is used when the platfrom supports unicode
	//1 is to indicate Mac
	//2 is reserved and not used
	//3 is Microsoft encoding

	//when the platformID is 0 the platfromSpecificID can have the following values:
	//
	//0, 1, 2 - these indicate Unicode standard version 1 and 1.1 respectively and a value of 2 is deprecated.
	//3 - Unicode 2.0 with Basic Multilingual Plane only (BMP this is called Plane 0 which contains the languages and a lot of symbols),
	//4 - Unicode 2.0 with non-BMP allowed
	//5 - Unicode Variation Sequences, these represent variation of a Glyph as an example (from unicode.com FAQ) ≠ which is a variation of =.
	//6 - very similar to (4) but can use more formats for the cmap than number 4
	if _, err := fontFile.Seek(int64(offsetToCMAP), 0); err != nil {
		return nil, err
	}

	dataToRead := make([]byte, 2+2)
	nBytesRead, err := fontFile.Read(dataToRead)
	if err != nil {
		return nil, err
	}
	if nBytesRead != 2+2 {
		return nil, fmt.Errorf("unexpected input")
	}
	//log("version:", printBytes(dataToRead[:2]))
	//log("number of sub tables:", bytesToInt(dataToRead[2:]))
	offsetTo03 := -1 // the most common platformID and platformSpecificID combination is (0, 3), so we are going to focus on that combination only, supporting the other combinations is just a practice in being more complete and doesn't contribute too much to the overall approach
	for i := 0; i < bytesToInt(dataToRead[2:]); i++ {
		dataToRead := make([]byte, 2+2+4)
		nBytesRead, err = fontFile.Read(dataToRead)
		if err != nil {
			return nil, err
		}
		if nBytesRead != 2+2+4 {
			return nil, fmt.Errorf("unexpected input")
		}
		//log("platform id:", printBytes(dataToRead[:2]))
		//log("platform specific id:", bytesToInt(dataToRead[2:4]))
		//log("offset:", bytesToInt(dataToRead[4:]))
		if bytesToInt(dataToRead[:2]) == 0 && bytesToInt(dataToRead[2:4]) == 3 {
			offsetTo03 = bytesToInt(dataToRead[4:])
		}
	}
	if offsetTo03 < 0 {
		return nil, fmt.Errorf("unsupported format")
	}
	if _, err := fontFile.Seek(int64(offsetToCMAP+offsetTo03), 0); err != nil {
		return nil, err
	}

	// cmap format 4
	//
	// format4 is one of the most common formats and is used when you support the Unicode BMP
	// but not all the Codepoints in the entire plane range (0x0000 to 0xFFFF)
	//
	//typedef struct {
	//	u16  format;
	//	u16  length;
	//	u16  language;
	//	u16  segCountX2;
	//	u16  searchRange;
	//	u16  entrySelector;
	//	u16  rangeShift; array of segCountX2 / 2
	//	u16  reservedPad;
	//	u16  *endCode; End characterCode for each segment, last=0xFFFF.
	//	u16  *startCode;  Start character code for each segment.
	//	u16  *idDelta; Delta for all character codes in segment
	//	u16  *idRangeOffset; Offsets into glyphIdArray or 0
	//	u16  *glyphIdArray; Glyph index array (arbitrary length)
	//} format4;
	bytesRead := 0
	dataToRead = make([]byte, 2+2+2+2+2+2+2)
	nBytesRead, err = fontFile.Read(dataToRead)
	if err != nil {
		return nil, err
	}
	if nBytesRead != 2+2+2+2+2+2+2 {
		return nil, fmt.Errorf("unexpected input")
	}
	bytesRead += nBytesRead
	format4TableLength := bytesToInt(dataToRead[2:4])
	//log("format:", printBytes(dataToRead[:2]), bytesToInt(dataToRead[:2]))
	//log("length:", printBytes(dataToRead[2:4]), bytesToInt(dataToRead[2:4]))
	//log("language:", printBytes(dataToRead[4:6]), bytesToStr(dataToRead[4:6]))
	//log("seg_count_x2:", printBytes(dataToRead[6:8]), bytesToInt(dataToRead[6:8]))
	//log("search_range:", printBytes(dataToRead[8:10]), bytesToInt(dataToRead[8:10]))
	//log("entry_selector:", printBytes(dataToRead[10:12]), bytesToInt(dataToRead[10:12]))
	//log("range_shift:", printBytes(dataToRead[12:14]), bytesToInt(dataToRead[12:14]))
	segCountX2 := bytesToInt(dataToRead[6:8])
	uSegCountX2 := bytesToUint(dataToRead[6:8])
	endCodeArr := make([]byte, segCountX2)
	for i := 0; i < segCountX2; i += 2 {
		dataToRead = make([]byte, 2)
		nBytesRead, err = fontFile.Read(dataToRead)
		if err != nil {
			return nil, err
		}
		if nBytesRead != 2 {
			return nil, fmt.Errorf("unexpected input")
		}
		bytesRead += nBytesRead
		endCodeArr[i] = dataToRead[0]
		endCodeArr[i+1] = dataToRead[1]
	}
	dataToRead = make([]byte, 2)
	nBytesRead, err = fontFile.Read(dataToRead)
	if err != nil {
		return nil, err
	}
	if nBytesRead != 2 {
		return nil, fmt.Errorf("unexpected input")
	}
	bytesRead += nBytesRead
	//log("reserved_pad:", printBytes(dataToRead), bytesToStr(dataToRead))
	startCodeArr := make([]byte, segCountX2)
	for i := 0; i < segCountX2; i += 2 {
		dataToRead = make([]byte, 2)
		nBytesRead, err = fontFile.Read(dataToRead)
		if err != nil {
			return nil, err
		}
		if nBytesRead != 2 {
			return nil, fmt.Errorf("unexpected input")
		}
		bytesRead += nBytesRead
		startCodeArr[i] = dataToRead[0]
		startCodeArr[i+1] = dataToRead[1]
	}
	idDelta := make([]byte, segCountX2)
	for i := 0; i < segCountX2; i += 2 {
		dataToRead = make([]byte, 2)
		nBytesRead, err = fontFile.Read(dataToRead)
		if err != nil {
			return nil, err
		}
		if nBytesRead != 2 {
			return nil, fmt.Errorf("unexpected input")
		}
		bytesRead += nBytesRead
		idDelta[i] = dataToRead[0]
		idDelta[i+1] = dataToRead[1]
	}

	offset, err := fontFile.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, err
	}
	idRangeOffsetsStart := uint(offset)
	idRangeOffset := make([]byte, segCountX2)
	for i := 0; i < segCountX2; i += 2 {
		dataToRead = make([]byte, 2)
		nBytesRead, err = fontFile.Read(dataToRead)
		if err != nil {
			return nil, err
		}
		if nBytesRead != 2 {
			return nil, fmt.Errorf("unexpected input")
		}
		bytesRead += nBytesRead
		idRangeOffset[i] = dataToRead[0]
		idRangeOffset[i+1] = dataToRead[1]
	}
	//log("format 4 table length is", format4TableLength, "and we read", bytesRead, "bytesRemaining represent glyphIdArray", format4TableLength-bytesRead)
	glyphIdArray := make([]byte, 0)
	offset, err = fontFile.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, err
	}
	glyphIdArrayStart := uint(offset)
	if format4TableLength-bytesRead > 0 && (format4TableLength-bytesRead)%2 == 0 {
		glyphIdArray = make([]byte, format4TableLength-bytesRead)
		for i := 0; i < format4TableLength-bytesRead; i++ {
			dataToRead = make([]byte, 1)
			nBytesRead, err = fontFile.Read(dataToRead)
			if err != nil {
				return nil, err
			}
			if nBytesRead != 1 {
				return nil, fmt.Errorf("unexpected input")
			}
			glyphIdArray[i] = dataToRead[0]
		}
	}
	f := font{
		segCount: uSegCountX2 / 2,
	}
	f.startCodeArr = make([]uint, f.segCount)
	f.endCodeArr = make([]uint, f.segCount)
	f.idDelta = make([]uint, f.segCount)
	f.idRangeOffset = make([]uint, f.segCount)
	f.glyphIdArr = make([]uint, 0)
	for i := 0; i < segCountX2/2; i += 2 {
		//log("start code: ", printBytes(startCodeArr[i:i+2]), bytesToInt(startCodeArr[i:i+2]), " |  end code:", printBytes(endCodeArr[i:i+2]), bytesToInt(endCodeArr[i:i+2]), " |  id delta:", printBytes(idDelta[i:i+2]), bytesToInt(idDelta[i:i+2]), " |  id range offset:", printBytes(idRangeOffset[i:i+2]), bytesToInt(idRangeOffset[i:i+2]))
		f.startCodeArr[i] = bytesToUint(startCodeArr[i : i+2])
		f.endCodeArr[i] = bytesToUint(endCodeArr[i : i+2])
		f.idDelta[i] = bytesToUint(idDelta[i : i+2])
		f.idRangeOffset[i] = bytesToUint(idRangeOffset[i : i+2])
	}
	if format4TableLength-bytesRead > 0 && (format4TableLength-bytesRead)%2 == 0 {
		f.glyphIdArr = make([]uint, (format4TableLength-bytesRead)/2)
		for i := 0; i < (format4TableLength-bytesRead)/2; i += 2 {
			//log("glyph Id Array : for ", i, printBytes(glyphIdArray[i:i+2]), bytesToInt(glyphIdArray[i:i+2]))
			f.glyphIdArr[i] = bytesToUint(glyphIdArray[i : i+2])
		}
	}
	f.idRangeOffsetsStart = idRangeOffsetsStart
	f.glyphIdArrStart = glyphIdArrayStart
	return &f, nil
}

func loadHEAD(f *font, fontFile *os.File, offsetToHEAD, lenHEAD int) error {
	//typedef struct {
	//	Fixed version; Fixed: this is a fixed-point 32 bit integer, 16 bits for the whole part (the upper 16 bits) and 16 bit is used for the fractional part (lower 16 bits)[1].
	//	Fixed fontRevision;
	//
	//	u32 checkSumAdjustment;
	//	u32 magicNumber;
	//
	//	u16 flags;
	//	u16 unitsPerEm;
	//
	//	i64 created;
	//	i64 modified;
	//
	//	FWord XMin;
	//	FWord YMin;
	//	FWord XMax;
	//	FWord YMax;
	//
	//	u16 macStyle;
	//	u16 lowestRecPPEM;
	//	i16 fontDirectionHint;
	//	i16 indexToLocFormat;
	//	i16 glyphDataFormat;
	//}
	if _, err := fontFile.Seek(int64(offsetToHEAD), 0); err != nil {
		return err
	}
	dataToRead := make([]byte, 4+4+4+4+2+2+8+8+2+2+2+2+2+2+2+2+2)
	nBytesRead, err := fontFile.Read(dataToRead)
	if err != nil {
		return err
	}
	if nBytesRead != 4+4+4+4+2+2+8+8+2+2+2+2+2+2+2+2+2 {
		return fmt.Errorf("unexpected input")
	}
	//log("font version:", printBytes(dataToRead[:4]), bytesToInt(dataToRead[:2]), ".", bytesToInt(dataToRead[2:4]))   // 4
	//log("fontRevision:", printBytes(dataToRead[4:8]), bytesToInt(dataToRead[4:6]), ".", bytesToInt(dataToRead[6:8])) // 8
	//log("magic number:", printBytes(dataToRead[12:16]))                                                              // 16
	// 16:18
	//log("units per em:", printBytes(dataToRead[18:20]), bytesToInt(dataToRead[18:20]))
	// 20:28
	// 28:36
	//log("XMin:", printBytes(dataToRead[36:38]), bytesToInt16(dataToRead[36:38]))
	//log("YMin:", printBytes(dataToRead[38:40]), bytesToInt16(dataToRead[38:40]))
	//log("XMax:", printBytes(dataToRead[40:42]), bytesToInt16(dataToRead[40:42]))
	//log("YMax:", printBytes(dataToRead[42:44]), bytesToInt16(dataToRead[42:44]))
	f.Xmin = bytesToInt16(dataToRead[36:38])
	f.Ymin = bytesToInt16(dataToRead[38:40])
	f.Xmax = bytesToInt16(dataToRead[40:42])
	f.Ymax = bytesToInt16(dataToRead[42:44])
	// 44:46
	// 46:48
	// 48:50
	//log("indexToLocFormat:", printBytes(dataToRead[50:52]), bytesToInt(dataToRead[50:52]))
	//log("glyphDataFormat:", printBytes(dataToRead[52:54]), bytesToInt(dataToRead[52:54]))
	return nil
}

func loadLOCA(f *font, fontFile *os.File, offsetToLOCA int, lenLOCA int) error {
	if _, err := fontFile.Seek(int64(offsetToLOCA), 0); err != nil {
		return err
	}

	f.loca = make([]byte, lenLOCA)
	for i := 0; i < lenLOCA; i++ {
		dataToRead := make([]byte, 1)
		nBytesRead, err := fontFile.Read(dataToRead)
		if err != nil {
			return err
		}
		if nBytesRead != 1 {
			return fmt.Errorf("unexpected input")
		}
		f.loca[i] = dataToRead[0]
	}

	return nil
}

func loadGLYF(f *font, fontFile *os.File, offsetToGLYF int, lenGLYF int) error {
	if _, err := fontFile.Seek(int64(offsetToGLYF), 0); err != nil {
		return err
	}
	offset, err := fontFile.Seek(0, io.SeekCurrent)
	if err != nil {
		return err
	}
	f.glyfTableStart = uint(offset)
	f.glyf = make([]byte, lenGLYF)
	for i := 0; i < lenGLYF; i++ {
		dataToRead := make([]byte, 1)
		nBytesRead, err := fontFile.Read(dataToRead)
		if err != nil {
			return err
		}
		if nBytesRead != 1 {
			return fmt.Errorf("unexpected input")
		}
		f.glyf[i] = dataToRead[0]
	}

	return nil
}

// https://learn.microsoft.com/en-us/typography/opentype/spec/hhea
func loadHHEA(f *font, fontFile *os.File, offsetToHHEA, lenHHEA int) error {
	if _, err := fontFile.Seek(int64(offsetToHHEA), 0); err != nil {
		return err
	}
	dataToRead := make([]byte, 4+2+2+2)
	nBytesRead, err := fontFile.Read(dataToRead)
	if err != nil {
		return err
	}
	if nBytesRead != 4+2+2+2 {
		return fmt.Errorf("unexpected input")
	}
	//log("version:", printBytes(dataToRead[:4]), bytesToInt(dataToRead[0:2]), ".", bytesToInt(dataToRead[2:4])) // 4
	//log("ascender:", printBytes(dataToRead[4:6]), bytesToInt16(dataToRead[4:6]))
	//log("descender:", printBytes(dataToRead[6:8]), bytesToInt16(dataToRead[6:8]))
	//log("lineGap:", printBytes(dataToRead[8:10]), bytesToInt16(dataToRead[8:10]))
	f.ascent = bytesToInt16(dataToRead[4:6])
	f.descent = bytesToInt16(dataToRead[6:8])
	f.lineGap = bytesToInt16(dataToRead[8:10])

	return nil
}

func (f *font) GetAscent() int16 {
	return f.ascent
}
func (f *font) GetDescent() int16 {
	return f.descent
}
func (f *font) GetLineGap() int16 {
	return f.lineGap
}

func (f *font) GetMaxX() int16 {
	return f.Xmax
}

func (f *font) GetMaxY() int16 {
	return f.Ymax
}

func (f *font) GetMinX() int16 {
	return f.Xmin
}

func (f *font) GetMinY() int16 {
	return f.Ymin
}

func (f *font) GetAvailableCharCodes() []uint32 {
	result := make([]uint32, 0)
	for i, endCode := range f.endCodeArr {
		for c := f.startCodeArr[i]; c <= endCode; c++ {
			if f.GetGlyphIndex(uint(c)) > 0 {
				result = append(result, uint32(c))
			}
		}
	}

	return result
}

func (f *font) GetGlyphIndex(charCode uint) uint {
	idx := -1
	for i, endCode := range f.endCodeArr {
		if endCode >= charCode {
			//log(endCode, ">", charCode)
			idx = i
			break
		}
	}
	if idx < 0 {
		//log("< 0")
		return 0
	}
	startCode := f.startCodeArr[idx]
	if startCode > charCode {
		//log(charCode, "<startCode", startCode)
		return 0
	}
	//log("startCode", startCode)
	if f.idRangeOffset[idx] == 0 { // If the idRangeOffset value for the segment is not 0, the mapping of character codes relies on glyphIdArray
		return (f.idDelta[idx] + charCode) % 65536
	}
	i := int64(f.idRangeOffsetsStart) + int64(idx) + int64(f.idRangeOffset[idx]) + int64(charCode) - int64(startCode) + int64(f.idDelta[idx]) - int64(f.glyphIdArrStart)
	i = i % 65536
	if i < 0 {
		return 0
	}
	i *= 2
	if uint(i) > uint(len(f.glyphIdArr)-1) {
		return 0
	}
	//log("glyph index", f.glyphIdArr[i])

	return f.glyphIdArr[i]
}

func (f *font) GetGlyphOffset(charCode uint) uint {
	idx := f.GetGlyphIndex(charCode)
	//log(idx)

	return f.GetGlyphOffsetByIndex(idx)
}

func (f *font) GetGlyphOffsetByIndex(idx uint) uint {
	idx *= 2
	//log(idx, f.indexToLocFormat)
	//log(printBytes(f.loca, 2))
	offset := uint(0)
	if f.indexToLocFormat != 0 {
		offset = bytesToUint([]byte{f.loca[idx], f.loca[idx+1], f.loca[idx+2], f.loca[idx+3]}) * 2
	} else {
		offset = bytesToUint([]byte{f.loca[idx], f.loca[idx+1]}) * 2
	}

	return offset
}

func (f *font) GetGlyphDataByIndex(indx uint) *GlyphData {
	g := &GlyphData{}
	offset := f.GetGlyphOffsetByIndex(indx)
	idx := int(offset)
	//log("offset:", offset, "| glyf table starts at:", f.glyfTableStart, "| diff:", offset-f.glyfTableStart, "| glyf table size:", len(f.glyf))
	//offset -= f.glyfTableStart
	//log("number of contours", printBytes([]byte{f.glyf[offset], f.glyf[offset+1]}), bytesToInt16([]byte{f.glyf[offset], f.glyf[offset+1]}))
	numberOfContours := bytesToInt16([]byte{f.glyf[offset], f.glyf[offset+1]})
	//log("XMin", printBytes([]byte{f.glyf[offset+2], f.glyf[offset+3]}), bytesToInt16([]byte{f.glyf[offset+2], f.glyf[offset+3]}))
	g.XMin = bytesToInt16([]byte{f.glyf[offset+2], f.glyf[offset+3]})
	//log("YMin", printBytes([]byte{f.glyf[offset+4], f.glyf[offset+5]}), bytesToInt16([]byte{f.glyf[offset+4], f.glyf[offset+5]}))
	g.YMin = bytesToInt16([]byte{f.glyf[offset+4], f.glyf[offset+5]})
	//log("XMax", printBytes([]byte{f.glyf[offset+6], f.glyf[offset+7]}), bytesToInt16([]byte{f.glyf[offset+6], f.glyf[offset+7]}))
	g.XMax = bytesToInt16([]byte{f.glyf[offset+6], f.glyf[offset+7]})
	//log("YMax", printBytes([]byte{f.glyf[offset+8], f.glyf[offset+9]}), bytesToInt16([]byte{f.glyf[offset+8], f.glyf[offset+9]}))
	g.YMax = bytesToInt16([]byte{f.glyf[offset+8], f.glyf[offset+9]})
	g.NumberOfContours = uint16(numberOfContours)

	if numberOfContours > 0 { // If the number of contours is positive or zero, it is a single glyph;
		return f.getSimpleGlyphData(int(numberOfContours), idx, g)
	}
	// If the number of contours less than zero, the glyph is compound

	return f.getCompoundGlyphData(idx+10, g)
}

func (f *font) GetGlyphData(charCode uint) *GlyphData {
	idx := f.GetGlyphIndex(charCode)
	//log(string(rune(charCode)), charCode, "glyph index is", idx)

	return f.GetGlyphDataByIndex(idx)
}

func (f *font) getSimpleGlyphData(numberOfContours int, idx int, g *GlyphData) *GlyphData {
	endOfContours := make([]uint16, 0)
	// uint16	endPtsOfContours[n]	Array of last Points of each contour; n is the number of contours; array entries are point indices
	// uint16	instructionLength	Total number of bytes needed for instructions
	// uint8	instructions[instructionLength]	Array of instructions for this glyph
	// uint8	flags[variable]	Array of flags
	// uint8 or int16	xCoordinates[]	Array of x-coordinates; the first is relative to (0,0), others are relative to previous point
	// uint8 or int16	yCoordinates[]	Array of y-coordinates; the first is relative to (0,0), others are relative to previous point
	numPoints := uint(0)
	indx := 0
	for i := 0; i < numberOfContours*2; i += 2 {
		//log(i, printBytes([]byte{f.glyf[idx+10+i], f.glyf[idx+10+i+1]}), bytesToUint([]byte{f.glyf[idx+10+i], f.glyf[idx+10+i+1]}))
		endOfContours = append(endOfContours, uint16(bytesToUint([]byte{f.glyf[idx+10+i], f.glyf[idx+10+i+1]})))
		numPoints = bytesToUint([]byte{f.glyf[idx+10+i], f.glyf[idx+10+i+1]})
		indx = idx + 10 + i + 1 + 1
	}
	numPoints += 1
	//log("instructionLength", printBytes([]byte{f.glyf[indx], f.glyf[indx+1]}), bytesToUint([]byte{f.glyf[indx], f.glyf[indx+1]}))
	instructionLen := bytesToInt([]byte{f.glyf[indx], f.glyf[indx+1]})
	indx += instructionLen + 1 + 1
	//log("flags start at", indx)
	// we know how many flags are there
	flags := make([][]byte, numPoints)

	for i := 0; i < int(numPoints); i++ {
		flag := getBits(f.glyf[indx])
		flags[i] = flag
		if flag[3] > 0 { // repeat
			// read next byte
			repeatTimes := f.glyf[indx+1]
			for j := 0; j < int(repeatTimes); j++ {
				flags[i+j+1] = flag
			}
			indx += 2
			i += int(repeatTimes)
			continue
		} else {
			indx += 1
		}
	}

	xCoordinates := make([]int16, numPoints)
	value := int16(0)
	for i := 0; i < int(numPoints); i++ {
		if flags[i][1] > 0 { // short
			if flags[i][4] > 0 {
				value += int16(f.glyf[indx])
			} else {
				value -= int16(f.glyf[indx])
			}
			indx += 1
		} else { // long
			if flags[i][4] > 0 {
				// we don't need to read
			} else {
				value += bytesToInt16([]byte{f.glyf[indx], f.glyf[indx+1]})
				indx += 2
			}
		}
		xCoordinates[i] += value
	}

	yCoordinates := make([]int16, numPoints)
	value = int16(0)
	for i := 0; i < int(numPoints); i++ {
		if flags[i][2] > 0 { // short
			if flags[i][5] > 0 {
				value += int16(f.glyf[indx])
			} else {
				value -= int16(f.glyf[indx])
			}
			indx += 1
		} else { // long
			if flags[i][5] > 0 {
				// we don't need to read
			} else {
				value += bytesToInt16([]byte{f.glyf[indx], f.glyf[indx+1]})
				indx += 2
			}
		}
		yCoordinates[i] += value
	}

	g.Points = make([][]GlyphPoint, 0)
	contour := make([]GlyphPoint, 0)
	isAppended := false
	for i := 0; i < int(numPoints); i++ {
		isAppended = false
		gp := GlyphPoint{}
		gp.OnCurve = flags[i][0] > 0

		gp.X = xCoordinates[i]
		gp.Y = yCoordinates[i]
		//if !gp.OnCurve && flags[i-1][0] < 1 {
		//	missingPoint := GlyphPoint{
		//		OnCurve: true,
		//		x:       (xCoordinates[i] + xCoordinates[i-1]) / 2,
		//		y:       (yCoordinates[i] + yCoordinates[i-1]) / 2,
		//	}
		//	contour = append(contour, missingPoint)
		//}
		for _, end := range endOfContours {
			if int(end)+1 == i {
				//log(end, i)
				g.Points = append(g.Points, contour)
				contour = make([]GlyphPoint, 0)
				isAppended = true
				break
			}
		}
		contour = append(contour, gp)

	}
	if !isAppended {
		g.Points = append(g.Points, contour)
	}

	return g
}

func (f *font) getCompoundGlyphData(idx int, g *GlyphData) *GlyphData {
	g.Points = make([][]GlyphPoint, 0)
	// https://learn.microsoft.com/en-us/typography/opentype/spec/glyf#composite-glyph-description
	// A composite, or compound, glyph describes an outline indirectly by referencing other glyphs that
	// get incorporated into the composite glyph as components.
	// This is useful when the same contours are needed for multiple glyphs as it provides consistency for contours
	// that are repeated across multiple glyphs and can also provide significant size reduction.
	//
	// Composite glyphs may be nested within other composite glyphs—that is, a composite glyph parent may include other
	// composite glyphs as child components.
	// Thus, a composite glyph description is a directed graph.
	//
	// The data block for each child component starts with two uint16 values:
	// a flags field, and a glyph ID. These are followed by two argument fields,
	// though the size and interpretation of the arguments varies according to the flags that are set.
	// Optional fields describing a transformation can follow the arguments, depending on the flags.
	//
	// u16	flags	component flag
	// u16	glyphIndex	glyph index of component
	// u8, i8, u16 or i16	argument1	x-offset for component or point number; type depends on bits 0 and 1 in component flags
	// u8, i8, u16 or i16	argument2	y-offset for component or point number; type depends on bits 0 and 1 in component flags
	// [transform data]		optional transform data—see below
	//
	// The C pseudo-code fragment below shows how the sequence of component glyph records is stored and parsed;
	// definitions for the flag bits follow this fragment:
	//
	// do {
	//    uint16 flags;
	//    uint16 glyphIndex;
	//    if ( flags & ARG_1_AND_2_ARE_WORDS) {
	//    (int16 or FWORD) argument1;
	//    (int16 or FWORD) argument2;
	//    } else {
	//        uint16 arg1and2; /* (arg1 << 8) | arg2 */
	//    }
	//    if ( flags & WE_HAVE_A_SCALE ) {
	//        F2DOT14  scale;    /* Format 2.14 */
	//    } else if ( flags & WE_HAVE_AN_X_AND_Y_SCALE ) {
	//        F2DOT14  xscale;    /* Format 2.14 */
	//        F2DOT14  yscale;    /* Format 2.14 */
	//    } else if ( flags & WE_HAVE_A_TWO_BY_TWO ) {
	//        F2DOT14  xscale;    /* Format 2.14 */
	//        F2DOT14  scale01;   /* Format 2.14 */
	//        F2DOT14  scale10;   /* Format 2.14 */
	//        F2DOT14  yscale;    /* Format 2.14 */
	//    }
	//} while ( flags & MORE_COMPONENTS )
	//if (flags & WE_HAVE_INSTRUCTIONS){
	//    uint16 numInstr
	//    uint8 instr[numInstr]
	//}
	//
	// The F2DOT14 format consists of a signed, 2’s complement integer and an unsigned fraction.
	// To compute the actual value, take the integer and add the fraction.
	// 16-bit signed fixed number with the low 14 bits of fraction (2.14).
	//
	// Decimal Value	Hex Value	Bits                         Integer    Fraction
	//  1.999939	    0x7fff	    01 11 1111 1111 1111         1          16383/16384
	//  1.75	        0x7000	    01 11 0000 0000 0000         1          12288/16384
	//  0.000061	    0x0001	    00 00 0000 0000 0001         0	        1/16384
	//  0.0	            0x0000	    00 00 0000 0000 0000         0	        0/16384
	// -0.000061	    0xffff	    11 11 1111 1111 1111        -1	        16383/16384
	// -2.0	            0x8000	    10 00 0000 0000 0000        -2	        0/16384
	//
	// max value 1.99993896484375
	// min value -2.0
	//
	// Mask	Name	Description
	// 0x0001	ARG_1_AND_2_ARE_WORDS	Bit 0: If this is set, the arguments are 16-bit (uint16 or int16); otherwise, they are bytes (uint8 or int8).
	// 0x0002	ARGS_ARE_XY_VALUES	Bit 1: If this is set, the arguments are signed xy values; otherwise, they are unsigned point numbers.
	// 0x0004	ROUND_XY_TO_GRID	Bit 2: If set and ARGS_ARE_XY_VALUES is also set, the xy values are rounded to the nearest grid line. Ignored if ARGS_ARE_XY_VALUES is not set.
	// 0x0008	WE_HAVE_A_SCALE	Bit 3: This indicates that there is a simple scale for the component. Otherwise, scale = 1.0.
	// 0x0020	MORE_COMPONENTS	Bit 5: Indicates at least one more glyph after this one.
	// 0x0040	WE_HAVE_AN_X_AND_Y_SCALE	Bit 6: The x direction will use a different scale from the y direction.
	// 0x0080	WE_HAVE_A_TWO_BY_TWO	Bit 7: There is a 2 by 2 transformation that will be used to scale the component.
	// 0x0100	WE_HAVE_INSTRUCTIONS	Bit 8: Following the last component are instructions for the composite glyph.
	// 0x0200	USE_MY_METRICS	Bit 9: If set, this forces the aw and lsb (and rsb) for the composite to be equal to those from this component glyph. This works for hinted and unhinted glyphs.
	// 0x0400	OVERLAP_COMPOUND	Bit 10: If set, the components of the compound glyph overlap. Use of this flag is not required — that is, component glyphs may overlap without having this flag set. When used, it must be set on the flag word for the first component. Some rasterizer implementations may require fonts to use this flag to obtain correct behavior — see additional remarks, above, for the similar OVERLAP_SIMPLE flag used in simple-glyph descriptions.
	// 0x0800	SCALED_COMPONENT_OFFSET	Bit 11: The composite is designed to have the component offset scaled. Ignored if ARGS_ARE_XY_VALUES is not set.
	// 0x1000	UNSCALED_COMPONENT_OFFSET	Bit 12: The composite is designed not to have the component offset scaled. Ignored if ARGS_ARE_XY_VALUES is not set.
	// 0xE010	Reserved	Bits 4, 13, 14 and 15 are reserved: set to 0.
	//
	// The argument1 and argument2 fields of the component glyph record are used to determine the placement
	// of the child component glyph within the parent composite glyph.
	// They are interpreted either as an offset vector or as points from the parent and the child,
	// according to whether the ARGS_ARE_XY_VALUES flag is set.
	// This flag must always be set for the first component of a composite glyph.
	//
	const ARG_1_AND_2_ARE_WORDS = 0x0001
	const ARGS_ARE_XY_VALUES = 0x0002
	const ROUND_XY_TO_GRID = 0x0004
	const WE_HAVE_A_SCALE = 0x0008
	const MORE_COMPONENTS = 0x0020
	const WE_HAVE_AN_X_AND_Y_SCALE = 0x0040
	const WE_HAVE_A_TWO_BY_TWO = 0x0080
	const WE_HAVE_INSTRUCTIONS = 0x0100
	const USE_MY_METRICS = 0x0200
	const OVERLAP_COMPOUND = 0x0400
	const SCALED_COMPONENT_OFFSET = 0x0800
	const UNSCALED_COMPONENT_OFFSET = 0x1000
	for {
		//log("flags:", printBytes([]byte{f.glyf[idx], f.glyf[idx+1]}), bytesToUint([]byte{f.glyf[idx], f.glyf[idx+1]}))
		flags := bytesToUint([]byte{f.glyf[idx], f.glyf[idx+1]})
		glyphIdx := bytesToUint([]byte{f.glyf[idx+2], f.glyf[idx+3]})
		childGlyph := f.GetGlyphDataByIndex(glyphIdx)
		//log("child idx", glyphIdx)
		if childGlyph == nil {
			//log("child data was not found")
			break
		}
		arg1 := int16(0)
		arg2 := int16(0)
		if flags&ARG_1_AND_2_ARE_WORDS != 0 {
			//log("ARG_1_AND_2_ARE_WORDS")
			arg1 = bytesToInt16([]byte{f.glyf[idx+4], f.glyf[idx+5]})
			arg1 = bytesToInt16([]byte{f.glyf[idx+6], f.glyf[idx+7]})
			idx = idx + 7 + 1
		} else {
			//log("ARG_1_AND_2_ARE_NOT_WORDS")
			arg1 = bytesToInt16([]byte{0x0, f.glyf[idx+4]})
			arg2 = bytesToInt16([]byte{0x0, f.glyf[idx+5]})
			idx = idx + 5 + 1
		}

		//log("arg1:", arg1)
		//log("arg2:", arg2)
		deltax := float64(0)
		deltay := float64(0)
		if flags&ARGS_ARE_XY_VALUES != 0 {
			//log("ARGS_ARE_XY_VALUES")
			// The argument1 and argument2 fields of the component glyph record are used to determine the placement of
			// the child component glyph within the parent composite glyph.
			// They are interpreted either as an offset vector or as points from the parent and the child,
			// according to whether the ARGS_ARE_XY_VALUES flag is set.
			// This flag must always be set for the first component of a composite glyph.
			//
			// If ARGS_ARE_XY_VALUES is set, then argument1 and argument2 are interpreted as units in the design
			// coordinate system
			// and an offset vector (x, y) = (argument1, argument2) is added to the coordinates of each control point
			// of the component glyph.
			// In a variable font, the offset vector can be modified by deltas in the 'gvar' table; see Point numbers and processing for composite glyphs in the 'gvar' chapter for details.
			// If a scale or transform matrix is provided, the offset vector might or might not be subject to the transformation;
			// see the discussion below of the SCALED_COMPONENT_OFFSET and UNSCALED_COMPONENT_OFFSET flags for details.
			//
			//	If ARGS_ARE_XY_VALUES is set and the ROUND_XY_TO_GRID flag is also set, the offset vector
			//	(after any transformation and variation deltas are applied) is grid-fitted, with the x and y values rounded to the nearest pixel grid line.
			//
			deltax = float64(arg1)
			deltay = float64(arg2)
		} else {
			//log("ARGS_ARE_NOT_XY_VALUES")
			//	If ARGS_ARE_XY_VALUES is not set, then argument1 is a point number in the parent glyph
			//	(from contours incorporated and re-numbered from previous component glyphs);
			//	and argument2 is a point number (prior to re-numbering) from the child component glyph.
			//	Phantom points from the parent or the child may be referenced.
			//	The child component glyph is positioned within the parent glyph by aligning the two points.
			//	If a scale or transform matrix is provided, the transformation is applied to the child’s point before the points are aligned.
			//
			// 	In a variable font, when a component is positioned by alignment of points,
			//	deltas are applied to component glyphs before this alignment is done.
			//	Any deltas specified for the parent composite glyph to be applied to components positioned by point
			//	alignment are ignored. See Point numbers and processing for composite glyphs in the 'gvar' chapter for details.
			i := 0
			cx := int16(0)
			cy := int16(0)
			foundChildPoint := false
			for _, c := range childGlyph.Points {
				for _, gp := range c {
					if i == int(arg2) {
						cx = gp.X
						cy = gp.Y
						foundChildPoint = true
						break
					}
					i++
				}
				if foundChildPoint {
					break
				}
			}
			//if !foundChildPoint {
			//	switch int(arg2) - i {
			//	case 0: // pp0
			//		foundChildPoint = true
			//		cx = 0
			//		cy = 0
			//	case 1: // pp1
			//		foundChildPoint = true
			//		cx = 0
			//		cy = 0
			//	case 2: // pp2
			//		foundChildPoint = true
			//		cx = 0
			//		cy = 0
			//	case 3: // pp4
			//		foundChildPoint = true
			//		cx = 0
			//		cy = 0
			//	}
			//}
			if foundChildPoint {
				//log("cx, cy", cx, cy, arg2, "found?", foundChildPoint)
				i = 0
				x := int16(0)
				y := int16(0)
				foundParentPoint := false
				for _, c := range g.Points {
					for _, gp := range c {
						if i == int(arg1) {
							x = gp.X
							y = gp.Y
							foundParentPoint = true
							break
						}
						i++
					}
					if foundParentPoint {
						break
					}
				}
				if foundParentPoint {
					//log("x, y", x, y, arg1, "found?", foundParentPoint)
					deltax = float64(x - cx)
					deltay = float64(y - cy)
				} else {
					//log("did not find parent point", "all points", i, "the needed point was", arg1)
				}
			} else {
				//log("did not find child point", "all points", i, "the needed point was", arg2)
			}
		}

		xscale := float64(1)
		scale01 := float64(0)
		scale10 := float64(0)
		yscale := float64(1)
		if flags&WE_HAVE_A_SCALE != 0 {
			//log("WE_HAVE_A_SCALE")
			xscale = get2Dot14([]byte{f.glyf[idx], f.glyf[idx+1]})
			idx = idx + 1 + 1
		} else if flags&WE_HAVE_AN_X_AND_Y_SCALE != 0 {
			//log("WE_HAVE_AN_X_AND_Y_SCALE")
			xscale = get2Dot14([]byte{f.glyf[idx], f.glyf[idx+1]})
			yscale = get2Dot14([]byte{f.glyf[idx+2], f.glyf[idx+3]})
			idx = idx + 3 + 1
		} else if flags&WE_HAVE_A_TWO_BY_TWO != 0 {
			//log("WE_HAVE_A_TWO_BY_TWO")
			xscale = get2Dot14([]byte{f.glyf[idx], f.glyf[idx+1]})
			scale01 = get2Dot14([]byte{f.glyf[idx+2], f.glyf[idx+3]})
			scale10 = get2Dot14([]byte{f.glyf[idx+4], f.glyf[idx+5]})
			yscale = get2Dot14([]byte{f.glyf[idx+6], f.glyf[idx+7]})
			idx = idx + 7 + 1
		}

		//log("xscale", xscale)
		//log("scale01", scale01)
		//log("scale10", scale10)
		//log("yscale", yscale)
		//log("deltax", deltax)
		//log("deltay", deltay)
		// x′ = xscale * x + scale10 * y + deltax
		// y′ = scale01 * x + yscale * y + deltay
		for _, c := range childGlyph.Points {
			contour := make([]GlyphPoint, 0)
			for _, gp := range c {
				contour = append(contour, GlyphPoint{
					OnCurve: gp.OnCurve,
					X:       int16(xscale*float64(gp.X) + scale10*float64(gp.Y) + deltax),
					Y:       int16(scale01*float64(gp.X) + yscale*float64(gp.Y) + deltay),
				})
			}
			g.Points = append(g.Points, contour)
		}
		if flags&MORE_COMPONENTS == 0 {
			break
		}
	}

	return g
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

func get2Dot14(bytes []byte) float64 {
	bitsPart1 := getBits(bytes[0])
	bitsPart2 := getBits(bytes[1])
	b2 := bitsToInt([]byte{bitsPart1[7], bitsPart1[6]})
	b14 := bitsToInt([]byte{bitsPart1[5], bitsPart1[4], bitsPart1[3], bitsPart1[2], bitsPart1[1], bitsPart1[0], bitsPart2[7], bitsPart2[6], bitsPart2[5], bitsPart2[4], bitsPart2[3], bitsPart2[2], bitsPart2[1], bitsPart2[0]})

	return float64(b2) + float64(b14/16384)
}

func bitsToInt(bits []byte) int {
	result := 0
	for i, b := range reverse(bits) {
		if b == 0 {
			continue
		}
		//result += 2^i
		if i == 0 {
			result += 1
			continue
		}

		p := 2
		for j := 1; j < i; j++ {
			p *= 2
		}
		result += p
	}

	return result
}

func reverse(b []byte) []byte {
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return b
}

func getBits(b byte) []byte {
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

	for i, j := 0, len(chunk)-1; i < j; i, j = i+1, j-1 {
		chunk[i], chunk[j] = chunk[j], chunk[i]
	}

	return chunk
}

func bytesToInt16(bytes []byte) int16 {
	return int16(binary.BigEndian.Uint16(bytes))
}

func bytesToStr(bytes []byte) string {
	return string(bytes)
}

func printBytes(bytes []byte, n ...int) string {
	if len(n) > 0 {
		msg := ""
		for i := 0; i < len(bytes); i += n[0] {
			msg += fmt.Sprintln(i, ":", printBytes(bytes[i:i+n[0]]), bytesToInt(bytes[i:i+n[0]]), bytesToInt(bytes[i:i+n[0]])*2)
		}
		return msg
	}
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

func bytesToUint(bytes []byte) uint {
	result := uint(0)
	for i := 0; i < len(bytes); i++ {
		if i > 3 {
			break
		}
		result = result << 8
		result += uint(bytes[i])

	}

	return result
}

func areEqual(a, b []byte, n int) bool {
	for i := 0; i < n; i++ {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func log(msgs ...interface{}) {
	_, f, line, _ := runtime.Caller(1)
	baseDir := ""
	pwd, err := os.Getwd()
	if err == nil {
		baseDir = pwd
	}

	relativePath := strings.Replace(f, baseDir+"/", "", -1)
	formatedMsg := fmt.Sprintln(time.Now().UTC().Format("15:04:05.999 02-01-2006"), fmt.Sprintf("%s:%d", relativePath, line), msgs)
	fmt.Print(formatedMsg)
}

func logError(msgs ...interface{}) {
	_, f, line, _ := runtime.Caller(1)
	baseDir := ""
	pwd, err := os.Getwd()
	if err == nil {
		baseDir = pwd
	}

	relativePath := strings.Replace(f, baseDir+"/", "", -1)
	formatedMsg := fmt.Sprintln("*********ERROR", time.Now().UTC().Format("15:04:05.999 02-01-2006"), fmt.Sprintf("%s:%d", relativePath, line), msgs)
	fmt.Print(formatedMsg)
}
