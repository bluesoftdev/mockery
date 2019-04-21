[![Go Lang Version](https://img.shields.io/badge/go-1.9-blue.svg?style=plastic)](http://golang.com)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=plastic)](https://godoc.org/github.com/bluesoftdev/mockery)
[![Go Report Card](https://goreportcard.com/badge/github.com/bluesoftdev/mockery)](https://goreportcard.com/report/github.com/bluesoftdev/mockery)
[![codecov](https://codecov.io/gh/bluesoftdev/mockery/branch/master/graph/badge.svg)](https://codecov.io/gh/bluesoftdev/mockery)
[![CircleCI](https://circleci.com/gh/bluesoftdev/mockery/tree/master.svg?style=svg)](https://circleci.com/gh/bluesoftdev/mockery/tree/master)

# Mockery
Mockery is a go library that enables programmers to create mock http
servers for the purpose of testing their integrations in isolation.  It
is particularly good at doing performance testing since one instance can
handle a very large number of tps.  I have tested a basic mockery
handling 100,000 tps without using more than 20% CPU on an 8 core
system.

# Getting Started

Here is an example mockery.

``` golang
package mockery_test

import (
	"log"
	"net/http"
	. "github.com/bluesoftdev/mockery/httpmock"
	. "github.com/bluesoftdev/go-http-matchers/extractor"
	. "github.com/bluesoftdev/go-http-matchers/predicate"
)

func main() {
	mockery := Mockery(func() {
		Endpoint("/foo/bar", func() {
			Method("GET", func() {
				Header("Content-Type", "application/json")
				Header("FOO", "BAR")
				RespondWithFile(500, "./error.json")
				FixedDelay("10ms")
			})
		})
		Endpoint("/foo/bar/", func() {
			Method("GET", func() {
				Header("Content-Type", "application/json")
				Header("FOO", "BAR")
				RespondWithFile(200, "./ok.json")
				NormalDelay("10s", "5s", "20s")
			})
		})
		Endpoint("/snafu/", func() {
			Method("GET", func() {
				Header("Content-Type", "application/xml")
				Header("Cache-Control", "no-cache")
				Header("Access-Control-Allow-Origin", "*")
				Switch(ExtractQueryParameter("foo"), func() {
					Case(StringEquals("bar"), func() {
						RespondWithFile(http.StatusOK, "response.xml")
					})
					Default(func() {
						RespondWithFile(http.StatusBadRequest, "error.xml")
					})
				})
			})
			NormalDelay("300ms", "120ms", "5s")
		})
	})

	log.Fatal(http.ListenAndServe(":8080", mockery))
}
```

# Contributing

see [Contributing](CONTRIBUTING.md)
