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

func (d *decl) LocalScope(name string) (*Scope, error) {
	local := NewScope(name)
	for i, alias := range d.args {
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
