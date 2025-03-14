package main

import (
	"fmt"
	"math"
	"os"
	"runtime"
	"strings"
	"time"
)

func main() {
	alpha := math.Pi / 4
	s := float64(3)
	x1 := float64(1)
	y1 := float64(1)
	xo := float64(2)
	yo := float64(3)
	r := float64(1.5)
	//x := float64(2)
	//y := float64(2)

	ta := math.Tan(alpha)
	a := ta*ta + 1
	if a == 0 {
		log(ta*ta+1, "=", 0)
		return
	}
	c1 := y1 - yo - ta*x1
	b := 2*ta*c1 - 2*xo
	c := c1*c1 + xo*xo - r*r
	ds := b*b - 4*a*c
	if ds < 0 {
		log(ds, "<", 0, "b=", b, "4*a*c=", 4*a*c)
		return
	}
	ds = math.Sqrt(ds)

	xx1 := (-b + ds) / 2 / a
	xx2 := (-b - ds) / 2 / a
	log(xx1, xx2, ta*(xx1-x1)+y1, ta*(xx2-x1)+y1)
	if xx1 < x1 || xx1 > x1+s {
		if xx2 < x1 || xx2 > x1+s {
			log(xx2, "<", x1, "||", xx2, ">", x1+s)
		} else {
			log("INTERSECTS! at", xx2, ta*(xx2-x1)+y1)
		}
	} else {
		log("INTERSECTS! at", xx1, ta*(xx1-x1)+y1)
	}
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
		return ""
	}

	return pwd
}
