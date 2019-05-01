package backend

import (
	"imp/errors"
)

//
// Aliases
//

type (
	Alias interface {
		Alias() string
		Pos()  int
	}
	regAlias struct {
		name string
		line int
	}
	numAlias struct {
		name string
		line int
	}
	cmdAlias struct {
		name string
		line int
	}
)

func RegAlias(name string, line int) regAlias {
	return regAlias{
		name: name,
		line: line,
	}
}

func NumAlias(name string, line int) numAlias {
	return numAlias{
		name: name,
		line: line,
	}
}

func CmdAlias(name string, line int) cmdAlias {
	return cmdAlias{
		name: name,
		line: line,
	}
}

func (r regAlias) Alias() string { return r.name }
func (n numAlias) Alias() string { return n.name }
func (c cmdAlias) Alias() string { return c.name }

func (r regAlias) Pos() int { return r.line }
func (n numAlias) Pos() int { return n.line }
func (c cmdAlias) Pos() int { return c.line }

//
// Statements
//

type (
	Stmt interface {
		Gen(*[]Ins, *Scope) error
		Pos() int
	}
	call struct {
		cmd  cmdAlias
		args []Alias
	}
	decl struct {
		cmd  cmdAlias
		args []Alias
		body []Stmt
	}
)

func Call(cmd cmdAlias, args []Alias) Stmt {
	return Stmt(call{
		cmd:  cmd,
		args: args,
	})
}

func Decl(cmd cmdAlias, args []Alias, body []Stmt) Stmt {
	return Stmt(decl{
		cmd:  cmd,
		args: args,
		body: body,
	})
}

func (c call) Pos() int { return c.cmd.Pos() }
func (d decl) Pos() int { return d.cmd.Pos() }

// Generates psuedo-instructions for a call.
func (c call) Gen(out *[]Ins, local *Scope) (err error) {
	defer func() {
		if err != nil {
			err = errors.New("call to %s on line %d: %s", c.cmd.Alias(), c.Pos(), err)
		}
		return
	}()

	// Look for Cmd in local scope.
	if entry, err := local.Lookup(c.cmd); err == nil {
		cmd, ok := entry.(Cmd)
		if !ok {
			return errors.New("cmd lookup returned non-cmd entry")
		}
		args, err := local.Typecheck(c.args, cmd.Params)
		if err != nil {
			return err
		}
		*out = append(*out, genProcCall(c.cmd.Alias(), cmd, args)...)
		return nil
	}

	// Look for Cmd as builtin.
	if fn, ok := Builtin[c.cmd.Alias()]; ok {
		args, err := local.Typecheck(c.args, nil)
		if err != nil {
			return err
		}
		ins, err := fn(args...)
		if err != nil {
			return err
		}
		*out = append(*out, ins...)
		return nil
	}

	return errors.Undefined(c.cmd.Alias())
}

// Generates psuedo-instructions for a declaration.
func (d decl) Gen(out *[]Ins, local *Scope) (err error) {
	defer func() {
		if err != nil {
			err = errors.New("decl of %s on line %d: %s", d.cmd.Alias(), d.Pos(), err)
		}
		return
	}()

	// Create parameter template for type checking call arguments.
	params := make([]Psuedo, len(d.args))
	for i, arg := range d.args {
		switch arg.(type) {
		case regAlias:
			params[i] = Reg(0)
		case numAlias:
			params[i] = Num(0)
		case cmdAlias:
			return errors.Unsupported("cmds as arguments")
		}
	}

	// Create entry and add to current scope.
	cmd := Cmd{
		Addr:   Num(len(*out)+1),
		Params: params,
	}
	err = local.Insert(d.cmd, cmd)
	if err != nil {
		return err
	}

	// Create inner scope for declaration body.
	inner, err := d.LocalScope(d.cmd.Alias())
	if err != nil {
		return err
	}

	// Generate psuedo-instructions for declaration body.
	body, err := flatten(d.body, inner)
	if err != nil {
		return err
	}
	*out = append(*out, Ins{
		Name: "JUMP_I",
		Args: []Psuedo{ Num(len(body)+len(*out)+1) },
	})
	*out = append(*out, body[0].WithComment("start of " + d.cmd.Alias()))
	*out = append(*out, body[1:]...)
	return nil
}
