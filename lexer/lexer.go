package lexer

import (
	"fmt"
	"strings"
)

type LexerState int

const (
	LexerStateInitial LexerState = iota
	LexerStateIdentifier
	LexerStatePunctuation
	LexerStateNumber
	LexerStateString
)

type Lexer struct {
	input  []rune
	state  LexerState
	from   int
	next   int
	tokens []Token
	err    error
}

func NewLexer(input string) *Lexer {
	return &Lexer{input: []rune(input)}
}

func (l *Lexer) nextRune() (rune, bool) {
	if l.next > len(l.input) {
		return 0, false
	}
	if l.next == len(l.input) {
		return ' ', true
	}
	r := l.input[l.next]
	return r, true
}

func (l *Lexer) Lex() ([]Token, error) {

	for {
		r, ok := l.nextRune()
		if !ok {
			return l.tokens, nil
		}
		switch l.state {
		case LexerStateInitial:
			l.parseForInitialState(r)
		case LexerStateIdentifier:
			l.parseForIdentifier(r)
		case LexerStateNumber:
			l.parseForNumber(r)
		case LexerStateString:
			l.parseForString(r)
		case LexerStatePunctuation:
			l.parseForPunctuation(r)
		default:
			l.err = fmt.Errorf("unexpected state %d", l.state)
		}
		if l.err != nil {
			return nil, l.err
		}
		l.next++
	}
}

func (l *Lexer) parseForPunctuation(r rune) {
	switch {
	case isPunctuation(r):
	case isDigitOrDot(r):
		l.tokenForPunctuation()
		l.changeState(LexerStateNumber)
	case isQuote(r):
		l.tokenForPunctuation()
		l.changeState(LexerStateString)
	case isWordStartWithCharacter(r):
		l.tokenForPunctuation()
		l.changeState(LexerStateIdentifier)
	case isWhitespace(r):
		l.tokenForPunctuation()
		l.changeState(LexerStateInitial)
	default:
		l.errorf(l.next, "unexpected character '%c'", r)

	}
}

func (l *Lexer) parseForString(r rune) {
	switch {
	case isQuote(r):
		l.tokenForString()
		l.changeState(LexerStateInitial)
	default:

	}
}

func (l *Lexer) parseForNumber(r rune) {
	switch {
	case isDigitOrDot(r):
	case isWhitespace(r):
		l.tokenForNumber()
		l.changeState(LexerStateInitial)
	case isPunctuation(r):
		l.tokenForNumber()
		l.changeState(LexerStatePunctuation)
	default:
		l.errorf(l.next, "unexpected character '%c'", r)
	}
}

func (l *Lexer) parseForIdentifier(r rune) {
	switch {
	case isWordCharacter(r):
	case isWhitespace(r):
		l.tokenForWord()
		l.changeState(LexerStateInitial)
	case isPunctuation(r):
		l.tokenForWord()
		l.changeState(LexerStatePunctuation)
	default:
		l.errorf(l.next, "unexpected character '%c'", r)
	}
}

func (l *Lexer) parseForInitialState(r rune) {
	switch {
	case isWhitespace(r):
	case isWordStartWithCharacter(r):
		l.changeState(LexerStateIdentifier)
	case isDigitOrDot(r):
		l.changeState(LexerStateNumber)
	case isQuote(r):
		l.changeState(LexerStateString)
	case isPunctuation(r):
		l.changeState(LexerStatePunctuation)
	default:
		l.errorf(l.next, "unexpected character '%c'", r)
	}
}

func (l *Lexer) changeState(s LexerState) {
	l.state = s
	l.from = l.next
}

func (l *Lexer) tokenForWord() {
	word := string(l.input[l.from:l.next])
	tokenType := TokenTypeIdentifier
	if typ, ok := KeywordMap[strings.ToLower(word)]; ok {
		tokenType = typ
	}
	l.produceToken(l.from, l.next, tokenType)
}

func (l *Lexer) produceToken(from, to int, typ TokenType) {
	token := Token{Value: string(l.input[from:to]), Type: typ, From: from, To: to}
	l.tokens = append(l.tokens, token)
}

func (l *Lexer) tokenForNumber() {
	l.produceToken(l.from, l.next, TokenTypeNumber)
}

func (l *Lexer) tokenForString() {
	from := l.from + 1
	to := l.next
	l.produceToken(from, to, TokenTypeString)
}

func (l *Lexer) tokenForPunctuation() {
	val := string(l.input[l.from:l.next])
	typ, ok := SymbolMap[val]
	if !ok {
		l.errorf(l.from, "invalid SQL: %q", val)
		return
	}
	l.produceToken(l.from, l.next, typ)
}

func (l *Lexer) errorf(pos int, s string, r ...any) {
	l.err = &SyntaxError{
		Position: pos,
		Message:  fmt.Sprintf(s, r...),
	}
}
