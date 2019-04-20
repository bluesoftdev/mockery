package mockery_test

import (
	. "github.com/bluesoftdev/go-http-matchers/extractor"
	. "github.com/bluesoftdev/go-http-matchers/predicate"
	. "github.com/bluesoftdev/mockery/httpmock"
	"log"
	"net/http"
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
