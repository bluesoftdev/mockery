package httpmock

import (
	. "github.com/bluesoftdev/go-http-matchers/predicate"
	. "github.com/bluesoftdev/go-http-matchers/extractor"
	"net/http"
)

type when struct {
	predicate     Predicate
	trueResponse  http.Handler
	falseResponse http.Handler
}

// When can be used within a Method's config function to conditionally choose one Response or another.
func When(predicate Predicate, trueResponseBuilder func(), falseResponseBuilder func()) {

	outerMockMethodHandler := currentMockHandler
	trueResponseBuilder()
	trueMockMethod := currentMockHandler

	currentMockHandler = outerMockMethodHandler
	falseResponseBuilder()
	falseMockMethod := currentMockHandler

	currentMockHandler = &when{predicate, trueMockMethod, falseMockMethod}
}

func (wh *when) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	if wh.predicate.Accept(request) {
		wh.trueResponse.ServeHTTP(w, request)
	} else {
		wh.falseResponse.ServeHTTP(w, request)
	}
}

type switchCase struct {
	predicate Predicate
	response  http.Handler
}

type switchCaseSet struct {
	keySupplier    Extractor
	switchCases    []*switchCase
	defaultHandler http.Handler
}

func (scs *switchCaseSet) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	key := scs.keySupplier.Extract(request)
	for _, sc := range scs.switchCases {
		if sc.predicate.Accept(key) {
			sc.response.ServeHTTP(w, request)
			return
		}
	}
	scs.defaultHandler.ServeHTTP(w, request)
}

var currentSwitch *switchCaseSet

// Switch can be used within a Method's config function to conditionally choose one of many possible responses.  The
// first Case whose predicate returns true will be selected.  Otherwise the Response defined in the Default is used.
// If there is no Default, then 404 is returned with an empty Body.
func Switch(keySupplier Extractor, cases func()) {
	handler := currentMockHandler
	sw := &switchCaseSet{
		keySupplier: keySupplier,
		switchCases: make([]*switchCase, 0, 10),
		defaultHandler: http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
			handler.ServeHTTP(w, request)
			w.WriteHeader(404)
		}),
	}
	outerSwitch := currentSwitch
	currentSwitch = sw
	cases()
	currentMockHandler = currentSwitch
	currentSwitch = outerSwitch
}

// Case used within a Switch to define a Response that will be returned if the case's predicate is true.  The order of
// the case calls matter as the first to match will be used.
func Case(predicate Predicate, responseBuilder func()) {
	outerMockMethodHandler := currentMockHandler
	responseBuilder()
	responseMockMethod := currentMockHandler
	if predicate != nil {
		currentSwitch.switchCases = append(currentSwitch.switchCases, &switchCase{predicate, responseMockMethod})
	} else {
		currentSwitch.defaultHandler = responseMockMethod
	}
	currentMockHandler = outerMockMethodHandler
}

// Default used to define the Response that will be returned when no other case is triggered.  The default can be placed
// anywhere but there can only be one.
func Default(responseBuilder func()) {
	Case(nil, responseBuilder)
}
