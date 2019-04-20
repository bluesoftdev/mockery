package httpmock

import (
	"github.com/bluesoftdev/go-http-matchers/predicate"
)

// Method is a DSL element that is used within an Endpoint element to define a method handler.
func Method(method string, configFunc func()) {
	Case(predicate.StringEquals(method), configFunc)
}
