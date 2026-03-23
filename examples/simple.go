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
	m.Create(2, 1024)

	//Create arbitrary coordinates
	c := []uint32{1023, 9}

	//Encode aforementioned coordinates
	e, err := m.Encode(c)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Coordinates: %v\n", c)
	fmt.Printf("Encoded Coordinates: %v\n", e)
	fmt.Printf("Decoded Coordinates: %v\n", m.Decode(e))

	// 4D example
	m4 := morton.New(4, 512)
	c4 := []uint32{511, 472, 103, 7}
	e4, err := m4.Encode(c4)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("\n4D Coordinates: %v\n", c4)
	fmt.Printf("4D Encoded Coordinates: %v\n", e4)
	fmt.Printf("4D Decoded Coordinates: %v\n", m4.Decode(e4))
}
