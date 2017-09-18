package httpmock

import "regexp"

// Endpoint defines an endpoint that uses the http.ServeMux to dispatch requests.  The content of the configureFunc
// should be Method elements which may contain
func Endpoint(url string, configureFunc func()) {
	outerCurrentMockHandler := currentMockHandler
	Switch(MethodExtractor, configureFunc)
	currentMockery.Handle(url, currentMockHandler)
	currentMockHandler = outerCurrentMockHandler
}

const DEFAULT_PRIORITY = 100

// EndpointPattern creates an endpoint that is selected by comparing the URL path with the pattern provided.
func EndpointPattern(urlPattern string, configFunc func()) {
	pathRegex := regexp.MustCompile(urlPattern)
	EndpointForCondition(PathMatches(pathRegex), configFunc)
}

// EndpointForCondition creates an endpoint that is selected by the predicate passed.
func EndpointForCondition(predicate Predicate, configFunc func()) {
	outerCurrentMockHandler := currentMockHandler
	configFunc()
	currentMockery.HandleForCondition(DEFAULT_PRIORITY, predicate, currentMockHandler)
	currentMockHandler = outerCurrentMockHandler
}

// EndpointForConditionWithPriority defines an endpoint that is selected by the predicate given with the priority
// provided.
func EndpointForConditionWithPriority(priority int, predicate Predicate, configFunc func()) {
	outerCurrentMockHandler := currentMockHandler
	configFunc()
	currentMockery.HandleForCondition(priority, predicate, currentMockHandler)
	currentMockHandler = outerCurrentMockHandler
}
