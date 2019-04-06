.PHONY: test clean

parser.go: parser.y
	go generate

clean:
	$(RM) parser.go y.output
