package httpmock

import (
	"github.com/bluesoftdev/go-http-matchers/extractor"
	"github.com/bluesoftdev/go-http-matchers/predicate"
	"regexp"
)

// Endpoint defines an endpoint that uses the http.ServeMux to dispatch requests.  The content of the configureFunc
// should be Method elements which may contain
func Endpoint(url string, configureFunc func()) {
	outerCurrentMockHandler := currentMockHandler
	Switch(extractor.ExtractMethod(), configureFunc)
	currentMockery.Handle(url, currentMockHandler)
	currentMockHandler = outerCurrentMockHandler
}

// DefaultPriority is the default priority for endpoint consideration.
const DefaultPriority = 100

// EndpointPattern creates an endpoint that is selected by comparing the URL path with the pattern provided.
func EndpointPattern(urlPattern string, configFunc func()) {
	pathRegex := regexp.MustCompile(urlPattern)
	EndpointForCondition(predicate.PathMatches(pathRegex), configFunc)
}

// EndpointForCondition creates an endpoint that is selected by the predicate passed.
func EndpointForCondition(predicate predicate.Predicate, configFunc func()) {
	outerCurrentMockHandler := currentMockHandler
	configFunc()
	currentMockery.HandleForCondition(DefaultPriority, predicate, currentMockHandler)
	currentMockHandler = outerCurrentMockHandler
}

// EndpointForConditionWithPriority defines an endpoint that is selected by the predicate given with the priority
// provided.
func EndpointForConditionWithPriority(priority int, predicate predicate.Predicate, configFunc func()) {
	outerCurrentMockHandler := currentMockHandler
	configFunc()
	currentMockery.HandleForCondition(priority, predicate, currentMockHandler)
	currentMockHandler = outerCurrentMockHandler
}
