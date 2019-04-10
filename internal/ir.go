package internal

type (
	Reg struct { Val int }
	Num struct { Val int64 }

	RegAlias struct { Name string }
	NumAlias struct { Name string }

	Call struct {
		Cmd  string
		Args []Arg
	}
	Decl struct {
		Cmd  string
		Params []Param
		Body []Stmt
	}
)

type (
	Arg   interface{ Arg() }
	Param interface{ Param() }
	Stmt  interface{ Stmt() }
)

func (r Reg) Arg()      {}
func (n Num) Arg()      {}
func (r RegAlias) Arg() {}
func (n NumAlias) Arg() {}

func (r RegAlias) Param() {}
func (n NumAlias) Param() {}

func (c Call) Stmt() {}
func (d Decl) Stmt() {}
