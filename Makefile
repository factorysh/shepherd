GIT_VERSION?=$(shell git describe --tags --always --abbrev=42 --dirty)

build: bin
	go build \
		-o bin/shepherd \
		-ldflags "-X github.com/factorysh/shepherd/version.version=$(GIT_VERSION)" \
		.

bin:
	mkdir -p bin

test:
	go test -cover \
		github.com/factorysh/shepherd/du \
		github.com/factorysh/shepherd/todo \
		github.com/factorysh/shepherd/shepherd

pull:
	docker pull bearstech/golang-dev

linux:
	GOOS=linux make build
	upx bin/shepherd

docker-build:
	mkdir -p .cache/go-pkg
	docker run -ti --rm \
		-v `pwd`:/src \
		-w /src \
		-v `pwd`/.cache:/.cache \
		-v `pwd`/.cache/go-pkg:/go/pkg \
		-u `id -u` \
		bearstech/golang-dev \
		make build

docker-upx:
	docker run -ti --rm \
		-u `id -u` \
		-v `pwd`/bin:/upx \
		-w /upx \
		bearstech/upx \
		upx shepherd
