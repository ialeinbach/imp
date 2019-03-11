package imp

//go:generate goyacc -l -o parser.go parser.y

import "fmt"

var LexerVerbosity int

const (
	MAX_IDENT_LENGTH int  = 8
	MAX_NUM_LENGTH   int  = 21 // can hold all int64 values (incl. sign)
	NUM_PREFIX       rune = '#'
	REG_PREFIX       rune = '@'
)

type lexer struct {
	start     int
	curr      int
	input     string
	err       error
}

func Lexer(input string) *lexer {
	return &lexer{input: input}
}

// No lexerPred closures allowed.
type lexerPred func(rn rune) bool

var (
	predNumber lexerPred = func(rn rune) bool {
		return rn >= '0' && rn <= '9'
	}
	predAlpha lexerPred = func(rn rune) bool {
		return rn >= 'a' && rn <= 'z'
	}
	predIdent lexerPred = func(rn rune) bool {
		return predAlpha(rn) || rn == '_' || rn == '?'
	}
	predRegister lexerPred = func(rn rune) bool {
		return predAlpha(rn) || predNumber(rn)
	}
	predRegPrefix lexerPred = func(rn rune) bool {
		return rn == REG_PREFIX
	}
	predNumPrefix lexerPred = func(rn rune) bool {
		return rn == NUM_PREFIX
	}
)

func (l *lexer) emit(tokName string, lval *yySymType) {
	lexeme := l.lexeme()
	debugLexer(1, true, "Emitting %s(%s)\n", tokName, lexeme)
	if LexerVerbosity == 1 {
		fmt.Printf("\n")
	}
	lval.str = lexeme
	l.start = l.curr
	debugLexer(2, true, "====================\n\n")
}

func (l *lexer) lexPred(pred lexerPred, max int) (n int) {
	debugLexer(2, true, "--------------------\n")
	debugLexer(2, true, "Reading '%s'\t", repr(rune(l.input[l.curr])))
	for len(l.input[l.curr:]) > 0 && n < max && pred(rune(l.input[l.curr])) {
		debugLexer(2, false, "  OK\n")
		l.curr++
		debugLexer(2, true, "Reading '%s'\t", repr(rune(l.input[l.curr])))
		n++
	}
	debugLexer(2, false, "FAIL\n")
	debugLexer(2, true, "--------------------\n")
	return
}

func (l *lexer) lexPrefixed(prefix lexerPred, suffix lexerPred, max int) (n int) {
	// ignore lone prefix
	defer func() {
		if n > 1 {
			l.start++
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
	debugLexer(2, true, "Lexing CMD\n")
	return l.lexPred(predIdent, MAX_IDENT_LENGTH)
}

func (l *lexer) lexReg() int {
	debugLexer(2, true, "Lexing REG\n")
	return l.lexPrefixed(predRegPrefix, predRegister, MAX_IDENT_LENGTH)
}

func (l *lexer) lexNum() int {
	debugLexer(2, true, "Lexing NUM\n")
	return l.lexPrefixed(predNumPrefix, predNumber, MAX_NUM_LENGTH)
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
			debugLexer(1, true, "Emitting CR\n\n")
			l.curr++
			l.start++
			return CR
		case ':', ',', '{', '}':
			debugLexer(1, true, "Emitting SYM(%c)\n\n", ch)
			l.curr++
			l.start++
			return int(ch)
		}
		debugLexer(2, true, "====================\n")
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
		l.err = ErrorPrefix("Unrecognized input encountered:", repr(rune(ch)), "\n")
		return int(ch)
	}
	return int(ch)
}

// Satisfies yyLexer.
func (l *lexer) Error(str string) {
	l.err = ErrorPrefix(str)
}

func Parse(input string) (ret int) {
	l := Lexer(input)
	ret = yyParse(l)
	if l.err != nil {
		debugLexer(0, true, "%s\n", l.err)
	}
	return
}
