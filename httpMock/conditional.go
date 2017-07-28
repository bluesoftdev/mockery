package httpMock

import "net/http"

type RequestPredicate func(*http.Request) bool


type when struct {
	predicate RequestPredicate
	trueResponse http.Handler
	falseResponse http.Handler
}

func When(predicate RequestPredicate, trueResponseBuilder func(), falseResponseBuilder func()) {

	outerMockMethod := currentMockMethod
	Method(outerMockMethod.method,trueResponseBuilder)
	trueMockMethod := currentMock.methods[outerMockMethod.method]
	Method(outerMockMethod.method,falseResponseBuilder)
	falseMockMethod := currentMock.methods[outerMockMethod.method]

	currentMockMethodHandler = &when{predicate,trueMockMethod,falseMockMethod}
}

func (wh *when) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	if wh.predicate(request) {
		wh.trueResponse.ServeHTTP(w, request)
	} else {
		wh.falseResponse.ServeHTTP(w, request)
	}
}