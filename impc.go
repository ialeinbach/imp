package main

//go:generate goyacc -l -o frontend/parser.go -v frontend/y.output frontend/parser.y

import (
	"flag"
	"io/ioutil"
	"os"

	"imp/frontend"
	"imp/backend"
	"imp/errors"
)

const (
	// Lexer Verbosity: -lexer-verbosity, -lv
	lexerVerbosityUsage     string = "level of lexer debugging information to print"

	// Parser Verbosity: -parser-verbosity, -pv
	parserVerbosityUsage    string = "level of parser debugging information to print"

	// Backend Verbosity: -backend-verbosity, -bv
	backendVerbosityUsage   string = "level of backend debugging information to print"

	// Target Architecture: -target-architecture, -arch
	targetArchitectureUsage string = "target architecture for code generation"
)

func init() {
	var lexerVerbosityLong, lexerVerbosityShort int
	flag.IntVar(&lexerVerbosityLong, "lexer-verbosity", 0, lexerVerbosityUsage)
	flag.IntVar(&lexerVerbosityShort, "lv", 0, lexerVerbosityUsage)

	var parserVerbosityLong, parserVerbosityShort int
	flag.IntVar(&parserVerbosityLong, "parser-verbosity", 0, parserVerbosityUsage)
	flag.IntVar(&parserVerbosityShort, "pv", 0, parserVerbosityUsage)

	var backendVerbosityLong, backendVerbosityShort int
	flag.IntVar(&backendVerbosityLong, "backend-verbosity", 0, backendVerbosityUsage)
	flag.IntVar(&backendVerbosityShort, "bv", 0, backendVerbosityUsage)

	var targetArchitectureLong, targetArchitectureShort string
	flag.StringVar(&targetArchitectureLong, "target-architechture", "", targetArchitectureUsage)
	flag.StringVar(&targetArchitectureShort, "arch", "", targetArchitectureUsage)


	flag.Parse()

	configLexerVerbosity(lexerVerbosityLong, lexerVerbosityShort)
	configParserVerbosity(parserVerbosityLong, parserVerbosityShort)
	configTargetArchitecture(targetArchitectureLong, targetArchitectureShort)
	configBackendVerbosity(backendVerbosityLong, backendVerbosityShort)
}

func main() {
	if flag.NArg() == 0 {
		errors.Print(errors.NoSourceFiles())
		flag.PrintDefaults()
		os.Exit(0)
	}
	for _, filename := range flag.Args() {
		src, err := ioutil.ReadFile(filename)
		if err != nil {
			errors.Print(errors.BadSourceFile(filename, err))
			os.Exit(1)
		}

		ast, err := frontend.Parse(string(src))
		if err != nil {
			errors.Print(err)
			os.Exit(1)
		}

		_, err = backend.Flatten(ast)
		if err != nil {
			errors.Print(err)
			os.Exit(1)
		}

		// Currently, by default, the compiler outputs nothing.
		// However, with the various verbosity flags, each stage can be printed
		// out.
		errors.Ok(filename)
	}
}
