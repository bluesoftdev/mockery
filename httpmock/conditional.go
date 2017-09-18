package httpmock

import (
	"gopkg.in/xmlpath.v2"
	"net/http"
	"regexp"
	"strings"
)

// Predicate is a class that can accept or reject a value based on some condition.
type Predicate interface {
	Accept(interface{}) bool
}

// PredicateFunc is an implementation of Predicate that is a function and calls itself on a call to Accept
type PredicateFunc func(interface{}) bool

func (pf PredicateFunc) Accept(v interface{}) bool {
	return pf(v)
}

// And returns a predicate that is true if all of the passed predicates are true for the input.
func And(predicates ...Predicate) Predicate {
	return PredicateFunc(func(v interface{}) bool {
		for _, p := range predicates {
			if !p.Accept(v) {
				return false
			}
		}
		return true
	})
}

// Or returns a predicate that is true if any of the passed predicates are true.  Furthermore, it
// stops executing predicates after the first true one.
func Or(predicates ...Predicate) Predicate {
	return PredicateFunc(func(v interface{}) bool {
		for _, p := range predicates {
			if p.Accept(v) {
				return true
			}
		}
		return false
	})
}

// Not returns a predicate that negates the condition defined by the passed predicate.
func Not(predicate Predicate) Predicate {
	return PredicateFunc(func(v interface{}) bool {
		return !predicate.Accept(v)
	})
}

// TruePredicate is a predicate that returns true for all inputs.
var TruePredicate Predicate = PredicateFunc(func(v interface{}) bool { return true })

// FalsePredicate is a predicate that returns false for all inputs.
var FalsePredicate Predicate = PredicateFunc(func(v interface{}) bool { return false })

// Extractor can extract a value from another value by calling the Extract method.
type Extractor interface {
	Extract(interface{}) interface{}
}

// ExtractorFunc is a function that calls itself when it's Extract emthod is called.
type ExtractorFunc func(interface{}) interface{}

func (ef ExtractorFunc) Extract(v interface{}) interface{} {
	return ef(v)
}

// A predicate that returns true if the value passed is a string and is equal to the value of 'value'
func StringEquals(value string) Predicate {
	return PredicateFunc(func(s interface{}) bool {
		return s.(string) == value
	})
}

// A predicate that returns true if the value passed contains a substring matching 'value'.
func StringContains(value string) Predicate {
	return PredicateFunc(func(s interface{}) bool {
		return strings.Contains(s.(string), value)
	})
}

// A predicate that returns true if the value passed starts with a substring matching 'value'.
func StringStartsWith(value string) Predicate {
	return PredicateFunc(func(s interface{}) bool {
		return strings.HasPrefix(s.(string), value)
	})
}

// A predicate that returns true if the value passed ends with a substring matching 'value'.
func StringEndsWith(value string) Predicate {
	return PredicateFunc(func(s interface{}) bool {
		return strings.HasSuffix(s.(string), value)
	})
}

// A predicate that returns true if the regex matches 'value'.
func StringMatches(regex *regexp.Regexp) Predicate {
	return PredicateFunc(func(s interface{}) bool {
		return regex.MatchString(s.(string))
	})
}

// ExtractedValueAccepted returns A predicate that extracts a value using Extractor and passes that value to the
// provided predicate
func ExtractedValueAccepted(extractor Extractor, predicate Predicate) Predicate {
	return PredicateFunc(func(v interface{}) bool {
		return predicate.Accept(extractor.Extract(v))
	})
}

// RequestKeyIdentity is a Extractor that returns the entire Request.
var IdentityExtractor Extractor = ExtractorFunc(func(r interface{}) interface{} {
	return r
})

// MethodExtractor is an extractor that returns the method of an http.Request.
var MethodExtractor Extractor = ExtractorFunc(func(r interface{}) interface{} {
	return r.(*http.Request).Method
})

// PathMatches returns a predicate that returns true if the path matches the pathRegex.
func PathMatches(pathRegex *regexp.Regexp) Predicate {
	return ExtractedValueAccepted(ExtractPath, StringMatches(pathRegex))
}

// PathEquals returns a predicate that returns true if the path equals 'path'
func PathEquals(path string) Predicate {
	return ExtractedValueAccepted(ExtractPath, StringEquals(path))
}

// PathStartsWith returns a predicate that returns true if the path starts with 'path'
func PathStartsWith(path string) Predicate {
	return ExtractedValueAccepted(ExtractPath, StringStartsWith(path))
}

// HeaderMatches returns a predicate that returns true if the header named 'name' matches 'regex'
func HeaderMatches(name string, regex *regexp.Regexp) Predicate {
	return ExtractedValueAccepted(ExtractHeader(name), StringMatches(regex))
}

// HeaderEquals returns a predicate that returns true if the header named 'name' equals 'value'
func HeaderEquals(name string, value string) Predicate {
	return ExtractedValueAccepted(ExtractHeader(name), StringEquals(value))
}

// HeaderEqualsIgnoreCase returns a predicate that returns true if if the header named 'name' equals 'value', ignoring case.
func HeaderEqualsIgnoreCase(name string, path string) Predicate {
	return ExtractedValueAccepted(UpperCaseExtractor(ExtractHeader(name)), StringEquals(strings.ToUpper(path)))
}

func HeaderContains(name string, path string) Predicate {
	return ExtractedValueAccepted(ExtractHeader(name), StringContains(path))
}

func HeaderContainsIgnoreCase(name string, path string) Predicate {
	return ExtractedValueAccepted(UpperCaseExtractor(ExtractHeader(name)), StringContains(strings.ToUpper(path)))
}

func HeaderStartsWith(name string, path string) Predicate {
	return ExtractedValueAccepted(ExtractHeader(name), StringStartsWith(path))
}

func RequestURIMatches(pathRegex *regexp.Regexp) Predicate {
	return ExtractedValueAccepted(ExtractRequestURI, StringMatches(pathRegex))
}

func RequestURIEquals(path string) Predicate {
	return ExtractedValueAccepted(ExtractRequestURI, StringEquals(path))
}

func RequestURIStartsWith(path string) Predicate {
	return ExtractedValueAccepted(ExtractRequestURI, StringStartsWith(path))
}

// MethodIs returns a predicate that takes a request, extracts the method, and returns true if it equals the method
// provided, ignoring case.
func MethodIs(method string) Predicate {
	return ExtractedValueAccepted(UpperCaseExtractor(MethodExtractor), StringEquals(strings.ToUpper(method)))
}

// QueryParamContainsIgnoreCase returns a Predicate that takes a request, extracts the query parameter specified and
// returns true if it equals the value provided.
func QueryParamEquals(name, value string) Predicate {
	return ExtractedValueAccepted(ExtractQueryParameter(name), StringEquals(value))
}

// QueryParamContainsIgnoreCase returns a Predicate that takes a request, extracts the query parameter specified and
// returns true if it equals the value provided, ignoring case.
func QueryParamEqualsIgnoreCase(name, value string) Predicate {
	return ExtractedValueAccepted(UpperCaseExtractor(ExtractQueryParameter(name)), StringEquals(strings.ToUpper(value)))
}

// QueryParamContainsIgnoreCase returns a Predicate that takes a request, extracts the query parameter specified and
// returns true if it contains the value provided.
func QueryParamContains(name, value string) Predicate {
	return ExtractedValueAccepted(ExtractQueryParameter(name), StringContains(value))
}

// QueryParamContainsIgnoreCase returns a Predicate that takes a request, extracts the query parameter specified and
// returns true if it contains the value provided, ignoring case.
func QueryParamContainsIgnoreCase(name, value string) Predicate {
	return ExtractedValueAccepted(UpperCaseExtractor(ExtractQueryParameter(name)), StringContains(strings.ToUpper(value)))
}

// QueryParamMatches returns a Predicate that takes a request, extracts the query parameter specified and
// returns true if the value matches the pattern provided.
func QueryParamMatches(name string, pattern *regexp.Regexp) Predicate {
	return ExtractedValueAccepted(ExtractQueryParameter(name), StringMatches(pattern))
}

// QueryParamStartsWith returns a Predicate that takes a request, extracts the query parameter specified and
// returns true if the value starts with the prefix provided.
func QueryParamStartsWith(name string, prefix string) Predicate {
	return ExtractedValueAccepted(ExtractQueryParameter(name), StringStartsWith(prefix))
}

// ExtractPath is an Extractor that returns the request URL's Path property.  This is just path path portion of the URI.
var ExtractPath Extractor = ExtractorFunc(func(r interface{}) interface{} {
	return r.(*http.Request).URL.Path
})

// ExtractRequestURI is an Extractor that returns the request URL's RequestURI property.  This is the path and the query
// portions of the URI.
var ExtractRequestURI Extractor = ExtractorFunc(func(r interface{}) interface{} {
	return r.(*http.Request).URL.RequestURI()
})

// ExtractHeader returns a Extractor that returns the value of the header named 'name'
func ExtractHeader(name string) Extractor {
	return ExtractorFunc(func(r interface{}) interface{} {
		return r.(*http.Request).Header.Get(name)
	})
}

// UpperCaseExctractor returns an Extractor that decorates the passed extractor by applying strings.ToUpper to the
// value returned.
func UpperCaseExtractor(extractor Extractor) Extractor {
	return ExtractorFunc(func(v interface{}) interface{} {
		return strings.ToUpper(extractor.Extract(v).(string))
	})
}

// ExtractXPathString returns a Extractor that uses XPATH expression to extract a string from the Body of the
// Request.
func ExtractXPathString(xpath string) Extractor {
	path := xmlpath.MustCompile(xpath)
	return ExtractorFunc(func(r interface{}) interface{} {
		str := ""
		root, err := xmlpath.Parse(r.(*http.Request).Body)
		if err == nil {
			str, _ = path.String(root)
		}
		return str
	})
}

// ExtractPathElementByIndex returns a Extractor that extracts the path element at the given position.  A
// negative number denotes a position from the end (starting at 1 e.g. -1 is the last element in the path).  For
// positive inputs, the counting starts at 1 as well.
func ExtractPathElementByIndex(idx int) Extractor {
	return ExtractorFunc(func(r interface{}) interface{} {
		elements := strings.Split(r.(*http.Request).URL.Path, "/")
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
	})
}

// ExtractQueryParameter returns a Extractor that extracts a query parameters value.
func ExtractQueryParameter(name string) Extractor {
	return ExtractorFunc(func(r interface{}) interface{} {
		return r.(*http.Request).URL.Query().Get(name)
	})
}

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

// Switch con be used within a Method's config function to conditionally choose one of many possible responses.  The
// first Case whose predicate returns true will be selected.  Otherwise the Response defined in the Default is used.
// If there is no Default, then 404 is returned with an empty Body.
func Switch(keySupplier Extractor, cases func()) {
	handler := currentMockHandler
	currentSwitch = &switchCaseSet{
		keySupplier: keySupplier,
		switchCases: make([]*switchCase, 0, 10),
		defaultHandler: http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
			handler.ServeHTTP(w, request)
			w.WriteHeader(404)
		}),
	}
	cases()
	currentMockHandler = currentSwitch
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
