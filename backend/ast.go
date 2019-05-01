package backend

import (
	"encoding/json"
)

type (
	Alias interface {
		Alias() string
		Pos()  int
	}
	regAlias struct {
		name string
		line int
	}
	numAlias struct {
		name string
		line int
	}
	cmdAlias struct {
		name string
		line int
	}
)

func (r regAlias) Alias() string { return r.name }
func (n numAlias) Alias() string { return n.name }
func (c cmdAlias) Alias() string { return c.name }

func (r regAlias) Pos() int { return r.line }
func (n numAlias) Pos() int { return n.line }
func (c cmdAlias) Pos() int { return c.line }

func RegAlias(name string, line int) regAlias {
	return regAlias{
		name: name,
		line: line,
	}
}

func NumAlias(name string, line int) numAlias {
	return numAlias{
		name: name,
		line: line,
	}
}

func CmdAlias(name string, line int) cmdAlias {
	return cmdAlias{
		name: name,
		line: line,
	}
}

type (
	Stmt interface {
		Gen(*[]Ins, *Scope) error
		Pos() int
	}
	call struct {
		cmd  cmdAlias
		args []Alias
	}
	decl struct {
		cmd  cmdAlias
		args []Alias
		body []Stmt
	}
)

func (c call) Pos() int { return c.cmd.Pos() }
func (d decl) Pos() int { return d.cmd.Pos() }

func Call(cmd cmdAlias, args []Alias) Stmt {
	return Stmt(call{
		cmd:  cmd,
		args: args,
	})
}

func Decl(cmd cmdAlias, args []Alias, body []Stmt) Stmt {
	return Stmt(decl{
		cmd:  cmd,
		args: args,
		body: body,
	})
}

func DumpAst(ast []Stmt) string {
	a, err := json.MarshalIndent(ast, "", "\t")
	if err != nil {
		panic(err)
	}
	return string(a)
}
