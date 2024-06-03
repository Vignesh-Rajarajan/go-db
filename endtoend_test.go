package mini_sql_db

import (
	"github.com/Vignesh-Rajarajan/go-db/parser"
	"github.com/Vignesh-Rajarajan/go-db/planner"
	"github.com/Vignesh-Rajarajan/go-db/sql"
	"github.com/Vignesh-Rajarajan/go-db/storage"
	"github.com/Vignesh-Rajarajan/go-db/types"
	"reflect"
	"testing"
)

func TestAll(t *testing.T) {
	sampleData := storage.GetSampleData()
	query := "SELECT films.title, people.name FROM films JOIN people ON films.director = people.id"
	want := &types.Relation{
		Schema: types.TableSchema{
			Columns: []types.ColumnSchema{
				{Name: "films.title", Type: types.TypeText},
				{Name: "people.name", Type: types.TypeText},
			},
		},
		Rows: [][]types.Value{
			{types.NewText("The Shawshank Redemption"), types.NewText("Frank Darabont")},
			{types.NewText("The Godfather"), types.NewText("Francis Ford Coppola")},
			{types.NewText("The Dark Knight"), types.NewText("Frank Darabont")},
		},
	}

	stmt, err := parser.Parse(query)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	selectStmt, ok := stmt.(*sql.SelectStatement)
	if !ok {
		t.Fatalf("unexpected SelectStatement: %T", stmt)
	}

	plan, err := planner.Plan(selectStmt, sampleData.Database)
	if err != nil {
		t.Fatalf("planner.Plan unexpected error: %v", err)
	}

	got := plan.Run(sampleData.Database)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}
