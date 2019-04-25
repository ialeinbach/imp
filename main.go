package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"imp/frontend"
	"imp/internal"
	"imp/backend"
	"imp/errors"
	"os"
)

func init() {
	flag.IntVar(&errors.LexerVerbosityFlag, "lexer-verbosity", 0, errors.LexerVerbosityUsage)
	flag.IntVar(&errors.LexerVerbosityFlag, "lv", 0, errors.LexerVerbosityUsage)
	flag.IntVar(&errors.ParserVerbosityFlag, "parser-verbosity", 0, errors.ParserVerbosityUsage)
	flag.IntVar(&errors.ParserVerbosityFlag, "pv", 0, errors.ParserVerbosityUsage)
	flag.StringVar(&backend.TargetArchitectureFlag, "arch", "x86", backend.TargetArchitectureUsage)
	flag.Parse()
}

func main() {
	if flag.NArg() == 0 {
		fmt.Println("Nothing to do.\n")
		flag.PrintDefaults()
		return
	}
	for _, filename := range flag.Args() {
		src, err := ioutil.ReadFile(filename)
		if err == nil {
			ast, err := frontend.Parse(string(src))
			if err != nil {
				errors.Print(err)
				os.Exit(1)
			}
			psuedo, err := internal.Flatten(ast)
			if err != nil {
				errors.Print(err)
				os.Exit(1)
			}
			fmt.Printf("%s:\n", filename)
			fmt.Println(backend.DumpPsuedo(psuedo))
		} else {
			errors.Print(errors.BadSourceFile(filename, err))
			os.Exit(1)
		}
	}
}
