package backend

import (
	"fmt"
	"strings"
)

func DumpScope(s *scope) string {
	var b strings.Builder

	for k, v := range s.regs {
		b.WriteString(fmt.Sprintf("%v: %v\n", k, v))
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
		b.WriteString("\n")
	}
	if out := b.String(); len(out) > 0 {
		return out[:len(out)-1]
	} else {
		return out
	}
}
