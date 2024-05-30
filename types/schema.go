package types

type TableSchema struct {
	Columns []ColumnSchema
}

func (s *TableSchema) GetColumn(name string) (int, Type) {
	for i, column := range s.Columns {
		if column.Name == name {
			return i, column.Type
		}
	}
	return -1, -1
}

type ColumnSchema struct {
	Name string
	Type Type
}
