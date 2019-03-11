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
%type <strlist> arg args delim term decl stmt program main

%start main

%%

main
	: program {
		debugParser(1, true, "main -> program\n\n")
		fmt.Printf("Parsed successfully.\n")
	}

program
	: program term {debugParser(1, true, "program -> program delim term\n\n")}
	| term {debugParser(1, true, "program -> term\n\n")}

term
	: decl delim {debugParser(1, true, "term -> decl \n\n")}
	| stmt delim {debugParser(1, true, "term -> stmt \n\n")}

decl
	: ':' CMD args '{' delim program '}' {debugParser(1, true, "decl -> :CMD args delim { program }\n\n")}

stmt
	: CMD args {debugParser(1, true, "stmt -> CMD args\n\n")}

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
