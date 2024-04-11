package lexer

import "testing"

var tokenList = TokenList{
	Input: "select * from foo",
	Tokens: []Token{
		{Type: TokenTypeSelect, Value: "select"},
		{Type: TokenTypeStar, Value: "*"},
		{Type: TokenTypeFrom, Value: "from"},
		{Type: TokenTypeIdentifier, Value: "foo"},
	},
}

func TestTokenList_Peek(t *testing.T) {
	l := TokenList{}
	_, err := l.Peek()
	if err == nil {
		t.Errorf("expected custom_error, but got nil")
	}
	l = tokenList

	got, err := l.Peek(TokenTypeSelect)
	if err != nil {
		t.Errorf("unexpected custom_error: %v", err)
	}
	if got != tokenList.Tokens[0] {
		t.Errorf("expected %v, got %v", tokenList.Tokens[0], got)
	}
	l = tokenList
	_, err = l.Peek(TokenTypeIdentifier)
	if err == nil {
		t.Errorf("expected custom_error, but got nil")
	}
	l = tokenList
	_, err = l.Peek(TokenTypeGte, TokenTypeEq, TokenTypeIdentifier)
	if err == nil {
		t.Errorf("expected custom_error, but got nil")
	}
}

func TestTokenList_Get(t *testing.T) {
	l := TokenList{}
	_, err := l.Get()
	if err == nil {
		t.Errorf("expected custom_error, but got nil")
	}
	l = tokenList

	got, err := l.Get(TokenTypeSelect)
	if err != nil {
		t.Errorf("unexpected custom_error: %v", err)
	}
	if got != tokenList.Tokens[0] {
		t.Errorf("expected %v, got %v", tokenList.Tokens[0], got)
	}
}
