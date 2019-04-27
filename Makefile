.PHONY: clean test

impc: frontend/* backend/* errors/*
	go build -o impc

frontend/lexer.go: frontend/parser.go

frontend/parser.go: frontend/parser.y
	go generate -x

clean:
	$(RM) frontend/parser.go frontend/y.output impc

test: impc
	./impc examples/*.imp
