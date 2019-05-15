package backend

import (
	"imp/errors"
	"imp/frontend"
	"strings"
	"fmt"
	"strconv"
)

type scope struct {
	name string
	cmds map[string]Cmd
	regs map[string]Reg
	nums map[string]Reg
}

func newScope(name string) *scope {
	return &scope{
		name: name,
		cmds: make(map[string]Cmd),
		regs: make(map[string]Reg),
		nums: make(map[string]Reg),
	}
}

func (s scope) String() string {
	var b strings.Builder

	b.WriteString("====================\n")
	b.WriteString(fmt.Sprintf("  Scope: %s\n", s.name))
	b.WriteString("--------------------\n")

	b.WriteString("  Registers\n")
	for k, v := range s.regs {
		b.WriteString(fmt.Sprintf("    @%s = %v\n", k, v))
	}

	b.WriteString("--------------------\n")

	b.WriteString("  Commands\n")
	for k, v := range s.cmds {
		b.WriteString(fmt.Sprintf("    :%s = %v\n", k, v))
	}

	b.WriteString("====================\n")

	return b.String()
}

func globalScope() *scope {
	return &scope{
		name: "__global__",
		cmds: make(map[string]Cmd),
		regs: map[string]Reg{
			"0": Reg(0),
			"1": Reg(1),
			"2": Reg(2),
			"3": Reg(3),
			"4": Reg(4),
			"5": Reg(5),
			"6": Reg(6),
			"7": Reg(7),
		},
		nums: make(map[string]Reg),
	}
}

func innerScope(context frontend.Decl) (*scope, error) {
	local := newScope(context.String())
	for i, param := range context.Params {
		switch param := param.(type) {
		case frontend.RegAlias:
			local.regs[param.String()] = Reg(i)
		case frontend.NumAlias:
			local.nums[param.String()] = Reg(i)
		default:
			return nil, errors.Unsupported("%s arguments", param.Type())
		}
	}
	return local, nil
}

func (s *scope) lookup(alias frontend.Alias) (Psuedo, error) {
	switch alias := alias.(type) {
	case frontend.CmdAlias:
		if cmd, ok := s.cmds[alias.String()]; ok {
			return cmd, nil
		}
	case frontend.RegAlias:
		if reg, ok := s.regs[alias.String()]; ok {
			return reg, nil
		}
	case frontend.NumAlias:
		// Always treat parseable numbers as numbers.
		num, err := strconv.ParseInt(alias.String(), 0, 0)
		if err == nil {
			return Num(num), nil
		}

		// Otherwise, check if it's a saved alias.
		if reg, ok := s.nums[alias.String()]; ok {
			return reg, nil
		}
	}
	return nil, errors.Undefined(alias)
}

func (s *scope) define(name string, cmd Cmd) {
	delete(s.cmds, name)
	s.cmds[name] = cmd
}

// Checks args for proper typing according to params. If type checking succeeds,
// returns slice of values associated with aliases in some local scope. If
// params == nil, there are no type restrictions.
func (s *scope) typecheck(args []frontend.Alias, params []Psuedo) ([]Psuedo, error) {
	out := make([]Psuedo, len(args))

	// No type restrictions imposed, so just fetch values from local scope.
	if params == nil {
		for i, arg := range args {
			psuedo, err := s.lookup(arg)
			if err != nil {
				return nil, errors.Undefined(arg)
			}
			out[i] = psuedo
		}
		return out, nil
	}

	// Check argument count.
	if len(params) != len(args) {
		return nil, errors.CountMismatch(len(params), len(args))
	}

	// Check argument types against param types and fetch values from local
	// scope.
	for i, param := range params {
		switch param := param.(type) {
		case Reg:
			switch arg := args[i].(type) {
			case frontend.RegAlias:
				psuedo, err := s.lookup(arg)
				if err != nil {
					return nil, errors.Undefined(arg)
				}
				out[i] = psuedo
			default:
				return nil, errors.TypeMismatch(param, arg)
			}
		case Num:
			switch arg := args[i].(type) {
			case frontend.RegAlias, frontend.NumAlias:
				psuedo, err := s.lookup(args[i])
				if err != nil {
					return nil, errors.Undefined(arg)
				}
				out[i] = psuedo
			default:
				return nil, errors.TypeMismatch(param, arg)
			}
		default:
			return nil, errors.Unsupported("%s arguments", param.Type())
		}
	}

	return out, nil
}
