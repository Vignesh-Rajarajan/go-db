package query

import (
	"fmt"
	"github.com/Vignesh-Rajarajan/go-db/storage"
	"github.com/Vignesh-Rajarajan/go-db/types"
)

type QueryPlan interface {
	Schema() types.TableSchema
	Run(db *storage.Database) *types.Relation
	Print(printer *Printer)
}

type Load struct {
	TableName   string
	TableSchema types.TableSchema
}

func NewLoad(tableName string, tableSchema types.TableSchema) *Load {
	return &Load{TableName: tableName, TableSchema: tableSchema.Prefix(tableName)}
}

func (l *Load) Schema() types.TableSchema {
	return l.TableSchema

}

func (l *Load) Run(db *storage.Database) *types.Relation {
	r, err := db.GetTable(l.TableName)
	if err != nil {
		panic(fmt.Errorf("table %s not found", l.TableName))
	}
	return r

}

func (l *Load) Print(printer *Printer) {
	printer.Println("Load {")
	printer.Indent()
	printer.Println("Table: %q", l.TableName)
	printer.Println("Schema: %s", l.TableSchema)
	printer.Dedent()
	printer.Println("}")

}

type Select struct {
	From      QueryPlan
	Condition Expression
}

func NewSelect(from QueryPlan, condition Expression) (*Select, error) {
	if condition.Type() != types.TypeBoolean {
		return nil, fmt.Errorf("condition must be a boolean expression %v", condition)
	}

	if err := condition.Check(from.Schema()); err != nil {
		return nil, err
	}

	return &Select{From: from, Condition: condition}, nil
}

func (s *Select) Schema() types.TableSchema {
	return s.From.Schema()
}

func (s *Select) Run(db *storage.Database) *types.Relation {
	from := s.From.Run(db)

	var rows [][]types.Value
	for i := range from.Rows {
		row := from.Row(i)
		res := s.Condition.Evaluate(row).(types.Boolean)
		if res.Bool() {
			rows = append(rows, row.Values)
		}
	}
	return &types.Relation{
		Schema: from.Schema,
		Rows:   rows,
	}
}

func (s *Select) Print(printer *Printer) {
	printer.Println("Select {")
	printer.Indent()
	printer.Println("From:")
	s.From.Print(printer)
	printer.Println("Condition: %s", s.Condition)
	printer.Dedent()
	printer.Println("}")
}

type OutputColumn struct {
	Name       string
	Expression Expression
}

func SimpleColumn(name string, index int, t types.Type) OutputColumn {
	return OutputColumn{
		Name:       name,
		Expression: NewColumnReference(index, t),
	}
}

func ComputedColumn(name string, expression Expression) OutputColumn {
	return OutputColumn{
		Name:       name,
		Expression: expression,
	}
}

func (oc OutputColumn) Schema() types.ColumnSchema {
	return types.ColumnSchema{
		Name: oc.Name,
		Type: oc.Expression.Type(),
	}
}

type Project struct {
	From    QueryPlan
	Columns []OutputColumn
}

func NewProject(from QueryPlan, columns []OutputColumn) (*Project, error) {
	names := make(map[string]bool)
	for _, c := range columns {
		if names[c.Name] {
			return nil, fmt.Errorf("duplicate column name %s", c.Name)
		}
		if err := c.Expression.Check(from.Schema()); err != nil {
			return nil, err
		}
		names[c.Name] = true
	}

	for _, c := range columns {
		if err := c.Expression.Check(from.Schema()); err != nil {
			return nil, err
		}
	}

	return &Project{From: from, Columns: columns}, nil
}

func (p *Project) Schema() types.TableSchema {
	var columns []types.ColumnSchema
	for _, c := range p.Columns {
		columns = append(columns, c.Schema())
	}
	return types.TableSchema{Columns: columns}
}

func (p *Project) Run(db *storage.Database) *types.Relation {
	from := p.From.Run(db)
	var rows [][]types.Value
	for i := range from.Rows {
		row := make([]types.Value, len(p.Columns))
		for j := range p.Columns {
			row[j] = p.Columns[j].Expression.Evaluate(from.Row(i))
		}
		rows[i] = row
	}
	return &types.Relation{
		Schema: p.Schema(),
		Rows:   rows,
	}
}

func (p *Project) Print(printer *Printer) {
	printer.Println("Project {")
	printer.Indent()
	printer.Println("From:")
	p.From.Print(printer)
	printer.Println("Columns:")
	printer.Indent()
	for _, c := range p.Columns {
		printer.Println("%s: %s", c.Name, c.Expression)
	}
	printer.Dedent()
	printer.Dedent()
	printer.Println("}")
}
