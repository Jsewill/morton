/*

Morton library example

*/
package main

import (
	".."
	"fmt"
)

func main() {
	//Create a new Morton
	m := new(morton.Morton)
	//Generate Tables and Magic bits
	m.Create(3, 512)

	//Create arbitrary coordinates
	c := []uint32{5, 9, 1}

	//Encode aforementioned coordinates
	e, err := m.Encode(c)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Coordinates: ", c)
	fmt.Println("Encoded Coordinates: ", e)
	fmt.Println("Decoded Coordinates: ", m.Decode(e))
}
