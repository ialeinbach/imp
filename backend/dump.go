package backend

import (
	"fmt"
	"imp/errors"
	"strings"
)

func indent(s string) string {
	return errors.PrefixLines(s, "\t")
}

func DumpArg(arg Alias) string {
	return fmt.Sprintf(
		"name: \"%s\"\n" +
		"type: %s\n" +
		"line: %d\n",
		arg,
		arg.Type(),
		arg.Pos(),
	)
}

func DumpArgs(args []Alias) string {
	if len(args) == 0 {
		return ""
	}

	var b strings.Builder

	b.WriteString("\n")
	for i, arg := range args {
		b.WriteString(DumpArg(arg))
		if i < len(args) - 1 {
			b.WriteString("\n")
		}
	}

	return b.String()
}

func DumpStmt(stmt Stmt) string {
	switch stmt := stmt.(type) {
	case call:
		return fmt.Sprintf(
			"name: \"%s\"\n" +
			"type: %s\n" +
			"line: %d\n" +
			"args: [%s]\n",
			stmt,
			stmt.Type(),
			stmt.Pos(),
			indent(DumpArgs(stmt.args)),
		)
	case decl:
		return fmt.Sprintf(
			"name: \"%s\"\n" +
			"type: %s\n" +
			"line: %d\n" +
			"params: [%s]\n" +
			"body: [\n%s]\n",
			stmt,
			stmt.Type(),
			stmt.Pos(),
			indent(DumpArgs(stmt.args)),
			indent(DumpAst(stmt.body)),
		)
	}
	return ""
}

func DumpAst(ast []Stmt) string {
	var b strings.Builder

	for _, stmt := range ast {
		b.WriteString(DumpStmt(stmt))
		b.WriteString("\n")
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
