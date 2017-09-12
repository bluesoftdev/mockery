package httpMock

import "regexp"

func Endpoint(url string, configureFunc func()) {
	outerCurrentMockHandler := currentMockHandler
	Switch(MethodExtractor, configureFunc)
	currentMockery.Handle(url, currentMockHandler)
	currentMockHandler = outerCurrentMockHandler
}

const DEFAULT_PRIORITY = 100

func EndpointPattern(urlPattern string, configFunc func()) {
	pathRegex := regexp.MustCompile(urlPattern)
	EndpointForCondition(PathMatches(pathRegex),configFunc)
}

func EndpointForCondition(predicate RequestPredicate, configFunc func()) {
	outerCurrentMockHandler := currentMockHandler
	configFunc()
	currentMockery.HandleForCondition(DEFAULT_PRIORITY, predicate, currentMockHandler)
	currentMockHandler = outerCurrentMockHandler
}

func EndpointForConditionWithPriority(priority int,predicate RequestPredicate, configFunc func()) {
	outerCurrentMockHandler := currentMockHandler
	configFunc()
	currentMockery.HandleForCondition(priority, predicate, currentMockHandler)
	currentMockHandler = outerCurrentMockHandler
}