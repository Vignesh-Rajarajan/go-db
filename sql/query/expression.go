package query

import (
	"fmt"
	"github.com/Vignesh-Rajarajan/go-db/types"
)

// Expression represents an expression of column, constant and operation
type Expression interface {
	Type() types.Type
	Evaluate(t *types.Row) types.Value
	Check(schema types.TableSchema) error
}

type Constant struct {
	Value types.Value
}

func NewConstant(value types.Value) Constant {
	return Constant{Value: value}
}

func (c Constant) Type() types.Type {
	return c.Value.Type()
}

func (c Constant) Evaluate(t *types.Row) types.Value {
	return c.Value
}

func (c Constant) Check(schema types.TableSchema) error {
	return nil
}

type ColumnReference struct {
	Index int
	T     types.Type
}

func NewColumnReference(index int, t types.Type) *ColumnReference {
	return &ColumnReference{Index: index, T: t}
}

func (c *ColumnReference) Type() types.Type {
	return c.T
}

func (c *ColumnReference) Evaluate(t *types.Row) types.Value {
	return t.Values[c.Index]
}

func (c *ColumnReference) Check(schema types.TableSchema) error {
	if c.Index < 0 || c.Index >= len(schema.Columns) {
		return fmt.Errorf("column index out of range: %d", c.Index)
	}

	got := schema.Columns[c.Index].Type
	expected := c.T
	if got != expected {
		return fmt.Errorf("mismatched types: %v and %v", got, expected)
	}
	return nil
}

type BinaryOperation struct {
	Left     Expression
	Right    Expression
	Operator BinaryOperator
}

func NewBinaryOperation(left, right Expression, operator BinaryOperator) (*BinaryOperation, error) {
	if left.Type() != right.Type() {
		return nil, fmt.Errorf("mismatched types: %v and %v", left.Type(), right.Type())
	}
	return &BinaryOperation{Left: left, Right: right, Operator: operator}, nil
}

func (b *BinaryOperation) Type() types.Type {
	return types.TypeBoolean
}

func (b *BinaryOperation) Evaluate(t *types.Row) types.Value {
	left := b.Left.Evaluate(t)
	right := b.Right.Evaluate(t)

	var result bool

	switch left.Compare(right) {
	case types.ComparisonEqual:
		result = b.Operator == BinaryOperatorLe || b.Operator == BinaryOperatorGe || b.Operator == BinaryOperatorEq
	case types.ComparisonLess:
		result = b.Operator == BinaryOperatorLt || b.Operator == BinaryOperatorLe
	case types.ComparisonGreater:
		result = b.Operator == BinaryOperatorGt || b.Operator == BinaryOperatorGe
	}
	return types.NewBoolean(result)
}

func (b *BinaryOperation) Check(schema types.TableSchema) error {
	if err := b.Left.Check(schema); err != nil {
		return err
	}
	if err := b.Right.Check(schema); err != nil {
		return err
	}
	return nil
}

type BinaryOperator int

const (
	BinaryOperatorEq BinaryOperator = iota
	BinaryOperatorNe
	BinaryOperatorLt
	BinaryOperatorLe
	BinaryOperatorGt
	BinaryOperatorGe
)
