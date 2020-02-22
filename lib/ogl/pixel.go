package ogl

func PutPixel(x, y int, color byte) {
	if y >= height {
		return
	}
	i := (x + yTable[y]) * 3
	if i < 0 {
		return
	}

	if !doubleBuffering {
		if i+2 > len(pixelArr)-1 {
			return
		}
		pixelArr[i] = color
		pixelArr[i+1] = color
		pixelArr[i+2] = color
		return
	}

	if index == 0 {
		if i+2 > len(pixelArr1)-1 {
			return
		}
		pixelArr1[i] = color
		pixelArr1[i+1] = color
		pixelArr1[i+2] = color
		return
	}

	if i+2 > len(pixelArr2)-1 {
		return
	}
	pixelArr2[i] = color
	pixelArr2[i+1] = color
	pixelArr2[i+2] = color
}
