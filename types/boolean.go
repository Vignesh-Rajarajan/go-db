package types

type Boolean struct {
	value bool
}

func NewBoolean(value bool) Boolean {
	return Boolean{value: value}
}

func (b Boolean) Type() Type {
	return TypeBoolean
}

func (b Boolean) Bool() bool {
	return b.value
}

func (b Boolean) Compare(next Value) Comparison {
	nextBoolean, ok := next.(Boolean)
	if !ok {
		return ComparisonIncomparable
	}
	if b.value == nextBoolean.value {
		return ComparisonEqual
	}
	if b.value {
		return ComparisonGreater
	}
	return ComparisonLess
}
