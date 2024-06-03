package planner

import (
	"fmt"
	"github.com/Vignesh-Rajarajan/go-db/sql"
	"github.com/Vignesh-Rajarajan/go-db/sql/query"
	"github.com/Vignesh-Rajarajan/go-db/storage"
)

func Plan(stmt *sql.SelectStatement, db *storage.Database) (query.QueryPlan, error) {
	plan, err := convertTableReference(stmt.From, db)
	if err != nil {
		return nil, err
	}
	switch what := stmt.What.(type) {
	case sql.Star:
	// Do nothing
	case sql.ExpressionList:
		schema := plan.Schema()
		columns := make([]query.OutputColumn, len(what.Expressions))
		for i, e := range what.Expressions {
			converted, name, err := ConvertExpression(e, schema)
			if err != nil {
				return nil, err
			}
			columns[i].Expression = converted
			columns[i].Name = name
		}
		plan, err = query.NewProject(plan, columns)
		if err != nil {
			return nil, err

		}
	default:
		return nil, fmt.Errorf("plan:: not implemented: %T", what)
	}
	return plan, nil
}

func convertTableReference(ref sql.TableReference, db *storage.Database) (query.QueryPlan, error) {
	switch f := ref.(type) {
	case sql.TableName:
		table, err := db.GetTable(f.Name)
		if err != nil {
			return nil, err
		}
		return query.NewLoad(f.Name, table.Schema), nil
	case *sql.Join:
		joinType := convertJoinType(f.Type)
		left, err := convertTableReference(f.Left, db)
		if err != nil {
			return nil, err
		}
		right, err := convertTableReference(f.Right, db)
		if err != nil {
			return nil, err
		}

		schema := query.CombineSchemas(left.Schema(), right.Schema())
		condition, _, err := ConvertExpression(f.Condition, schema)
		if err != nil {
			return nil, err
		}
		return query.NewJoin(joinType, left, right, condition)
	}

	return nil, fmt.Errorf("plan:: not implemented: %T", ref)
}

func convertJoinType(t sql.JoinType) query.JoinType {
	switch t {
	case sql.JoinTypeInner:
		return query.JoinTypeInner
	case sql.JoinTypeLeft:
		return query.JoinTypeLeftOuter
	case sql.JoinTypeRight:
		return query.JoinTypeRightOuter
	}
	return query.JoinTypeInner
}
