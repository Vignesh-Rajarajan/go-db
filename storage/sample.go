package storage

import "github.com/Vignesh-Rajarajan/go-db/types"

type SampleData struct {
	Database      *Database
	Films, People *types.Relation
}

func GetSampleData() *SampleData {
	filmSchema := types.TableSchema{
		Columns: []types.ColumnSchema{
			{Name: "id", Type: types.TypeDecimal},
			{Name: "title", Type: types.TypeText},
			{Name: "director", Type: types.TypeDecimal},
			{Name: "release_date", Type: types.TypeDate},
		},
	}

	filmRows := [][]types.Value{
		{
			types.NewDecimal("1"),
			types.NewText("The Shawshank Redemption"),
			types.NewDecimal("1"),
			types.NewDate(1994, 9, 23),
		},
		{
			types.NewDecimal("2"),
			types.NewText("The Godfather"),
			types.NewDecimal("2"),
			types.NewDate(1972, 3, 24),
		}, {
			types.NewDecimal("3"),
			types.NewText("The Dark Knight"),
			types.NewDecimal("1"),
			types.NewDate(2008, 7, 18),
		},
	}

	films := &types.Relation{
		Schema: filmSchema,
		Rows:   filmRows,
	}

	personSchema := types.TableSchema{
		Columns: []types.ColumnSchema{
			{Name: "id", Type: types.TypeDecimal},
			{Name: "name", Type: types.TypeText},
		},
	}
	personRows := [][]types.Value{
		{types.NewDecimal("1"), types.NewText("Frank Darabont")},
		{types.NewDecimal("2"), types.NewText("Francis Ford Coppola")},
	}

	people := &types.Relation{
		Schema: personSchema,
		Rows:   personRows,
	}

	db := &Database{
		tables: map[string]*types.Relation{
			"films":  films,
			"people": people,
		},
	}

	return &SampleData{
		Database: db,
		Films:    films,
		People:   people,
	}
}
