/*

Morton library example

*/
package main

import (
	"fmt"

	"github.com/Jsewill/morton"
)

func main() {
	//Create a new Morton
	m := new(morton.Morton)
	//Generate Tables and Magic bits
	m.Create(4, 512)

	//Create arbitrary coordinates
	c := []uint32{511, 472, 103, 7}

	//Encode aforementioned coordinates
	e, err := m.Encode(c)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Coordinates: %v\n", c)
	fmt.Printf("Decoded Coordinates: %v\n", m.Decode(e))
}
