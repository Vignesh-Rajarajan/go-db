package query

import (
	"github.com/Vignesh-Rajarajan/go-db/types"
	"testing"
)

func sampleSchema() types.TableSchema {
	return types.TableSchema{
		Columns: []types.ColumnSchema{
			{Name: "foo", Type: types.TypeBoolean},
			{Name: "bar", Type: types.TypeText},
		},
	}
}

func sampleRow() *types.Row {
	return &types.Row{
		Schema: sampleSchema(),
		Values: []types.Value{
			types.NewBoolean(true),
			types.NewText("hello"),
		},
	}
}

func TestColumnReference_Type(t *testing.T) {
	cases := []struct {
		index   int
		t       types.Type
		wantErr bool
	}{
		{
			index:   1,
			t:       types.TypeText,
			wantErr: false,
		},
		{
			index:   1,
			t:       types.TypeBoolean,
			wantErr: true,
		},
		{
			index:   2,
			t:       types.TypeBoolean,
			wantErr: true,
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			err := NewColumnReference(c.index, c.t).Check(sampleSchema())
			if err != nil && !c.wantErr {
				t.Errorf("NewColumnReference(%d, %v) returned error, want nil", c.index, c.t)
			}
			if err == nil && c.wantErr {
				t.Errorf("NewColumnReference(%d, %v) returned nil, want error", c.index, c.t)
			}
		})
	}
}
