package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ialeinbach/imp/errors"
	"github.com/ialeinbach/imp/frontend"
	"github.com/ialeinbach/imp/backend"
)

var interactiveMode bool

const interactiveModeUsage string = "interpreter blocks on each pseudo-instruction with options for querying internal state"

func init() {
	flag.BoolVar(&interactiveMode, "i", false, interactiveModeUsage)
	flag.Parse()
}

func main() {
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

		psuedo, err := backend.Flatten(ast)
		if err != nil {
			errors.Print(err)
			os.Exit(1)
		}

		imptwerpreter := NewTwerp(psuedo, backend.MaxRegCount)
		ret, err := imptwerpreter.Exec(interactiveMode)
		if err != nil {
			errors.Print(err)
			os.Exit(1)
		}

		fmt.Printf("Imptwerpreter returned successfully with %v.\n", ret)
	}
}
