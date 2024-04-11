package parser

type Statement interface {
}

type StatementSelect struct {
	What  []Expression
	From  []FromExpression
	Where *Condition
}
type Star struct{}

type Column struct {
	Name string
}

type Expression interface{}
type FromExpression interface{}
type Condition interface{}

type TableName string
