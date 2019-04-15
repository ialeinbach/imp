%{

package frontend

import (
	"fmt"
	"strconv"
	"imp/internal"
)

const MAX_ARG_COUNT int = 6

var ParserVerbosity int

%}

%union{
	tok       Token
	arglist   []internal.Arg
	stmtlist  []internal.Stmt
}

%token <tok> CMD REG NUM CR

%type <arglist> arg args
%type <stmtlist> decl call stmt program main

%start main

%%

main:
	program {
		debugParser(1, true, "main -> program\n\n")
		$$ = $1
		fmt.Println("Parsed successfully.\n")

		// pretty print struct
		if ParserVerbosity >= 2 {
			pprintAst($$)
			fmt.Print("\n")
		}
	}

program:
	program stmt {
		debugParser(1, true, "program -> program delim stmt\n\n")
		$$ = append($1, $2...)
	}
|
	stmt {
		debugParser(1, true, "program -> stmt\n\n")
		$$ = $1
	}

stmt:
	decl delim {
		debugParser(1, true, "stmt -> decl \n\n")
		$$ = $1
	}
|
	call delim {
		debugParser(1, true, "stmt -> call \n\n")
		$$ = $1
	}

decl:
	':' CMD args '{' delim program '}' {
		debugParser(1, true, "decl -> :CMD args delim { program }\n\n")
		$$ = []internal.Stmt{
			internal.Stmt(internal.Decl{
				Cmd: $2.Lexeme,
				Args: $3,
				Body: $6,
				Line: $2.Line,
			}),
		}
	}

call:
	CMD args {
		debugParser(1, true, "call -> CMD args\n\n")
		$$ = []internal.Stmt{
			internal.Stmt(internal.Call{
				Cmd: $1.Lexeme,
				Args: $2,
				Line: $1.Line,
			}),
		}
	}

args:
	/* nullable */ {
		debugParser(1, true, "args -> EPSILON\n\n")
		$$ = make([]internal.Arg, 0, 0)
	}
|
	arg ',' args {
		debugParser(1, true, "args -> arg, args\n\n")
		$$ = append($1, $3...)
	}
|
	arg {
		debugParser(1, true, "args -> arg\n\n")
		$$ = $1
	}

arg:
	REG {
		debugParser(1, true, "arg -> REG\n\n")
		$$ = []internal.Arg{
			internal.Arg(internal.Reg{
				Alias: $1.Lexeme,
				Line:  $1.Line,
			}),
		}
	}
|
	NUM {
		debugParser(1, true, "arg -> NUM\n\n")
		n, err := strconv.ParseInt($1.Lexeme, 0, 64)
		if err != nil {
			panic(err)
		}
		$$ = []internal.Arg{
			internal.Arg(internal.Num{
				Value: n,
				Line:  $1.Line,
			}),
		}
	}

delim:
	CR delim {
		debugParser(1, true, "delim -> CR delim\n\n")
	}
|
	CR {
		debugParser(1, true, "delim -> CR\n\n")
	}
