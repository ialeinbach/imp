package frontend

type (
	Reg struct {
		Alias string
		Line  int
	}
	Num struct {
		Value int64
		Line  int
	}
	Call struct {
		Cmd  string
		Args []Arg
		Line int
	}
	Decl struct {
		Cmd  string
		Args []Arg
		Body []Stmt
		Line int
	}
)

type (
	Arg interface {
		Arg()
		Pos() int
	}
	Stmt interface {
		Stmt()
		Pos() int
	}
)

func (r Reg) Arg() {}
func (n Num) Arg() {}

func (c Call) Stmt() {}
func (d Decl) Stmt() {}

func (r Reg) Pos() int {
	return r.Line
}

func (n Num) Pos() int {
	return n.Line
}

func (c Call) Pos() int {
	return c.Line
}

func (d Decl) Pos() int {
	return d.Line
}
