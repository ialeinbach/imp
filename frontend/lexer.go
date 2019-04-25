package frontend

//go:generate goyacc -l -o parser.go parser.y

import (
	"imp/errors"
	"imp/internal"
)

const (
	MaxCmdLength int = 16
	MaxRegLength int = 8
	MaxNumLength int = 21 // can hold all int64 values (incl. sign)

	NumPrefix rune = '#'
	RegPrefix rune = '@'
)

type lexer struct {
	start int
	curr  int
	input string
	line  int
	err   error
	ret   []internal.Stmt
}

type Token struct {
	Lexeme string
	Line   int
}

func Lexer(input string) *lexer {
	return &lexer{ input: input, line: 1 }
}

// No lexerPred closures allowed.
type lexerPred func(rn rune) bool

var (
	predDec lexerPred = func(rn rune) bool {
		return rn >= '0' && rn <= '9'
	}
	predHex lexerPred = func(rn rune) bool {
		return (
			(rn >= 'A' && rn <= 'F') ||
			(rn >= 'a' && rn <= 'f') ||
			predDec(rn)              )
	}
	predAlpha lexerPred = func(rn rune) bool {
		return (
			(rn >= 'a' && rn <= 'z') ||
			(rn >= 'A' && rn <= 'Z') )
	}
	predCmdLiteral lexerPred = func(rn rune) bool {
		return predAlpha(rn) || predDec(rn) || rn == '_' || rn == '?'
	}
	predRegLiteral lexerPred = func(rn rune) bool {
		return predAlpha(rn) || predHex(rn)
	}
	predNumLiteral lexerPred = func(rn rune) bool {
		return predAlpha(rn) || predHex(rn)
	}
	predRegPrefix lexerPred = func(rn rune) bool {
		return rn == RegPrefix
	}
	predNumPrefix lexerPred = func(rn rune) bool {
		return rn == NumPrefix
	}
)

func (l *lexer) emit(tokName string, lval *yySymType) {
	lexeme := l.lexeme()

	errors.DebugLexer(1, true, "Emitting %s(%s)\n", tokName, lexeme)

//	if LexerVerbosityFlag == 1 {
//		fmt.Printf("\n")
//	}
	lval.tok = Token{
		Lexeme: lexeme,
		Line:   l.line,
	}
	l.start = l.curr

	errors.DebugLexer(2, true, "====================\n\n")
}

func (l *lexer) lexPred(pred lexerPred, max int) (n int) {
	errors.DebugLexer(2, true, "--------------------\n")
	errors.DebugLexer(2, true, "Reading '%s'\t", errors.Repr(rune(l.input[l.curr])))

	for len(l.input[l.curr:]) > 0 && n < max && pred(rune(l.input[l.curr])) {
		errors.DebugLexer(2, false, "  OK\n")
		l.curr++
		errors.DebugLexer(2, true, "Reading '%s'\t", errors.Repr(rune(l.input[l.curr])))
		n++
	}

	errors.DebugLexer(2, false, "FAIL\n")
	errors.DebugLexer(2, true, "--------------------\n")

	return
}

func (l *lexer) lexPrefixed(prefix lexerPred, suffix lexerPred, max int) (n int) {
	defer func() {
		// ignore prefix
		switch {
		case n > 1:
			l.start++
			n -= 1
		case n == 1:
			l.curr--
			n = 0
		}
	}()

	// declaration must be referencable from definition
	var pred lexerPred

	// lex prefix then pass baton to suffix
	pred = func(rn rune) (ok bool) {
		defer func() { pred = suffix }()
		return prefix(rn)
	}

	return l.lexPred(func(rn rune) bool {
		return pred(rn)
	}, max)
}

func (l *lexer) lexCmd() int {
	errors.DebugLexer(2, true, "Lexing CMD\n")
	return l.lexPred(predCmdLiteral, MaxCmdLength)
}

func (l *lexer) lexReg() int {
	errors.DebugLexer(2, true, "Lexing REG\n")
	return l.lexPrefixed(predRegPrefix, predRegLiteral, MaxRegLength)
}

func (l *lexer) lexNum() int {
	errors.DebugLexer(2, true, "Lexing NUM\n")
	return l.lexPrefixed(predNumPrefix, predNumLiteral, MaxNumLength)
}

func (l *lexer) lexeme() string {
	return string(l.input[l.start:l.curr])
}

// Satisfies yyLexer.
func (l *lexer) Lex(lval *yySymType) int {
	var ch int
	for len(l.input[l.curr:]) > 0 {
		ch = int(l.input[l.curr])
		switch ch {
		case ' ', '\t':
			l.curr++
			l.start++
			continue
		case '\n':
			errors.DebugLexer(1, true, "Emitting CR\n\n")
			l.curr++
			l.start++
			l.line++
			return CR
		case ':', ',', '{', '}':
			errors.DebugLexer(1, true, "Emitting SYM(%c)\n\n", ch)
			l.curr++
			l.start++
			return ch
		}

		errors.DebugLexer(2, true, "====================\n")

		switch {
		case l.lexCmd() > 0:
			l.emit("CMD", lval)
			return CMD
		case l.lexReg() > 0:
			l.emit("REG", lval)
			return REG
		case l.lexNum() > 0:
			l.emit("NUM", lval)
			return NUM
		}

		l.Error(errors.UnrecognizedInput(rune(ch)))
		return ch
	}
	return ch
}

// Satisfies yyLexer.
func (l *lexer) Error(str string) {
	l.err = errors.New(str)
}
