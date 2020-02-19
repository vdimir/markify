FROM golang:1.13-alpine as build-backend

ADD . /build
WORKDIR /build

RUN go get -v -t -d ./... && \
    go get github.com/rakyll/statik

# RUN go install -v ./...

RUN GOPATH=$(go env GOPATH) go generate ./... && \
    version=$(date +%Y%m%d_%H%M%S) && \
    echo "version=$version" && \
    go build -o markify -ldflags "-X main.revision=${version} -s -w" ./

RUN go test -timeout=60s ./...

FROM alpine:3.11

WORKDIR /srv

COPY --from=build-backend /build/markify /srv/markify

EXPOSE 8080

CMD ["/srv/markify"]
