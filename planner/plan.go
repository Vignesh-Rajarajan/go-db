package planner

import (
	"fmt"
	"github.com/Vignesh-Rajarajan/go-db/sql"
	"github.com/Vignesh-Rajarajan/go-db/sql/query"
	"github.com/Vignesh-Rajarajan/go-db/storage"
)

func Plan(stmt *sql.SelectStatement, db *storage.Database) (query.QueryPlan, error) {
	var plan query.QueryPlan
	var err error
	switch f := stmt.From.(type) {

	case sql.TableName:
		table, err := db.GetTable(f.Name)
		if err != nil {
			return nil, err
		}
		plan = query.NewLoad(f.Name, table.Schema)
	case *sql.Join:
		return nil, fmt.Errorf("not implemented: %T", f)

	default:
		return nil, fmt.Errorf("not implemented: %T", f)
	}

	if stmt.Where != nil {
		schema := plan.Schema()
		condition, err := ConvertExpression(stmt.Where, schema)
		if err != nil {
			return nil, err
		}
		plan, err = query.NewSelect(plan, condition)
		if err != nil {
			return nil, err
		}
	}

	switch what := stmt.What.(type) {
	case sql.Star:
	// Do nothing
	case sql.ExpressionList:
		schema := plan.Schema()
		columns := make([]query.OutputColumn, len(what.Expressions))
		for i, e := range what.Expressions {
			converted, err := ConvertExpression(e, schema)
			if err != nil {
				return nil, err
			}
			columns[i].Expression = converted
			columns[i].Name = schema.Columns[i].Name
		}
		plan, err = query.NewProject(plan, columns)
		if err != nil {
			return nil, err

		}
	default:
		return nil, fmt.Errorf("Plan:: not implemented: %T", what)
	}
	return plan, nil
}
