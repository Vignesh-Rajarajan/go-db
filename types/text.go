package types

type Text struct {
	value string
}

func NewText(value string) Text {
	return Text{value: value}
}

func (t Text) Type() Type {
	return TypeText
}

func (t Text) Compare(next Value) Comparison {
	nextText, ok := next.(Text)
	if !ok {
		return ComparisonIncomparable
	}
	if t.value < nextText.value {
		return ComparisonLess
	}
	if t.value > nextText.value {
		return ComparisonGreater
	}
	return ComparisonEqual
}
