package main

import (
	. "github.homedepot.com/dxp8048/mockery/httpMock"
	"log"
	"net/http"
)

func main() {

	mockery := Mockery(func() {
		Endpoint("/foo/bar", func() {
			Method("GET", func() {
				Header("Content-Type", "application/json")
				Header("FOO", "BAR")
				RespondWithFile(500, "./httpMock/error.json")
				FixedDelay("10ms")
			})
		})
		Endpoint("/foo/bar/", func() {
			Method("GET", func() {
				Header("Content-Type", "application/json")
				Header("FOO", "BAR")
				RespondWithFile(200, "./httpMock/ok.json")
				NormalDelay("10s", "5s", "20s")
			})
		})
		Endpoint("/snafu/", func() {
			Method("GET", func() {
				Switch(ExtractQueryParameter("foo"), func() {
					Case(RequestKeyStringEquals("bar"), func() {
						Header("Content-Type", "application/xml")
						Header("Cache-Control", "no-cache")
						Header("Access-Control-Allow-Origin", "*")
						RespondWithFile(http.StatusOK, "snafu_foo_bar_response.xml")
					})
				})
			})
			NormalDelay("300ms", "120ms", "5s")
		})
	})

	log.Fatal(http.ListenAndServe(":8080", mockery))
}
