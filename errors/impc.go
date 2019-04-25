package errors

import (
	"fmt"
	"os"
)

func Print(err error) {
	fmt.Fprintf(os.Stderr, "impc: %s\n", err)
}

func BadSourceFile(filename string, err error) error {
	return New(fmt.Sprintf("error opening %s: %v", filename, err))
}
