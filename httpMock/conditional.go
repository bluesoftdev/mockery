package httpMock

import (
	"gopkg.in/xmlpath.v2"
	"net/http"
	"regexp"
	"strings"
)

// RequestPredicate is a function that takes a request and returns true or false.
type RequestPredicate func(*http.Request) bool

// RequestPredicateTrue is a function that takeas all requests.
var RequestPredicateTrue RequestPredicate = func(r *http.Request) bool {
	return true
}

// RequestKeySupplier is a function that extracts a key from a request.  For instance a RequestKeySupplier could
// return the value of a query parameter or a header.
type RequestKeySupplier func(*http.Request) interface{}

// RequestKeyIdentity is a RequestKeySupplier that returns the entire request.
var RequestKeyIdentity RequestKeySupplier = func(r *http.Request) interface{} {
	return r
}

// ExtractXPathString returns a RequestKeySupplier that uses XPATH expression to extract a string from the body of the
// request.
func ExtractXPathString(xpath string) RequestKeySupplier {
	path := xmlpath.MustCompile(xpath)
	return func(r *http.Request) interface{} {
		str := ""
		root, err := xmlpath.Parse(r.Body)
		if err == nil {
			str, _ = path.String(root)
		}
		return str
	}
}

// ExtractPathElementByIndex returns a RequestKeySupplier that extracts the path element at the given position.  A
// negative number denotes a position from the end (starting at 1 e.g. -1 is the last element in the path).  For
// positive inputs, the counting starts at 1 as well.
func ExtractPathElementByIndex(idx int) RequestKeySupplier {
	return func(r *http.Request) interface{} {
		elements := strings.Split(r.URL.Path, "/")
		var i int
		if idx < 0 {
			i = len(elements) + idx
		} else {
			i = idx
		}
		if i < 0 || i >= len(elements) {
			return ""
		}
		return elements[i]
	}
}

// ExtractQueryParameter returns a RequestKeySupplier that extracts a query parameters value.
func ExtractQueryParameter(name string) RequestKeySupplier {
	return func(r *http.Request) interface{} {
		return r.URL.Query().Get(name)
	}
}

// RequestKeyPredicate is a function that takes a Request Key provided by a RequestKeyProvider and returns either true
// or false.
type RequestKeyPredicate func(interface{}) bool

// RequestKeyStringEquals returns a function that will compare the a RequestKey to the string provided and return true
// if the strings are equal.
func RequestKeyStringEquals(str string) RequestKeyPredicate {
	return func(key interface{}) bool {
		return key.(string) == str
	}
}

// RequestKeyStringEquals returns a function that will compare the a RequestKey to the regex provided and return true
// if the string matches.
func RequestKeyStringMatches(regexStr string) RequestKeyPredicate {
	regex := regexp.MustCompile(regexStr)
	return func(key interface{}) bool {
		return regex.MatchString(key.(string))
	}
}

type when struct {
	predicate     RequestPredicate
	trueResponse  http.Handler
	falseResponse http.Handler
}

// When can be used within a Method's config function to conditionally choose one response or another.
func When(predicate RequestPredicate, trueResponseBuilder func(), falseResponseBuilder func()) {

	outerMockMethodHandler := currentMockMethodHandler
	trueResponseBuilder()
	trueMockMethod := currentMockMethodHandler

	currentMockMethodHandler = outerMockMethodHandler
	falseResponseBuilder()
	falseMockMethod := currentMockMethodHandler

	currentMockMethodHandler = &when{predicate, trueMockMethod, falseMockMethod}
}

func (wh *when) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	if wh.predicate(request) {
		wh.trueResponse.ServeHTTP(w, request)
	} else {
		wh.falseResponse.ServeHTTP(w, request)
	}
}

type switchCase struct {
	predicate RequestKeyPredicate
	response  http.Handler
}

type switchCaseSet struct {
	keySupplier    RequestKeySupplier
	switchCases    []*switchCase
	defaultHandler http.Handler
}

func (scs *switchCaseSet) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	key := scs.keySupplier(request)
	for _, sc := range scs.switchCases {
		if sc.predicate(key) {
			sc.response.ServeHTTP(w, request)
			return
		}
	}
	scs.defaultHandler.ServeHTTP(w, request)
}

var currentSwitch *switchCaseSet

// Switch con be used within a Method's config function to conditionally choose one of many possible responses.  The
// first Case whose predicate returns true will be selected.  Otherwise the response defined in the Default is used.
// If there is no Default, then 404 is returned with an empty body.
func Switch(keySupplier RequestKeySupplier, cases func()) {
	handler := currentMockMethodHandler
	currentSwitch = &switchCaseSet{
		keySupplier:    keySupplier,
		switchCases:    make([]*switchCase, 0, 10),
		defaultHandler: http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
			handler.ServeHTTP(w,request)
			w.WriteHeader(404)
		}),
	}
	cases()
	currentMockMethodHandler = currentSwitch
}

// Case used within a Switch to define a response that will be returned if the case's predicate is true.  The order of
// the case calls matter as the first to match will be used.
func Case(predicate RequestKeyPredicate, responseBuilder func()) {
	if currentMockMethod == nil {
		panic("Switch must be inside a method.")
	}
	outerMockMethodHandler := currentMockMethodHandler
	responseBuilder()
	responseMockMethod := currentMockMethodHandler
	if predicate != nil {
		currentSwitch.switchCases = append(currentSwitch.switchCases, &switchCase{predicate, responseMockMethod})
	} else {
		currentSwitch.defaultHandler = responseMockMethod
	}
	currentMockMethodHandler = outerMockMethodHandler
}

// Default used to define the response that will be returned when no other case is triggered.  The default can be placed
// anywhere but there can only be one.
func Default(responseBuilder func()) {
	Case(nil, responseBuilder)
}
