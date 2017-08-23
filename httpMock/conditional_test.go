package httpMock

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"io/ioutil"
	"strings"
)

func TestWhen(t *testing.T) {
	handler := Mockery(func() {
		Endpoint("/foo/bar", func() {
			Method("GET", func() {
				When(func(r *http.Request) bool {
					return r.Header.Get("Accept") == "application/json"
				}, func() {
					RespondWithFile(200, "./ok.json")
					Header("Content-Type", "application/json")
				}, func() {
					Respond(406)
				})
			})
		})
	})

	assert.NotNil(t, handler)
	if assert.IsType(t, &http.ServeMux{}, handler, "mockery is not an http.ServeMux") {
		serveMux, _ := handler.(*http.ServeMux)
		testReq, err := http.NewRequest("GET", "http://localhost/foo/bar", nil)

		assert.NoError(t, err)
		pathHandler, pattern := serveMux.Handler(testReq)
		assert.NotEmpty(t, pattern, "pattern should not be empty: %s", pattern)
		assert.NotNil(t, pathHandler, "path handler should be defined")
		if assert.IsType(t, &mock{}, pathHandler, "path handler is not a mock") {
			pathMock, _ := pathHandler.(*mock)
			getHandler, ok := pathMock.methods["GET"]
			assert.True(t, ok, "No GET method found")
			if assert.IsType(t, &when{}, getHandler, "handler is not a when") {
				mockWriter := httptest.NewRecorder()
				handler.ServeHTTP(mockWriter, testReq)
				result := mockWriter.Result()
				assert.Equal(t, 406, result.StatusCode)

				mockWriter = httptest.NewRecorder()
				testReq.Header.Add("Accept", "application/json")
				handler.ServeHTTP(mockWriter, testReq)
				result = mockWriter.Result()
				assert.Equal(t, 200, result.StatusCode)
			}
		}
	}
}

func TestSwitch(t *testing.T) {
	handler := Mockery(func() {
		Endpoint("/foo/bar", func() {
			Method("GET", func() {
				Switch(func(r *http.Request) interface{} {
					return r.Header.Get("Accept")
				}, func() {
					Case(func(acceptHeader interface{}) bool {
						return acceptHeader == "application/json"
					}, func() {
						RespondWithFile(200, "./ok.json")
						Header("Content-Type", "application/json")
					})
					Case(func(acceptHeader interface{}) bool {
						return acceptHeader == "application/xml"
					}, func() {
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
	if assert.IsType(t, &http.ServeMux{}, handler, "mockery is not an http.ServeMux") {
		serveMux, _ := handler.(*http.ServeMux)
		testReq, err := http.NewRequest("GET", "http://localhost/foo/bar", nil)

		assert.NoError(t, err)
		pathHandler, pattern := serveMux.Handler(testReq)
		assert.NotEmpty(t, pattern, "pattern should not be empty: %s", pattern)
		assert.NotNil(t, pathHandler, "path handler should be defined")
		if assert.IsType(t, &mock{}, pathHandler, "path handler is not a mock") {
			pathMock, _ := pathHandler.(*mock)
			getHandler, ok := pathMock.methods["GET"]
			assert.True(t, ok, "No GET method found")
			if assert.IsType(t, &switchCaseSet{}, getHandler, "handler is not a when") {
				mockWriter := httptest.NewRecorder()
				handler.ServeHTTP(mockWriter, testReq)
				result := mockWriter.Result()
				assert.Equal(t, 406, result.StatusCode)

				mockWriter = httptest.NewRecorder()
				testReq.Header.Set("Accept", "application/json")
				handler.ServeHTTP(mockWriter, testReq)
				result = mockWriter.Result()
				assert.Equal(t, 200, result.StatusCode)
				assert.Equal(t, "application/json", result.Header.Get("Content-Type"))

				mockWriter = httptest.NewRecorder()
				testReq.Header.Set("Accept", "application/xml")
				handler.ServeHTTP(mockWriter, testReq)
				result = mockWriter.Result()
				assert.Equal(t, 200, result.StatusCode)
				assert.Equal(t, "application/xml", result.Header.Get("Content-Type"))

				mockWriter = httptest.NewRecorder()
				testReq.Header.Set("Accept", "application/pdf")
				handler.ServeHTTP(mockWriter, testReq)
				result = mockWriter.Result()
				assert.Equal(t, 406, result.StatusCode)
			}
		}
	}
}

func TestExtractXPathString(t *testing.T) {
	xml := `<foo><bar snafu="fubar"></bar></foo>`
	path := "/foo/bar/@snafu"
	result := ExtractXPathString(path)(&http.Request{Body: ioutil.NopCloser(strings.NewReader(xml))})
	assert.Equal(t,"fubar",result)
}
