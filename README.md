# [markify.dev](https://markify.dev)

[![Uptime Robot status](https://img.shields.io/uptimerobot/status/m784426448-74cf1366e423303b7543bc6f?label=markify.dev)](https://markify.dev)
[![test](https://github.com/vdimir/markify/workflows/test/badge.svg)](https://github.com/vdimir/markify/actions?query=workflow%3Atest)
[![Go Report Card](https://goreportcard.com/badge/github.com/vdimir/markify)](https://goreportcard.com/report/github.com/vdimir/markify)
[![Coverage Status](https://coveralls.io/repos/github/vdimir/markify/badge.svg)](https://coveralls.io/github/vdimir/markify)

Simple and minimalistic markdown sharing service.

Features:
* Keep clean and simple
* Free to use, open-source code
* Supports some handy extensions like Table of Contents, social media embedding and others

Read more at [markify.dev/about](https://markify.dev/about)

## Development

### Build & Run

Build dockerized app and run:
```
docker-compose up --build
```
Tests checked in docker build.

Run naively on machine (golang installed required):
```
make run
```
or
```
make run-debug # run with --debug flag to reload assets on fly
```

Run tests:
```
make test
```
