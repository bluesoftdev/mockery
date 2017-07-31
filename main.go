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
				NormalDelay("100ms", "20ms", "500ms")
			})
		})
	})

	log.Fatal(http.ListenAndServe(":8080", mockery))
}
