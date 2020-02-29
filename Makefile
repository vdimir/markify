EXE_NAME:=$(shell basename $(CURDIR))

.PHONY: all generate test build run run-debug clean

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

clean:
	rm -f ${EXE_NAME}
