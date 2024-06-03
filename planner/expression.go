package planner

import (
	"fmt"
	"github.com/Vignesh-Rajarajan/go-db/lexer"
	"github.com/Vignesh-Rajarajan/go-db/sql"
	"github.com/Vignesh-Rajarajan/go-db/sql/query"
	"github.com/Vignesh-Rajarajan/go-db/types"
	"strings"
)

func ConvertExpression(input sql.Expression, schema types.TableSchema) (query.Expression, string, error) {
	switch e := input.(type) {
	case sql.StringLiteral:
		return query.NewConstant(types.NewText(e.Value)), "", nil
	case sql.Boolean:
		return query.NewConstant(types.NewBoolean(e.Value)), "", nil
	case sql.NumberLiteral:
		return query.NewConstant(e.Value), "", nil
	case sql.ColumnReference:
		return convertColumnReference(e, schema)
	case *sql.BinaryOperation:
		return convertBinaryOperation(*e, schema)
	}
	return nil, "", fmt.Errorf("ConvertExpression:: not implemented: %T", input)
}

func convertColumnReference(e sql.ColumnReference, schema types.TableSchema) (query.Expression, string, error) {
	if e.Relation == "" {
		index, name, err := FindColumnIndex(e.Name, schema)
		if err != nil {
			return nil, "", err
		}
		return query.NewColumnReference(index, schema.Columns[index].Type), name, nil
	}
	name := fmt.Sprintf("%s.%s", e.Relation, e.Name)
	index, t, ok := schema.GetColumn(name)
	if !ok {
		return nil, "", fmt.Errorf("column %s not found", name)
	}
	return query.NewColumnReference(index, t), name, nil
}

func ConvertBinaryOperator(input lexer.BinaryOperator) query.BinaryOperator {
	switch input {
	case lexer.BinaryOperatorEq:
		return query.BinaryOperatorEq
	case lexer.BinaryOperatorNotEq:
		return query.BinaryOperatorNe
	case lexer.BinaryOperatorGt:
		return query.BinaryOperatorGt
	case lexer.BinaryOperatorGte:
		return query.BinaryOperatorGe
	case lexer.BinaryOperatorLt:
		return query.BinaryOperatorLt
	case lexer.BinaryOperatorLte:
		return query.BinaryOperatorLe
	}

	panic(fmt.Errorf("ConvertBinaryOperator:: not implemented: %v", input))
}

func convertBinaryOperation(input sql.BinaryOperation, schema types.TableSchema) (*query.BinaryOperation, string, error) {
	left, _, err := ConvertExpression(input.Left, schema)
	if err != nil {
		return nil, "", err
	}
	right, _, err := ConvertExpression(input.Right, schema)
	if err != nil {
		return nil, "", err
	}
	operator := ConvertBinaryOperator(input.Operator)
	expr, err := query.NewBinaryOperation(left, right, operator)
	return expr, "", err
}

func FindColumnIndex(name string, schema types.TableSchema) (int, string, error) {
	suffix := fmt.Sprintf(".%s", name)
	var index int
	var resName string
	var found bool
	for i, column := range schema.Columns {
		if strings.HasSuffix(column.Name, suffix) {
			if found {
				return 0, "", fmt.Errorf("unkown column name: %s", name)
			}
			index = i
			resName = column.Name
			found = true
		}
	}
	if !found {
		return 0, "", fmt.Errorf("column not found: %s", name)
	}
	return index, resName, nil
}
