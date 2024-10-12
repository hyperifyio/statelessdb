.PHONY: all test bench

all: test bench

test:
	grc go test -tags prod -v ./...

bench: 
	grc go test -bench=.  ./...
