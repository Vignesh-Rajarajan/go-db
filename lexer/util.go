package lexer

import (
	"fmt"
	"unicode"
)

const (
	TokenTypeIdentifier TokenType = iota
	TokenTypeString
	TokenTypeNumber

	TokenTypeComma
	TokenTypeDot
	TokenTypeStar
	TokenTypeSemicolon
	TokenTypeOpenParen
	TokenTypeCloseParen
	TokenTypeEq
	TokenTypeNotEq
	TokenTypeLt
	TokenTypeLte
	TokenTypeGt
	TokenTypeGte

	TokenTypeSelect
	TokenTypeFrom
	TokenTypeWhere
	TokenTypeAnd
	TokenTypeOr
	TokenTypeNot
	TokenTypeIn
	TokenTypeIs
	TokenTypeNull
	TokenTypeLeft
	TokenTypeRight
	TokenTypeInner
	TokenTypeOuter
	TokenTypeJoin
)

func (t TokenType) String() string {
	switch t {
	case TokenTypeIdentifier:
		return "Identifier"
	case TokenTypeString:
		return "String"
	case TokenTypeNumber:
		return "Number"
	case TokenTypeComma:
		return "Comma"
	case TokenTypeDot:
		return "Dot"
	case TokenTypeStar:
		return "Star"
	case TokenTypeSemicolon:
		return "Semicolon"
	case TokenTypeOpenParen:
		return "OpenParen"
	case TokenTypeCloseParen:
		return "CloseParen"
	case TokenTypeEq:
		return "Eq"
	case TokenTypeNotEq:
		return "NotEq"
	case TokenTypeLt:
		return "Lt"
	case TokenTypeLte:
		return "Lte"
	case TokenTypeGt:
		return "Gt"
	case TokenTypeGte:
		return "Gte"
	case TokenTypeSelect:
		return "Select"
	case TokenTypeFrom:
		return "From"
	case TokenTypeWhere:
		return "Where"
	case TokenTypeAnd:
		return "And"
	case TokenTypeOr:
		return "Or"
	case TokenTypeNot:
		return "Not"
	case TokenTypeIn:
		return "In"
	case TokenTypeIs:
		return "Is"
	case TokenTypeNull:
		return "Null"
	case TokenTypeLeft:
		return "Left"
	case TokenTypeRight:
		return "Right"
	case TokenTypeInner:
		return "Inner"
	case TokenTypeOuter:
		return "Outer"
	case TokenTypeJoin:
		return "Join"
	}
	return fmt.Sprintf("Unknown token type %d", t)
}

var KeywordMap = map[string]TokenType{
	"select": TokenTypeSelect,
	"from":   TokenTypeFrom,
	"where":  TokenTypeWhere,
	"and":    TokenTypeAnd,
	"or":     TokenTypeOr,
	"not":    TokenTypeNot,
	"in":     TokenTypeIn,
	"is":     TokenTypeIs,
	"null":   TokenTypeNull,
	"left":   TokenTypeLeft,
	"right":  TokenTypeRight,
	"inner":  TokenTypeInner,
	"outer":  TokenTypeOuter,
	"join":   TokenTypeJoin,
}

var SymbolMap = map[string]TokenType{
	",":  TokenTypeComma,
	".":  TokenTypeDot,
	"*":  TokenTypeStar,
	";":  TokenTypeSemicolon,
	"(":  TokenTypeOpenParen,
	")":  TokenTypeCloseParen,
	"=":  TokenTypeEq,
	"!=": TokenTypeNotEq,
	"<>": TokenTypeNotEq,
	"<":  TokenTypeLt,
	"<=": TokenTypeLte,
	">":  TokenTypeGt,
	">=": TokenTypeGte,
}

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
