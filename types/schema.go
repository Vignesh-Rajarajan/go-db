package types

import (
	"fmt"
	"strings"
)

type TableSchema struct {
	Columns []ColumnSchema
}

func (s *TableSchema) GetColumn(name string) (i int, t Type, ok bool) {
	for i, column := range s.Columns {
		if column.Name == name {
			return i, column.Type, true
		}
	}
	return
}

func (s *TableSchema) Check(row []Value) error {
	if len(row) != len(s.Columns) {
		return fmt.Errorf("wrong number of values: expected %d values, got %d", len(s.Columns), len(row))
	}
	for i := range s.Columns {
		if err := s.Columns[i].Check(row[i]); err != nil {
			return err
		}
	}
	return nil
}

func (s *TableSchema) Prefix(name string) TableSchema {
	var columns []ColumnSchema
	for _, column := range s.Columns {
		columns = append(columns, ColumnSchema{Name: fmt.Sprintf("%s.%s", name, column.Name), Type: column.Type})
	}
	return TableSchema{Columns: columns}
}

func (s *TableSchema) String() string {
	var columns []string
	for _, column := range s.Columns {
		columns = append(columns, column.String())
	}
	return fmt.Sprintf("TableSchema(%s)", strings.Join(columns, ", "))
}

type ColumnSchema struct {
	Name string
	Type Type
}

func (c *ColumnSchema) Check(value Value) error {
	if c.Type != value.Type() {
		return fmt.Errorf("mismatched types: column %s is of type %v, got %v", c.Name, c.Type, value.Type())
	}
	return nil
}

func (c *ColumnSchema) String() string {
	return fmt.Sprintf("%s %v", c.Name, c.Type)
}
