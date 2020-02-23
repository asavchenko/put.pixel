package ogl

// #include <string.h>
import "C"

import (
	"assa.com/put.pixel/lib/mlib"
	"unsafe"
)

func Line(xa, ya, xb, yb int, color byte) int {
	yb = Ym - yb
	x1 := xa
	y1 := ya

	x2 := xb
	y2 := yb

	visible := 0
	code := 0x0

	if x2 < 0 {
		code += 1
	} else {
		if x2 > Xm {
			code += 2
		}
	}

	if y2 < 0 {
		code += 8
	} else {
		if y2 > Ym {
			code += 4
		}
	}

	if x1 < 0 {
		code += 16
	} else {
		if x1 > Xm {
			code += 32
		}
	}

	if y1 < 0 {
		code += 128
	} else {
		if y1 > Ym {
			code += 64
		}
	}

	switch code {
	case 0x98:
	case 0x89:
	case 0x9A:
	case 0xA9:
	case 0x95:
	case 0x59:
	case 0x54:
	case 0x45:
	case 0x46:
	case 0x64:
	case 0x6A:
	case 0xA6:
	case 0x8A:
	case 0xA8:
	case 0x88:
	case 0x99:
	case 0xAA:
	case 0x55:
	case 0x44:
	case 0x66:
	case 0x11:
	case 0x22:
	case 0x91:
	case 0x19:
	case 0x51:
	case 0x15:
	case 0xA2:
	case 0x2A:
	case 0x26:
	case 0x62:
	case 0x56:
	case 0x65:
	case 0x00:
		visible++
	case 0x10:
		visible++
		y1 = ya - xa*(yb-ya)/(xb-xa)
		x1 = 0
	case 0x01:
		visible++
		y2 = ya - xa*(yb-ya)/(xb-xa)
		x2 = 0

	case 0x80:
		visible++
		x1 = xa - ya*(xb-xa)/(yb-ya)
		y1 = 0
	case 0x08:
		visible++
		x2 = xa - ya*(xb-xa)/(yb-ya)
		y2 = 0

	case 0x02:
		visible++
		y2 = ya + (yb-ya)*(Xm-xa)/(xb-xa)
		x2 = Xm
	case 0x20:
		visible++
		y1 = ya + (yb-ya)*(Xm-xa)/(xb-xa)
		x1 = Xm
	case 0x04:
		visible++
		x2 = xa + (xb-xa)*(Ym-ya)/(yb-ya)
		y2 = Ym
	case 0x40:
		visible++
		x1 = xa + (xb-xa)*(Ym-ya)/(yb-ya)
		y1 = Ym
	case 0x84:
		visible++
		x1 = xa - ya*(xb-xa)/(yb-ya)
		y1 = 0
		x2 = x1 + (x2-x1)*Ym/yb
		y2 = Ym
	case 0x48:
		visible++
		x2 = xa - ya*(xb-xa)/(yb-ya)
		y2 = 0
		x1 = xa - (xb-xa)*(Ym-ya)/ya
		y1 = Ym
	case 0x12:
		visible++
		y1 = ya - xa*(yb-ya)/(xb-xa)
		x1 = 0
		y2 = ya + (yb-ya)*Xm/xb
		x2 = Xm
	case 0x21:
		visible++
		y2 = ya - xa*(yb-ya)/(xb-xa)
		x2 = 0
		y1 = ya - (yb-ya)*(Xm-xa)/xa
		x1 = Xm
	case 0x90:
		visible++
		y1 = ya - xa*(yb-ya)/(xb-xa)
		x1 = 0
		if y1 < 0 {
			x1 = -y1 * x2 / (y2 - y1)
			y1 = 0
		}
	case 0x09:
		visible++
		y2 = ya - xa*(yb-ya)/(xb-xa)
		x2 = 0
		if y2 < 0 {
			x2 = xa + ya*xa/(y2-ya)
			y2 = 0
		}
	case 0x0A:
		visible++
		y2 = ya + (yb-ya)*(Xm-xa)/(xb-xa)
		x2 = Xm
		if y2 < 0 {
			x2 = xa - ya*(xb-xa)/(yb-ya)
			y2 = 0
		}
	case 0xA0:
		visible++
		y1 = ya + (yb-ya)*(Xm-xa)/(xb-xa)
		x1 = Xm
		if y1 < 0 {
			x1 = xa - ya*(xb-xa)/(yb-ya)
			y1 = 0
		}
	case 0x05:
		visible++
		y2 = ya - xa*(yb-ya)/(xb-xa)
		x2 = 0
		if y2 > Ym {
			x2 = xa + -xa*(Ym-ya)/(y2-ya)
			y2 = Ym
		}
	case 0x50:
		visible++
		y1 = ya - xa*(yb-ya)/(xb-xa)
		x1 = 0
		if y1 > Ym {
			x1 = xb * (Ym - y1) / (yb - y1)
			y1 = Ym
		}
	case 0x06:
		visible++
		y2 = ya + (yb-ya)*(Xm-xa)/(xb-xa)
		x2 = Xm
		if y2 > Ym {
			x2 = x1 + (x2-x1)*(Ym-y1)/(y2-y1)
			y2 = Ym
		}
	case 0x60:
		visible++
		y1 = ya + (yb-ya)*(Xm-xa)/(xb-xa)
		x1 = Xm
		if y1 > Ym {
			x1 = x1 + (x2-x1)*(Ym-y1)/(y2-y1)
			y1 = Ym
		}
	case 0x96:
		y1 = ya - xa*(yb-ya)/(xb-xa)
		x1 = 0
		if y1 > Ym {
			break
		}
		y2 = y1 + (y2-y1)*Xm/x2
		x2 = Xm
		if y2 < 0 {
			break
		}
		if y1 < 0 {
			x1 = -y1 * x2 / (y2 - y1)
			y1 = 0
		}
		if y2 > Ym {
			x2 = x2 * (Ym - y1) / (y2 - y1)
			y2 = Ym
		}
		visible++
	case 0x69:
		y2 = ya - xa*(yb-ya)/(xb-xa)
		x2 = 0
		if y2 > Ym {
			break
		}
		y1 = y1 - (y2-y1)*(Xm-x1)/x1
		x1 = Xm
		if y1 < 0 {
			break
		}
		if y2 < 0 {
			x2 = x1 - y1*(x2-x1)/(y2-y1)
			y2 = 0
		}
		if y1 > Ym {
			x1 = x1 + (x2-x1)*(Ym-y1)/(y2-y1)
			y1 = Ym
		}
		visible++
	case 0x5A:
		y1 = ya - xa*(yb-ya)/(xb-xa)
		x1 = 0
		if y1 < 0 {
			break
		}
		y2 = y1 + (y2-y1)*(Xm-x1)/(x2-x1)
		x2 = Xm
		if y2 > Ym {
			break
		}
		if y1 > Ym {
			x1 = x1 + (x2-x1)*(Ym-y1)/(y2-y1)
			y1 = Ym
		}
		if y2 < 0 {
			x2 = x1 - y1*(x2-x1)/(y2-y1)
			y2 = 0
		}
		visible++
	case 0xA5:
		y2 = ya - xa*(yb-ya)/(xb-xa)
		x2 = 0
		if y2 < 0 {
			break
		}
		y1 = ya + (y2-y1)*(Xm-x1)/(x2-x1)
		x1 = Xm
		if y1 > Ym {
			break
		}
		if y2 > Ym {
			x2 = x1 + (x2-x1)*(Ym-y1)/(y2-y1)
			y2 = Ym
		}
		if y1 < 0 {
			x1 = x1 - y1*(x2-x1)/(y2-y1)
			y1 = 0
		}
		visible++
	case 0xA4:
		x2 = xa + (xb-xa)*(Ym-ya)/(yb-ya)
		y2 = Ym
		if x2 > Xm {
			break
		}
		y1 = y1 + (y2-y1)*(Xm-x1)/(x2-x1)
		x1 = Xm
		if y1 < 0 {
			x1 = x1 - y1*(x2-x1)/(y2-y1)
			y1 = 0
		}
		visible++
	case 0x4A:
		x1 = x1 + (xb-xa)*(Ym-ya)/(yb-ya)
		y1 = Ym
		if x1 > Xm {
			break
		}
		y2 = y1 + (y2-y1)*(Xm-x1)/(x2-x1)
		x2 = Xm
		if y2 < 0 {
			x2 = x1 - y1*(x2-x1)/(y2-y1)
			y2 = 0
		}
		visible++
	case 0x94:
		x2 = xa + (xb-xa)*(Ym-ya)/(yb-ya)
		y2 = Ym
		if x2 < 0 {
			break
		}
		y1 = y1 - x1*(y2-y1)/(x2-x1)
		x1 = 0
		if y1 < 0 {
			x1 = x1 - y1*(x2-x1)/(y2-y1)
			y1 = 0
		}
		visible++
	case 0x49:
		x1 = xa + (xb-xa)*(Ym-ya)/(yb-ya)
		y1 = Ym
		if x1 < 0 {
			break
		}
		y2 = y1 - x1*(y2-y1)/(x2-x1)
		x2 = 0
		if y2 < 0 {
			x2 = x1 - y1*(x2-x1)/(y2-y1)
			y2 = 0
		}
		visible++
	case 0x58:
		x2 = xa - ya*(xb-xa)/(yb-ya)
		y2 = 0
		if x2 < 0 {
			break
		}
		x1 = x1 + (x2-x1)*(Ym-ya)/(y2-y1)
		y1 = Ym
		if x1 < 0 {
			y1 = y1 - x1*(y2-y1)/(x2-x1)
			x1 = 0
		}
		visible++
	case 0x85:
		x1 = xa - ya*(xb-xa)/(yb-ya)
		y1 = 0
		if x1 < 0 {
			break
		}
		x2 = x1 + (x2-x1)*(Ym-y1)/(y2-y1)
		y2 = Ym
		if x2 < 0 {
			y2 = y1 - x1*(y2-y1)/(x2-x1)
			x2 = 0
		}
		visible++
	case 0x68:
		x2 = xa - ya*(xb-xa)/(yb-ya)
		y2 = 0
		if x2 > Xm {
			break
		}
		y1 = y1 + (y2-y1)*(Xm-x1)/(x2-x1)
		x1 = Xm
		if y1 > Ym {
			x1 = x1 + (x2-x1)*(Ym-y1)/(y2-y1)
			y1 = Ym
		}
		visible++
	case 0x86:
		x1 = xa - ya*(xb-xa)/(yb-ya)
		y1 = 0
		if x1 > Xm {
			break
		}
		y2 = y1 + (y2-y1)*(Xm-x1)/(x2-x1)
		x2 = Xm
		if y2 > Ym {
			x2 = x1 + (x2-x1)*(Ym-y1)/(y2-y1)
			y2 = Ym
		}
		visible++
	case 0xA1:
		y2 = ya - xa*(yb-ya)/(xb-xa)
		x2 = 0
		if y2 < 0 {
			break
		}
		x1 = x1 - y1*(x2-x1)/(y2-y1)
		y1 = 0
		if x1 > Xm {
			y1 = y1 + (y2-y1)*(Xm-x1)/(x2-x1)
			x1 = Xm
		}
		visible++
	case 0x1A:
		y1 = ya - xa*(yb-ya)/(xb-xa)
		x1 = 0
		if y1 < 0 {
			break
		}
		x2 = x1 - y1*(x2-x1)/(y2-y1)
		y2 = 0
		if x2 > Xm {
			y2 = y1 + (y2-y1)*(Xm-x1)/(x2-x1)
			x2 = Xm
		}
		visible++
	case 0x92:
		y2 = ya + (yb-ya)*(Xm-xa)/(xb-xa)
		x2 = Xm
		if y2 < 0 {
			break
		}
		x1 = x1 - y1*(x2-x1)/(y2-y1)
		y1 = 0
		if x1 < 0 {
			y1 = y1 - x1*(y2-y1)/(x2-x1)
			x1 = 0
		}
		visible++
	case 0x29:
		y1 = ya + ((yb-ya)*(Xm-xa))/(xb-xa)
		x1 = Xm
		if y1 < 0 {
			break
		}
		x2 = x1 - y1*(x2-x1)/(y2-y1)
		y2 = 0
		if x2 < 0 {
			y2 = y1 - x1*(y2-y1)/(x2-x1)
			x2 = 0
		}
		visible++
	case 0x61:
		y2 = ya - xa*(yb-ya)/(xb-xa)
		x2 = 0
		if y2 > Ym {
			break
		}
		x1 = x1 + (x2-x1)*(Ym-y1)/(y2-y1)
		y1 = Ym
		if x1 > Xm {
			y1 = y1 + (y2-y1)*(Xm-x1)/(x2-x1)
			x1 = Xm
		}
		visible++
	case 0x16:
		y1 = ya - xa*(yb-ya)/(xb-xa)
		x1 = 0
		if y1 > Ym {
			break
		}
		x2 = x1 + (x2-x1)*(Ym-y1)/(y2-y1)
		y2 = Ym
		if x2 > Xm {
			y2 = y1 + (y2-y1)*(Xm-x1)/(x2-x1)
			x2 = Xm
		}
		visible++
	case 0x52:
		y2 = ya + (yb-ya)*(Xm-xa)/(xb-xa)
		x2 = Xm
		if y2 > Ym {
			break
		}
		x1 = x1 + (x2-x1)*(Ym-y1)/(y2-y1)
		y1 = Ym
		if x1 < 0 {
			y1 = y1 - x1*(y2-y1)/(x2-x1)
			x1 = 0
		}
		visible++
	case 0x25:
		y1 = ya + (yb-ya)*(Xm-xa)/(xb-xa)
		x1 = Xm
		if y1 > Ym {
			break
		}
		x2 = x1 + (x2-x1)*(Ym-y1)/(y2-y1)
		y2 = Ym
		if x2 < 0 {
			y2 = y1 - x1*(y2-y1)/(x2-x1)
			x2 = 0
		}
		visible++
	case 0x14:
		y1 = ya - xa*(yb-ya)/(xb-xa)
		x1 = 0
		if y1 > Ym {
			break
		}
		x2 = x1 + (Ym-y1)*(x2-x1)/(y2-y1)
		y2 = Ym
		visible++
	case 0x41:
		y2 = ya - xa*(yb-ya)/(xb-xa)
		x2 = 0
		if y2 > Ym {
			break
		}
		x1 = x1 + (x2-x1)*(Ym-y1)/(y2-y1)
		y1 = Ym
		visible++
	case 0x18:
		y1 = ya - xa*(yb-ya)/(xb-xa)
		x1 = 0
		if y1 < 0 {
			break
		}
		x2 = x1 - y1*(x2-x1)/(y2-y1)
		y2 = 0
		visible++
	case 0x81:
		y2 = ya - xa*(yb-ya)/(xb-xa)
		x2 = 0
		if y2 < 0 {
			break
		}
		x1 = xa - y1*(x2-x1)/(y2-y1)
		y1 = 0
		visible++
	case 0x24:
		y1 = ya + (yb-ya)*(Xm-xa)/(xb-xa)
		x1 = Xm
		if y1 > Ym {
			break
		}
		x2 = x1 + (x2-x1)*(Ym-y1)/(y2-y1)
		y2 = Ym
		visible++
	case 0x42:
		y2 = ya + (yb-ya)*(Xm-xa)/(xb-xa)
		x2 = Xm
		if y2 > Ym {
			break
		}
		x1 = x1 + (x2-x1)*(Ym-y1)/(y2-y1)
		y1 = Ym
		visible++
	case 0x28:
		y1 = ya + (yb-ya)*(Xm-xa)/(xb-xa)
		x1 = Xm
		if y1 < 0 {
			break
		}
		x2 = x1 - y1*(x2-x1)/(y2-y1)
		y2 = 0
		visible++
	case 0x82:
		y2 = ya + (yb-ya)*(Xm-xa)/(xb-xa)
		x2 = Xm
		if y2 < 0 {
			break
		}
		x1 = x1 - y1*(x2-x1)/(y2-y1)
		y1 = 0
		visible++
	default:
		visible = 0
	}
	if visible > 0 {
		__line(x1, y1, x2, y2, color)
	}

	return visible
}

func _line(x1, y1, x2, y2 int, color byte) {
	x := x1
	y := y1
	dx := (mlib.AbsInt(x2 - x1)) << 1
	dy := (mlib.AbsInt(y2 - y1)) << 1

	sx := mlib.SignInt(x2 - x1)
	sy := mlib.SignInt(y2 - y1)

	swp := 0
	if dy > dx {
		dy, dx = dx, dy
		swp = 1
	}
	dx_ := dx >> 1
	e := dy - dx_

	for k := 0; k < dx_; k++ {
		putPixel(x, y, color)
		for {
			if e < 0 {
				break
			}
			if swp > 0 {
				x += sx
			} else {
				y += sy
			}
			e -= dx
		}
		if swp > 0 {
			y += sy
		} else {
			x += sx
		}
		e += dy
	}
}

func __line(x1, y1, x2, y2 int, color byte) {
	if x2 == x1 {
		if y1 < y2 {
			for k := y1; k < y2+1; k++ {
				putPixel(x1, k, color)
			}
		} else {
			for k := y2; k < y1+1; k++ {
				putPixel(x1, k, color)
			}
		}
		return
	}

	if y2 == y1 {
		if x1 < x2 {
			C.memset(unsafe.Pointer(&(screen[(x1+yTable[y1])*3])), C.int(color), C.ulong((x2-x1+1)*3))
			return
		}
		C.memset(unsafe.Pointer(&(screen[(x2+yTable[y1])*3])), C.int(color), C.ulong((x1-x2+1)*3))
		return
	}
	x := x1
	y := y1

	dx := (mlib.AbsInt(x2 - x1)) << 1
	dy := (mlib.AbsInt(y2 - y1)) << 1

	if x2 > x1 {
		if y2 > y1 {
			if dy > dx {
				dy_ := dy >> 1
				e := dx - dy_
				for k := 0; k < dy_; k++ {
					putPixel(x, y, color)
					for {
						if e < 0 {
							break
						}
						x++
						e -= dy
					}
					y++
					e += dx
				}
			} else {
				dx_ := dx >> 1
				e := dy - dx_
				for k := 0; k < dx_; k++ {
					putPixel(x, y, color)
					for {
						if e < 0 {
							break
						}
						y++
						e -= dx
					}
					x++
					e += dy
				}
			}
		} else if y2 < y1 {
			if dy > dx {
				dy_ := dy >> 1
				e := dx - dy_
				for k := 0; k < dy_; k++ {
					putPixel(x, y, color)
					for {
						if e < 0 {
							break
						}
						x++
						e -= dy
					}
					y--
					e += dx
				}
			} else {
				dx_ := dx >> 1
				e := dy - dx_
				for k := 0; k < dx_; k++ {
					putPixel(x, y, color)
					for {
						if e < 0 {
							break
						}
						y--
						e -= dx
					}
					x++
					e += dy
				}
			}
		}
	} else if x2 < x1 {
		if y2 > y1 {
			if dy > dx {
				dy_ := dy >> 1
				e := dx - dy_
				for k := 0; k < dy_; k++ {
					putPixel(x, y, color)
					for {
						if e < 0 {
							break
						}
						x++
						e -= dy
					}
					y++
					e += dx
				}
			} else {
				dx_ := dx >> 1
				e := dy - dx_
				for k := 0; k < dx_; k++ {
					putPixel(x, y, color)
					for {
						if e < 0 {
							break
						}
						y++
						e -= dx
					}
					x--
					e += dy
				}
			}
		} else if y2 < y1 {
			if dy > dx {
				dy_ := dy >> 1
				e := dx - dy_
				for k := 0; k < dy_; k++ {
					putPixel(x, y, color)
					for {
						if e < 0 {
							break
						}
						x++
						e -= dy
					}
					y--
					e += dx
				}
			} else {
				dx_ := dx >> 1
				e := dy - dx_
				for k := 0; k < dx_; k++ {
					putPixel(x, y, color)
					for {
						if e < 0 {
							break
						}
						y--
						e -= dx
					}
					x--
					e += dy
				}
			}
		}
	}
}
