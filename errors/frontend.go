package errors

import (
	"fmt"
)

// Flag-configurables.
var (
	LexerVerbosityFlag  int
	ParserVerbosityFlag int
)

const (
	LexerPrefix  string = "[LEXER  ] "
	ParserPrefix string = "[PARSER ] "
)

func DebugLexer(verbosity int, prefix bool, format string, a ...interface{}) {
	if LexerVerbosityFlag >= verbosity {
		out := fmt.Sprintf(format, a...)
		if prefix {
			out = PrefixLines(out, LexerPrefix)
		}
		fmt.Print(out)
	}
}

func DebugParser(verbosity int, prefix bool, format string, a ...interface{}) {
	if ParserVerbosityFlag >= verbosity {
		out := fmt.Sprintf(format, a...)
		if prefix {
			out = PrefixLines(out, ParserPrefix)
		}
		fmt.Print(out)
	}
}

func UnrecognizedInput(rn rune) error {
	return New(fmt.Sprintf("Unrecognized input: %s\n", Repr(rn)))
}
