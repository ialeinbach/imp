.PHONY: test clean

parser.go: parser.y
	go generate

test: parser.go
	go run impc/main.go -lv 2 -pv 1 impc/f.imp

clean:
	$(RM) parser.go y.output
