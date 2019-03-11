package imp

import (
	"errors"
	"fmt"
)

func repr(rn rune) string {
	switch rn {
	case '\n':
		return "CR"
	case '\t':
		return "TAB"
	default:
		return string(rn)
	}
}

func ErrorPrefix(strs ...string) error {
	s := make([]interface{}, len(strs)+1)
	for i, str := range append([]string{"imp:"}, strs...) {
		s[i] = str
	}
	return errors.New(fmt.Sprint(s...))
}

func lexerPrefix(str string) string {
	return "[LEXER] " + str
}

func parserPrefix(str string) string {
	return "[PARSER] " + str
}

func debugLexerStr(v int, prefix bool, format string, args ...interface{}) (ret string) {
	if LexerVerbosity >= v {
		ret = fmt.Sprintf(format, args...)
		if prefix {
			ret = lexerPrefix(ret)
		}
	}
	return
}

func debugParserStr(v int, prefix bool, format string, args ...interface{}) (ret string) {
	if ParserVerbosity >= v {
		ret = fmt.Sprintf(format, args...)
		if prefix {
			ret = parserPrefix(ret)
		}
	}
	return
}

func debugLexer(v int, prefix bool, format string, args ...interface{}) {
	fmt.Print(debugLexerStr(v, prefix, format, args...))
}

func debugParser(v int, prefix bool, format string, args ...interface{}) {
	fmt.Print(debugParserStr(v, prefix, format, args...))
}
