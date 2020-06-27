NAME=go-scheme-handler

.PHONY: build clean

all: build

build:
	go build
	mv $(NAME) GoHandler.app/Contents/MacOS/bin

clean:
	-rm $(NAME)
