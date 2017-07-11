/*

  Morton implements Z-Order Curve encoding and decoding for N-dimensions, using lookup tables and magic bits respectively.

  In order to supply for N-dimensions, this library generates the magic bits used in decoding.  While this library does supply for N-dimensions, because this type of ordering uses bit interleaving for encoding it is limited by the width of the uint64 type divided by the number of dimensions (i.e., uint64/3 for 3 dimensions).

*/
package morton

import (
	"errors"
	"fmt"
	"sort"
)

type Table struct {
	Index  uint8
	Length uint32
	Encode []Bit
}

// Sortable Table slice type to satisfy the sort package interface
type ByTable []Table

func (t ByTable) Len() int {
	return len(t)
}

func (t ByTable) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t ByTable) Less(i, j int) bool {
	return t[i].Index < t[j].Index
}

type Bit struct {
	Index uint32
	Value uint64
}

// Sortable Table slice type to satisfy the sort package interface
type ByBit []Bit

func (b ByBit) Len() int {
	return len(b)
}

func (b ByBit) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b ByBit) Less(i, j int) bool {
	return b[i].Index < b[j].Index
}

type Morton struct {
	Dimensions uint8
	Tables     []Table
	Magic      []uint64
}

// Convenience function
func NewMorton(dimensions uint8, size uint32) *Morton {
	m := new(Morton)
	m.Create(dimensions, size)
	return m
}

func (m *Morton) Create(dimensions uint8, size uint32) {
	table_ops := make(chan bool)
	go m.CreateTables(dimensions, size, table_ops)
	<-table_ops
}

func (m *Morton) CreateTables(dimensions uint8, length uint32, done chan<- bool) {
	ch := make(chan *Table)

	m.Dimensions = dimensions

	for i := uint8(0); i < dimensions; i++ {
		go createTable(i, dimensions, length, ch)
	}

	mch := make(chan []uint64)

	go makeMagic(dimensions, mch)

	for i := uint8(0); i < dimensions; i++ {
		t := <-ch
		m.Tables = append(m.Tables, *t)
	}

	close(ch)

	sort.Sort(ByTable(m.Tables))

	m.Magic = <-mch

	done <- true
}

func makeMagic(dimensions uint8, mch chan<- []uint64) {
	// Generate nth and ith bits variables
	d := uint64(dimensions)
	limit := 64/d + 1
	nth := []uint64{0, 0, 0, 0, 0}
	for i := uint64(0); i < 64; i++ {
		if i < limit {
			nth[0] |= 1 << (i * (d))
		}

		nth[1] |= 3 << (i * (d << 1))
		nth[2] |= 0xf << (i * (d << 2))
		nth[3] |= 0xff << (i * (d << 3))
		nth[4] |= 0xffff << (i * (d << 4))
	}

	mch <- nth
	close(mch)
}

func (m *Morton) Decode(code uint64) (result []uint32) {
	if m.Dimensions == 0 {
		return
	}

	d := uint64(m.Dimensions)
	r := make([]uint64, d)

	// Process each dimension
	for i := uint64(0); i < d; i++ {
		r[i] = (code >> i) & m.Magic[0]

		r[i] = (r[i] ^ (r[i] >> (1 << (d - 2)))) & m.Magic[1]
		r[i] = (r[i] ^ (r[i] >> (2 << (d - 2)))) & m.Magic[2]
		r[i] = (r[i] ^ (r[i] >> (4 << (d - 2)))) & m.Magic[3]
		r[i] = (r[i] ^ (r[i] >> (8 << (d - 2)))) & m.Magic[4]

		result = append(result, uint32(r[i]))
	}

	return
}

func (m *Morton) Encode(vector []uint32) (result uint64, err error) {
	length := len(m.Tables)
	if length == 0 {
		err = errors.New("No lookup tables.  Please generate them via CreateTables().")
		return
	}

	if len(vector) > length {
		err = errors.New("Input vector slice length exceeds the number of lookup tables.  Please regenerate them via CreateTables()")
		return
	}

	//sort.Sort(sort.Reverse(ByUint32Index(vector)))

	for k, v := range vector {
		if v > uint32(len(m.Tables[k].Encode)-1) {
			err = errors.New(fmt.Sprint("Input vector component, ", k, " length exceeds the corresponding lookup table's size.  Please regenerate them via CreateTables() and specify the appropriate table length"))
			return
		}

		result |= m.Tables[k].Encode[v].Value
	}

	return
}

func createTable(index, dimensions uint8, length uint32, ch chan<- *Table) {
	t := &Table{Index: index, Length: length}
	bch := make(chan *Bit)

	// Build interleave queue
	for i := uint32(0); i < length; i++ {
		go interleaveBits(i, uint32(index), uint32(dimensions-1), bch)
	}

	// Pull from interleave queue
	for i := uint32(0); i < length; i++ {
		ib := <-bch
		t.Encode = append(t.Encode, *ib)
	}

	close(bch)

	sort.Sort(ByBit(t.Encode))

	ch <- t
}

func interleaveBits(value, offset, spread uint32, bch chan<- *Bit) {
	ib := &Bit{value, 0}

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
		ib.Value |= (v & (1 << i)) << (i * s)
	}
	ib.Value = ib.Value << o

	bch <- ib
}
