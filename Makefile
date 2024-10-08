.PHONY: build run clean tidy

STATELESSDB_TAGS := prod
STATELESSDB_SOURCES := $(shell find ./*.go ./cmd ./internal -type f -iname '*.go' ! -iname '*_test.go')

all: build

tidy:
	go mod tidy

build: statelessdb

statelessdb: $(STATELESSDB_SOURCES) Makefile
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -tags $(STATELESSDB_TAGS) -gcflags -m -ldflags '-extldflags "-static"' -o statelessdb ./cmd/statelessdb

test: Makefile
	go test -tags $(STATELESSDB_TAGS) -v ./...

color-test: Makefile
	grc go test -tags $(STATELESSDB_TAGS) -v ./...

color-bench: Makefile
	grc go test -bench=. -tags $(STATELESSDB_TAGS) -v ./...

clean:
	rm -f ./statelessdb
