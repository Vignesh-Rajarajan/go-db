package parser

import (
	"reflect"
	"testing"
)

func TestLexer(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  []Token
	}{
		{
			name:  "runs tests successfully with empty input",
			input: "",
			want:  nil,
		},
		{
			name:  "runs tests successfully with multiple tokens",
			input: "select foo, bar, baz from temptable;",
			want: []Token{
				{Type: TokenTypeSelect, Value: "select"},
				{Type: TokenTypeIdentifier, Value: "foo"},
				{Type: TokenTypeComma, Value: ","},
				{Type: TokenTypeIdentifier, Value: "bar"},
				{Type: TokenTypeComma, Value: ","},
				{Type: TokenTypeIdentifier, Value: "baz"},
				{Type: TokenTypeFrom, Value: "from"},
				{Type: TokenTypeIdentifier, Value: "temptable"},
				{Type: TokenTypeSemicolon, Value: ";"},
			},
		},
		{
			name:  "runs tests successfully with whitespace and newlines",
			input: "select\tfoo  ,bar\t,baz from \t \n temptable;",
			want: []Token{
				{Type: TokenTypeSelect, Value: "select"},
				{Type: TokenTypeIdentifier, Value: "foo"},
				{Type: TokenTypeComma, Value: ","},
				{Type: TokenTypeIdentifier, Value: "bar"},
				{Type: TokenTypeComma, Value: ","},
				{Type: TokenTypeIdentifier, Value: "baz"},
				{Type: TokenTypeFrom, Value: "from"},
				{Type: TokenTypeIdentifier, Value: "temptable"},
				{Type: TokenTypeSemicolon, Value: ";"},
			},
		},
		{
			name:  "runs tests successfully with single quotes and special characters",
			input: "select _foo, bar123, ab$c FROM temptable",
			want: []Token{
				{Type: TokenTypeSelect, Value: "select"},
				{Type: TokenTypeIdentifier, Value: "_foo"},
				{Type: TokenTypeComma, Value: ","},
				{Type: TokenTypeIdentifier, Value: "bar123"},
				{Type: TokenTypeComma, Value: ","},
				{Type: TokenTypeIdentifier, Value: "ab$c"},
				{Type: TokenTypeFrom, Value: "FROM"},
				{Type: TokenTypeIdentifier, Value: "temptable"},
			},
		},
		{
			name:  "runs tests successfully with numbers and parentheses",
			input: "select foo from bar where (x=123.45 or y<0) and z>= .4",
			want: []Token{
				{Type: TokenTypeSelect, Value: "select"},
				{Type: TokenTypeIdentifier, Value: "foo"},
				{Type: TokenTypeFrom, Value: "from"},
				{Type: TokenTypeIdentifier, Value: "bar"},
				{Type: TokenTypeWhere, Value: "where"},
				{Type: TokenTypeOpenParen, Value: "("},
				{Type: TokenTypeIdentifier, Value: "x"},
				{Type: TokenTypeEq, Value: "="},
				{Type: TokenTypeNumber, Value: "123.45"},
				{Type: TokenTypeOr, Value: "or"},
				{Type: TokenTypeIdentifier, Value: "y"},
				{Type: TokenTypeLt, Value: "<"},
				{Type: TokenTypeNumber, Value: "0"},
				{Type: TokenTypeCloseParen, Value: ")"},
				{Type: TokenTypeAnd, Value: "and"},
				{Type: TokenTypeIdentifier, Value: "z"},
				{Type: TokenTypeGte, Value: ">="},
				{Type: TokenTypeNumber, Value: ".4"},
			},
		},
		{
			name:  "runs tests successfully with null identifier",
			input: "select * from temptable where x is not null",
			want: []Token{
				{Type: TokenTypeSelect, Value: "select"},
				{Type: TokenTypeStar, Value: "*"},
				{Type: TokenTypeFrom, Value: "from"},
				{Type: TokenTypeIdentifier, Value: "temptable"},
				{Type: TokenTypeWhere, Value: "where"},
				{Type: TokenTypeIdentifier, Value: "x"},
				{Type: TokenTypeIs, Value: "is"},
				{Type: TokenTypeNot, Value: "not"},
				{Type: TokenTypeNull, Value: "null"},
			},
		},
		{
			name:  "runs tests successfully with not equal operator",
			input: "select * from temptable where x != 123 or y <> 'hello'",
			want: []Token{
				{Type: TokenTypeSelect, Value: "select"},
				{Type: TokenTypeStar, Value: "*"},
				{Type: TokenTypeFrom, Value: "from"},
				{Type: TokenTypeIdentifier, Value: "temptable"},
				{Type: TokenTypeWhere, Value: "where"},
				{Type: TokenTypeIdentifier, Value: "x"},
				{Type: TokenTypeNotEq, Value: "!="},
				{Type: TokenTypeNumber, Value: "123"},
				{Type: TokenTypeOr, Value: "or"},
				{Type: TokenTypeIdentifier, Value: "y"},
				{Type: TokenTypeNotEq, Value: "<>"},
				{Type: TokenTypeString, Value: "hello"},
			},
		},
		{
			name:  "runs tests successfully with join clause",
			input: "select foo.x, bar.y from foo left outer join bar",
			want: []Token{
				{Type: TokenTypeSelect, Value: "select"},
				{Type: TokenTypeIdentifier, Value: "foo"},
				{Type: TokenTypeDot, Value: "."},
				{Type: TokenTypeIdentifier, Value: "x"},
				{Type: TokenTypeComma, Value: ","},
				{Type: TokenTypeIdentifier, Value: "bar"},
				{Type: TokenTypeDot, Value: "."},
				{Type: TokenTypeIdentifier, Value: "y"},
				{Type: TokenTypeFrom, Value: "from"},
				{Type: TokenTypeIdentifier, Value: "foo"},
				{Type: TokenTypeLeft, Value: "left"},
				{Type: TokenTypeOuter, Value: "outer"},
				{Type: TokenTypeJoin, Value: "join"},
				{Type: TokenTypeIdentifier, Value: "bar"},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			lexer := NewLexer(c.input)
			got, err := lexer.Lex()
			if err != nil {
				t.Fatalf("Lex(%q) returned error: %v", c.input, err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Lex(%q) == %v, want %v", c.input, got, c.want)
			}
		})
	}
}

func TestInvalidCases(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		{
			name:  "should return error for invalid input",
			input: "select foo,, bar, baz from temptable",
		},
		{
			name:  "runs tests successfully with invalid input",
			input: "select % from temptable",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			lexer := NewLexer(c.input)
			_, err := lexer.Lex()
			if err == nil {
				t.Fatalf("Lex(%q) returned nil error", c.input)
			}
		})
	}
}
