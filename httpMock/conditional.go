package httpMock

import (
	"gopkg.in/xmlpath.v2"
	"net/http"
	"regexp"
	"strings"
)

type Predicate interface {
	Accept(interface{}) bool
}

type PredicateFunc func(interface{}) bool

func (pf PredicateFunc) Accept(v interface{}) bool {
	return pf(v)
}

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

func Not(predicate Predicate) Predicate {
	return PredicateFunc(func(v interface{}) bool {
		return !predicate.Accept(v)
	})
}

var TruePredicate Predicate = PredicateFunc(func(v interface{}) bool { return true })
var FalsePredicate Predicate = PredicateFunc(func(v interface{}) bool { return false })

type Extractor interface {
	Extract(interface{}) interface{}
}

type ExtractorFunc func(interface{}) interface{}

func (ef ExtractorFunc) Extract(v interface{}) interface{} {
	return ef(v)
}

// RequestPredicate is a function that takes a Request and returns true or false.
type RequestPredicate Predicate

func StringEquals(value string) Predicate {
	return PredicateFunc(func(s interface{}) bool {
		return s.(string) == value
	})
}

func StringContains(value string) Predicate {
	return PredicateFunc(func(s interface{}) bool {
		return strings.Contains(s.(string),value)
	})
}

func StringStartsWith(value string) Predicate {
	return PredicateFunc(func(s interface{}) bool {
		return strings.HasPrefix(s.(string), value)
	})
}

func StringEndsWith(value string) Predicate {
	return PredicateFunc(func(s interface{}) bool {
		return strings.HasSuffix(s.(string), value)
	})
}

func StringMatches(regex *regexp.Regexp) Predicate {
	return PredicateFunc(func(s interface{}) bool {
		return regex.MatchString(s.(string))
	})
}

func ExtractedValueAccepted(extractor Extractor, predicate Predicate) Predicate {
	return PredicateFunc(func(v interface{}) bool {
		return predicate.Accept(extractor.Extract(v))
	})
}

// RequestKeySupplier is a function that extracts a key from a Request.  For instance a RequestKeySupplier could
// return the value of a query parameter or a header.
type RequestKeySupplier Extractor

// RequestKeyIdentity is a RequestKeySupplier that returns the entire Request.
var IdentityExtractor RequestKeySupplier = ExtractorFunc(func(r interface{}) interface{} {
	return r
})

var MethodExtractor RequestKeySupplier = ExtractorFunc(func(r interface{}) interface{} {
	return r.(*http.Request).Method
})

func PathMatches(pathRegex *regexp.Regexp) RequestPredicate {
	return ExtractedValueAccepted(ExtractPath, StringMatches(pathRegex))
}

func PathEquals(path string) RequestPredicate {
	return ExtractedValueAccepted(ExtractPath, StringEquals(path))
}

func PathStartsWith(path string) RequestPredicate {
	return ExtractedValueAccepted(ExtractPath, StringStartsWith(path))
}

func HeaderMatches(name string, pathRegex *regexp.Regexp) RequestPredicate {
	return ExtractedValueAccepted(ExtractHeader(name), StringMatches(pathRegex))
}

func HeaderEquals(name string, path string) RequestPredicate {
	return ExtractedValueAccepted(ExtractHeader(name), StringEquals(path))
}

func HeaderEqualsIgnoreCase(name string, path string) RequestPredicate {
	return ExtractedValueAccepted(UpperCaseExtractor(ExtractHeader(name)), StringEquals(strings.ToUpper(path)))
}

func HeaderContains(name string, path string) RequestPredicate {
	return ExtractedValueAccepted(ExtractHeader(name), StringContains(path))
}

func HeaderContainsIgnoreCase(name string, path string) RequestPredicate {
	return ExtractedValueAccepted(UpperCaseExtractor(ExtractHeader(name)), StringContains(strings.ToUpper(path)))
}

func HeaderStartsWith(name string, path string) RequestPredicate {
	return ExtractedValueAccepted(ExtractHeader(name), StringStartsWith(path))
}

func RequestURIMatches(pathRegex *regexp.Regexp) RequestPredicate {
	return ExtractedValueAccepted(ExtractRequestURI, StringMatches(pathRegex))
}

func RequestURIEquals(path string) RequestPredicate {
	return ExtractedValueAccepted(ExtractRequestURI, StringEquals(path))
}

func RequestURIStartsWith(path string) RequestPredicate {
	return ExtractedValueAccepted(ExtractRequestURI, StringStartsWith(path))
}

func MethodIs(method string) RequestPredicate {
	return ExtractedValueAccepted(MethodExtractor, StringEquals(method))
}

func QueryParamEquals(name, value string) RequestPredicate {
	return ExtractedValueAccepted(ExtractQueryParameter(name), StringEquals(value))
}

func QueryParamEqualsIgnoreCase(name, value string) RequestPredicate {
	return ExtractedValueAccepted(UpperCaseExtractor(ExtractQueryParameter(name)), StringEquals(strings.ToUpper(value)))
}

func QueryParamContains(name, value string) RequestPredicate {
	return ExtractedValueAccepted(ExtractQueryParameter(name), StringContains(value))
}

func QueryParamContainsIgnoreCase(name, value string) RequestPredicate {
	return ExtractedValueAccepted(UpperCaseExtractor(ExtractQueryParameter(name)), StringContains(strings.ToUpper(value)))
}

func QueryParamMatches(name string, pattern *regexp.Regexp) RequestPredicate {
	return ExtractedValueAccepted(ExtractQueryParameter(name), StringMatches(pattern))
}
func QueryParamStartsWith(name string, prefix string) RequestPredicate {
	return ExtractedValueAccepted(ExtractQueryParameter(name), StringStartsWith(prefix))
}

var ExtractPath RequestKeySupplier = ExtractorFunc(func(r interface{}) interface{} {
	return r.(*http.Request).URL.Path
})

var ExtractRequestURI RequestKeySupplier = ExtractorFunc(func(r interface{}) interface{} {
	return r.(*http.Request).URL.RequestURI()
})

func ExtractHeader(name string) RequestKeySupplier {
	return ExtractorFunc(func(r interface{}) interface{} {
		return r.(*http.Request).Header.Get(name)
	})
}
func UpperCaseExtractor(extractor Extractor) Extractor {
	return ExtractorFunc(func(v interface{}) interface{} {
		return strings.ToUpper(extractor.Extract(v).(string))
	})
}

// ExtractXPathString returns a RequestKeySupplier that uses XPATH expression to extract a string from the Body of the
// Request.
func ExtractXPathString(xpath string) RequestKeySupplier {
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

// ExtractPathElementByIndex returns a RequestKeySupplier that extracts the path element at the given position.  A
// negative number denotes a position from the end (starting at 1 e.g. -1 is the last element in the path).  For
// positive inputs, the counting starts at 1 as well.
func ExtractPathElementByIndex(idx int) RequestKeySupplier {
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

// ExtractQueryParameter returns a RequestKeySupplier that extracts a query parameters value.
func ExtractQueryParameter(name string) RequestKeySupplier {
	return ExtractorFunc(func(r interface{}) interface{} {
		return r.(*http.Request).URL.Query().Get(name)
	})
}

// RequestKeyPredicate is a function that takes a Request Key provided by a RequestKeyProvider and returns either true
// or false.
type RequestKeyPredicate Predicate

// RequestKeyStringEquals returns a function that will compare the a RequestKey to the string provided and return true
// if the strings are equal.
func RequestKeyStringEquals(str string) RequestKeyPredicate {
	return PredicateFunc(func(key interface{}) bool {
		return key.(string) == str
	})
}

// RequestKeyStringEquals returns a function that will compare the a RequestKey to the regex provided and return true
// if the string matches.
func RequestKeyStringMatches(regexStr string) RequestKeyPredicate {
	regex := regexp.MustCompile(regexStr)
	return PredicateFunc(func(key interface{}) bool {
		return regex.MatchString(key.(string))
	})
}

type when struct {
	predicate     RequestPredicate
	trueResponse  http.Handler
	falseResponse http.Handler
}

// When can be used within a Method's config function to conditionally choose one Response or another.
func When(predicate RequestPredicate, trueResponseBuilder func(), falseResponseBuilder func()) {

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
	predicate RequestKeyPredicate
	response  http.Handler
}

type switchCaseSet struct {
	keySupplier    RequestKeySupplier
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
func Switch(keySupplier RequestKeySupplier, cases func()) {
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
func Case(predicate RequestKeyPredicate, responseBuilder func()) {
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
