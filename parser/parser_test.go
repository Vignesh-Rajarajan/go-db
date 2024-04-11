package parser

import (
	"github.com/Vignesh-Rajarajan/go-db/custom_error"
	"reflect"
	"testing"
)

func TestParseValid(t *testing.T) {
	tt := []struct {
		input string
		want  Statement
		err   error
	}{
		{
			input: "select x,y from foo",
			want: &StatementSelect{
				What: []Expression{
					&Column{Name: "x"},
					&Column{Name: "y"},
				},
				From: []FromExpression{
					TableName("foo"),
				},
			},
			err: nil,
		},
		{
			input: "select * from foo",
			want: &StatementSelect{
				What: []Expression{
					&Star{},
				},
				From: []FromExpression{
					TableName("foo"),
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.input, func(t *testing.T) {
			stmt, err := Parse(tc.input)
			if err != nil {
				t.Fatalf("Parse(%q) returned custom_error: %v", tc.input, err)
			}
			if !reflect.DeepEqual(stmt, tc.want) {
				t.Fatalf("Parse(%q) == %v, want %v", tc.input, stmt, tc.want)
			}
		})
	}
}

func TestParseErrorCase(t *testing.T) {
	cases := []struct {
		input string
		want  custom_error.SyntaxError
	}{
		{
			input: "",
			want: custom_error.SyntaxError{
				Position: 0,
				Message:  "unexpected end of input",
			},
		},
		{
			input: "select % from temptable",
			want: custom_error.SyntaxError{
				Position: 7,
				Message:  "unexpected character '%'",
			},
		},
		{
			input: "hello x from y",
			want: custom_error.SyntaxError{
				Position: 0,
				Message:  `expected Select, got "hello"`,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			_, err := Parse(c.input)
			if err.Error() != c.want.Error() {
				t.Fatalf("got %v, want %v", err, c.want)
			}
		})
	}
}
