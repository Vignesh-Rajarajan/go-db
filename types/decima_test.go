package types

import (
	"testing"
)

func TestParseDecimal(t *testing.T) {
	cases := []struct {
		input    string
		expected Decimal
		wantErr  bool
	}{
		{
			"0",
			Decimal{},
			false,
		},
		{
			"000",
			Decimal{},
			false,
		},
		{
			"123",
			Decimal{digits: []uint8{1, 2, 3}, scale: 3},
			false,
		},
		{
			"000123",
			Decimal{digits: []uint8{1, 2, 3}, scale: 3},
			false,
		},
		{
			"123.456",
			Decimal{digits: []uint8{1, 2, 3, 4, 5, 6}, scale: 3},
			false,
		},
		{
			"123.456000",
			Decimal{digits: []uint8{1, 2, 3, 4, 5, 6}, scale: 3},
			false,
		},
		{
			"-123",
			Decimal{negative: true, digits: []uint8{1, 2, 3}, scale: 3},
			false,
		},
		{
			"-123.456",
			Decimal{negative: true, digits: []uint8{1, 2, 3, 4, 5, 6}, scale: 3},
			false,
		}, {
			"123.456.789",
			Decimal{},
			true,
		},
		{
			"123.",
			Decimal{digits: []uint8{1, 2, 3}, scale: 3},
			false,
		},
		{
			".456",
			Decimal{digits: []uint8{4, 5, 6}, scale: 0},
			false,
		},
		{
			"123.000456",
			Decimal{digits: []uint8{1, 2, 3, 0, 0, 0, 4, 5, 6}, scale: 3},
			false,
		},
		{
			".000456",
			Decimal{digits: []uint8{0, 0, 0, 4, 5, 6}, scale: 0},
			false,
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			d, err := ParseDecimal(c.input)
			if c.wantErr && err == nil {
				t.Fatalf("ParseDecimal(%q) returned nil, want error", c.input)
			}
			if !c.wantErr && err != nil {
				t.Fatalf("ParseDecimal(%q) returned error, want nil: %v", c.input, err)
			}

			//if !reflect.DeepEqual(d, c.expected) {
			if d.String() != c.expected.String() {
				t.Fatalf("ParseDecimal(%q) = %v, want %v", c.input, d, c.expected)
			}
		})
	}
}
