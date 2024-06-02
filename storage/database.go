package storage

import (
	"fmt"
	"github.com/Vignesh-Rajarajan/go-db/types"
)

type Database struct {
	tables map[string]*types.Relation
}

func NewDatabase() *Database {
	return &Database{
		tables: make(map[string]*types.Relation),
	}
}

func (db *Database) GetTable(name string) (*types.Relation, error) {
	t, ok := db.tables[name]
	if !ok {
		return nil, fmt.Errorf("table %s not found", name)
	}
	return t, nil
}

func (db *Database) CreateTable(name string, schema types.TableSchema) (*types.Relation, error) {
	_, ok := db.tables[name]
	if ok {
		return nil, fmt.Errorf("table %s already exists", name)
	}
	table := &types.Relation{
		Schema: schema,
	}
	db.tables[name] = table
	return table, nil
}
