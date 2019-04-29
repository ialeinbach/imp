.PHONY: clean test

impc: frontend/* backend/* errors/*
	go build -o impc

twerp: impc interp/*
	go build -o interp/twerp interp/*.go

frontend/lexer.go: frontend/parser.go

frontend/parser.go: frontend/parser.y
	go generate -x

clean:
	$(RM) frontend/parser.go frontend/y.output impc interp/twerp

test: impc
	./impc examples/*.imp
