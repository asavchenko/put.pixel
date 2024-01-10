package characters

func GetCharacterWidth(size int) int {
	switch size {
	case 14:
		return 17
	default:
		if size > 14 {
			return 17 + size - 14
		}
		if size < 0 {
			return 3
		}
		return 17 - (14 - size)
	}
}

func GetCharacterHeight(size int) int {
	switch size {
	case 14:
		return 20
	default:
		if size > 14 {
			return 20 + size - 14
		}
		if size < 0 {
			return 5
		}
		return 20 - (14 - size)
	}
}

func GetSpaceSizeBtwCharacters(size int) int {
	return GetCharacterWidth(size) / 9
}

func GetLineSpaceSize(size int) int {
	return GetCharacterWidth(size) / 6
}
