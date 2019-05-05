%{

package frontend

import (
	"imp/errors"
)

// yyParse can only return an int. Final grammar rule places AST here so that
// Parse() can return it.
var abstractSyntaxTree []Stmt

func Parse(input string) ([]Stmt, error) {
	l := Lexer(input)
	yyParse(l)
	if l.err != nil {
		return nil, errors.Line(l.line, l.err)
	}
	return abstractSyntaxTree, nil
}

%}

%union{
	tok      Token
	arglist  []Alias
	stmtlist []Stmt
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
		errors.DebugParser(1, false, "\n")
		errors.DebugParser(2, true, DumpAst($1))
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
		cmd := CmdAlias{$2.Lexeme, $2.Line}
		decl := Decl{cmd, $3, $6}
		$$ = []Stmt{decl}
		errors.DebugParser(1, true, "decl -> :CMD args delim { program }\n")
	}

call:
	CMD args {
		cmd := CmdAlias{$1.Lexeme, $1.Line}
		call := Call{cmd, $2}
		$$ = []Stmt{call}
		errors.DebugParser(1, true, "call -> CMD args\n")
	}

args:
	/* nullable */ {
		$$ = make([]Alias, 0, 0)
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
		reg := RegAlias{$1.Lexeme, $1.Line}
		$$ = []Alias{reg}
		errors.DebugParser(1, true, "arg -> REG\n")
	}
|
	NUM {
		num := NumAlias{$1.Lexeme, $1.Line}
		$$ = []Alias{num}
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
