/*
Morton implements Z-Order Curve encoding and decoding for N-dimensions, using lookup tables and magic bits respectively.

In order to supply for N-dimensions, this library generates the magic bits used in decoding.  While this library does supply for N-dimensions, because this type of ordering uses bit interleaving for encoding it is limited by the width of the uint64 type divided by the number of dimensions (i.e., uint64/3 for 3 dimensions).
*/
package morton

import (
	"errors"
	"fmt"
)

// A type for working with Morton lookup tables, and subsequent encoding and decoding.
type Morton struct {
	Dimensions uint8
	Table      []uint64 // flat lookup: Table[dim*tableSize + val]
	Magic      []uint64
	tableSize  uint32
}

// Convenience function for creating a new Morton.
func New(dimensions uint8, size uint32) *Morton {
	m := new(Morton)
	m.Create(dimensions, size)
	return m
}

// Creates lookup table and magic bits.
func (m *Morton) Create(dimensions uint8, size uint32) {
	m.CreateTable(dimensions, size)
}

// Creates the flat lookup table.
func (m *Morton) CreateTable(dimensions uint8, length uint32) {
	m.Dimensions = dimensions
	m.Magic = MakeMagic(dimensions)
	m.tableSize = length
	m.Table = make([]uint64, uint32(dimensions)*length)
	for i := uint8(0); i < dimensions; i++ {
		start := uint32(i) * length
		for j := uint32(0); j < length; j++ {
			m.Table[start+j] = InterleaveBitsMagic(j, uint32(i), uint32(dimensions-1), m.Magic)
		}
	}
}

// Makes magic bits.
func MakeMagic(dimensions uint8) []uint64 {
	// Generate nth and ith bits variables
	d := uint64(dimensions)
	limit := 64/d + 1
	nth := []uint64{0, 0, 0, 0, 0, 0}
	for i := uint64(0); i < limit; i++ {

		switch {
		case i <= 32:
			//32
			nth[0] |= 1 << (i * d)
			fallthrough
		case i <= 16:
			//16
			nth[1] |= 3 << (i * (d << 1))
			fallthrough
		case i <= 8:
			//8
			nth[2] |= 0xf << (i * (d << 2))
			fallthrough
		case i <= 4:
			//4
			nth[3] |= 0xff << (i * (d << 3))
			fallthrough
		case i <= 2:
			//2
			nth[4] |= 0xffff << (i * (d << 4))
			fallthrough
		case i <= 1:
			//1
			nth[5] |= 0xffffff << (i * (d << 5))
		}
	}

	return nth
}

// Decodes a Morton number.
func (m *Morton) Decode(code uint64) []uint32 {
	if m.Dimensions == 0 || len(m.Magic) == 0 {
		return nil
	}
	result := make([]uint32, m.Dimensions)
	m.DecodeInto(code, result)
	return result
}

// Decodes a Morton number into a caller-provided buffer (zero allocations).
func (m *Morton) DecodeInto(code uint64, result []uint32) {
	if m.Dimensions == 0 || len(m.Magic) == 0 {
		return
	}
	d := uint64(m.Dimensions)
	for i := uint64(0); i < d; i++ {
		r := (code >> i) & m.Magic[0]
		for j := 0; j < len(m.Magic)-1; j++ {
			r = (r ^ (r >> ((d - 1) * (1 << uint64(j))))) & m.Magic[j+1]
		}
		result[i] = uint32(r)
	}
}

// Encodes a Morton number via lookup tables.
func (m *Morton) Encode(vector []uint32) (result uint64, err error) {
	if m.tableSize == 0 {
		err = errors.New("no lookup table, please generate one via CreateTable()")
		return
	}

	if len(vector) > int(m.Dimensions) {
		err = errors.New("input vector slice length exceeds the number of dimensions, please regenerate via CreateTable()")
		return
	}

	for k, v := range vector {
		if v >= m.tableSize {
			err = errors.New(fmt.Sprint("input vector component, ", k, " length exceeds the lookup table's size, please regenerate via CreateTable() and specify the appropriate table length"))
			return
		}

		result |= m.Table[uint32(k)*m.tableSize+v]
	}

	return
}

// Interleave bits of a uint32.
func InterleaveBits(value, offset, spread uint32) uint64 {
	var result uint64

	// Determine the minimum number of single shifts required. There's likely a better, and more efficient, way to do this.
	n := value
	limit := uint64(0)
	for i := uint32(0); n != 0; i++ {
		n = n >> 1
		limit++
	}

	// Offset value for interleaving and reconcile types
	v, o, s := uint64(value), uint64(offset), uint64(spread)
	for i := uint64(0); i < limit; i++ {
		// Interleave bits, bit by bit.
		result |= (v & (1 << i)) << (i * s)
	}
	result = result << o

	return result
}

// Interleave bits of a uint32 by magic.
func InterleaveBitsMagic(value, offset, spread uint32, magic []uint64) uint64 {
	v, o, s := uint64(value)&magic[len(magic)-1], uint64(offset), uint64(spread)
	for i := len(magic) - 2; i >= 0; i-- {
		j := uint64(i)
		v = (v ^ (v << (s * (1 << j)))) & magic[j]
	}
	return v << o
}
