package backend

import (
	"imp/errors"
	"strings"
	"fmt"
	"strconv"
)

type Scope struct {
	Name string
	Cmds map[string]Cmd
	Regs map[string]Reg
	Nums map[string]Reg
}

func NewScope(name string) *Scope {
	return &Scope{
		Name: name,
		Cmds: make(map[string]Cmd),
		Regs: make(map[string]Reg),
		Nums: make(map[string]Reg),
	}
}

func (s Scope) String() string {
	var b strings.Builder
	b.WriteString("====================\n")

	b.WriteString(fmt.Sprintf("  Scope: %s\n", s.Name))

	b.WriteString("--------------------\n")

	b.WriteString("  Registers\n")
	for k, v := range s.Regs {
		b.WriteString(fmt.Sprintf("    @%s = %v\n", k, v))
	}

	b.WriteString("--------------------\n")

	b.WriteString("  Commands\n")
	for k, v := range s.Cmds {
		b.WriteString(fmt.Sprintf("    :%s = %v\n", k, v))
	}

	b.WriteString("====================\n")
	return b.String()
}

func GlobalScope() *Scope {
	return &Scope{
		Name: "Global",
		Cmds: make(map[string]Cmd),
		Regs: map[string]Reg{
			"0": Reg(0),
			"1": Reg(1),
			"2": Reg(2),
			"3": Reg(3),
			"4": Reg(4),
			"5": Reg(5),
			"6": Reg(6),
			"7": Reg(7),
		},
		Nums: make(map[string]Reg),
	}
}

func LocalScope(context decl) (*Scope, error) {
	local := NewScope(context.cmd.Alias())
	for i, alias := range context.args {
		switch alias := alias.(type) {
		case cmdAlias:
			return nil, errors.Unsupported("cmds as arguments")
		case regAlias:
			local.Regs[alias.Alias()] = Reg(i)
		case numAlias:
			local.Nums[alias.Alias()] = Reg(i)
		}
	}
	return local, nil
}

func (s *Scope) Lookup(alias Alias) (Psuedo, error) {
	switch alias := alias.(type) {
	case cmdAlias:
		if cmd, ok := s.Cmds[alias.Alias()]; ok {
			return cmd, nil
		}
	case regAlias:
		if reg, ok := s.Regs[alias.Alias()]; ok {
			return reg, nil
		}
	case numAlias:
		// Always treat parseable numbers as numbers.
		num, err := strconv.ParseInt(alias.Alias(), 0, 0)
		if err == nil {
			return Num(num), nil
		}

		// Otherwise, check if it's a saved alias.
		if reg, ok := s.Nums[alias.Alias()]; ok {
			return reg, nil
		}
	default:
		return nil, errors.New("unrecognized alias type: %t", alias)
	}
	return nil, errors.New("undefined alias: %s", alias)
}

func (s *Scope) Insert(alias Alias, psuedo Psuedo) error {
	switch alias := alias.(type) {
	case cmdAlias:
		switch psuedo := psuedo.(type) {
		case Cmd:
			delete(s.Cmds, alias.Alias())
			s.Cmds[alias.Alias()] = psuedo
		case Reg:
			return errors.TypeMismatch("CmdAlias", "Reg")
		case Num:
			return errors.TypeMismatch("CmdAlias", "Num")
		}
	case regAlias:
		switch psuedo := psuedo.(type) {
		case Cmd:
			return errors.TypeMismatch("RegAlias", "Cmd")
		case Reg:
			delete(s.Regs, alias.Alias())
			s.Regs[alias.Alias()] = psuedo
		case Num:
			return errors.TypeMismatch("RegAlias", "Num")
		}
	case numAlias:
		return errors.Unsupported("nums in scopes")
	}
	return nil
}

// Checks args for proper typing according to params. If type checking succeeds,
// returns slice of values associated with aliases in some local scope. If
// params == nil, there are no type restrictions.
func (s *Scope) typecheck(args []Alias, params []Psuedo) ([]Psuedo, error) {
	out := make([]Psuedo, len(args))

	// No type restrictions imposed, so just fetch values from local scope.
	if params == nil {
		for i, arg := range args {
			psuedo, err := s.Lookup(arg)
			if err != nil {
				return nil, errors.Undefined(arg.Alias())
			}
			out[i] = psuedo
		}
		return out, nil
	}

	// Check argument count.
	if len(params) != len(args) {
		return nil, errors.New("argument count: expected %d but found %d\n" +
		                       "params: %v\n" +
		                       "args:   %v\n", len(params), len(args), params, args)
	}

	// Check argument types against param types and fetch values from local
	// scope.
	for i, param := range params {
		switch param.(type) {
		case Reg:
			switch args[i].(type) {
			case regAlias:
				psuedo, err := s.Lookup(args[i])
				if err != nil {
					return nil, errors.Undefined(args[i].Alias())
				}
				out[i] = psuedo
			default:
				return nil, errors.TypeExpected("register")
			}
		case Num:
			switch args[i].(type) {
			case regAlias, numAlias:
				psuedo, err := s.Lookup(args[i])
				if err != nil {
					return nil, errors.Undefined(args[i].Alias())
				}
				out[i] = psuedo
			default:
				return nil, errors.TypeExpected("register or number")
			}
		case Cmd:
			return nil, errors.Unsupported("cmds as arguments")
		}
	}

	return out, nil
}
