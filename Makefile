EXE_NAME:=$(shell basename $(CURDIR))

.PHONY: all generate test build run run-debug clean

all: build

generate:
	go generate -x ./...

test: generate
	go test ./...

build: test generate
	go build -o ${EXE_NAME} ./

run: generate
	go run ./

run-debug: generate
	go run ./ --debug --host=localhost

clean:
	rm -f ${EXE_NAME}