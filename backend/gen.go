package backend

import (
	"imp/errors"
)

//
// AST -> Psuedo-instructions
//

// Wraps flatten() and drops in the global scope.
func Flatten(prog []Stmt) ([]Ins, error) {
	out, err := flatten(prog, GlobalScope())
	errors.DebugBackend(1, true, DumpPsuedo(out) + "\n")
	return out, err
}

// Converts an AST into a list of psuedo-instructions.
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

//
// Procedure Call Generation
//

func genProcCall(name string, cmd Cmd, args []Psuedo) []Ins {
	out := []Ins{}

	out = append(out, genProcCallProlog(args)...)
	out = append(out, Ins{
		Name: "CALL_I",
		Args: []Psuedo{ cmd.Addr },
	}.WithComment("call " + name))
	out = append(out, genProcCallEpilog(args)...)

	return out
}

func genProcCallProlog(args []Psuedo) []Ins {
	out := []Ins{}

	// See DepSeqs definition for info about dependency sequences.
	regSeqs, numSeqs := DepSeqs(args)

	// Generate psuedo-instructions for handling dep seqs that start with
	// numbers.
	for num, seq := range numSeqs {
		i := len(seq) - 1
		out = append(out, Ins{
			Name: "PUSH_R",
			Args: []Psuedo{ Reg(seq[i]) },
		})
		for i--; i >= 0; i-- {
			out = append(out, Ins{
				Name: "MOVE_R",
				Args: []Psuedo{ Reg(seq[i]), Reg(seq[i+1]) },
			})
		}
		out = append(out, Ins{
			Name: "MOVE_I",
			Args: []Psuedo{ Num(num), Reg(seq[0]) },
		})
	}

	// Generate psuedo-instructions for handling dep seqs that start with
	// registers.
	for reg, seq := range regSeqs {
		i := len(seq) - 1
		out = append(out, Ins{
			Name: "PUSH_R",
			Args: []Psuedo{ Reg(seq[i]) },
		})
		for i--; i >= 0; i-- {
			out = append(out, Ins{
				Name: "MOVE_R",
				Args: []Psuedo{ Reg(seq[i]), Reg(seq[i+1]) },
			})
		}

		// Handle cyclic dep seqs.
		if seq[len(seq)-1] == reg {
			out = append(out, Ins{
				Name: "POP_R",
				Args: []Psuedo{ Reg(seq[0]) },
			})
		} else {
			out = append(out, Ins{
				Name: "MOVE_R",
				Args: []Psuedo{ Reg(reg), Reg(seq[0]) },
			})
		}
	}

	return out
}

func genProcCallEpilog(args []Psuedo) []Ins {
	out := []Ins{}

	// See DepSeqs definition for info about dependency sequences.
	regSeqs, numSeqs := DepSeqs(args)

	// Generate psuedo-instructions for handling dep seqs that start with
	// numbers.
	for _, seq := range numSeqs {
		for i := 1; i < len(seq); i++ {
			out = append(out, Ins{
				Name: "MOVE_R",
				Args: []Psuedo{ Reg(seq[i]), Reg(seq[i-1]) },
			})
		}
		out = append(out, Ins{
			Name: "POP_R",
			Args: []Psuedo{ Reg(seq[len(seq)-1]) },
		})
	}

	// Generate psuedo-instructions for handling dep seqs that start with
	// registers.
	for reg, seq := range regSeqs {
		// Handle cyclic dep seqs.
		if i := len(seq)-1; reg == seq[i] {
			out = append(out, Ins{
				Name: "PUSH_R",
				Args: []Psuedo{ Reg(seq[0]) },
			})
		} else {
			out = append(out, Ins{
				Name: "MOVE_R",
				Args: []Psuedo{ Reg(seq[0]), Reg(reg) },
			})
		}

		for i := 1; i < len(seq); i++ {
			out = append(out, Ins{
				Name: "MOVE_R",
				Args: []Psuedo{ Reg(seq[i]), Reg(seq[i-1]) },
			})
		}
		out = append(out, Ins{
			Name: "POP_R",
			Args: []Psuedo{ Reg(seq[len(seq)-1]) },
		})
	}

	return out
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
func DepSeqs(args []Psuedo) (regSeqs map[int][]int, numSeqs map[int][]int) {
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
		for {
			dst, free[dst] = regs[dst], false
			if !free[dst] { break }
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
