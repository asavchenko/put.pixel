package main

import (
	"fmt"
	"os"

	"assa.com/put.pixel/lib/bitreader"
)

func main() {
	pngFile, err := os.Open("sample.png")
	if err != nil {
		panic(err)
	}

	defer pngFile.Close()

	r := bitreader.GetNew(pngFile)
	header, err := r.GetBytes(8)
	if err != nil {
		panic(err)
	}

	fmt.Println(printBytes(header))
}

func printBytes(bytes []byte) string {
	msg := ""
	for _, b := range bytes {
		msg += fmt.Sprintf("0x%02X ", b)
	}

	return msg
}
