package backend

import (
	"imp/frontend"
)

type gen struct{
	scopes []*scope
	code   []Ins
}

func (g *gen) localScope() *scope {
	if len(g.scopes) == 0 {
		return nil
	}
	return g.scopes[len(g.scopes)-1]
}

func (g *gen) enterScope(s *scope) {
	g.scopes = append(g.scopes, s)
}

func (g *gen) exitScope() {
	g.scopes = g.scopes[:len(g.scopes)-1]
}

func (g *gen) lookup(alias frontend.Alias) (Psuedo, error) {
	return g.localScope().lookup(alias)
}

func (g *gen) define(name string, cmd Cmd) {
	delete(g.localScope().cmds, name)
	g.localScope().cmds[name] = cmd
}

func (g *gen) typecheck(args []frontend.Alias, params []Psuedo) ([]Psuedo, error) {
	return g.localScope().typecheck(args, params)
}

func (g *gen) emit(i ...Ins) int {
	g.code = append(g.code, i...)
	return len(i)
}
