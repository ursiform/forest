# forest

[![Coverage Status](https://coveralls.io/repos/github/ursiform/forest/badge.svg?branch=master)](https://coveralls.io/github/ursiform/forest?branch=master)

[![API documentation](https://godoc.org/github.com/ursiform/forest?status.svg)](https://godoc.org/github.com/ursiform/forest)

`forest` is a micro-framework for building REST services that talk JSON. Its
core unit is a [`forest.App`](https://godoc.org/github.com/ursiform/forest#App)
that is built upon a [`bear`](https://github.com/ursiform/bear) multiplexer for
URL routing. It outputs responses using
[`forest.Response`](https://godoc.org/github.com/ursiform/forest#Response)
and provides utility methods for many common tasks required by web services.

## Install
    go get github.com/ursiform/forest

## Test
    go test -cover github.com/ursiform/forest

## API

[![API documentation](https://godoc.org/github.com/ursiform/forest?status.svg)](https://godoc.org/github.com/ursiform/forest)

## License
[MIT License](LICENSE)
