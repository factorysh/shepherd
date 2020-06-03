GIT_VERSION?=$(shell git describe --tags --always --abbrev=42 --dirty)

build: bin
	go build \
		-o bin/shepherd \
		-ldflags "-X github.com/factorysh/shepherd/version.version=$(GIT_VERSION)" \
		.

bin:
	mkdir -p bin

test:
	go test -v -cover github.com/factorysh/shepherd/todo
	go test -v -cover github.com/factorysh/shepherd/shepherd

pull:
	docker pull bearstech/golang-dev

docker-build:
	mkdir -p .cache
	docker run -ti --rm \
		-v `pwd`:/src \
		-w /src \
		-v `pwd`/.cache:/.cache \
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