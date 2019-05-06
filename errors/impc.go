package errors

import (
	"fmt"
	"os"
)

// Prints and imp error message.
func Print(err error) {
	fmt.Fprintf(os.Stderr, "imp: %s\n", err)
}

// Prints a message to report successful compilation.
func Ok(filename string) {
	fmt.Printf("Source file \"%s\" compiled with no errors.\n", filename)
}

func BadSourceFile(filename string, err error) error {
	return New(fmt.Sprintf("error opening %s: %v", filename, err))
}
