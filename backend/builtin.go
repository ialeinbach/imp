package backend

import (
	"github.com/ialeinbach/imp/errors"
)

type genFn func(*gen, ...Psuedo) (int, error)

var builtins map[string]genFn

func init() {
	builtins = map[string]genFn{
		"add": (*gen).add,
		"sub": (*gen).sub,
		"mov": (*gen).mov,
		"ret": (*gen).ret,
		"rec": (*gen).rec,
	}
}

func (g *gen) rec(args ...Psuedo) (int, error) {
	if len(args) != 0 && len(args) != 2 {
		return 0, errors.New("ret expects either 0 or 2 arguments")
	}

	var n int
	if len(args) == 2 {
		right, ok := args[1].(Reg)
		if !ok {
			return 0, errors.New("right argument of ret must be a register")
		}
		switch left := args[0].(type) {
		case Reg:
			n = g.emit(Ins{
				Name: "BNE_R",
				Args: []Psuedo{ left, right, g.here()+2 },
			})
		case Num:
			n = g.emit(Ins{
				Name: "BNE_I",
				Args: []Psuedo{ left, right, g.here()+2 },
			})
		default:
			return 0, errors.New("left argument of ret must be a register or number")
		}
	}
	n += g.emit(Ins{
		Name: "CALL_I",
		Args: []Psuedo{ g.context().Addr },
	})

	return n, nil
}

func (g *gen) ret(args ...Psuedo) (int, error) {
	if len(args) == 0 {
		return g.emit(Ins{ Name: "RET" }), nil
	}
	if len(args) != 2 {
		return 0, errors.New("ret expects either 0 or 2 arguments")
	}

	right, ok := args[1].(Reg)
	if !ok {
		return 0, errors.New("right argument of ret must be a register")
	}

	var n int
	switch left := args[0].(type) {
	case Reg:
		n = g.emit(Ins{
			Name: "BNE_R",
			Args: []Psuedo{ left, right, g.here()+2 },
		})
	case Num:
		n = g.emit(Ins{
			Name: "BNE_I",
			Args: []Psuedo{ left, right, g.here()+2 },
		})
	default:
		return 0, errors.New("left argument of ret must be a register or number")
	}
	n += g.emit(Ins{ Name: "RET" })
	return n, nil
}

func (g *gen) mov(args ...Psuedo) (int, error) {
	if len(args) != 2 {
		return 0, errors.New("mov expects 2 arguments")
	}

	dst, ok := args[1].(Reg)
	if !ok {
		return 0, errors.New("dst argument of mov must be a register")
	}

	var n int
	switch src := args[0].(type) {
	case Reg:
		n = g.emit(Ins{
			Name: "MOVE_R",
			Args: []Psuedo{ src, dst },
		})
	case Num:
		n = g.emit(Ins{
			Name: "MOVE_I",
			Args: []Psuedo{ src, dst },
		})
	default:
		return 0, errors.New("src argument of mov must be a register or number")
	}
	return n, nil
}

func (g *gen) add(args ...Psuedo) (int, error) {
	if len(args) != 2 {
		return 0, errors.New("add expects 2 arguments")
	}

	dst, ok := args[1].(Reg)
	if !ok {
		return 0, errors.New("dst argument of add must be a register")
	}

	var n int
	switch src := args[0].(type) {
	case Reg:
		n = g.emit(Ins{
			Name: "ADD_R",
			Args: []Psuedo{ src, dst },
		})
	case Num:
		n = g.emit(Ins{
			Name: "ADD_I",
			Args: []Psuedo{ src, dst },
		})
	default:
		return 0, errors.New("src argument of add must be a register or number")
	}
	return n, nil
}

func (g *gen) sub(args ...Psuedo) (int, error) {
	if len(args) != 2 {
		return 0, errors.New("sub expects 2 arguments")
	}

	dst, ok := args[1].(Reg)
	if !ok {
		return 0, errors.New("dst argument of sub must be a register")
	}

	var n int
	switch src := args[0].(type) {
	case Reg:
		n = g.emit(Ins{
			Name: "SUB_R",
			Args: []Psuedo{ src, dst },
		})
	case Num:
		n = g.emit(Ins{
			Name: "SUB_I",
			Args: []Psuedo{ src, dst },
		})
	default:
		return 0, errors.New("src argument of sub must be a register or number")
	}
	return n, nil
}
