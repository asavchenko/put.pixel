package freetype

type GlyphPoint struct {
	OnCurve bool
	X       int16
	Y       int16
}

type GlyphData struct {
	NumberOfContours uint16 // the NumberOfContours tells us how many contours this shape has,
	// it also gives us the size of endPtsOfContours,
	// the endPtsOfContours array gives us the indices of the end Points of the contours
	// i.e: endPtsOfContours[0] gives us the index into xCoordinates and yCoordinates array
	// where the first contour ends and endPtsOfContours[1] gives us the index where the 2nd contour ends.
	//
	// This means the last value in endPtsContour array will give us the number of Points
	// all the contours combined have thus we can allocate enough memory for all of them.
	XMin int16
	YMin int16
	XMax int16
	YMax int16

	// the intructionLength and instructions together gives us the instructions that are needed
	// to do grid fitting (hinting) for this Glyph.
	// We won't be doing this as it requires us to code a complete virtual machine to run the
	//	instructions but its a fun project to do so maybe for another time.
	//instructionLength uint16
	//instructions      []byte

	// On Curve                                 0	If set, the point is on the curve;
	//                                              Otherwise, it is off the curve.
	//
	// x-Short Vector                           1	If set, the corresponding x-coordinate is 1 byte long;
	//                                              Otherwise, the corresponding x-coordinate is 2 bytes long
	//
	// y-Short Vector                           2	If set, the corresponding y-coordinate is 1 byte long;
	//                                              Otherwise, the corresponding y-coordinate is 2 bytes long
	//
	//Repeat	                                3	If set, the next byte specifies the number of additional times this set of flags is to be repeated.
	//                                              In this way, the number of flags listed can be smaller than the number of Points in a character.
	//
	//This x is same (Positive x-Short vector)	4	This flag has one of two meanings, depending on how the x-Short Vector flag is set.
	//                                              If the x-Short Vector bit is set, this bit describes the sign of the value,
	//                                              with a value of 1 equalling positive and a zero value negative.
	//
	//                                              If the x-short Vector bit is not set, and this bit is set,
	//                                              then the current x-coordinate is the same as the previous x-coordinate.
	//
	//                                              If the x-short Vector bit is not set, and this bit is not set,
	//                                              the current x-coordinate is a signed 16-bit delta vector. In this case, the delta vector is the change in x
	//
	//This y is same (Positive y-Short vector)	5	This flag has one of two meanings, depending on how the y-Short Vector flag is set.
	//                                              If the y-Short Vector bit is set, this bit describes the sign of the value,
	//                                              with a value of 1 equalling positive and a zero value negative.
	//
	//                                              If the y-short Vector bit is not set, and this bit is set, then
	//                                              the current y-coordinate is the same as the previous y-coordinate.
	//
	//If the y-short Vector bit is not set, and this bit is not set, the current y-coordinate is a signed 16-bit delta vector. In this case, the delta vector is the change in y

	//Reserved	                                6 - 7	Set to zero
	//flags []byte

	//xCoordinates     []int16 // Array of x-coordinates; the first is relative to (0,0), others are relative to previous point
	//yCoordinates     []int16 // Array of y-coordinates; the first is relative to (0,0), others are relative to previous point
	Points [][]GlyphPoint
	//endPtsOfContours []uint16 // Array of last Points of each contour; n is the number of contours; array entries are point indices
}

func (d GlyphData) GetWithMiddlePoints() [][]GlyphPoint {
	newPoints := make([][]GlyphPoint, 0)
	for _, contour := range d.Points {
		newContour := make([]GlyphPoint, 0)
		var prev GlyphPoint
		for pi, p := range contour {
			if p.OnCurve {
				newContour = append(newContour, p)
				prev = p
				continue
			}
			if prev.OnCurve {
				newContour = append(newContour, p)
				prev = p
				continue
			}

			mp := GlyphPoint{
				OnCurve: true,
				X:       (contour[pi].X + contour[pi-1].X) / 2,
				Y:       (contour[pi].Y + contour[pi-1].Y) / 2,
			}
			newContour = append(newContour, mp)
			prev = p
			newContour = append(newContour, p)
		}
		if len(newContour) > 0 {
			newContour = append(newContour, newContour[0])
		}
		newPoints = append(newPoints, newContour)
	}

	return newPoints
}

func (d GlyphData) ContainsPoint(x, y int) bool {
	points := d.GetWithMiddlePoints()
	winding := 0
	for _, contour := range points {
		i := 0
		for {
			if i+1 >= len(contour) {
				break
			}
			var startY float64
			var endY float64
			var startX float64
			var endX float64
			if contour[i].OnCurve && contour[i+1].OnCurve { // simple line case
				startX = float64(contour[i].X)
				startY = float64(contour[i].Y)
				endX = float64(contour[i+1].X)
				endY = float64(contour[i+1].Y)
				if startY <= float64(y) && endY > float64(y) { // Crossing from top to bottom
					dir := (startX-float64(x))*(endY-float64(y)) - (endX-float64(x))*(startY-float64(y))
					if dir > 0 { // On Right
						winding -= 1
					}
				} else if startY > float64(y) && endY <= float64(y) { // Crossing from bottom to top
					dir := (startX-float64(x))*(endY-float64(y)) - (endX-float64(x))*(startY-float64(y))
					if dir < 0 { // On Left
						winding += 1
					}
				}
				i += 1
				continue
			}
			if contour[i].OnCurve && !contour[i+1].OnCurve { // Bezier curve case
				startX = float64(contour[i].X)
				startY = float64(contour[i].Y)
				endX = float64(contour[i+2].X)
				endY = float64(contour[i+2].Y)
				if startY <= float64(y) && endY > float64(y) { // Crossing from top to bottom
					dir := (startX-float64(x))*(endY-float64(y)) - (endX-float64(x))*(startY-float64(y))
					if dir > 0 { // On Right
						winding -= 1
					}
				} else if startY > float64(y) && endY <= float64(y) { // Crossing from bottom to top
					dir := (startX-float64(x))*(endY-float64(y)) - (endX-float64(x))*(startY-float64(y))
					if dir < 0 { // On Left
						winding += 1
					}
				}

				xa := float64(contour[i].X)
				ya := float64(contour[i].Y)
				xb := float64(contour[i+1].X)
				yb := float64(contour[i+1].Y)
				xc := float64(contour[i+2].X)
				yc := float64(contour[i+2].Y)

				b := -(-float64(y)*xa + yc*xa - yc*float64(x) + ya*float64(x) + float64(y)*xc - ya*xc) / (xb*yc - xb*ya - yc*xa - yb*xc + yb*xa + ya*xc)
				c := (xb*float64(y) - xb*ya - yb*float64(x) + yb*xa - float64(y)*xa + ya*float64(x)) / (xb*yc - xb*ya - yc*xa - yb*xc + yb*xa + ya*xc)
				a := 1.0 - b - c

				if a >= 0.0 && a <= 1.0 && b >= 0.0 && b <= 1.0 && c >= 0.0 && c <= 1.0 {
					controlToStartX := xa - xb
					controlToStartY := ya - yb
					controlToEndX := xc - xb
					controlToEndY := yc - yb
					crossZ := controlToStartX*controlToEndY - controlToStartY*controlToEndX

					uvPointX := a + 0.5*b
					uvPointY := a
					uvValue := uvPointX*uvPointX - uvPointY
					//if crossZ < 0 {
					//	uvValue = -uvValue
					//}

					if uvValue < 0.0 {
						if crossZ < 0 {
							winding -= 1
						} else {
							winding += 1
						}
					}
				}

				i += 2
			}
		}
	}

	return winding != 0
}

type Font interface {
	GetAvailableCharCodes() []uint32
	GetGlyphIndex(charCode uint) uint
	GetGlyphOffset(charCode uint) uint
	GetGlyphData(charCode uint) *GlyphData
	GetMaxX() int16
	GetMaxY() int16
	GetMinY() int16
	GetMinX() int16
	GetAscent() int16
	GetDescent() int16
	GetLineGap() int16
}
