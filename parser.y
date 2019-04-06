%{

package imp

import (
	"fmt"
	"strconv"
)

const (
	MAX_ARG_COUNT int = 6
)

var ParserVerbosity int

%}

%union{
	str       string
	arglist   []Arg
	paramlist []Param
	stmtlist  []Stmt
}

%token <str> CMD REG REG_ALIAS NUM NUM_ALIAS CR

%type <arglist> arg args
%type <paramlist> param params
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
	':' CMD params '{' delim program '}' {
		debugParser(1, true, "decl -> :CMD params delim { program }\n\n")
		$$ = []Stmt{ Stmt(Decl{Cmd: $2, Params: $3, Body: $6}) }
	}

call:
	CMD args {
		debugParser(1, true, "call -> CMD args\n\n")
		$$ = []Stmt{ Stmt(Call{Cmd: $1, Args: $2}) }
	}

args:
	/* nullable */ {
		debugParser(1, true, "args -> EPSILON\n\n")
		$$ = make([]Arg, 0, 0)
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
		r, err := strconv.ParseInt($1, 0, 64)
		if err != nil {
			panic(err)
		}
		$$ = []Arg{ Arg(Reg{Val: int(r)}) }
	}
|
	NUM {
		debugParser(1, true, "arg -> NUM\n\n")
		n, err := strconv.ParseInt($1, 0, 64)
		if err != nil {
			panic(err)
		}
		$$ = []Arg{ Arg(Num{Val: n}) }
	}
|
	REG_ALIAS {
		debugParser(1, true, "arg -> REG_ALIAS\n\n")
		$$ = []Arg{ Arg(RegAlias{Name: $1}) }
	}
|
	NUM_ALIAS {
		debugParser(1, true, "arg -> NUM_ALIAS\n\n")
		$$ = []Arg{ Arg(NumAlias{Name: $1}) }
	}

params:
	/* nullable */ {
		debugParser(1, true, "params -> EPSILON\n\n")
		$$ = make([]Param, 0, 0)
	}
|
	param ',' params {
		debugParser(1, true, "params -> param, params\n\n")
		$$ = append($1, $3...)
	}
|
	param {
		debugParser(1, true, "params -> param\n\n")
		$$ = $1
	}



param:
	REG_ALIAS {
		debugParser(1, true, "arg -> REG_ALIAS\n\n")
		$$ = []Param{ Param(RegAlias{Name: $1}) }
	}
|
	NUM_ALIAS {
		debugParser(1, true, "arg -> NUM_ALIAS\n\n")
		$$ = []Param{ Param(NumAlias{Name: $1}) }
	}

delim
	: CR delim { debugParser(1, true, "delim -> CR delim\n\n") }
	| CR { debugParser(1, true, "delim -> CR\n\n") }
