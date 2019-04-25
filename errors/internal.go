package errors

import (
	"fmt"
)

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
