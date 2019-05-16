package backend

import (
	"imp/frontend"
)

type gen struct{
	scopes []*scope
	code   []Ins
}

func (g *gen) here() Num {
	return Num(len(g.code))
}

func (g *gen) localScope() *scope {
	if len(g.scopes) == 0 {
		return nil
	}
	return g.scopes[len(g.scopes)-1]
}

func (g *gen) context() Cmd {
	return g.localScope().cmds[g.localScope().name]
}

func (g *gen) enterScope(context frontend.Decl) error {
	local, err := innerScope(context)
	if err != nil {
		return err
	}
	g.scopes = append(g.scopes, local)
	return nil
}

func (g *gen) exitScope() {
	g.scopes = g.scopes[:len(g.scopes)-1]
}

func (g *gen) lookup(alias frontend.Alias) (ps Psuedo, err error) {
	for i := len(g.scopes)-1; i >= 0; i-- {
		ps, err = g.scopes[i].lookup(alias)
		if err == nil {
			return
		}
	}
	return // failed lookup returns same error for any scope
}

func (g *gen) define(name string, cmd Cmd) {
	g.localScope().define(name, cmd)
}

func (g *gen) typecheck(args []frontend.Alias, params []Psuedo) ([]Psuedo, error) {
	return g.localScope().typecheck(args, params)
}

func (g *gen) emit(i ...Ins) int {
	g.code = append(g.code, i...)
	return len(i)
}
