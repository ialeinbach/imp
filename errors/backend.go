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

func Undefined(v interface{String() string}) error {
	return New(fmt.Sprintf("undefined: %s", v))
}

func Unsupported(feature string) error {
	return New(fmt.Sprintf("unsupported feature: %s", feature))
}

func TypeMismatch(expected, found interface{Type() string}) error {
	return New(fmt.Sprintf(
		"type mismatch: expected %s but found %s",
		expected.Type(), found.Type(),
	))
}

func CountMismatch(expected, found int) error {
	return New(fmt.Sprintf(
		"count mismatch: expected %d but found %d",
		expected, found,
	))
}
