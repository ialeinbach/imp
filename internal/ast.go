package internal

import (
	"encoding/json"
	"imp/backend"
)

type (
	Alias interface {
		Alias() string
		Pos()  int
	}
	RegAlias struct {
		Name string
		Line int
	}
	NumAlias struct {
		Name string
		Line int
	}
	CmdAlias struct {
		Name string
		Line int
	}
)

func (r RegAlias) Alias() string { return r.Name }
func (n NumAlias) Alias() string { return n.Name }
func (c CmdAlias) Alias() string { return c.Name }

func (r RegAlias) Pos() int { return r.Line }
func (n NumAlias) Pos() int { return n.Line }
func (c CmdAlias) Pos() int { return c.Line }

type (
	Stmt interface {
		Gen(*[]backend.Ins, *Scope) error
		Pos() int
	}
	Call struct {
		Cmd  CmdAlias
		Args []Alias
	}
	Decl struct {
		Cmd  CmdAlias
		Args []Alias
		Body []Stmt
	}
)

func (c Call) Pos() int { return c.Cmd.Line }
func (d Decl) Pos() int { return d.Cmd.Line }

// DEBUGGING
func DumpAst(ast []Stmt) string {
	a, err := json.MarshalIndent(ast, "", "\t")
	if err != nil {
		panic(err)
	}
	return string(a)
}
