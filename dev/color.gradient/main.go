package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello, playground")
	//r2 := 0x42
	g2 := 0x82
	//b2 := 0x6c

	//r1 := 0x38
	g1 := 0x70
	//b1 := 0x5d

	dg := (g2 - g1) / 10
	//dr := int(float64(0x38) / float64(0x70) * (float64(g2) - float64(g1)) / 10)
	//db := int(float64(0x5d) / float64(0x70) * (float64(g2) - float64(g1)) / 10)
	for i := 0; i < 100; i++ {
		g := g1 - dg*i
		r := int(float64(0x38) / float64(0x70) * (-(float64(g2)-float64(g1))/10*float64(i) + float64(g1)))
		b := int(float64(0x2d) / float64(0x70) * (-(float64(g2)-float64(g1))/10*float64(i) + float64(g1)))
		//fmt.Println(i, ":", fmt.Sprintf("0x%02X", r1+dr*i), fmt.Sprintf("0x%02X", g1+dg*i), fmt.Sprintf("0x%02X", b1+db*i))
		if i%10 == 0 || i == 99 {
			fmt.Println(i, ":", fmt.Sprintf("0x%02X", r), fmt.Sprintf("0x%02X", g), fmt.Sprintf("0x%02X", b))
		}
	}
}
