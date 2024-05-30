package storage

import (
	"fmt"
	"github.com/Vignesh-Rajarajan/go-db/types"
)

type Database struct {
	name   string
	tables map[string]*types.Relation
}

func (db *Database) GetTable(name string) (*types.Relation, error) {
	t, ok := db.tables[name]
	if !ok {
		return nil, fmt.Errorf("table %s not found", name)
	}
	return t, nil
}
