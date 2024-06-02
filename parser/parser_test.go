package parser

import (
	"github.com/Vignesh-Rajarajan/go-db/lexer"
	"github.com/Vignesh-Rajarajan/go-db/sql"
	"reflect"
	"testing"
)

func checkParser[T any](t *testing.T, name string, parse Parser[T], input string, want T) {
	t.Helper()
	tokens, err := lexer.Tokenize(input)
	if err != nil {
		t.Fatalf("unexpected error while tokenizing: %v", err)
	}
	tokensList := &lexer.TokenList{Input: input, Tokens: tokens}
	got, remaining, err := parse(tokensList)
	if err != nil {
		t.Fatalf("unexpected error while parsing %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	if remaining != nil && remaining.Len() != 0 {
		t.Fatalf("unexpected remaining tokens, did not get fully parsed: %v", remaining)
	}

}

func checkParserInvalid[T any](t *testing.T, name string, parse Parser[T], input string) {
	t.Helper()
	tokens, err := lexer.Tokenize(input)
	if err != nil {
		t.Fatalf("unexpected error while tokenizing: %v", err)
	}
	tokensList := &lexer.TokenList{Input: input, Tokens: tokens}
	_, _, err = parse(tokensList)
	if err == nil {
		t.Fatalf("expected error, but got nil")
	}
}

func TestParseValidValue(t *testing.T) {
	tt := []struct {
		input string
		want  sql.Expression
		err   error
	}{
		{
			input: `"hello"`,
			want:  sql.StringLiteral{Value: "hello"},
		},
		{
			input: "foo",
			want:  sql.ColumnReference{Name: "foo"},
		},
		{
			input: "foo.bar",
			want:  sql.ColumnReference{Relation: "foo", Name: "bar"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.input, func(t *testing.T) {
			checkParser(t, "ParseValue", ParseValue, tc.input, tc.want)
		})
	}
}

func TestParseValidExpression(t *testing.T) {
	cases := []struct {
		input string
		want  sql.Expression
	}{
		{
			input: "'hello'",
			want:  sql.StringLiteral{Value: "hello"},
		},
		{
			input: "'hello' = 'world'",
			want: &sql.BinaryOperation{
				Left:     sql.StringLiteral{Value: "hello"},
				Right:    sql.StringLiteral{Value: "world"},
				Operator: lexer.BinaryOperatorEq,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			checkParser(t, "ParseExpression", ParseExpression, c.input, c.want)
		})

	}
}

func TestParseSelectList(t *testing.T) {
	cases := []struct {
		input string
		want  sql.SelectList
	}{
		{
			input: "*",
			want:  sql.Star{},
		},
		{
			input: "",
			want:  sql.ExpressionList{},
		},
		{
			input: "foo",
			want: sql.ExpressionList{
				Expressions: []sql.Expression{sql.ColumnReference{Name: "foo"}},
			},
		},
		{
			input: "'a', 'b', 'c'",
			want: sql.ExpressionList{
				Expressions: []sql.Expression{
					sql.StringLiteral{Value: "a"},
					sql.StringLiteral{Value: "b"},
					sql.StringLiteral{Value: "c"},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			checkParser(t, "ParseSelectList", ParseSelectList, c.input, c.want)
		})
	}

	invalid := []string{
		"foo,",
	}
	for _, c := range invalid {
		t.Run(c, func(t *testing.T) {
			checkParserInvalid(t, "ParseSelectList", ParseSelectList, c)
		})
	}
}

func TestParserValueInvalid(t *testing.T) {
	cases := []struct {
		input string
		want  sql.Expression
	}{
		{
			input: "'hello'",
			want:  sql.StringLiteral{Value: "hello"},
		},
		{
			input: "foo",
			want:  sql.ColumnReference{Name: "foo"},
		},
		{
			input: "foo.bar",
			want:  sql.ColumnReference{Relation: "foo", Name: "bar"},
		},
	}

	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			checkParser(t, "ParseValue", ParseValue, c.input, c.want)
		})
	}
}

func TestParseTableReference(t *testing.T) {
	condition := sql.BinaryOperation{
		Left:     sql.ColumnReference{Name: "x", Relation: "foo"},
		Operator: lexer.BinaryOperatorEq,
		Right:    sql.ColumnReference{Name: "y", Relation: "bar"},
	}

	cases := []struct {
		input string
		want  sql.TableReference
	}{
		{
			input: "foo",
			want:  sql.TableName{Name: "foo"},
		},
		{
			input: "foo join bar on foo.x = bar.y",
			want: sql.Join{
				Left:      sql.TableName{Name: "foo"},
				Right:     sql.TableName{Name: "bar"},
				Condition: &condition,
				Type:      sql.JoinTypeInner,
			},
		},
		{
			input: "foo left outer join bar on foo.x = bar.y",
			want: sql.Join{
				Left:      sql.TableName{Name: "foo"},
				Right:     sql.TableName{Name: "bar"},
				Condition: &condition,
				Type:      sql.JoinTypeLeft,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			checkParser(t, "ParseTableReference", ParseTableReference, c.input, c.want)
		})
	}
}

func TestParseSelectStatement(t *testing.T) {
	condition := sql.BinaryOperation{
		Left:     sql.ColumnReference{Name: "x", Relation: "foo"},
		Operator: lexer.BinaryOperatorEq,
		Right:    sql.ColumnReference{Name: "y", Relation: "bar"},
	}

	cases := []struct {
		input string
		want  sql.SelectStatement
	}{
		{
			input: "select x,y from foo",
			want: sql.SelectStatement{
				What: sql.ExpressionList{
					Expressions: []sql.Expression{sql.ColumnReference{Name: "x"}, sql.ColumnReference{Name: "y"}},
				},
				From: sql.TableName{Name: "foo"},
			},
		},
		{
			input: "select * from foo join bar on foo.x = bar.y",
			want: sql.SelectStatement{
				What: sql.Star{},
				From: sql.Join{
					Left:      sql.TableName{Name: "foo"},
					Right:     sql.TableName{Name: "bar"},
					Condition: &condition,
					Type:      sql.JoinTypeInner,
				},
			},
		},
	}
	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			checkParser(t, "ParseSelectStatement", ParseSelectStatement, c.input, &c.want)
		})
	}
}

func TestParse(t *testing.T) {
	_, err := Parse("select * from foo;")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = Parse("select * from foo; where x = 1")
	if err == nil {
		t.Fatalf("expected error, but got nil")
	}
}
