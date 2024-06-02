package query

import (
	"fmt"
	"strings"
)

type Printer struct {
	builder     *strings.Builder
	indentation int
}

func NewPrinter() *Printer {
	return &Printer{builder: new(strings.Builder)}
}

func (p *Printer) Indent() {
	p.indentation++
}

func (p *Printer) Dedent() {
	if p.indentation == 0 {
		return
	}
	p.indentation--
}

func (p *Printer) Println(format string, a ...any) {
	for i := 0; i < p.indentation; i++ {
		fmt.Fprint(p.builder, "   ")
	}
	fmt.Fprintf(p.builder, format, a...)
	fmt.Fprintln(p.builder)

}
