package httpmock_test

import (
	. "github.com/bluesoftdev/mockery/httpmock"
	. "github.com/bluesoftdev/go-http-matchers/predicate"
	. "github.com/bluesoftdev/go-http-matchers/extractor"

	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"
)

func TestWhen(t *testing.T) {
	handler := Mockery(func() {
		Endpoint("/foo/bar", func() {
			Method("GET", func() {
				When(PredicateFunc(func(r interface{}) bool {
					return r.(*http.Request).Header.Get("Accept") == "application/json"
				}), func() {
					RespondWithFile(200, "./ok.json")
					Header("Content-Type", "application/json")
				}, func() {
					Respond(406)
				})
				Header("Foo", "Bar")
			})
		})
	})

	assert.NotNil(t, handler)
	testReq := httptest.NewRequest("GET", "http://localhost/foo/bar", nil)

	mockWriter := httptest.NewRecorder()
	handler.ServeHTTP(mockWriter, testReq)
	result := mockWriter.Result()
	assert.Equal(t, 406, result.StatusCode)
	assert.Equal(t, "Bar", result.Header.Get("FOO"))

	mockWriter = httptest.NewRecorder()
	testReq.Header.Add("Accept", "application/json")
	handler.ServeHTTP(mockWriter, testReq)
	result = mockWriter.Result()
	assert.Equal(t, 200, result.StatusCode)
	assert.Equal(t, "application/json", result.Header.Get("Content-Type"))
	assert.Equal(t, "", result.Trailer.Get("Content-Type"))
	assert.Equal(t, "Bar", result.Header.Get("FOO"))
}

func TestSwitch(t *testing.T) {
	handler := Mockery(func() {
		Endpoint("/foo/bar", func() {
			Method("GET", func() {
				Header("Foo", "Bar")
				Switch(ExtractorFunc(func(r interface{}) interface{} {
					return r.(*http.Request).Header.Get("Accept")
				}), func() {
					Case(PredicateFunc(func(acceptHeader interface{}) bool {
						return acceptHeader.(string) == "application/json"
					}), func() {
						RespondWithFile(200, "./ok.json")
						Header("Content-Type", "application/json")
					})
					Case(PredicateFunc(func(acceptHeader interface{}) bool {
						return acceptHeader.(string) == "application/xml"
					}), func() {
						RespondWithFile(200, "./ok.xml")
						Header("Content-Type", "application/xml")
					})
					Default(func() {
						Respond(406)
					})
				})
			})
		})
	})

	assert.NotNil(t, handler)

	testReq, err := http.NewRequest("GET", "http://localhost/foo/bar", nil)

	assert.NoError(t, err)

	mockWriter := httptest.NewRecorder()
	testReq.Header.Set("Accept", "application/json")
	handler.ServeHTTP(mockWriter, testReq)
	result := mockWriter.Result()
	assert.Equal(t, 200, result.StatusCode)
	assert.Equal(t, "application/json", result.Header.Get("Content-Type"))
	assert.Equal(t, "Bar", result.Header.Get("FOO"))

	mockWriter = httptest.NewRecorder()
	testReq.Header.Set("Accept", "application/xml")
	handler.ServeHTTP(mockWriter, testReq)
	result = mockWriter.Result()
	assert.Equal(t, 200, result.StatusCode)
	assert.Equal(t, "application/xml", result.Header.Get("Content-Type"))
	assert.Equal(t, "Bar", result.Header.Get("FOO"))

	mockWriter = httptest.NewRecorder()
	testReq.Header.Set("Accept", "application/pdf")
	handler.ServeHTTP(mockWriter, testReq)
	result = mockWriter.Result()
	assert.Equal(t, 406, result.StatusCode)
	assert.Equal(t, "Bar", result.Header.Get("FOO"))
}

func TestExtractXPathString(t *testing.T) {
	xml := `<foo><bar snafu="fubar"></bar></foo>`
	path := "/foo/bar/@snafu"
	result := ExtractXPathString(path).Extract(&http.Request{Body: ioutil.NopCloser(strings.NewReader(xml))})
	assert.Equal(t, "fubar", result)
}

func TestExtractQueryParameter(t *testing.T) {
	request := &http.Request{URL: &url.URL{RawQuery: "foo=bar&snafu=fubar"}}
	result := ExtractQueryParameter("foo").Extract(request)
	assert.Equal(t, "bar", result)
	result = ExtractQueryParameter("snafu").Extract(request)
	assert.Equal(t, "fubar", result)
}

func TestExtractPathElementByIndex(t *testing.T) {

	url, _ := url.Parse("http://localhost/foo/bar/snafu")
	request := &http.Request{URL: url}
	result := ExtractPathElementByIndex(-1).Extract(request)
	assert.Equal(t, "snafu", result)

	result = ExtractPathElementByIndex(-2).Extract(request)
	assert.Equal(t, "bar", result)

	result = ExtractPathElementByIndex(-3).Extract(request)
	assert.Equal(t, "foo", result)

	result = ExtractPathElementByIndex(-4).Extract(request)
	assert.Equal(t, "", result)

	result = ExtractPathElementByIndex(4).Extract(request)
	assert.Equal(t, "", result)

	result = ExtractPathElementByIndex(3).Extract(request)
	assert.Equal(t, "snafu", result)

	result = ExtractPathElementByIndex(2).Extract(request)
	assert.Equal(t, "bar", result)

	result = ExtractPathElementByIndex(1).Extract(request)
	assert.Equal(t, "foo", result)
}

func TestRequestKeyStringMatches(t *testing.T) {
	key := "foo"
	assert.False(t, StringMatches(regexp.MustCompile("\\d+")).Accept(key))
	assert.True(t, StringMatches(regexp.MustCompile("[a-z]+")).Accept(key))
}
