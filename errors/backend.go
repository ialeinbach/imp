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
			out = prefixLines(out, BackendPrefix)
		}
		fmt.Print(out)
	}
}

func Undefined(name string) error {
	return New(fmt.Sprintf("undefined: %s", name))
}

func Unsupported(feature string) error {
	return New(fmt.Sprintf("unsupported feature: %s", feature))
}

func TypeExpected(expected string) error {
	return New(fmt.Sprintf("type expected: %s", expected))
}

func TypeMismatch(expected, found string) error {
	return New(fmt.Sprintf("type mismatch: expected %s but found %s", expected, found))
}

func CountMismatch(expected, found int) error {
	return New(fmt.Sprintf("count mismatch: expected %d but found %d", expected, found))
}
