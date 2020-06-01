GIT_VERSION?=$(shell git describe --tags --always --abbrev=42 --dirty)

build: bin
	go build .
	go build \
		-o bin/janitor \
		-ldflags "-X github.com/factorysh/janitor-go/version.version=$(GIT_VERSION)" \
		.

bin:
	mkdir -p bin

test:
	go test -v -cover github.com/factorysh/janitor-go/todo
	go test -v -cover github.com/factorysh/janitor-go/janitor
