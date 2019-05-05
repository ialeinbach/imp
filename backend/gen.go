package backend

import (
	"imp/errors"
	"imp/frontend"
)

type gen struct {
	scopes []*scope
	code   []Ins
}

func (g *gen) local() *scope {
	return g.scopes[len(g.scopes)-1]
}

func (g *gen) enterScope(inner *scope) {
	g.scopes = append(g.scopes, inner)
}

func (g *gen) exitScope() {
	g.scopes = g.scopes[:len(g.scopes)-1]
}

func (g *gen) lookup(alias frontend.Alias) (Psuedo, error) {
	return g.local().lookup(alias)
}

func (g *gen) typecheck(args []frontend.Alias, params []Psuedo) ([]Psuedo, error) {
	return g.local().typecheck(args, params)
}

func (g *gen) emit(i ...Ins) int {
	g.code = append(g.code, i...)
	return len(i)
}

// Returns a flat slice of psuedo-instructions generated from an AST.
func Flatten(prog []frontend.Stmt) ([]Ins, error) {
	g := &gen{
		scopes: []*scope{globalScope()},
		code:   []Ins{},
	}
	_, err := g.flatten(prog)
	if err != nil {
		return nil, err
	}
	errors.DebugBackend(1, true, DumpPsuedo(g.code))
	errors.DebugBackend(1, false, "\n\n")
	return g.code, nil
}

// Converts an AST into a list of psuedo-instructions.
func (g *gen) flatten(prog []frontend.Stmt) (n int, err error) {
	var i int
	for _, stmt := range prog {
		switch stmt := stmt.(type) {
		case frontend.Call:
			if i, err = g.genCall(stmt); err != nil {
				return
			}
		case frontend.Decl:
			if i, err = g.genDecl(stmt); err != nil {
				return
			}
		}
		n += i
	}
	return
}

// Generates psuedo-instructions for a call.
func (g *gen) genCall(call frontend.Call) (n int, err error) {
	// Look for Cmd in local scope.
	if entry, err := g.lookup(call.Cmd); err == nil {
		cmd := entry.(Cmd) // ensured by lookup()
		args, err := g.typecheck(call.Args, cmd.Params)
		if err != nil {
			return 0, errors.Wrap(err, call)
		}
		return g.genProcCall(cmd, args), nil
	}

	// Look for Cmd as builtin.
	if fn, ok := Builtin[call.String()]; ok {
		args, err := g.typecheck(call.Args, nil)
		if err != nil {
			return 0, errors.Wrap(err, call)
		}
		ins, err := fn(args...)
		if err != nil {
			return 0, errors.Wrap(err, call)
		}
		return g.emit(ins...), nil
	}

	return 0, errors.Undefined(call)
}

// Generates psuedo-instructions for a declaration.
func (g *gen) genDecl(decl frontend.Decl) (n int, err error) {
	// Create parameter template for type checking call arguments.
	params := make([]Psuedo, len(decl.Params))
	for i, param := range decl.Params {
		switch param.(type) {
		case frontend.RegAlias:
			params[i] = Reg(0)
		case frontend.NumAlias:
			params[i] = Num(0)
		default:
			err = errors.Unsupported("%s parameters", param.Type())
			return 0, errors.Wrap(err, decl)
		}
	}

	// Addr to be backfilled after decl body size known.
	n += g.emit(Ins{
		Name: "JUMP_I",
	})

	// Create entry and add to current scope.
	cmd := Cmd{
		Addr:   Num(len(g.code)),
		Params: params,
	}
	g.local().define(decl.String(), cmd)

	// Create inner scope for declaration body.
	inner, err := localScope(decl)
	if err != nil {
		return
	}
	g.enterScope(inner)
	defer g.exitScope()

	// Generate psuedo-instructions for declaration body.
	i, err := g.flatten(decl.Body)
	if err != nil {
		return 0, errors.Wrap(err, decl)
	}
	n += i

	// Backfill jump over declaration body.
	g.code[len(g.code)-1-i].Args = []Psuedo{Num(len(g.code))}

	return
}

func (g *gen) genProcCall(cmd Cmd, args []Psuedo) (n int) {
	n += g.genProcCallProlog(args)
	n += g.emit(Ins{
		Name: "CALL_I",
		Args: []Psuedo{ cmd.Addr },
	})
	n += g.genProcCallEpilog(args)
	return
}

func (g *gen) genProcCallProlog(args []Psuedo) (n int) {
	// See depSeqs definition for info about dependency sequences.
	regSeqs, numSeqs := depSeqs(args)

	// Generate psuedo-instructions for handling dep seqs that start with
	// numbers.
	for num, seq := range numSeqs {
		i := len(seq) - 1
		n += g.emit(Ins{
			Name: "PUSH_R",
			Args: []Psuedo{ Reg(seq[i]) },
		})
		for i--; i >= 0; i-- {
			n += g.emit(Ins{
				Name: "MOVE_R",
				Args: []Psuedo{ Reg(seq[i]), Reg(seq[i+1]) },
			})
		}
		n += g.emit(Ins{
			Name: "MOVE_I",
			Args: []Psuedo{ Num(num), Reg(seq[0]) },
		})
	}

	// Generate psuedo-instructions for handling dep seqs that start with
	// registers.
	for reg, seq := range regSeqs {
		i := len(seq) - 1
		n += g.emit(Ins{
			Name: "PUSH_R",
			Args: []Psuedo{ Reg(seq[i]) },
		})
		for i--; i >= 0; i-- {
			n += g.emit(Ins{
				Name: "MOVE_R",
				Args: []Psuedo{ Reg(seq[i]), Reg(seq[i+1]) },
			})
		}

		// Handle cyclic dep seqs.
		if seq[len(seq)-1] == reg {
			n += g.emit(Ins{
				Name: "POP_R",
				Args: []Psuedo{ Reg(seq[0]) },
			})
		} else {
			n += g.emit(Ins{
				Name: "MOVE_R",
				Args: []Psuedo{ Reg(reg), Reg(seq[0]) },
			})
		}
	}

	return
}

func (g *gen) genProcCallEpilog(args []Psuedo) (n int) {
	// See depSeqs definition for info about dependency sequences.
	regSeqs, numSeqs := depSeqs(args)

	// Generate psuedo-instructions for handling dep seqs that start with
	// numbers.
	for _, seq := range numSeqs {
		for i := 1; i < len(seq); i++ {
			n += g.emit(Ins{
				Name: "MOVE_R",
				Args: []Psuedo{ Reg(seq[i]), Reg(seq[i-1]) },
			})
		}
		n += g.emit(Ins{
			Name: "POP_R",
			Args: []Psuedo{ Reg(seq[len(seq)-1]) },
		})
	}

	// Generate psuedo-instructions for handling dep seqs that start with
	// registers.
	for reg, seq := range regSeqs {
		// Handle cyclic dep seqs.
		if i := len(seq)-1; reg == seq[i] {
			n += g.emit(Ins{
				Name: "PUSH_R",
				Args: []Psuedo{ Reg(seq[0]) },
			})
		} else {
			n += g.emit(Ins{
				Name: "MOVE_R",
				Args: []Psuedo{ Reg(seq[0]), Reg(reg) },
			})
		}

		for i := 1; i < len(seq); i++ {
			n += g.emit(Ins{
				Name: "MOVE_R",
				Args: []Psuedo{ Reg(seq[i]), Reg(seq[i-1]) },
			})
		}
		n += g.emit(Ins{
			Name: "POP_R",
			Args: []Psuedo{ Reg(seq[len(seq)-1]) },
		})
	}

	return
}

// Returns "dependency sequences" for generating instructions to perform
// maximally in-place, stack-assisted register reorderings that occurs in call
// prologs/epilogs. A dependency sequence A, B, C means:
//
//   (a) for call prologs, A->B must happen after B->C must happen after C is
//       pushed on the stack, and
//
//   (b) for call epilogs, A<-B must happen before B<-C must happen before the
//       stack is popped into C.
//
// A cyclic dependency sequence like A, B, A is valid and results in a circular
// shift of contents of registers in the sequence. The length of a dependency
// sequence is always at least 2.
func depSeqs(args []Psuedo) (regSeqs map[int][]int, numSeqs map[int][]int) {
	if len(args) == 0 {
		return make(map[int][]int), make(map[int][]int)
	}

	// These helper data structures encode the register transfers that must
	// result from the call prologs/epilogs that are to be generated. For
	// example, if reg R is passed as argument A, then:
	//
	//   (a) after the prolog executes, the contents of reg R must have been
	//       placed into reg A, and
	//
	//   (b) after the epilog executes, the contents of reg A must have been
	//       placed back into reg R.
	//
	// The requirements of this example would be represented in the helper data
	// structures by the fact that reg[R]=A and free[R]=true.
	var (
		// The key-value pair (A, B) means we need to move num A into reg B.
		nums = make(map[int]int)

		// The combination of regs[A]=B and free[A]=true means we need to move
		// reg A's contents into reg B.
		regs = make([]int, MaxRegCount)
		free = make([]bool, MaxRegCount)
	)
	for dst, src := range args {
		switch src := src.(type) {
		case Reg:
			regs[int(src)] = dst
			free[int(src)] = true
		case Num:
			nums[int(src)] = dst
		}
	}

	// The key-value pair (A, B) represents the dependency sequence A, B[0],
	// B[1], ..., B[N] where N == len(B)-1. The keys of numSeqs represent
	// numbers whereas the keys of regSeqs represent registers.
	numSeqs = make(map[int][]int)
	regSeqs = make(map[int][]int)

	// Find dependency sequences starting with numbers. Order does not matter
	// (hence nums is a map) because numbers are trivially guaranteed to start
	// dependency sequences. Put another way, we don't have to:
	//
	//   (a) in the prolog, worry about something needing to be "moved into a
	//       number" (as opposed to "moved into a register"), or
	//
	//   (b) in the epilog, worry about "restoring the contents of a number" (as
	//       opposed to "restoring the contents of a register")
	//
	// because neither of those things even make sense (hence trivial).
	for num, dst := range nums {
		numSeqs[num] = []int{dst}
		for free[dst] {
			dst, free[dst] = regs[dst], false
			numSeqs[num] = append(numSeqs[num], dst)
		}
	}

	// Find dependency sequences starting with registers.
	var reg, dst int
	for {
		// It is important that we look for dependency sequences starting with the
		// highest register as it guarantees we never start in the middle of a
		// dependency sequence. Formal proof of this in the works. Numbers are
		// trivially guaranteed to start dependency sequences because they are
		// just values and cannot be written into.
		for reg = MaxRegCount-1; !free[reg]; reg-- {
			if reg <= 0 { return }
		}
		dst, free[reg] = regs[reg], false

		regSeqs[reg] = []int{dst}
		for free[dst] {
			dst, free[dst] = regs[dst], false
			regSeqs[reg] = append(regSeqs[reg], dst)
		}
	}
}
