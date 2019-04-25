package internal

import (
	"imp/backend"
	"imp/errors"
	"strings"
	"fmt"
	"strconv"
)

type Scope struct {
	Name string
	Cmds map[string]backend.Cmd
	Regs map[string]backend.Reg
	Nums map[string]backend.Reg
}

func NewScope(name string) *Scope {
	return &Scope{
		Name: name,
		Cmds: make(map[string]backend.Cmd),
		Regs: make(map[string]backend.Reg),
		Nums: make(map[string]backend.Reg),
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
		Cmds: make(map[string]backend.Cmd),
		Regs: map[string]backend.Reg{
			"0": backend.Reg(0),
			"1": backend.Reg(1),
			"2": backend.Reg(2),
			"3": backend.Reg(3),
			"4": backend.Reg(4),
			"5": backend.Reg(5),
			"6": backend.Reg(6),
			"7": backend.Reg(7),
		},
		Nums: make(map[string]backend.Reg),
	}
}

func (d *Decl) LocalScope(name string) (*Scope, error) {
	local := NewScope(name)
	for i, alias := range d.Args {
		switch alias := alias.(type) {
		case CmdAlias:
			return nil, errors.Unsupported("cmds as arguments")
		case RegAlias:
			local.Regs[alias.Alias()] = backend.Reg(i)
		case NumAlias:
			local.Nums[alias.Alias()] = backend.Reg(i)
		}
	}
	return local, nil
}

func (s *Scope) Lookup(alias Alias) (backend.Psuedo, error) {
	switch alias := alias.(type) {
	case CmdAlias:
		if cmd, ok := s.Cmds[alias.Alias()]; ok {
			return cmd, nil
		}
	case RegAlias:
		if reg, ok := s.Regs[alias.Alias()]; ok {
			return reg, nil
		}
	case NumAlias:
		// Always treat parseable numbers as numbers.
		num, err := strconv.ParseInt(alias.Alias(), 0, 0)
		if err == nil {
			return backend.Num(num), nil
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

func (s *Scope) Insert(alias Alias, psuedo backend.Psuedo) error {
	switch alias := alias.(type) {
	case CmdAlias:
		switch psuedo := psuedo.(type) {
		case backend.Cmd:
			delete(s.Cmds, alias.Alias())
			s.Cmds[alias.Alias()] = psuedo
		case backend.Reg:
			return errors.TypeMismatch("CmdAlias", "Reg")
		case backend.Num:
			return errors.TypeMismatch("CmdAlias", "Num")
		}
	case RegAlias:
		switch psuedo := psuedo.(type) {
		case backend.Cmd:
			return errors.TypeMismatch("RegAlias", "Cmd")
		case backend.Reg:
			delete(s.Regs, alias.Alias())
			s.Regs[alias.Alias()] = psuedo
		case backend.Num:
			return errors.TypeMismatch("RegAlias", "Num")
		}
	case NumAlias:
		return errors.Unsupported("nums in scopes")
	}
	return nil
}
