%{

package imp

import "fmt"

var ParserVerbosity int

%}

%union{
    str     string
    strlist []string
}

%token <str> CMD REG NUM CR
%type <strlist> arg args delim call decl stmt program main

%start main

%%

main
	: program {
		debugParser(1, true, "main -> program\n\n")
		fmt.Printf("Parsed successfully.\n")
	}

program
	: program stmt {debugParser(1, true, "program -> program delim stmt\n\n")}
	| stmt {debugParser(1, true, "program -> stmt\n\n")}

stmt
	: decl delim {debugParser(1, true, "stmt -> decl \n\n")}
	| call delim {debugParser(1, true, "stmt -> call \n\n")}

decl
	: ':' CMD args '{' delim program '}' {debugParser(1, true, "decl -> :CMD args delim { program }\n\n")}

call
	: CMD args {debugParser(1, true, "call -> CMD args\n\n")}

args
	: /* nullable */ {debugParser(1, true, "args -> EPSILON\n\n")}
	| arg ',' args {debugParser(1, true, "args -> arg, args\n\n")}
	| arg {debugParser(1, true, "args -> arg\n\n")}

arg
	: REG {debugParser(1, true, "arg -> REG\n\n")}
	| NUM {debugParser(1, true, "arg -> NUM\n\n")}

delim
	: CR delim {debugParser(1, true, "delim -> CR delim\n\n")}
	| CR {debugParser(1, true, "delim -> CR\n\n")}
