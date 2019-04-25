.PHONY: clean test

imp: main.go frontend/parser.go
	go build

frontend/parser.go: frontend/parser.y frontend/lexer.go
	cd frontend/ && go generate && cd ..

clean:
	$(RM) frontend/parser.go frontend/y.output imp

test: imp
	./imp examples/*.imp
