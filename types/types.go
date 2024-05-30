package types

type Type int

const (
	TypeDate Type = iota
	TypeText
	TypeBoolean
	TypeDecimal
)

type Comparison int

const (
	ComparisonLess         Comparison = -1
	ComparisonEqual        Comparison = 0
	ComparisonGreater      Comparison = 1
	ComparisonIncomparable Comparison = 2
)

type Value interface {
	Compare(next Value) Comparison
	Type() Type
}
