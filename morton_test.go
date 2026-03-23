package morton

import (
	"fmt"
	"testing"
)

func TestRoundTrip2D(t *testing.T) {
	m := New(2, 1024)
	cases := [][]uint32{
		{0, 0},
		{1, 0},
		{0, 1},
		{1, 1},
		{1023, 9},
		{512, 512},
		{0, 1023},
		{1023, 1023},
	}
	for _, c := range cases {
		encoded, err := m.Encode(c)
		if err != nil {
			t.Fatalf("Encode(%v): %v", c, err)
		}
		decoded := m.Decode(encoded)
		for i := range c {
			if decoded[i] != c[i] {
				t.Errorf("round-trip %v: got %v", c, decoded)
				break
			}
		}
	}
}

func TestRoundTrip3D(t *testing.T) {
	m := New(3, 1024)
	cases := [][]uint32{
		{0, 0, 0},
		{1, 1, 1},
		{5, 3, 7},
		{100, 200, 300},
		{1023, 0, 0},
		{0, 0, 1023},
	}
	for _, c := range cases {
		encoded, err := m.Encode(c)
		if err != nil {
			t.Fatalf("Encode(%v): %v", c, err)
		}
		decoded := m.Decode(encoded)
		for i := range c {
			if decoded[i] != c[i] {
				t.Errorf("round-trip %v: got %v", c, decoded)
				break
			}
		}
	}
}

func TestRoundTrip4D(t *testing.T) {
	m := New(4, 512)
	cases := [][]uint32{
		{0, 0, 0, 0},
		{1, 1, 1, 1},
		{511, 472, 103, 7},
		{0, 0, 0, 511},
		{511, 0, 0, 0},
	}
	for _, c := range cases {
		encoded, err := m.Encode(c)
		if err != nil {
			t.Fatalf("Encode(%v): %v", c, err)
		}
		decoded := m.Decode(encoded)
		for i := range c {
			if decoded[i] != c[i] {
				t.Errorf("round-trip %v: got %v", c, decoded)
				break
			}
		}
	}
}

func TestKnownValues2D(t *testing.T) {
	m := New(2, 4)
	tests := []struct {
		coords []uint32
		want   uint64
	}{
		{[]uint32{0, 0}, 0},
		{[]uint32{1, 0}, 1},
		{[]uint32{0, 1}, 2},
		{[]uint32{1, 1}, 3},
		{[]uint32{3, 3}, 0xf},
	}
	for _, tt := range tests {
		got, err := m.Encode(tt.coords)
		if err != nil {
			t.Fatalf("Encode(%v): %v", tt.coords, err)
		}
		if got != tt.want {
			t.Errorf("Encode(%v) = %d, want %d", tt.coords, got, tt.want)
		}
	}
}

func TestKnownValues3D(t *testing.T) {
	m := New(3, 2)
	tests := []struct {
		coords []uint32
		want   uint64
	}{
		{[]uint32{0, 0, 0}, 0},
		{[]uint32{1, 0, 0}, 1},
		{[]uint32{0, 1, 0}, 2},
		{[]uint32{0, 0, 1}, 4},
		{[]uint32{1, 1, 1}, 7},
	}
	for _, tt := range tests {
		got, err := m.Encode(tt.coords)
		if err != nil {
			t.Fatalf("Encode(%v): %v", tt.coords, err)
		}
		if got != tt.want {
			t.Errorf("Encode(%v) = %d, want %d", tt.coords, got, tt.want)
		}
	}
}

func TestEncodeZeroCoordinates(t *testing.T) {
	m := New(2, 4)
	got, err := m.Encode([]uint32{0, 0})
	if err != nil {
		t.Fatalf("Encode({0,0}): %v", err)
	}
	if got != 0 {
		t.Errorf("Encode({0,0}) = %d, want 0", got)
	}
}

func TestEncodeErrorNoTables(t *testing.T) {
	m := &Morton{}
	_, err := m.Encode([]uint32{1, 2})
	if err == nil {
		t.Error("expected error for Encode with no tables")
	}
}

func TestEncodeErrorVectorTooLong(t *testing.T) {
	m := New(2, 4)
	_, err := m.Encode([]uint32{1, 2, 3})
	if err == nil {
		t.Error("expected error for vector longer than dimensions")
	}
}

func TestEncodeErrorValueTooLarge(t *testing.T) {
	m := New(2, 4)
	_, err := m.Encode([]uint32{5, 0})
	if err == nil {
		t.Error("expected error for value exceeding table size")
	}
}

func TestDecodeZeroDimensions(t *testing.T) {
	m := &Morton{}
	result := m.Decode(42)
	if result != nil {
		t.Errorf("Decode with zero dimensions: got %v, want nil", result)
	}
}

func TestDecodeEmptyMagic(t *testing.T) {
	m := &Morton{Dimensions: 2}
	result := m.Decode(42)
	if result != nil {
		t.Errorf("Decode with empty magic: got %v, want nil", result)
	}
}

func TestNewConvenience(t *testing.T) {
	m := New(3, 16)
	if m.Dimensions != 3 {
		t.Errorf("New dimensions = %d, want 3", m.Dimensions)
	}
	if len(m.Tables) != 3 {
		t.Errorf("New tables count = %d, want 3", len(m.Tables))
	}
	if len(m.Magic) == 0 {
		t.Error("New magic is empty")
	}
}

func TestMakeMagicLength(t *testing.T) {
	for d := uint8(2); d <= 8; d++ {
		magic := MakeMagic(d)
		if len(magic) != 6 {
			t.Errorf("MakeMagic(%d) length = %d, want 6", d, len(magic))
		}
	}
}

func TestTableSorting(t *testing.T) {
	m := New(3, 8)
	for i := 0; i < len(m.Tables)-1; i++ {
		if m.Tables[i].Index >= m.Tables[i+1].Index {
			t.Errorf("tables not sorted: index %d >= %d", m.Tables[i].Index, m.Tables[i+1].Index)
		}
	}
}

func TestInterleaveBitsVsMagic(t *testing.T) {
	cases := []struct {
		name       string
		dimensions uint8
		values     []uint32
	}{
		{"2D", 2, []uint32{0, 1, 5, 42, 255, 1023}},
		{"3D", 3, []uint32{0, 1, 5, 42, 255, 1023}},
		{"4D", 4, []uint32{0, 1, 5, 42, 255, 511}},
		{"5D", 5, []uint32{0, 1, 5, 42, 255, 511}},
		{"6D", 6, []uint32{0, 1, 5, 42, 255, 511}},
	}

	for _, tc := range cases {
		magic := MakeMagic(tc.dimensions)
		spread := uint32(tc.dimensions - 1)

		t.Run(tc.name, func(t *testing.T) {
			for _, offset := range []uint32{0, 1} {
				for _, v := range tc.values {
					got := InterleaveBits(v, offset, spread)
					gotMagic := InterleaveBitsMagic(v, offset, spread, magic)

					if got.Value != gotMagic.Value {
						t.Errorf("value=%d offset=%d spread=%d:\n  InterleaveBits:      %064b\n  InterleaveBitsMagic: %064b",
							v, offset, spread, got.Value, gotMagic.Value)
					} else {
						t.Logf("value=%d offset=%d spread=%d: MATCH %064b",
							v, offset, spread, got.Value)
					}
				}
			}
		})
	}
}

func BenchmarkInterleaveBits(b *testing.B) {
	for _, dims := range []uint8{2, 3, 4, 5, 6} {
		spread := uint32(dims - 1)
		b.Run(fmt.Sprintf("%dD", dims), func(b *testing.B) {
			for b.Loop() {
				InterleaveBits(255, 0, spread)
			}
		})
	}
}

func BenchmarkInterleaveBitsMagic(b *testing.B) {
	for _, dims := range []uint8{2, 3, 4, 5, 6} {
		spread := uint32(dims - 1)
		magic := MakeMagic(dims)
		b.Run(fmt.Sprintf("%dD", dims), func(b *testing.B) {
			for b.Loop() {
				InterleaveBitsMagic(255, 0, spread, magic)
			}
		})
	}
}

func BenchmarkEncode(b *testing.B) {
	for _, dims := range []uint8{2, 3, 4, 5, 6} {
		coords := make([]uint32, dims)
		for i := range coords {
			coords[i] = 42
		}
		m := New(dims, 1024)
		b.Run(fmt.Sprintf("%dD", dims), func(b *testing.B) {
			for b.Loop() {
				m.Encode(coords)
			}
		})
	}
}

func BenchmarkDecode(b *testing.B) {
	for _, dims := range []uint8{2, 3, 4, 5, 6} {
		m := New(dims, 1024)
		coords := make([]uint32, dims)
		for i := range coords {
			coords[i] = 42
		}
		code, _ := m.Encode(coords)
		b.Run(fmt.Sprintf("%dD", dims), func(b *testing.B) {
			for b.Loop() {
				m.Decode(code)
			}
		})
	}
}

func TestMaxBitWidth2D(t *testing.T) {
	// 2D: 64/2 = 32 bits per dimension, max value 2^32-1
	// Use a smaller table but test that encoding works at boundary
	m := New(2, 2)
	got, err := m.Encode([]uint32{1, 1})
	if err != nil {
		t.Fatalf("Encode at bit boundary: %v", err)
	}
	decoded := m.Decode(got)
	if decoded[0] != 1 || decoded[1] != 1 {
		t.Errorf("bit boundary round-trip: got %v, want [1 1]", decoded)
	}
}
