.PHONY: build test lint clean install

build:
	CGO_ENABLED=0 go build -o bin/supervisor ./cmd/supervisor

test:
	go test ./... -v -race -coverprofile=coverage.out

lint:
	golangci-lint run ./...

clean:
	rm -rf bin/ coverage.out

install: build
	cp bin/supervisor /usr/local/bin/

.DEFAULT_GOAL := build
