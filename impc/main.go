package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"imp"
)

const (
	lexerVerbosityUsage  string = "level of lexer debugging information to print"
	parserVerbosityUsage string = "level of parser debugging information to print"
)

func init() {
	flag.IntVar(&imp.LexerVerbosity, "lexer-verbosity", 0, lexerVerbosityUsage)
	flag.IntVar(&imp.LexerVerbosity, "lv", 0, lexerVerbosityUsage)
	flag.IntVar(&imp.ParserVerbosity, "parser-verbosity", 0, parserVerbosityUsage)
	flag.IntVar(&imp.ParserVerbosity, "pv", 0, parserVerbosityUsage)
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
			imp.Parse(string(src))
		} else {
			fmt.Fprintf(os.Stderr, imp.ErrorPrefix(fmt.Sprintf("error opening %s: %s\n", f, err)).Error())
		}
	}
	fmt.Println("Done.")
}
