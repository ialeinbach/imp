package internal

import (
	"imp/backend"
	"imp/errors"
)

// Wrapper for flatten().
func Flatten(prog []Stmt) ([]backend.Ins, error) {
	return flatten(prog, GlobalScope())
}

// Converts an AST into a flat list of psuedo-instructions.
func flatten(prog []Stmt, local *Scope) ([]backend.Ins, error) {
	out := make([]backend.Ins, 0)

	var (
		buf []backend.Ins
		err error
	)
	for _, stmt := range prog {
		switch stmt := stmt.(type) {
		case Call:
			err = stmt.Gen(&out, local)
			if err != nil {
				return buf, errors.New("call to %s on line %d: %s", stmt.Cmd.Alias(), stmt.Pos(), err)
			}
		case Decl:
			err = stmt.Gen(&out, local)
			if err != nil {
				return buf, errors.New("decl of %s on line %d: %s", stmt.Cmd.Alias(), stmt.Pos(), err)
			}
		}
		out = append(out, buf...)
	}

	return out, nil
}

// Checks argAliases for proper typing according to params. If type checking
// succeeds, returns slice of values associated with aliases in some local
// scope. If params == nil, there are no type restrictions.
func Typecheck(args []Alias, params []backend.Psuedo, local Scope) ([]backend.Psuedo, error) {
	out := make([]backend.Psuedo, len(args))

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
		case backend.Reg:
			switch args[i].(type) {
			case RegAlias:
				psuedo, err := local.Lookup(args[i])
				if err != nil {
					return nil, errors.Undefined(args[i].Alias())
				}
				out[i] = psuedo
			default:
				return nil, errors.TypeExpected("register")
			}
		case backend.Num:
			switch args[i].(type) {
			case RegAlias, NumAlias:
				psuedo, err := local.Lookup(args[i])
				if err != nil {
					return nil, errors.Undefined(args[i].Alias())
				}
				out[i] = psuedo
			default:
				return nil, errors.TypeExpected("register or number")
			}
		case backend.Cmd:
			return nil, errors.Unsupported("cmds as arguments")
		}
	}

	return out, nil
}

// Generates psuedo-instructions for a call.
func (c Call) Gen(out *[]backend.Ins, local *Scope) error {
	// Look for Cmd in local scope.
	if entry, err := local.Lookup(c.Cmd); err == nil {
		cmd, ok := entry.(backend.Cmd)
		if !ok {
			return errors.New("cmd lookup returned non-cmd entry")
		}
		args, err := Typecheck(c.Args, cmd.Params, *local)
		if err != nil {
			return err
		}
		*out = append(*out, backend.Call(c.Cmd.Alias(), cmd, args)...)
		return nil
	}

	// Look for Cmd as builtin.
	if fn, ok := backend.Builtin[c.Cmd.Alias()]; ok {
		args, err := Typecheck(c.Args, nil, *local)
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

	return errors.Undefined(c.Cmd.Alias())
}

// Generates psuedo-instructions for a declaration.
func (d Decl) Gen(out *[]backend.Ins, local *Scope) error {
	// Create parameter template for type checking call arguments.
	params := make([]backend.Psuedo, len(d.Args))
	for i, arg := range d.Args {
		switch arg.(type) {
		case RegAlias:
			params[i] = backend.Reg(0)
		case NumAlias:
			params[i] = backend.Num(0)
		case CmdAlias:
			return errors.Unsupported("cmds as arguments")
		}
	}

	// Create entry and add to current scope.
	cmd := backend.Cmd{
		Addr:   backend.Num(len(*out)+1),
		Params: params,
	}
	err := local.Insert(d.Cmd, cmd)
	if err != nil {
		return err
	}

	// Create inner scope for declaration body.
	inner, err := d.LocalScope(d.Cmd.Alias())
	if err != nil {
		return err
	}

	// Generate psuedo-instructions for declaration body.
	body, err := flatten(d.Body, inner)
	if err != nil {
		return err
	}
	*out = append(*out, backend.Ins{
		Name: "JUMP_I",
		Args: []backend.Psuedo{ backend.Num(len(body)+len(*out)+1) },
	})
	*out = append(*out, body[0].WithComment("start of " + d.Cmd.Alias()))
	*out = append(*out, body[1:]...)
	return nil
}
