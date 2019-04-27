package frontend

import (
	"imp/errors"
)

const (
	maxCmdLength int = 32
	maxRegLength int = 32
	maxNumLength int = 32

	numPrefix rune = '#'
	regPrefix rune = '@'
)

//
// Lexer
//

// Implements goyacc's yyLexer interface.
type lexer struct {
	// Source code being lexed.
	input string

	// Indices into input used to build a lexeme for tokenization.
	start int
	curr  int

	// Lexer state for debugging.
	line int
	err  error
}

type Token struct {
	Lexeme string
	Line   int
}

func Lexer(input string) *lexer {
	return &lexer{
		input: input,
		line:  1,
	}
}

func (l *lexer) head() rune {
	return rune(l.input[l.curr])
}

func (l *lexer) lexeme() string {
	return string(l.input[l.start:l.curr])
}

//
// Lexer Predicate Functions
//

type lexerPred func(rn rune) bool

// Primitives that are composed to build functions that define the textual
// representation of imp data.
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
)

var (
	// Cmds are alphanumeric with underscores.
	predCmdLiteral lexerPred = func(rn rune) bool {
		return predAlpha(rn) || predDec(rn) || rn == '_'
	}

	// Reg aliases are alphanumeric and Reg literals are numeric.
	predRegPrefix lexerPred = func(rn rune) bool {
		return rn == regPrefix
	}
	predRegBody lexerPred = func(rn rune) bool {
		return predAlpha(rn) || predHex(rn)
	}

	// Num aliases are alphanumeric and Num literals are numeric.
	predNumPrefix lexerPred = func(rn rune) bool {
		return rn == numPrefix
	}
	predNumBody lexerPred = func(rn rune) bool {
		return predAlpha(rn) || predHex(rn)
	}
)

//
// Lexer Driver Methods
//

// Emits a token to the parser. Called for non-trivial tokens i.e. ones that
// have a corresponding item on the goyacc value stack.
func (l *lexer) emit(name string, lval *yySymType) {
	errors.DebugLexer(1, true, "Emitting %s(%s)\n", name, l.lexeme())
	errors.DebugLexer(2, true, "====================\n")
	errors.DebugLexer(2, false, "\n")

	// Tokenize lexeme and emit.
	lval.tok = Token{
		Lexeme: l.lexeme(),
		Line:   l.line,
	}

	// Move past consumed lexeme.
	l.start = l.curr

	return
}

// Advances lexer according to pred for at most max characters. Returns the number of
// characters by which the lexer was advanced.
func (l *lexer) lexPred(pred lexerPred, max int) (n int) {
	errors.DebugLexer(2, true, "--------------------\n")
	errors.DebugLexer(2, true, "Reading '%s'\t", errors.Repr(l.head()))

	for len(l.input[l.curr:]) > 0 && n < max && pred(l.head()) {
		l.curr++
		n++

		errors.DebugLexer(2, false, "  OK\n")
		errors.DebugLexer(2, true, "Reading '%s'\t", errors.Repr(l.head()))
	}

	errors.DebugLexer(2, false, "FAIL\n")
	errors.DebugLexer(2, true, "--------------------\n")

	return
}

// Advances lexer with prefix predicate for first character, then uses body
// predicate for at most max characters. The prefix symbol does not count
// towards the character count as it pertains to max.
func (l *lexer) lexPrefixed(prefix lexerPred, body lexerPred, max int) (n int) {
	defer func() {
		switch {
		// Ignore prefix if successful.
		case n > 1:
			l.start++
			n -= 1
		// Rewind lexer if we only saw prefix.
		case n == 1:
			l.curr--
			n = 0
		}
	}()

	// Must be self-referentiable.
	var pred lexerPred

	// Lex prefix then replace self in order to lex body.
	pred = func(rn rune) (ok bool) {
		defer func() { pred = body }()
		return prefix(rn)
	}

	// Wrap in a closure so caller is none the wiser.
	return l.lexPred(func(rn rune) bool {
		return pred(rn)
	}, max)
}

func (l *lexer) lexCmd() int {
	errors.DebugLexer(2, true, "Lexing CMD\n")
	return l.lexPred(predCmdLiteral, maxCmdLength)
}

func (l *lexer) lexReg() int {
	errors.DebugLexer(2, true, "Lexing REG\n")
	return l.lexPrefixed(predRegPrefix, predRegBody, maxRegLength)
}

func (l *lexer) lexNum() int {
	errors.DebugLexer(2, true, "Lexing NUM\n")
	return l.lexPrefixed(predNumPrefix, predNumBody, maxNumLength)
}

//
// Handles For Goyacc
//

// Satisfies yyLexer.
func (l *lexer) Lex(lval *yySymType) int {
	var ch int

	for len(l.input[l.curr:]) > 0 {
		switch ch = int(l.head()); ch {

		// Ignore whitespace within a given line.
		case ' ', '\t':
			l.curr++
			l.start++
			continue

		// Newlines are slightly more important whitespace (for collecting line
		// info and imp syntax).
		case '\n':
			l.curr++
			l.start++
			l.line++

			errors.DebugLexer(1, true, "Emitting CR\n")
			return CR

		// Lex trivial tokens. Used as-is in parser i.e. no corresponding item
		// on goyacc value stack.
		case ':', ',', '{', '}':
			l.curr++
			l.start++

			errors.DebugLexer(1, true, "Emitting SYM(%c)\n", ch)
			return ch
		}

		errors.DebugLexer(2, false, "\n")
		errors.DebugLexer(2, true, "====================\n")

		// Lex non-trivial tokens (i.e. can be directly mapped to nodes on the
		// AST). Pass control of lexer to its methods.
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
