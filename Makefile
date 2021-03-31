EXE_NAME:=$(shell basename $(CURDIR))

.PHONY: all
all: build

.PHONY: test
test:
	go test -timeout=60s ./...

.PHONY: build
build: test
	go build -o ${EXE_NAME} ./

.PHONY: run
run:
	go run ./ --debug --host=localhost

.PHONY: docker
docker:
	export REVISION_INFO="$(shell git diff HEAD --exit-code --quiet || echo '*')$(shell git rev-parse --short HEAD)-$(shell date +%Y%m%d_%H%M%S)-dev"; \
	docker build --build-arg REVISION_INFO=${REVISION_INFO}  . -t markify

.PHONY: clean
clean:
	rm -f ${EXE_NAME}
