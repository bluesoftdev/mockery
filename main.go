package main

import (
	. "code.bluesoftdev.com/v1/repos/mockery/httpMock"
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
			})
		})
		Endpoint("/foo/bar/", func() {
			Method("GET", func() {
				Header("Content-Type", "application/json")
				Header("FOO", "BAR")
				RespondWithFile(200, "./httpMock/ok.json")
			})
		})
	})

	log.Fatal(http.ListenAndServe(":8080", mockery))
}
