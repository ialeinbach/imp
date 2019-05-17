.PHONY: clean test all

all: imp twerp

imp: frontend/* backend/* errors/*
	go build -o imp

twerp: imp interp/*.go
	go build -o twerp interp/*.go

frontend/lexer.go: frontend/parser.go

frontend/parser.go: frontend/parser.y
	go get -u golang.org/x/tools/cmd/goyacc
	go generate -x

test: imp
	@echo ""
	@echo "Compiling Examples"
	@echo "=================="
	@./imp examples/*.imp
	@echo ""
	$(MAKE) clean

clean:
	$(RM) frontend/y.output
	$(RM) imp twerp
