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
		name  string
	}{
		{
			input: sql.StringLiteral{Value: "hello"},
			want:  query.NewConstant(types.NewText("hello")),
			name:  "",
		}, {
			input: sql.Boolean{Value: true},
			want:  query.NewConstant(types.NewBoolean(true)),
			name:  "",
		}, {
			input: sql.NumberLiteral{Value: types.NewDecimal("123")},
			want:  query.NewConstant(types.NewDecimal("123")),
			name:  "",
		},
	}
	for _, c := range cases {
		t.Run(c.input.String(), func(t *testing.T) {
			got, _, err := ConvertExpression(c.input, schema)
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
	wantName := "films.release_date"
	got, resName, err := FindColumnIndex(name, schema)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != want {
		t.Errorf("Column index for %s is got %d, want %d", name, got, want)
	}
	if resName != wantName {
		t.Errorf("Column name for %s is got %s, want %s", name, name, wantName)

	}
	name = "name"
	_, _, err = FindColumnIndex(name, schema)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
