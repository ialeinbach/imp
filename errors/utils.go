package errors

import (
	"errors"
	"fmt"
	"strings"
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

