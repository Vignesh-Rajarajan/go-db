package types

type Text struct {
	Value string
}

func NewText(value string) Text {
	return Text{Value: value}
}

func (t Text) Compare(other Text) int {
	if t.Value == other.Value {
		return 0
	}
	if t.Value < other.Value {
		return -1
	}
	return 1
}
