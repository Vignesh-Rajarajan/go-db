package planner

import (
	"github.com/Vignesh-Rajarajan/go-db/sql"
	"github.com/Vignesh-Rajarajan/go-db/sql/query"
	"github.com/Vignesh-Rajarajan/go-db/storage"
	"github.com/Vignesh-Rajarajan/go-db/types"
	"reflect"
	"testing"
)

func TestConvertExpression(t *testing.T) {
	sampleData := storage.GetSampleData()
	schema := sampleData.Films.Schema

	cases := []struct {
		input sql.Expression
		want  query.Expression
	}{
		{
			input: sql.StringLiteral{Value: "hello"},
			want:  query.NewConstant(types.NewText("hello")),
		}, {
			input: sql.Boolean{Value: true},
			want:  query.NewConstant(types.NewBoolean(true)),
		}, {
			input: sql.NumberLiteral{Value: types.NewDecimal("123")},
			want:  query.NewConstant(types.NewDecimal("123")),
		}, {
			input: sql.ColumnReference{Relation: "films", Name: "title"},
			want:  query.NewColumnReference(1, types.TypeText),
		},
	}
	for _, c := range cases {
		t.Run(c.input.String(), func(t *testing.T) {
			got, err := ConvertExpression(c.input, schema)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("Expression for \n%v\n is \ngot %v, want %v\n", c.input, got, c.want)
			}
		})
	}
}

func TestFindColumnIndex(t *testing.T) {
	schema := types.TableSchema{
		Columns: []types.ColumnSchema{
			{Name: "films.name", Type: types.TypeText},
			{Name: "films.release_date", Type: types.TypeDate},
			{Name: "people.name", Type: types.TypeText},
		},
	}

	name := "release_date"
	want := 1
	got, err := FindColumnIndex(name, schema)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != want {
		t.Errorf("Column index for %s is got %d, want %d", name, got, want)
	}
	name = "name"
	_, err = FindColumnIndex(name, schema)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
