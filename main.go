package main

import (
	"fmt"

	xml "scripts/xml"
)

func main() {
	filename := "./test.xml"
	decoder, err := xml.Decode(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	// _ = decoder
	// bytes, _ := json.Marshal(decoder)
	// fmt.Println(string(bytes))

	root, err := xml.Parse(decoder)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(root)
}
