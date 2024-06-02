package types

type Relation struct {
	Schema TableSchema
	Rows   [][]Value
}

func (r *Relation) Row(i int) *Row {
	return &Row{Schema: r.Schema, Values: r.Rows[i]}
}

type Row struct {
	Schema TableSchema
	Values []Value
}

func (r *Relation) Insert(row []Value) error {
	if err := r.Schema.Check(row); err != nil {
		return err
	}
	r.Rows = append(r.Rows, row)
	return nil
}
