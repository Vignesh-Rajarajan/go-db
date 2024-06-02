package parser

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
