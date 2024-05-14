package types

type Boolean struct {
	value bool
}

func NewBoolean(value bool) Boolean {
	return Boolean{value: value}
}

func (b Boolean) Compare(other Boolean) int {
	if b.value == other.value {
		return 0
	}
	if b.value {
		return 1
	}
	return -1
}
