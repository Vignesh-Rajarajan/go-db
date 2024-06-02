package types

import "fmt"

type Date struct {
	year, month, day int
}

func NewDate(year, month, day int) Date {
	if year < 1 || year > 9999 || month < 1 || month > 12 || day < 1 || day > daysInMonth(year, month) {
		return Date{}
	}
	return Date{year: year, month: month, day: day}
}

func daysInMonth(year int, month int) int {
	switch month {
	case 1, 3, 5, 7, 8, 10, 12:
		return 31
	case 2:
		if isLeapYear(year) {
			return 29
		}
		return 28
	}
	return 30
}

func isLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

func (d Date) Type() Type {
	return TypeDate
}

func (d Date) Compare(v Value) Comparison {
	next, ok := v.(Date)
	if !ok {
		return ComparisonIncomparable

	}
	switch {
	case d.year < next.year:
		return ComparisonLess
	case d.year > next.year:
		return ComparisonGreater
	case d.month < next.month:
		return ComparisonLess
	case d.month > next.month:
		return ComparisonGreater
	case d.day < next.day:
		return ComparisonLess
	case d.day > next.day:
		return ComparisonGreater
	}
	return ComparisonEqual
}

func (d Date) String() string {
	return fmt.Sprintf("Date(%4d-%2d-%2d)", d.year, d.month, d.day)
}
