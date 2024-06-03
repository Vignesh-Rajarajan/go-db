package parser

import (
	"github.com/Vignesh-Rajarajan/go-db/lexer"
	"github.com/Vignesh-Rajarajan/go-db/sql"
	"github.com/Vignesh-Rajarajan/go-db/types"
)

type Parser[T any] func(*lexer.TokenList) (T, *lexer.TokenList, error)

// ParseExpression parses an expression from a list of tokens and returns the parsed expression and the remaining tokens
func ParseExpression(tokens *lexer.TokenList) (sql.Expression, *lexer.TokenList, error) {
	left, remTokens, err := ParseValue(tokens)
	if err != nil {
		return nil, nil, err
	}

	token, err := remTokens.Get(lexer.TokenTypeEq, lexer.TokenTypeNotEq, lexer.TokenTypeGt, lexer.TokenTypeLt, lexer.TokenTypeGte, lexer.TokenTypeLte)
	if err != nil {
		return left, remTokens, nil
	}

	op := lexer.TokenToBinaryOperator[token.Type]

	right, remTokens, err := ParseValue(remTokens)
	if err != nil {
		return nil, nil, err
	}
	result := &sql.BinaryOperation{
		Left:     left,
		Right:    right,
		Operator: op,
	}
	return result, remTokens, nil
}

// ParseValue parses a value from a list of tokens and returns the parsed value and the remaining tokens
func ParseValue(tokens *lexer.TokenList) (sql.Expression, *lexer.TokenList, error) {
	token, err := tokens.Peek(lexer.TokenTypeString, lexer.TokenTypeNumber, lexer.TokenTypeIdentifier, lexer.TokenTypeTrue, lexer.TokenTypeFalse)
	if err != nil {
		return nil, nil, err
	}

	switch token.Type {
	case lexer.TokenTypeString:
		return ParseString(tokens)
	case lexer.TokenTypeNumber:
		return ParseNumber(tokens)
	case lexer.TokenTypeTrue, lexer.TokenTypeFalse:
		_ = tokens.Consume()
		return sql.Boolean{Value: token.Type == lexer.TokenTypeTrue}, tokens, nil

	default:
		return ParseColumnReference(tokens)
	}
}

// ParseString parses a string from a list of tokens and returns the parsed string and the remaining tokens
func ParseString(tokens *lexer.TokenList) (sql.Expression, *lexer.TokenList, error) {
	token, err := tokens.Get(lexer.TokenTypeString)
	if err != nil {
		return sql.StringLiteral{}, nil, err
	}
	return sql.StringLiteral{Value: token.Value}, tokens, nil
}

// ParseNumber parses a number from a list of tokens and returns the parsed number and the remaining tokens
func ParseNumber(tokens *lexer.TokenList) (sql.Expression, *lexer.TokenList, error) {
	token, err := tokens.Get(lexer.TokenTypeNumber)
	if err != nil {
		return sql.NumberLiteral{}, nil, err
	}
	decimal, err := types.ParseDecimal(token.Value)
	if err != nil {
		return sql.NumberLiteral{}, nil, lexer.SyntaxError{
			Position: token.From,
			Message:  err.Error(),
		}
	}
	return sql.NumberLiteral{Value: decimal}, tokens, nil
}

// ParseColumnReference parses a column reference from a list of tokens and returns the parsed column reference and the remaining tokens
func ParseColumnReference(tokens *lexer.TokenList) (sql.Expression, *lexer.TokenList, error) {
	token, err := tokens.Get(lexer.TokenTypeIdentifier)
	if err != nil {
		return nil, nil, err
	}

	err = tokens.Consume(lexer.TokenTypeDot)
	if err != nil {
		return sql.ColumnReference{Name: token.Value}, tokens, nil
	}
	secondToken, err := tokens.Get(lexer.TokenTypeIdentifier)
	if err != nil {
		return nil, nil, err
	}
	return sql.ColumnReference{Relation: token.Value, Name: secondToken.Value}, tokens, nil
}

// ParseSelectList parses a select list from a list of tokens and returns the parsed select list and the remaining tokens
func ParseSelectList(tokens *lexer.TokenList) (sql.SelectList, *lexer.TokenList, error) {
	err := tokens.Consume(lexer.TokenTypeStar)
	if err == nil {
		return sql.Star{}, tokens, nil
	}

	expr, remTokens, err := ParseExpressionList(tokens)
	if err != nil {
		return nil, nil, err
	}
	return sql.ExpressionList{Expressions: expr}, remTokens, nil
}

// ParseExpressionList parses an expression list from a list of tokens and returns the parsed expression list and the remaining tokens
func ParseExpressionList(tokens *lexer.TokenList) ([]sql.Expression, *lexer.TokenList, error) {
	var result []sql.Expression

	first := true
	for {
		expr, remTokens, err := ParseExpression(tokens)
		if err != nil {
			if first {
				break
			}
			return nil, nil, err
		}
		first = false
		result = append(result, expr)

		err = remTokens.Consume(lexer.TokenTypeComma)
		if err != nil {
			break
		}
	}
	return result, tokens, nil
}

func ParseTableName(tokens *lexer.TokenList) (sql.TableName, *lexer.TokenList, error) {
	token, err := tokens.Get(lexer.TokenTypeIdentifier)
	if err != nil {
		return sql.TableName{}, nil, err
	}
	return sql.TableName{Name: token.Value}, tokens, nil
}

func ParseTableReference(tokens *lexer.TokenList) (sql.TableReference, *lexer.TokenList, error) {
	left, remTokens, err := ParseTableName(tokens)
	if err != nil {
		return nil, nil, err
	}

	token, err := tokens.Get(lexer.TokenTypeLeft, lexer.TokenTypeRight, lexer.TokenTypeJoin)
	if err != nil {
		// No join
		return left, remTokens, nil
	}

	join := &sql.Join{
		Left: left,
	}

	switch token.Type {
	case lexer.TokenTypeLeft:
		join.Type = sql.JoinTypeLeft
		_ = remTokens.Consume(lexer.TokenTypeOuter)
		_ = remTokens.Consume(lexer.TokenTypeJoin)
	case lexer.TokenTypeRight:
		join.Type = sql.JoinTypeRight
		_ = remTokens.Consume(lexer.TokenTypeOuter)
		_ = remTokens.Consume(lexer.TokenTypeJoin)
	default:
		join.Type = sql.JoinTypeInner
	}

	right, remTokens, err := ParseTableReference(remTokens)
	if err != nil {
		return nil, nil, err
	}
	join.Right = right

	err = remTokens.Consume(lexer.TokenTypeOn)
	if err != nil {
		return nil, nil, err
	}

	expression, remTokens, err := ParseExpression(remTokens)
	if err != nil {
		return nil, nil, err
	}
	join.Condition = expression

	return join, remTokens, nil
}

func ParseStatement(tokens *lexer.TokenList) (sql.Statement, *lexer.TokenList, error) {
	return ParseSelectStatement(tokens)
}

func ParseSelectStatement(tokens *lexer.TokenList) (*sql.SelectStatement, *lexer.TokenList, error) {
	err := tokens.Consume(lexer.TokenTypeSelect)
	if err != nil {
		return nil, nil, err
	}
	result := sql.SelectStatement{}

	what, remTokens, err := ParseSelectList(tokens)
	if err != nil {
		return nil, nil, err
	}
	result.What = what

	err = remTokens.Consume(lexer.TokenTypeFrom)
	if err != nil {
		return nil, nil, err
	}

	from, remTokens, err := ParseTableReference(remTokens)
	if err != nil {
		return nil, nil, err
	}
	result.From = from

	err = remTokens.Consume(lexer.TokenTypeWhere)
	if err == nil {
		result.Where, remTokens, err = ParseExpression(remTokens)
		if err != nil {
			return nil, nil, err
		}
	}
	return &result, remTokens, nil
}

func Parse(input string) (sql.Statement, error) {
	tokens, err := lexer.Tokenize(input)
	if err != nil {
		return nil, err
	}
	tokenList := &lexer.TokenList{
		Input:  input,
		Tokens: tokens,
	}
	stmt, remToken, err := ParseSelectStatement(tokenList)
	if err != nil {
		return nil, err
	}
	_ = remToken.Consume(lexer.TokenTypeSemicolon)

	if err = remToken.ExpectedEnd(); err != nil {
		return nil, err
	}
	return stmt, nil
}
