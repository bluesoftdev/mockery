package httpmock
import (
	. "gitlab.com/ComputersFearMe/go-http-matchers/predicate"
)
// Method is a DSL element that is used within an Endpoint element to define a method handler.
func Method(method string, configFunc func()) {
	Case(StringEquals(method), configFunc)
}
