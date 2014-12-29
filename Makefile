.PHONY: all build test cov install doc clean

all: test install

build:
	go build ./...

test:
	go test -test.v ./...

cov:
	go get -v github.com/axw/gocov/gocov
	gocov test | gocov report

install:
	go install ./...

doc:
	go get -v github.com/robertkrimen/godocdown/godocdown
	cp .readme.header README.md
	godocdown | tail -n +7 >> README.md

clean:
	go clean -i ./...
