.PHONY: clean test all

all: imp twerp

imp: frontend/* backend/* errors/*
	go build -o imp

twerp: imp interp/*.go
	go build -o twerp interp/*.go

frontend/lexer.go: frontend/parser.go

frontend/parser.go: frontend/parser.y
	go generate -x

test: imp
	./imp examples/*.imp

clean:
	$(RM) frontend/parser.go frontend/y.output imp twerp
