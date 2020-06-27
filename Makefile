NAME=go-scheme-handler

.PHONY: build clean

all: build

build:
	go build
	mv $(NAME) GoHandler.app/Contents/Resources/Scripts/

clean:
	-rm $(NAME)
