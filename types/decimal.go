package types

import (
	"fmt"
	"strings"
)

type Decimal struct {
	negative bool
	digits   []uint8
	scale    int // number of digits to the left of the decimal point
}

func DecimalZero() Decimal {
	return Decimal{}
}

func NewDecimal(input string) Decimal {
	result, err := ParseDecimal(input)
	if err != nil {
		return Decimal{}
	}
	return result
}

func ParseDecimal(input string) (Decimal, error) {
	runes := []rune(input)
	negative := false
	if len(runes) > 0 && runes[0] == '-' {
		negative = true
		runes = runes[1:]
	}
	if len(runes) == 0 {
		return Decimal{}, fmt.Errorf("empty input")
	}
	var n int
	var digits []uint8
	decimalFound := false

	for i, r := range runes {
		switch {
		case r >= '0' && r <= '9':
			digits = append(digits, uint8(r-'0'))
		case r == '.':
			if decimalFound {
				return Decimal{}, fmt.Errorf("multiple decimal points")
			}
			decimalFound = true
			n = i
		default:
			return Decimal{}, fmt.Errorf("invalid number format: %s", input)
		}
	}
	if !decimalFound {
		n = len(runes)
	}
	return normalise(negative, digits, n), nil
}

func (d Decimal) Type() Type {
	return TypeDecimal
}

func (d Decimal) Compare(next Value) Comparison {
	nextDecimal, ok := next.(Decimal)
	if !ok {
		return ComparisonIncomparable
	}
	gt, lt := ComparisonGreater, ComparisonLess

	switch {
	case d.negative && nextDecimal.negative:
		gt, lt = lt, gt
	case !d.negative && nextDecimal.negative:
		return gt
	case d.negative && !nextDecimal.negative:
		return lt
	}

	switch {
	case d.scale < nextDecimal.scale:
		return lt
	case d.scale > nextDecimal.scale:
		return gt
	}

	minDigits := len(d.digits)
	if len(nextDecimal.digits) < minDigits {
		minDigits = len(nextDecimal.digits)
	}

	for i := 0; i < minDigits; i++ {
		if d.digits[i] < nextDecimal.digits[i] {
			return lt
		}
		if d.digits[i] > nextDecimal.digits[i] {
			return gt
		}
	}

	switch {
	case len(d.digits) < len(nextDecimal.digits):
		return lt
	case len(d.digits) > len(nextDecimal.digits):
		return gt
	}
	return ComparisonEqual
}

func (d Decimal) String() string {
	builder := new(strings.Builder)
	if d.negative {
		fmt.Fprintf(builder, "-")
	}
	if d.scale == 0 {
		fmt.Fprintf(builder, "0")
	}

	for i, digit := range d.digits {
		if i == d.scale {
			fmt.Fprintf(builder, ".")
		}
		fmt.Fprintf(builder, "%d", digit)
	}
	return builder.String()
}

func normalise(negative bool, digits []uint8, scale int) Decimal {

	for scale > 0 && len(digits) > 0 && digits[0] == 0 {
		scale = scale - 1
		digits = digits[1:]
	}

	for len(digits) > scale && digits[len(digits)-1] == 0 {
		digits = digits[:len(digits)-1]
	}
	if len(digits) == 0 {
		return DecimalZero()
	}
	return Decimal{negative: negative, digits: digits, scale: scale}
}
