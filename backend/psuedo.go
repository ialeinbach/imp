package backend

import (
	"fmt"
	"strings"
)

//
// Psuedo-values
//

type (
	Psuedo interface {
		Psuedo()
		Type() string
	}
	Reg int
	Num int64
	Cmd struct {
		Addr   Num
		Params []Psuedo
	}
)

func (r Reg) Psuedo() {}
func (n Num) Psuedo() {}
func (c Cmd) Psuedo() {}

func (r Reg) Type() string { return "Reg" }
func (n Num) Type() string { return "Num" }
func (c Cmd) Type() string { return "Cmd" }

func (r Reg) String() string {
	return fmt.Sprint(int(r))
}

func (n Num) String() string {
	return fmt.Sprint(int64(n))
}

func (c Cmd) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("[Addr: %d, Params:", c.Addr))
	for _, typ := range c.Params {
		b.WriteString(fmt.Sprintf(" %v", typ))
	}
	b.WriteString("]\n")
	return b.String()
}

//
// Psuedo-instructions
//

type Ins struct {
	Name    string
	Args    []Psuedo
	Comment string
}

func (i Ins) String() string {
	var b strings.Builder
	b.WriteString(i.Name)
	for _, arg := range i.Args {
		b.WriteString(fmt.Sprintf(" %v", arg))
	}
	return b.String()
}

func (i Ins) WithComment(format string, a ...interface{}) Ins {
	i.Comment = fmt.Sprintf(format, a...)
	return i
}
