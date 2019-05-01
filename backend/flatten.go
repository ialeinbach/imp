package backend

import (
	"imp/errors"
)

// Wrapper for flatten().
func Flatten(prog []Stmt) (out []Ins, err error) {
	out, err = flatten(prog, GlobalScope())
	errors.DebugBackend(1, true, DumpPsuedo(out) + "\n")
	return
}

// Converts an AST into a flat list of psuedo-instructions.
func flatten(prog []Stmt, local *Scope) (out []Ins, err error) {
	out = []Ins{}

	for _, stmt := range prog {
		err = stmt.Gen(&out, local)
		if err != nil {
			return
		}
	}

	return
}

// Checks argAliases for proper typing according to params. If type checking
// succeeds, returns slice of values associated with aliases in some local
// scope. If params == nil, there are no type restrictions.
func Typecheck(args []Alias, params []Psuedo, local Scope) ([]Psuedo, error) {
	out := make([]Psuedo, len(args))

	// No type restrictions imposed, so just fetch values from local scope.
	if params == nil {
		for i, arg := range args {
			psuedo, err := local.Lookup(arg)
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
				psuedo, err := local.Lookup(args[i])
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
				psuedo, err := local.Lookup(args[i])
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
		args, err := Typecheck(c.args, cmd.Params, *local)
		if err != nil {
			return err
		}
		*out = append(*out, GenCall(c.cmd.Alias(), cmd, args)...)
		return nil
	}

	// Look for Cmd as builtin.
	if fn, ok := Builtin[c.cmd.Alias()]; ok {
		args, err := Typecheck(c.args, nil, *local)
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
