package types

import (
	"testing"
)

func TestNewDate(t *testing.T) {
	cases := []struct {
		year  int
		month int
		day   int
		want  Date
	}{
		{
			year:  2020,
			month: 1,
			day:   1,
			want:  Date{year: 2020, month: 1, day: 1},
		}, {
			year:  9999,
			month: 12,
			day:   31,
			want:  Date{year: 9999, month: 12, day: 31},
		},
	}

	for _, c := range cases {
		got := NewDate(c.year, c.month, c.day)
		if got != c.want {
			t.Errorf("NewDate(%d, %d, %d) == %v, want %v", c.year, c.month, c.day, got, c.want)
		}
	}
}

func TestDaysInMonth(t *testing.T) {
	cases := []struct {
		year  int
		month int
		want  int
	}{
		{
			year:  2020,
			month: 1,
			want:  31,
		},
		{
			year:  1999,
			month: 2,
			want:  28,
		}, {
			year:  2000,
			month: 2,
			want:  29,
		},
	}

	for _, c := range cases {
		got := daysInMonth(c.year, c.month)
		if got != c.want {
			t.Errorf("daysInMonth(%d, %d) == %d, want %d", c.year, c.month, got, c.want)
		}
	}
}

func TestDate_Compare(t *testing.T) {
	cases := []struct {
		d    Date
		next Date
		want Comparison
	}{
		{
			d:    Date{year: 2020, month: 1, day: 1},
			next: Date{year: 2020, month: 1, day: 1},
			want: ComparisonEqual,
		},
		{
			d:    Date{year: 2020, month: 1, day: 1},
			next: Date{year: 2020, month: 1, day: 2},
			want: ComparisonLess,
		},
		{
			d:    Date{year: 2020, month: 1, day: 1},
			next: Date{year: 2020, month: 2, day: 1},
			want: ComparisonLess,
		},
		{
			d:    Date{year: 2020, month: 1, day: 1},
			next: Date{year: 2021, month: 1, day: 1},
			want: ComparisonLess,
		},
	}

	for _, c := range cases {
		got := c.d.Compare(c.next)
		if got != c.want {
			t.Errorf("%v.Compare(%v) == %d, want %d", c.d, c.next, got, c.want)
		}
	}
}
