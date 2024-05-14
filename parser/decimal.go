package parser

import (
	"fmt"
	"github.com/Vignesh-Rajarajan/go-db/sql"
	"math"
	"strings"
)

// ParseDecimal parses a decimal from a string
func ParseDecimal(input string) (sql.Decimal, error) {
	if input == "" {
		return sql.Decimal{}, nil
	}

	before, after, _ := strings.Cut(input, ".")

	x, _, isNeg, ok := ParseInteger(before)
	if !ok {
		return sql.Decimal{}, fmt.Errorf("could not parse decimal %q", input)
	}

	after = strings.TrimRight(after, "0")

	y, digits, negativeDecimal, ok := ParseInteger(after)
	if !ok || negativeDecimal {
		return sql.Decimal{}, fmt.Errorf("could not parse decimal, not a valid number %q", input)
	}

	value := shift(x, digits) + y
	if isNeg {
		value = -value
	}

	return sql.Decimal{Value: value, Digits: digits}, nil

}

func shift(x int, digits int) int {
	return x * int(math.Pow10(digits))
}

// ParseInteger parses an integer from a string
func ParseInteger(before string) (result, digits int, isNeg, ok bool) {
	runes := []rune(before)
	if len(runes) > 0 && runes[0] == '-' {
		isNeg = true
		runes = runes[1:]
	}

	for i, digit := range runes {
		if digit < '0' || digit > '9' {
			return 0, 0, false, false
		}
		result = result*10 + int(digit-'0')
		digits = i + 1
	}
	ok = true
	return
}
