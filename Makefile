.PHONY: all deps build test install clean

all: test install

deps:
	go get -d -v -t ./...

build: deps
	go build ./...

test: deps
	go test -test.v ./...

install: deps
	go install ./...

clean:
	go clean -i ./...
