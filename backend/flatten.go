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
func flatten(prog []Stmt, local *Scope) ([]Ins, error) {
	out := make([]Ins, 0)

	var (
		buf []Ins
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
			case RegAlias:
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
			case RegAlias, NumAlias:
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
func (c Call) Gen(out *[]Ins, local *Scope) error {
	// Look for Cmd in local scope.
	if entry, err := local.Lookup(c.Cmd); err == nil {
		cmd, ok := entry.(Cmd)
		if !ok {
			return errors.New("cmd lookup returned non-cmd entry")
		}
		args, err := Typecheck(c.Args, cmd.Params, *local)
		if err != nil {
			return err
		}
		*out = append(*out, GenCall(c.Cmd.Alias(), cmd, args)...)
		return nil
	}

	// Look for Cmd as builtin.
	if fn, ok := Builtin[c.Cmd.Alias()]; ok {
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
func (d Decl) Gen(out *[]Ins, local *Scope) error {
	// Create parameter template for type checking call arguments.
	params := make([]Psuedo, len(d.Args))
	for i, arg := range d.Args {
		switch arg.(type) {
		case RegAlias:
			params[i] = Reg(0)
		case NumAlias:
			params[i] = Num(0)
		case CmdAlias:
			return errors.Unsupported("cmds as arguments")
		}
	}

	// Create entry and add to current scope.
	cmd := Cmd{
		Addr:   Num(len(*out)+1),
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
	*out = append(*out, Ins{
		Name: "JUMP_I",
		Args: []Psuedo{ Num(len(body)+len(*out)+1) },
	})
	*out = append(*out, body[0].WithComment("start of " + d.Cmd.Alias()))
	*out = append(*out, body[1:]...)
	return nil
}
