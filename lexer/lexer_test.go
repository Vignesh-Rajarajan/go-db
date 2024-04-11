package lexer

import (
	"github.com/Vignesh-Rajarajan/go-db/custom_error"
	"reflect"
	"testing"
)

func TestLexer(t *testing.T) {
	cases := []struct {
		name        string
		input, want string
	}{
		{
			name:  "runs tests successfully with empty input",
			input: "",
			want:  "",
		},
		{
			name:  "runs tests successfully with multiple Tokens",
			input: "select foo, bar, baz from temptable;",
			want:  `Select (Identifier "foo") Comma (Identifier "bar") Comma (Identifier "baz") From (Identifier "temptable") Semicolon`,
		},
		{
			name:  "runs tests successfully with whitespace and newlines",
			input: "select\tfoo  ,bar\t,baz from \t \n temptable;",
			want:  `Select (Identifier "foo") Comma (Identifier "bar") Comma (Identifier "baz") From (Identifier "temptable") Semicolon`,
		},
		{
			name:  "runs tests successfully with single quotes and special characters",
			input: "select _foo, bar123, ab$c FROM temptable;",
			want:  `Select (Identifier "_foo") Comma (Identifier "bar123") Comma (Identifier "ab$c") From (Identifier "temptable") Semicolon`,
		},
		{
			name:  "runs tests successfully with numbers and parentheses",
			input: "select foo from bar where (x=123.45 or y<0) and z>= .4",
			want:  `Select (Identifier "foo") From (Identifier "bar") Where OpenParen (Identifier "x") Eq (Number "123.45") Or (Identifier "y") Lt (Number "0") CloseParen And (Identifier "z") Gte (Number ".4")`,
		},
		{
			name:  "runs tests successfully with null identifier",
			input: "select * from temptable where x is not null",
			want:  `Select Star From (Identifier "temptable") Where (Identifier "x") Is Not Null`,
		},
		{
			name:  "runs tests successfully with not equal operator",
			input: "select * from temptable where x != 123 or y <> 'hello'",
			want:  `Select Star From (Identifier "temptable") Where (Identifier "x") NotEq (Number "123") Or (Identifier "y") NotEq (String "hello")`,
		},
		{
			name:  "runs tests successfully with join clause",
			input: "select foo.x, bar.y from foo left outer join bar",
			want:  `Select (Identifier "foo") Dot (Identifier "x") Comma (Identifier "bar") Dot (Identifier "y") From (Identifier "foo") Left Outer Join (Identifier "bar")`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			lexer := NewLexer(c.input)
			tokens, err := lexer.Lex()
			if err != nil {
				t.Fatalf("Lex(%q) returned custom_error: %v", c.input, err)
			}
			got := PrintTokens(tokens)
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
			name:  "should return custom_error for invalid input",
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
				t.Fatalf("Lex(%q) returned nil custom_error", c.input)
			}
		})
	}
}

func TestLexError(t *testing.T) {

	cases := []struct {
		input string
		want  custom_error.SyntaxError
	}{
		{
			input: "select foo,, bar, baz from temptable",
			want: custom_error.SyntaxError{
				Position: 10,
				Message:  `invalid SQL: ",,"`,
			},
		},
		{
			input: "select % from temptable",
			want: custom_error.SyntaxError{
				Position: 7,
				Message:  `unexpected character '%'`,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			_, err := Tokenize(c.input)
			if err.Error() != c.want.Error() {
				t.Fatalf("Lex(%q) == %v, want %v", c.input, err, c.want)
			}
		})
	}
}
