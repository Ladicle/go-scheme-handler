NAME=go-scheme-handler

.PHONY: build test clean

all: build

build:
	go build
	mv $(NAME) GoHandler.app/Contents/MacOS/bin

test:
	go test

clean:
	-rm $(NAME)
