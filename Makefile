
SOURCE_GO = $(shell fd -e go)

.PHONY: all
all: bin/textml

bin/textml: $(SOURCE_GO)
	mkdir -p bin
	go build -v -o ./bin/textml ./cmd/textml
