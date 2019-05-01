package backend

import (
	"encoding/json"
	"fmt"
	"strings"
)

func DumpAst(ast []Stmt) string {
	a, err := json.MarshalIndent(ast, "", "\t")
	if err != nil {
		panic(err)
	}
	return string(a)
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
