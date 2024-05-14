package sql

import (
	"fmt"
	"github.com/Vignesh-Rajarajan/go-db/lexer"
	"strings"
)

type Statement interface {
	String() string
}

// Decimal represents a decimal number
type Decimal struct {
	Value  int
	Digits int
}

// Column represents a column in a table
type Column struct {
	Name string
}

// Expression represents an expression in SQL
type Expression interface {
	String() string
}

// ExpressionList represents a list of expressions
type ExpressionList struct {
	Expressions []Expression
}

func (e ExpressionList) String() string {
	list := make([]string, len(e.Expressions))
	for i, expr := range e.Expressions {
		list[i] = expr.String()
	}
	return fmt.Sprintf("ExpressionList(%s)", strings.Join(list, ", "))
}

// StringLiteral represents a string literal
type StringLiteral struct {
	Value string
}

func (s StringLiteral) String() string {
	return fmt.Sprintf("String(%q)", s.Value)
}

// NumberLiteral represents a number
type NumberLiteral struct {
	Value Decimal
}

func (i NumberLiteral) String() string {
	return fmt.Sprintf("NumberLiteral(%v)", i.Value)
}

// BinaryOperation represents a binary equality operation with a left and right side
type BinaryOperation struct {
	Left     Expression
	Operator lexer.BinaryOperator
	Right    Expression
}

func (o BinaryOperation) String() string {
	return fmt.Sprintf("BinaryOperation(Left: %s, Operator: %s, Right: %s)", o.Left, o.Operator, o.Right)
}

// ColumnReference represents a reference to a column in a table
type ColumnReference struct {
	Relation string
	Name     string
}

func (c ColumnReference) String() string {
	return fmt.Sprintf("Column(Relation: %s, Name: %s)", c.Relation, c.Name)
}

// SelectList represents a list of columns to select
type SelectList interface {
	String() string
}

// Star represents a select all columns
type Star struct{}

func (s Star) String() string {
	return "Star"
}

// TableReference represents a reference to a table
type TableReference interface {
	String() string
}

// TableName represents a reference to a table by name
type TableName struct {
	Name string
}

func (t TableName) String() string {
	return fmt.Sprintf("Table(%s)", t.Name)
}

type Join struct {
	Left      TableReference
	Right     TableReference
	Condition Expression
	Type      JoinType
}

func (j Join) String() string {
	return fmt.Sprintf("Join(%s,%s,%s,%s)", j.Type, j.Left, j.Right, j.Condition)
}

type JoinType int

const (
	JoinTypeInner JoinType = iota
	JoinTypeLeft
	JoinTypeRight
)

func (j JoinType) String() string {
	switch j {
	case JoinTypeInner:
		return "INNER"
	case JoinTypeLeft:
		return "left outer"
	case JoinTypeRight:
		return "right outer"
	}
	return fmt.Sprintf("Unknown JoinType(%d)", j)
}

type SelectStatement struct {
	What  SelectList
	From  TableReference
	Where Expression
}

func (s SelectStatement) String() string {
	where := ""
	if s.Where != nil {
		where = fmt.Sprintf(", Where: %s", s.Where.String())
	}
	return fmt.Sprintf("SelectStatement(What: %s, From: %s, Where:%s)", s.What.String(), s.From.String(), where)
}
