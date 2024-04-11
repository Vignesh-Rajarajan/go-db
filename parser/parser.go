package parser

import "github.com/Vignesh-Rajarajan/go-db/lexer"

func Parse(input string) (Statement, error) {
	tokens, err := lexer.Tokenize(input)
	if err != nil {
		return nil, err
	}
	return ParseTokens(&lexer.TokenList{Input: input, Tokens: tokens})
}

func ParseTokens(tokens *lexer.TokenList) (Statement, error) {
	err := tokens.Consume(lexer.TokenTypeSelect)
	if err != nil {
		return nil, err
	}

	what, remTokens, err := parseExpressions(tokens)
	if err != nil {
		return nil, err
	}

	err = remTokens.Consume(lexer.TokenTypeFrom)
	if err != nil {
		return nil, err
	}

	from, remTokens, err := parseFromExpressions(tokens)
	if err != nil {
		return nil, err
	}

	return &StatementSelect{What: what, From: from}, nil
}

func parseExpressions(tokens *lexer.TokenList) ([]Expression, *lexer.TokenList, error) {
	var expressions []Expression

	expression, remainingToken, err := parseExpression(tokens)
	if err != nil {
		return nil, nil, err
	}
	expressions = append(expressions, expression)

	err = remainingToken.Consume(lexer.TokenTypeComma)
	if err != nil {
		return expressions, remainingToken, nil
	}

	moreExpressions, remainingToken, err := parseExpressions(remainingToken)
	if err != nil {
		return nil, nil, err
	}
	expressions = append(expressions, moreExpressions...)
	return expressions, remainingToken, nil
}

func parseExpression(tokens *lexer.TokenList) (Expression, *lexer.TokenList, error) {
	token, err := tokens.Get(lexer.TokenTypeStar, lexer.TokenTypeIdentifier)
	if err != nil {
		return nil, nil, err
	}
	if token.Type == lexer.TokenTypeStar {
		return new(Star), tokens, nil
	}
	return &Column{Name: token.Value}, tokens, nil
}

func parseFromExpressions(tokens *lexer.TokenList) ([]FromExpression, *lexer.TokenList, error) {

	token, err := tokens.Get(lexer.TokenTypeIdentifier)
	if err != nil {
		return nil, nil, err
	}
	tableName := TableName(token.Value)
	return []FromExpression{tableName}, tokens, nil
}
