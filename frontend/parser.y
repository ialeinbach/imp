%{

package frontend

import (
	"imp/errors"
	"imp/backend"
)

// yyParse can only return an int. Final grammar rule places AST here so that
// Parse() can return it.
var abstractSyntaxTree []backend.Stmt

func Parse(input string) ([]backend.Stmt, error) {
	l := Lexer(input)
	yyParse(l)
	if l.err != nil {
		return nil, l.err
	}
	return abstractSyntaxTree, nil
}

%}

%union{
	tok      Token
	arglist  []backend.Alias
	stmtlist []backend.Stmt
}

%token <tok> CMD REG NUM CR

%type <arglist> arg args
%type <stmtlist> decl call stmt program main

%start main

%%

main:
	program {
		// yyParse can only return an int. Place AST in package global so that
		// Parse() can return it.
		abstractSyntaxTree = $1

		errors.DebugParser(1, true, "main -> program\n")
		errors.DebugParser(2, false, "\n")
		errors.DebugParser(2, true, backend.DumpAst($1))
		errors.DebugParser(2, false, "\n\n")
	}

program:
	program stmt {
		$$ = append($1, $2...)
		errors.DebugParser(1, true, "program -> program delim stmt\n")
	}
|
	stmt {
		errors.DebugParser(1, true, "program -> stmt\n")
		$$ = $1
	}

stmt:
	decl delim {
		$$ = $1
		errors.DebugParser(1, true, "stmt -> decl \n")
	}
|
	call delim {
		$$ = $1
		errors.DebugParser(1, true, "stmt -> call \n")
	}

decl:
	':' CMD args '{' delim program '}' {
		$$ = []backend.Stmt{backend.Stmt(
			backend.Decl{
				Cmd: backend.CmdAlias{
					Name: $2.Lexeme,
					Line: $2.Line,
				},
				Args: $3,
				Body: $6,
			},
		)}
		errors.DebugParser(1, true, "decl -> :CMD args delim { program }\n")
	}

call:
	CMD args {
		$$ = []backend.Stmt{backend.Stmt(
			backend.Call{
				Cmd: backend.CmdAlias{
					Name: $1.Lexeme,
					Line: $1.Line,
				},
				Args: $2,
			},
		)}
		errors.DebugParser(1, true, "call -> CMD args\n")
	}

args:
	/* nullable */ {
		$$ = make([]backend.Alias, 0, 0)
		errors.DebugParser(1, true, "args -> EPSILON\n")
	}
|
	arg ',' args {
		$$ = append($1, $3...)
		errors.DebugParser(1, true, "args -> arg, args\n")
	}
|
	arg {
		$$ = $1
		errors.DebugParser(1, true, "args -> arg\n")
	}

arg:
	REG {
		$$ = []backend.Alias{backend.Alias(
			backend.RegAlias{
				Name: $1.Lexeme,
				Line: $1.Line,
			},
		)}
		errors.DebugParser(1, true, "arg -> REG\n")
	}
|
	NUM {
		$$ = []backend.Alias{backend.Alias(
			backend.NumAlias{
				Name: $1.Lexeme,
				Line: $1.Line,
			},
		)}
		errors.DebugParser(1, true, "arg -> NUM\n")
	}

delim:
	CR delim {
		errors.DebugParser(1, true, "delim -> CR delim\n")
	}
|
	CR {
		errors.DebugParser(1, true, "delim -> CR\n")
	}
