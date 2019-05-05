package frontend

type (
	Alias interface {
		Alias()
		String() string
		Type() string
		Pos() int
	}
	RegAlias struct {
		name string
		line int
	}
	NumAlias struct {
		name string
		line int
	}
	CmdAlias struct {
		name string
		line int
	}
)

func (r RegAlias) Alias() {}
func (n NumAlias) Alias() {}
func (c CmdAlias) Alias() {}

func (r RegAlias) String() string { return r.name }
func (n NumAlias) String() string { return n.name }
func (c CmdAlias) String() string { return c.name }

func (r RegAlias) Type() string { return "RegAlias" }
func (n NumAlias) Type() string { return "NumAlias" }
func (c CmdAlias) Type() string { return "CmdAlias" }

func (r RegAlias) Pos() int { return r.line }
func (n NumAlias) Pos() int { return n.line }
func (c CmdAlias) Pos() int { return c.line }

type (
	Stmt interface {
		Stmt()
		String() string
		Type() string
		Pos() int
	}
	Call struct {
		Cmd  CmdAlias
		Args []Alias
	}
	Decl struct {
		Cmd    CmdAlias
		Params []Alias
		Body   []Stmt
	}
)

func (c Call) Stmt() {}
func (d Decl) Stmt() {}

func (c Call) Pos() int { return c.Cmd.Pos() }
func (d Decl) Pos() int { return d.Cmd.Pos() }

func (c Call) String() string { return c.Cmd.String() }
func (d Decl) String() string { return d.Cmd.String() }

func (c Call) Type() string { return "Call" }
func (d Decl) Type() string { return "Decl" }
