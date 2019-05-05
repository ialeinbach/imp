package frontend

import (
	"fmt"
	"imp/errors"
	"strings"
)

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
	case Call:
		return fmt.Sprintf(
			"name: \"%s\"\n" +
			"type: %s\n" +
			"line: %d\n" +
			"args: [%s]\n",
			stmt,
			stmt.Type(),
			stmt.Pos(),
			errors.Indent(DumpArgs(stmt.Args)),
		)
	case Decl:
		return fmt.Sprintf(
			"name: \"%s\"\n" +
			"type: %s\n" +
			"line: %d\n" +
			"params: [%s]\n" +
			"body: [\n%s]\n",
			stmt,
			stmt.Type(),
			stmt.Pos(),
			errors.Indent(DumpArgs(stmt.Params)),
			errors.Indent(DumpAst(stmt.Body)),
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

