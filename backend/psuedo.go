package backend

import (
	"fmt"
	"strings"
)

type (
	Psuedo interface {
		Psuedo()
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

type Ins struct {
	Name    string
	Args    []Psuedo
	Comment string
}

func (i Ins) WithComment(comment string) Ins {
	i.Comment = comment
	return i
}

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

func (i Ins) String() string {
	var b strings.Builder
	b.WriteString(i.Name)
	for _, arg := range i.Args {
		b.WriteString(fmt.Sprintf(" %v", arg))
	}
	return b.String()
}

func DumpPsuedo(psuedo []Ins) string {
	var b strings.Builder
	for i, ins := range psuedo {
		b.WriteString(fmt.Sprintf("%2d: %s", i, ins))
		if len(ins.Comment) > 0 {
			b.WriteString(fmt.Sprintf("    # %s", ins.Comment))
		}
		b.WriteRune('\n')
	}
	return b.String()
}
