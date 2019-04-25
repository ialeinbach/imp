package errors

import (
	"errors"
	"fmt"
)

func New(format string, a ...interface{}) error {
	return errors.New(fmt.Sprintf(format, a...))
}

func Line(line int, err error) error {
	return New("line %d: %s", line, err)
}

func Repr(rn rune) string {
	switch rn {
	case '\n':
		return "CR"
	case '\t':
		return "TAB"
	default:
		return string(rn)
	}
}
