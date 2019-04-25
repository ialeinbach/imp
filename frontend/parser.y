%{

package frontend

import (
	"imp/errors"
	"imp/internal"
)

func Parse(input string) ([]internal.Stmt, error) {
	l := Lexer(input)
	yyParse(l)
	if l.err != nil {
		return nil, l.err
	}
	return l.ret, nil
}

%}

%union{
	tok      Token
	arglist  []internal.Alias
	stmtlist []internal.Stmt
}

%token <tok> CMD REG NUM CR

%type <arglist> arg args
%type <stmtlist> decl call stmt program main

%start main

%%

main:
	program {
		errors.DebugParser(1, true, "main -> program\n\n")
		errors.DebugParser(2, false, internal.DumpAst($1) + "\n\n")

		// Work around yyParse return value.
		yylex.(*lexer).ret = $1
	}

program:
	program stmt {
		errors.DebugParser(1, true, "program -> program delim stmt\n\n")
		$$ = append($1, $2...)
	}
|
	stmt {
		errors.DebugParser(1, true, "program -> stmt\n\n")
		$$ = $1
	}

stmt:
	decl delim {
		errors.DebugParser(1, true, "stmt -> decl \n\n")
		$$ = $1
	}
|
	call delim {
		errors.DebugParser(1, true, "stmt -> call \n\n")
		$$ = $1
	}

decl:
	':' CMD args '{' delim program '}' {
		errors.DebugParser(1, true, "decl -> :CMD args delim { program }\n\n")
		$$ = []internal.Stmt{internal.Stmt(
			internal.Decl{
				Cmd: internal.CmdAlias{
					Name: $2.Lexeme,
					Line: $2.Line,
				},
				Args: $3,
				Body: $6,
			},
		)}
	}

call:
	CMD args {
		errors.DebugParser(1, true, "call -> CMD args\n\n")
		$$ = []internal.Stmt{internal.Stmt(
			internal.Call{
				Cmd: internal.CmdAlias{
					Name: $1.Lexeme,
					Line: $1.Line,
				},
				Args: $2,
			},
		)}
	}

args:
	/* nullable */ {
		errors.DebugParser(1, true, "args -> EPSILON\n\n")
		$$ = make([]internal.Alias, 0, 0)
	}
|
	arg ',' args {
		errors.DebugParser(1, true, "args -> arg, args\n\n")
		$$ = append($1, $3...)
	}
|
	arg {
		errors.DebugParser(1, true, "args -> arg\n\n")
		$$ = $1
	}

arg:
	REG {
		errors.DebugParser(1, true, "arg -> REG\n\n")
		$$ = []internal.Alias{internal.Alias(
			internal.RegAlias{
				Name: $1.Lexeme,
				Line: $1.Line,
			},
		)}
	}
|
	NUM {
		errors.DebugParser(1, true, "arg -> NUM\n\n")
		$$ = []internal.Alias{internal.Alias(
			internal.NumAlias{
				Name: $1.Lexeme,
				Line: $1.Line,
			},
		)}
	}

delim:
	CR delim {
		errors.DebugParser(1, true, "delim -> CR delim\n\n")
	}
|
	CR {
		errors.DebugParser(1, true, "delim -> CR\n\n")
	}
