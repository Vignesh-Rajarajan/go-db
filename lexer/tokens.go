package lexer

import (
	"fmt"
	"github.com/Vignesh-Rajarajan/go-db/custom_error"
	"strings"
)

type TokenType int

type Token struct {
	Type     TokenType
	Value    string
	From, To int
}

func (t Token) String() string {
	if t.Type >= TokenTypeComma {
		return t.Type.String()
	}
	return fmt.Sprintf("(%s %q)", t.Type, t.Value)
}

func PrintTokens(tokens []Token) string {
	values := make([]string, len(tokens))
	for i, token := range tokens {
		values[i] = token.String()
	}
	return strings.Join(values, " ")
}

func Tokenize(input string) ([]Token, error) {
	return NewLexer(input).Lex()
}

type TokenList struct {
	Input  string
	Tokens []Token
}

func (t *TokenList) Peek(expected ...TokenType) (Token, error) {
	if err := t.checkEnd(); err != nil {
		return Token{}, err
	}
	if err := t.checkType(expected); err != nil {
		return Token{}, err
	}
	return t.Tokens[0], nil
}

func (t *TokenList) Get(expected ...TokenType) (Token, error) {
	if err := t.checkEnd(); err != nil {
		return Token{}, err
	}
	first, rest := t.Tokens[0], t.Tokens[1:]
	if err := t.checkType(expected); err != nil {
		return Token{}, err
	}
	t.Tokens = rest
	return first, nil
}

func (t *TokenList) Consume(expected ...TokenType) error {
	if err := t.checkEnd(); err != nil {
		return err
	}
	if err := t.checkType(expected); err != nil {
		return err
	}
	t.Tokens = t.Tokens[1:]
	return nil
}

func (t *TokenList) checkEnd() error {
	if len(t.Tokens) == 0 {
		return custom_error.SyntaxError{
			Position: len([]rune(t.Input)),
			Message:  "unexpected end of input",
		}
	}
	return nil
}

func (t *TokenList) checkType(expected []TokenType) error {
	token := t.Tokens[0]
	if len(expected) == 0 {
		return nil
	}
	for _, e := range expected {
		if token.Type == e {
			return nil
		}
	}
	if len(expected) == 1 {
		return custom_error.SyntaxError{
			Position: token.From,
			Message:  fmt.Sprintf("expected %s, got %q", expected[0], token.Value),
		}
	}
	return custom_error.SyntaxError{
		Position: token.From,
		Message:  fmt.Sprintf("expected one of %s, got %q", joinWith(expected), token.Value),
	}
}

func joinWith(expected []TokenType) string {
	builder := new(strings.Builder)
	last := len(expected) - 1
	for i, e := range expected {
		switch {
		case i == 0:
			fmt.Fprintf(builder, "%v", e)
		case i == last:
			fmt.Fprintf(builder, "or %v", e)
		default:
			fmt.Fprintf(builder, ", %v", e)
		}
	}
	return builder.String()
}
