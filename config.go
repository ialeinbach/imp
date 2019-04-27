package main

import (
	"imp/backend"
	"imp/errors"
)

func configLexerVerbosity(short, long int) {
	if short >= long {
		errors.LexerVerbosity = short
	} else {
		errors.LexerVerbosity = long
	}
}

func configParserVerbosity(short, long int) {
	if short >= long {
		errors.ParserVerbosity = short
	} else {
		errors.ParserVerbosity = long
	}
}

func configBackendVerbosity(short, long int) {
	if short >= long {
		errors.BackendVerbosity = short
	} else {
		errors.BackendVerbosity = long
	}
}

func configTargetArchitecture(short, long string) {
	if long == "" {
		backend.TargetArchitecture = short
	} else {
		backend.TargetArchitecture = long
	}
}
