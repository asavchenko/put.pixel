package characters

import "assa.com/put.pixel/src/characters/utf8"

var chrs []*Chr
var availableCharCodes = utf8.GetAvailableCharCodes()

func GetCharacterWidth(size int) int {
	switch size {
	case 14:
		return utf8.GetShapeWidth()
	default:
		return utf8.GetShapeWidth() * size / 14
	}
}

func GetCharacterHeight(size int) int {
	switch size {
	case 14:
		return utf8.GetShapeHeight()
	default:
		return utf8.GetShapeHeight() * size / 14
	}
}

func GetSpaceSizeBtwCharacters(size int) int {
	return 1
}

func GetLineSpaceSize(size int) int {
	return 2
}

func AddText(text string, x, y int, fontSize int, color int) {
	for _, r := range text {
		chrs = append(chrs, GetNew(r, x, y, uint32(color)).SetCharacterSize(fontSize))
		x += GetCharacterWidth(fontSize) + GetSpaceSizeBtwCharacters(fontSize)
	}
}

func ReplaceText(text string, x, y int, fontSize int, color int) {
	chrs = make([]*Chr, 0)
	for _, r := range text {
		chrs = append(chrs, GetNew(r, x, y, uint32(color)).SetCharacterSize(fontSize))
		x += GetCharacterWidth(fontSize) + GetSpaceSizeBtwCharacters(fontSize)
	}
}

func GetText() []*Chr {
	return chrs
}

func GetAvailableCharCodes() []uint32 {
	return utf8.GetAvailableCharCodes()
}
