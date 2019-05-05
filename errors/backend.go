package errors

import (
	"fmt"
)

// Flag-configurables.
var (
	BackendVerbosityFlag int
)

const BackendPrefix string = "[BACKEND] "

func DebugBackend(verbosity int, prefix bool, format string, a ...interface{}) {
	if BackendVerbosityFlag >= verbosity {
		out := fmt.Sprintf(format, a...)
		if prefix {
			out = PrefixLines(out, BackendPrefix)
		}
		fmt.Print(out)
	}
}

func Undefined(t Textual) error {
	return New(fmt.Sprintf("undefined: %s", t))
}

func Unsupported(format string, a ...interface{}) error {
	return New("unsupported feature: " + format, a...)
}

func TypeMismatch(expected, found Typed) error {
	return New(fmt.Sprintf(
		"type mismatch: can't get %s from %s",
		expected.Type(), found.Type(),
	))
}

func CountMismatch(expected, found int) error {
	return New(fmt.Sprintf(
		"count mismatch: expected %d but found %d",
		expected, found,
	))
}
