package main

import (
	"github.com/ialeinbach/imp/backend"
	"github.com/ialeinbach/imp/errors"
)

func configLexerVerbosity(short, long int) {
	if short >= long {
		errors.LexerVerbosityFlag = short
	} else {
		errors.LexerVerbosityFlag = long
	}
}

func configParserVerbosity(short, long int) {
	if short >= long {
		errors.ParserVerbosityFlag = short
	} else {
		errors.ParserVerbosityFlag = long
	}
}

func configBackendVerbosity(short, long int) {
	if short >= long {
		errors.BackendVerbosityFlag = short
	} else {
		errors.BackendVerbosityFlag = long
	}
}

func configTargetArchitecture(short, long string) {
	if long == "" {
		backend.TargetArchitectureFlag = short
	} else {
		backend.TargetArchitectureFlag = long
	}
}

func configHelp(short, long bool) {
	HelpFlag = short || long
}
