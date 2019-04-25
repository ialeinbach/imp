package backend

import (
	"imp/errors"
)

type BuiltinFn func(...Psuedo) ([]Ins, error)

var Builtin map[string]BuiltinFn

func init() {
	Builtin = map[string]BuiltinFn{
		"add": Add,
		"mov": Move,
		"ret": Ret,
	}
}

func Ret(args ...Psuedo) ([]Ins, error) {
	if len(args) != 0 {
		return nil, errors.New("ret expects 0 arguments")
	}
	return []Ins{Ins{ Name: "RET" }}, nil
}

func Move(args ...Psuedo) ([]Ins, error) {
	if len(args) != 2 {
		return nil, errors.New("mov expects 2 arguments")
	}
	dst, ok := args[1].(Reg)
	if !ok {
		return nil, errors.New("dst argument of mov must be a register")
	}
	switch src := args[0].(type) {
	case Reg:
		return []Ins{Ins{
			Name: "MOVE_R",
			Args: []Psuedo{ src, dst },
		}}, nil
	case Num:
		return []Ins{Ins{
			Name: "MOVE_I",
			Args: []Psuedo{ src, dst },
		}}, nil
	default:
		return nil, errors.New("src argument of mov must be a register or number")
	}
}

func Add(args ...Psuedo) ([]Ins, error) {
	if len(args) != 2 {
		return nil, errors.New("add expects 2 arguments")
	}
	dst, ok := args[1].(Reg)
	if !ok {
		return nil, errors.New("dst argument of add must be a register")
	}
	switch src := args[0].(type) {
	case Reg:
		return []Ins{Ins{
			Name: "ADD_R",
			Args: []Psuedo{ src, dst },
		}}, nil
	case Num:
		return []Ins{Ins{
			Name: "ADD_I",
			Args: []Psuedo{ src, dst },
		}}, nil
	default:
		return nil, errors.New("src argument of add must be a register or number")
	}
}
