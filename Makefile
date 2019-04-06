.PHONY: clean test

test: frontend/parser.go frontend/lexer.go
	go run impc/main.go -lv 2 -pv 2 examples/f.imp

clean:
	$(RM) frontend/parser.go frontend/y.output

frontend/parser.go:
	cd frontend/ && go generate && cd ..
