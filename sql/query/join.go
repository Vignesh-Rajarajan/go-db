package query

import (
	"fmt"
	"github.com/Vignesh-Rajarajan/go-db/storage"
	"github.com/Vignesh-Rajarajan/go-db/types"
)

type JoinType int

const (
	JoinTypeInner JoinType = iota
	JoinTypeLeftOuter
	JoinTypeRightOuter
)

func (j JoinType) String() string {
	switch j {
	case JoinTypeInner:
		return "inner"
	case JoinTypeLeftOuter:
		return "left outer"
	case JoinTypeRightOuter:
		return "right outer"
	}
	return fmt.Sprintf("Unknown JoinType(%d)", j)
}

type Join struct {
	Type           JoinType
	Left           QueryPlan
	Right          QueryPlan
	Condition      Expression
	combinedSchema types.TableSchema
}

func NewJoin(t JoinType, left QueryPlan, right QueryPlan, condition Expression) (*Join, error) {
	if t != JoinTypeInner {
		return nil, fmt.Errorf("only inner join is supported")
	}
	if condition.Type() != types.TypeBoolean {
		return nil, fmt.Errorf("condition must be a boolean expression %v", condition)
	}
	combinedSchema := CombineSchemas(left.Schema(), right.Schema())
	if err := condition.Check(combinedSchema); err != nil {
		return nil, err
	}
	return &Join{Type: t, Left: left, Right: right, Condition: condition, combinedSchema: combinedSchema}, nil
}

func (j *Join) Schema() types.TableSchema {
	return j.combinedSchema
}

func (j *Join) Run(db *storage.Database) *types.Relation {
	left := j.Left.Run(db)
	right := j.Right.Run(db)
	schema := j.Schema()

	var rows [][]types.Value

	for _, leftRow := range left.Rows {
		for _, rightRow := range right.Rows {
			row := &types.Row{
				Schema: schema,
				Values: combinedRows(leftRow, rightRow),
			}

			got := j.Condition.Evaluate(row).(types.Boolean)
			if got.Bool() {
				rows = append(rows, row.Values)
			}

		}
	}
	return &types.Relation{
		Schema: schema,
		Rows:   rows,
	}
}

func (j *Join) Print(printer *Printer) {
	printer.Println("Join {")
	printer.Indent()
	printer.Println("Type: %s", j.Type)
	printer.Println("Left:")
	j.Left.Print(printer)
	printer.Println("Right:")
	j.Right.Print(printer)
	printer.Println("Condition: %s", j.Condition)
	printer.Dedent()
	printer.Println("}")

}

func CombineSchemas(left, right types.TableSchema) types.TableSchema {
	var columns []types.ColumnSchema
	columns = append(columns, left.Columns...)
	columns = append(columns, right.Columns...)
	return types.TableSchema{Columns: columns}
}

func combinedRows(left, right []types.Value) []types.Value {
	var row []types.Value
	row = append(row, left...)
	row = append(row, right...)
	return row
}
