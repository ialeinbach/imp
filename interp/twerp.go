package main

import (
	"errors"
	"imp/backend"
	"fmt"
	"os"
	"strings"
	"strconv"
	"bufio"
)

// Represents a twerp instruction.
type ins func(*twerp,[]backend.Psuedo) error

const twerpUsage string = `Commands:
"h", "help": print this message
"n", "next": execute an instruction
"r", "regs": print the contents of the registers
"s", "stack": print the contents of the stack
"i", "ip": print the instruction pointer
"q", "quit": quit
"c", "continue": leave interactive mode and continue execution
"p", "prog": print the loaded program
`

type twerp struct {
	regs  []int64
	stack []int64
	prog  []backend.Ins
	ip    int64
}

func NewTwerp(prog []backend.Ins, regs int) *twerp {
	return &twerp{
		prog:  prog,
		regs:  make([]int64, regs),
		stack: []int64{},
	}
}

func (t *twerp) Reset() {
	*t = *NewTwerp(t.prog, len(t.regs))
}

func (t *twerp) Dump() string {
	var b strings.Builder

	b.WriteString(t.dumpIp())
	b.WriteRune('\n')
	b.WriteString(t.dumpProg())
	b.WriteRune('\n')
	b.WriteString(t.dumpStack())
	b.WriteRune('\n')
	b.WriteString(t.dumpRegs())
	b.WriteRune('\n')

	return b.String()
}

func (t *twerp) dumpIp() string {
	return fmt.Sprintln(t.ip)
}

func (t *twerp) dumpProg() string {
	return backend.DumpPsuedo(t.prog)
}

func (t *twerp) dumpStack() string {
	var b strings.Builder

	for _, item := range t.stack {
		b.WriteString(strconv.FormatInt(item, 10))
		b.WriteRune('\n')
	}

	return b.String()
}

func (t *twerp) dumpRegs() string {
	var b strings.Builder

	for _, reg := range t.regs {
		b.WriteString(strconv.FormatInt(reg, 10))
		b.WriteRune('\n')
	}

	return b.String()
}

func (t *twerp) fetch() backend.Ins {
	return t.prog[t.ip]
}

func (t *twerp) ret() int64 {
	return t.regs[0]
}

func (t *twerp) pop() (i int64, err error) {
	if top := len(t.stack) - 1; top >= 0 {
		i, t.stack = t.stack[top], t.stack[:top]
	} else {
		err = errors.New("cannot pop empty stack")
	}
	return
}

func (t *twerp) push(i int64) {
	t.stack = append(t.stack, i)
}

// Executes loaded program.
func (t *twerp) Exec(interactive bool) (int64, error) {
	var (
		decoded ins
		fetched backend.Ins

		ctrl *bufio.Scanner
	)

	if interactive {
		ctrl = bufio.NewScanner(os.Stdin)
	}
	for int(t.ip) < len(t.prog) {
InteractLoop:
		for interactive {
			if fmt.Print("> "); ctrl.Scan() {
				switch ctrl.Text() {
				case "n", "next":
					break InteractLoop
				case "r", "regs":
					fmt.Println(t.dumpRegs())
				case "s", "stack":
					fmt.Println(t.dumpStack())
				case "p", "prog", "program":
					fmt.Println(t.dumpProg())
				case "i", "ip":
					fmt.Println(t.dumpIp())
				case "q", "quit":
					return t.ret(), errors.New("forced quit")
				case "c", "continue":
					interactive = false
				case "h", "help":
					fmt.Println(twerpUsage)
				default:
					fmt.Println("Unrecognized command. Quit with \"q\" or continue with \"c\".")
				}
			} else {
				fmt.Printf("scanning error: %s\n", ctrl.Err())
			}
		}

		switch fetched = t.fetch(); fetched.Name {
		case "MOVE_I": decoded = (*twerp).MoveI
		case "MOVE_R": decoded = (*twerp).MoveR
		case "ADD_I":  decoded = (*twerp).AddI
		case "ADD_R":  decoded = (*twerp).AddR
		case "SUB_I":  decoded = (*twerp).SubI
		case "SUB_R":  decoded = (*twerp).SubR
		case "RET":    decoded = (*twerp).Ret
		case "JUMP_I": decoded = (*twerp).JumpI
		case "CALL_I": decoded = (*twerp).CallI
		case "PUSH_R": decoded = (*twerp).PushR
		case "POP_R":  decoded = (*twerp).PopR
		case "BNE_R":  decoded = (*twerp).BneR
		case "BNE_I":  decoded = (*twerp).BneI
		default:
			return t.ret(), errors.New("fetched not recognized: " + fetched.Name)
		}
		if err := decoded(t, fetched.Args); err != nil {
			return t.ret(), errors.New("error executing " + fetched.Name + ": " + err.Error())
		}
	}
	return t.ret(), nil
}

func (t *twerp) BneR(args []backend.Psuedo) (err error) {
	r0 := t.regs[int(args[0].(backend.Reg))]
	r1 := t.regs[int(args[1].(backend.Reg))]
	if r0 != r1 {
		t.ip = int64(args[2].(backend.Num))
	} else {
		t.ip++
	}
	return
}

func (t *twerp) BneI(args []backend.Psuedo) (err error) {
	n0 := int64(args[0].(backend.Num))
	r1 := t.regs[int(args[1].(backend.Reg))]
	if n0 != r1 {
		t.ip = int64(args[2].(backend.Num))
	} else {
		t.ip++
	}
	return
}

func (t *twerp) MoveR(args []backend.Psuedo) (err error) {
	t.regs[int(args[1].(backend.Reg))] = t.regs[int(args[0].(backend.Reg))]
	t.ip++
	return
}

func (t *twerp) MoveI(args []backend.Psuedo) (err error) {
	t.regs[int(args[1].(backend.Reg))] = int64(args[0].(backend.Num))
	t.ip++
	return
}

func (t *twerp) AddR(args []backend.Psuedo) (err error) {
	t.regs[int(args[1].(backend.Reg))] += int64(t.regs[int(args[0].(backend.Reg))])
	t.ip++
	return
}

func (t *twerp) AddI(args []backend.Psuedo) (err error) {
	t.regs[int(args[1].(backend.Reg))] += int64(args[0].(backend.Num))
	t.ip++
	return
}

func (t *twerp) SubR(args []backend.Psuedo) (err error) {
	t.regs[int(args[1].(backend.Reg))] -= int64(t.regs[int(args[0].(backend.Reg))])
	t.ip++
	return
}

func (t *twerp) SubI(args []backend.Psuedo) (err error) {
	t.regs[int(args[1].(backend.Reg))] -= int64(args[0].(backend.Num))
	t.ip++
	return
}

func (t *twerp) Ret(args []backend.Psuedo) (err error) {
	t.ip, err = t.pop()
	return
}

func (t *twerp) JumpI(args []backend.Psuedo) (err error) {
	t.ip = int64(args[0].(backend.Num))
	return
}

func (t *twerp) CallI(args []backend.Psuedo) (err error) {
	addr := int64(args[0].(backend.Num))
	t.push(t.ip + 1)
	t.ip = addr
	return
}

func (t *twerp) PushI(args []backend.Psuedo) (err error) {
	t.push(int64(args[0].(backend.Num)))
	t.ip++
	return
}

func (t *twerp) PushR(args []backend.Psuedo) (err error) {
	t.push(t.regs[int(args[0].(backend.Reg))])
	t.ip++
	return
}

func (t *twerp) PopR(args []backend.Psuedo) (err error) {
	t.regs[int(args[0].(backend.Reg))], err = t.pop()
	t.ip++
	return
}
