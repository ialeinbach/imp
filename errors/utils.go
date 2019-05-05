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

func Indent(s string) string {
	return PrefixLines(s, "\t")
}

func Wrap(err error, t Textual) error {
	return New(fmt.Sprintf("%s at %d: %s: %s", t.Type(), t.Pos(), t, err))
}
