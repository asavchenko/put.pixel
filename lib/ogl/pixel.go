package ogl

func PutPixel(x, y int, color byte) {
	if y >= height {
		return
	}
	i := (x + yTable[y]) * 3
	if i < 0 {
		return
	}

	if i+2 > len(screen)-1 {
		return
	}
	screen[i] = color
	screen[i+1] = color
	screen[i+2] = color
}

func putPixel(x, y int, color byte) {
	i := (x + yTable[y]) * 3
	screen[i] = color
	screen[i+1] = color
	screen[i+2] = color
}
