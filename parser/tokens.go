package parser

import "fmt"

type TokenType int

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

type Token struct {
	Type  TokenType
	Value string
}

func (t Token) String() string {
	if t.Type >= TokenTypeComma {
		return t.Type.String()
	}
	return fmt.Sprintf("{%s %q}", t.Type, t.Value)
}

func Tokenize(input string) ([]Token, error) {
	return NewLexer(input).Lex()
}
