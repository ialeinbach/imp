package errors

import (
	"fmt"
)

// Flag-configurables.
var (
	LexerVerbosity  int
	ParserVerbosity int
)

const (
	LexerPrefix  string = "[LEXER  ] "
	ParserPrefix string = "[PARSER ] "
)

func DebugLexer(verbosity int, prefix bool, format string, a ...interface{}) {
	if LexerVerbosity >= verbosity {
		out := fmt.Sprintf(format, a...)
		if prefix {
			out = prefixLines(out, LexerPrefix)
		}
		fmt.Print(out)
	}
}

func DebugParser(verbosity int, prefix bool, format string, a ...interface{}) {
	if ParserVerbosity >= verbosity {
		if prefix {
			format = prefixLines(format, ParserPrefix)
		}
		fmt.Printf(format, a...)
	}
}

func UnrecognizedInput(rn rune) string {
	return fmt.Sprintf("Unrecognized input: %s\n", Repr(rn))
}
