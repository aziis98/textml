
.PHONY: all
all: bin/textml

bin/textml:
	mkdir -p bin
	go build -v -o ./bin/textml ./cmd/textml
