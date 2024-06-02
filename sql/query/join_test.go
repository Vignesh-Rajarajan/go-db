package query

import (
	"github.com/Vignesh-Rajarajan/go-db/storage"
	"github.com/Vignesh-Rajarajan/go-db/types"
	"reflect"
	"testing"
)

func TestNewJoin(t *testing.T) {
	sampleData := storage.GetSampleData()

	wantSchema := types.TableSchema{
		Columns: []types.ColumnSchema{
			{Name: "films.id", Type: types.TypeDecimal},
			{Name: "films.title", Type: types.TypeText},
			{Name: "films.director", Type: types.TypeDecimal},
			{Name: "films.release_date", Type: types.TypeDate},
			{Name: "people.id", Type: types.TypeDecimal},
			{Name: "people.name", Type: types.TypeText},
		},
	}

	wantRows := [][]types.Value{
		{
			types.NewDecimal("1"),
			types.NewText("The Shawshank Redemption"),
			types.NewDecimal("1"),
			types.NewDate(1994, 9, 23),
			types.NewDecimal("1"),
			types.NewText("Frank Darabont"),
		}, {
			types.NewDecimal("2"),
			types.NewText("The Godfather"),
			types.NewDecimal("2"),
			types.NewDate(1972, 3, 24),
			types.NewDecimal("2"),
			types.NewText("Francis Ford Coppola"),
		}, {
			types.NewDecimal("3"),
			types.NewText("The Dark Knight"),
			types.NewDecimal("1"),
			types.NewDate(2008, 7, 18),
			types.NewDecimal("1"),
			types.NewText("Frank Darabont"),
		},
	}

	want := &types.Relation{
		Schema: wantSchema,
		Rows:   wantRows,
	}

	left := NewLoad("films", sampleData.Films.Schema)
	right := NewLoad("people", sampleData.People.Schema)
	col1 := NewColumnReference(2, types.TypeDecimal)
	col2 := NewColumnReference(4, types.TypeDecimal)
	condition, err := NewBinaryOperation(col1, col2, BinaryOperatorEq)
	if err != nil {
		t.Fatalf("unexpected error in NewBinaryOperation: %v", err)
	}
	join, err := NewJoin(JoinTypeInner, left, right, condition)
	if err != nil {
		t.Fatalf("unexpected error in NewJoin: %v", err)
	}

	gotSchema := join.Schema()
	if !reflect.DeepEqual(gotSchema, wantSchema) {
		t.Errorf("join.Schema() = %v, want %v", gotSchema, wantSchema)
	}

	got := join.Run(sampleData.Database)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("join.Run() = %v, want %v", got, want)
	}
}
