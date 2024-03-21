package parser

import "unicode"

func isQuote(r rune) bool {
	return r == '"' || r == '\''
}

func isDigitOrDot(r rune) bool {
	return isDigit(r) || r == '.'
}

func isDigit(r rune) bool {
	return unicode.IsDigit(r)
}

func isWordStartWithCharacter(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
}

func isWordCharacter(r rune) bool {
	return isWordStartWithCharacter(r) || isDigit(r) || r == '$'
}

func isWhitespace(r rune) bool {
	switch r {
	case ' ', '\t', '\n', '\r':
		return true
	}
	return false
}

func isPunctuation(r rune) bool {
	switch r {
	case ',', '.', ';', '(', ')', '=', '<', '>', '!', '*':
		return true
	}
	return false
}
