package query

import (
	"fmt"
	"github.com/Vignesh-Rajarajan/go-db/storage"
	"github.com/Vignesh-Rajarajan/go-db/types"
)

type Query interface {
	Schema() types.TableSchema
	Run(db *storage.Database) *types.Relation
}

type Select struct {
	From      Query
	Condition Expression
}

func NewSelect(from Query, condition Expression) (*Select, error) {
	if condition.Type() != types.TypeBoolean {
		return nil, fmt.Errorf("condition must be a boolean expression %v", condition)
	}

	if err := condition.Check(from.Schema()); err != nil {
		return nil, err
	}

	return &Select{From: from, Condition: condition}, nil
}

func (s *Select) Schema() types.TableSchema {
	return s.From.Schema()
}

func (s *Select) Run(db *storage.Database) *types.Relation {
	from := s.From.Run(db)

	var rows [][]types.Value
	for i := range from.Rows {
		row := from.Row(i)
		res := s.Condition.Evaluate(row).(types.Boolean)
		if res.Bool() {
			rows = append(rows, row.Values)
		}
	}
	return &types.Relation{
		Schema: from.Schema,
		Rows:   rows,
	}
}
