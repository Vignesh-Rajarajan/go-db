package query

import (
	"github.com/Vignesh-Rajarajan/go-db/storage"
	"github.com/Vignesh-Rajarajan/go-db/types"
	"reflect"
	"testing"
)

func sampleRows() [][]types.Value {
	return [][]types.Value{
		{types.NewBoolean(false), types.NewText("hello")},
		{types.NewBoolean(true), types.NewText("world")},
	}
}

func sampleDatabase(t *testing.T) *storage.Database {
	t.Helper()
	db := storage.NewDatabase()
	table, err := db.CreateTable("mytable", sampleSchema())
	if err != nil {
		t.Fatalf("db.CreateTable() = %v, want nil", err)
	}
	for _, rows := range sampleRows() {
		err := table.Insert(rows)
		if err != nil {
			t.Fatalf("table.Insert() = %v, want nil", err)
		}
	}
	return db
}

func TestNewLoad(t *testing.T) {
	schema := sampleSchema()
	db := sampleDatabase(t)
	l := NewLoad("mytable", schema)
	got := l.Run(db)
	want := &types.Relation{
		Schema: schema,
		Rows:   sampleRows(),
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestNewSelect(t *testing.T) {
	schema := sampleSchema()
	db := sampleDatabase(t)
	l := NewLoad("mytable", schema)

	condition, err := NewBinaryOperation(
		NewColumnReference(0, types.TypeBoolean),
		NewConstant(types.NewBoolean(true)),
		BinaryOperatorEq)
	if err != nil {
		t.Fatalf("NewBinaryOperation() = %v, want nil", err)
	}

	s, err := NewSelect(l, condition)
	if err != nil {
		t.Fatalf("NewSelect() = %v, want nil", err)
	}
	got := s.Run(db)
	want := &types.Relation{
		Schema: schema,
		Rows:   [][]types.Value{sampleRows()[1]},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}
