package errors

import (
	"fmt"
)

var (
	LexerVerbosityFlag  int
	ParserVerbosityFlag int
)

const (
	LexerVerbosityUsage  string = "level of lexer debugging information to print"
	ParserVerbosityUsage string = "level of parser debugging information to print"
)

func UnrecognizedInput(rn rune) string {
	return fmt.Sprintf("Unrecognized input: %s\n", Repr(rn))
}

func lexerPrefix(str string) string {
	return "[LEXER] " + str
}

func parserPrefix(str string) string {
	return "[PARSER] " + str
}

func DebugLexerStr(v int, prefix bool, format string, args ...interface{}) (ret string) {
	if LexerVerbosityFlag >= v {
		ret = fmt.Sprintf(format, args...)
		if prefix {
			ret = lexerPrefix(ret)
		}
	}
	return
}

func DebugParserStr(v int, prefix bool, format string, args ...interface{}) (ret string) {
	if ParserVerbosityFlag >= v {
		ret = fmt.Sprintf(format, args...)
		if prefix {
			ret = parserPrefix(ret)
		}
	}
	return
}

func DebugLexer(v int, prefix bool, format string, args ...interface{}) {
	fmt.Print(DebugLexerStr(v, prefix, format, args...))
}

func DebugParser(v int, prefix bool, format string, args ...interface{}) {
	fmt.Print(DebugParserStr(v, prefix, format, args...))
}
