EXE_NAME:=$(shell basename $(CURDIR))

.PHONY: all generate test build run run-debug clean docker

all: build

generate:
	go generate -x ./...

test: generate
	go test -timeout=60s ./...

build: test generate
	go build -o ${EXE_NAME} ./

run: generate
	go run ./ --host=localhost --data ./var

run-debug: generate
	go run ./ --debug --host=localhost --data ./var

docker:
	export REVISION_INFO="$(shell git diff HEAD --exit-code --quiet || echo '*')$(shell git rev-parse --short HEAD)-$(shell date +%Y%m%d_%H%M%S)-dev"; \
	docker build --build-arg REVISION_INFO=${REVISION_INFO}  . -t markify

clean:
	rm -f ${EXE_NAME}
