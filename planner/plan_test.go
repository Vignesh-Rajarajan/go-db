package planner

import (
	"github.com/Vignesh-Rajarajan/go-db/parser"
	"github.com/Vignesh-Rajarajan/go-db/sql"
	"github.com/Vignesh-Rajarajan/go-db/sql/query"
	"github.com/Vignesh-Rajarajan/go-db/storage"
	"github.com/Vignesh-Rajarajan/go-db/types"
	"reflect"
	"testing"
)

func parse(t *testing.T, stmt string) *sql.SelectStatement {
	t.Helper()
	parsed, err := parser.Parse(stmt)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	selectStatement, ok := parsed.(*sql.SelectStatement)
	if !ok {
		t.Fatalf("unexpected statement: %T", parsed)

	}
	return selectStatement
}

func TestPlan(t *testing.T) {
	sampleData := storage.GetSampleData()

	cases := []struct {
		stmt string
		want query.QueryPlan
	}{
		{
			stmt: "SELECT * FROM films",
			want: query.NewLoad("films", sampleData.Films.Schema),
		},
		{
			stmt: "SELECT id FROM films WHERE title = 'The Godfather'",
			want: &query.Project{
				From: &query.Select{
					From: query.NewLoad("films", sampleData.Films.Schema),
					Condition: &query.BinaryOperation{
						Left:     query.ColumnReference{Index: 1, T: types.TypeText},
						Right:    query.NewConstant(types.NewText("The Godfather")),
						Operator: query.BinaryOperatorEq,
					},
				},
				Columns: []query.OutputColumn{
					{
						"films.id",
						query.ColumnReference{Index: 0, T: types.TypeDecimal},
					},
				},
			},
		},
		//{
		//	stmt: "SELECT id, title, release_date, director FROM films",
		//	want: &query.Project{
		//		From: query.NewLoad("films", sampleData.Films.Schema),
		//		Columns: []query.OutputColumn{
		//			{
		//				"films.id",
		//				query.NewColumnReference(0, types.TypeDecimal),
		//			}, {
		//				"films.title",
		//				query.NewColumnReference(1, types.TypeText),
		//			},
		//			{
		//				"films.director",
		//				query.NewColumnReference(3, types.TypeDecimal),
		//			},
		//			{
		//				"films.release_date",
		//				query.NewColumnReference(2, types.TypeDate),
		//			},
		//		},
		//	},
		//},
	}

	for _, c := range cases {
		t.Run(c.stmt, func(t *testing.T) {
			stmt := parse(t, c.stmt)

			got, err := Plan(stmt, sampleData.Database)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("Query Plan for \n%s\n is \ngot %v, want %v\n", stmt, got, c.want)
			}
		})
	}
}
