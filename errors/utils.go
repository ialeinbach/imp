package errors

import (
	"errors"
	"fmt"
	"strings"
)

type Textual interface {
	String() string
	Pos() int
	Typed
}

type Typed interface {
	Type() string
}

// Wrapper for New() from Go standard library that adds string formatting.
func New(format string, a ...interface{}) error {
	return errors.New(fmt.Sprintf(format, a...))
}

// Wraps an error with the context of a line.
func Line(line int, err error) error {
	return New("line %d: %s", line, err)
}

// Returns a string representation of a rune that's easier to read.
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

// Returns a string with a prefix inserted at the beginning of each line.
func PrefixLines(s, prefix string) string {
	var b strings.Builder

	for next := 0; len(s) > 0; s = s[next:] {
		b.WriteString(prefix)
		next = strings.IndexRune(s, '\n') + 1
		if next == 0 {
			b.WriteString(s)
			break
		}
		b.WriteString(s[:next])
	}

	return b.String()
}

// Returns a string with each line prefixed by a tab (i.e. indented).
func Indent(s string) string {
	return PrefixLines(s, "\t")
}

// Wraps an error with the context of a Textual object.
func Wrap(err error, t Textual) error {
	return New(fmt.Sprintf("%s at %d (%s): %s", t.Type(), t.Pos(), t, err))
}
