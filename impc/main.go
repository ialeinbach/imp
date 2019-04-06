package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"imp/frontend"
)

const (
	lexerVerbosityUsage  string = "level of lexer debugging information to print"
	parserVerbosityUsage string = "level of parser debugging information to print"
)

func init() {
	flag.IntVar(&frontend.LexerVerbosity, "lexer-verbosity", 0, lexerVerbosityUsage)
	flag.IntVar(&frontend.LexerVerbosity, "lv", 0, lexerVerbosityUsage)
	flag.IntVar(&frontend.ParserVerbosity, "parser-verbosity", 0, parserVerbosityUsage)
	flag.IntVar(&frontend.ParserVerbosity, "pv", 0, parserVerbosityUsage)
	flag.Parse()
}

func main() {
	var (
		src []byte
		err error
	)
	if flag.NArg() == 0 {
		fmt.Println("Nothing to do.\n")
		flag.PrintDefaults()
		return
	}
	for _, f := range flag.Args() {
		src, err = ioutil.ReadFile(f)
		if err == nil {
			frontend.Parse(string(src))
		} else {
			fmt.Fprintf(os.Stderr, frontend.ErrorPrefix(fmt.Sprintf("error opening %s: %s\n", f, err)).Error())
		}
	}
	fmt.Println("Done.")
}
