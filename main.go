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
		Endpoint("/VPPService/rs/vpp/calculateVolumePricing/", func() {
			Method("POST", func() {
				Header("Content-Type", "application/xml")
				Header("Cache-Control", "no-cache")
				Header("Access-Control-Allow-Origin", "*")
				RespondWithFile(http.StatusOK, "usom_pricing_bidroom_service_response.xml")
				NormalDelay("300ms", "120ms", "5s")
			})
		})
	})

	log.Fatal(http.ListenAndServe(":8080", mockery))
}
