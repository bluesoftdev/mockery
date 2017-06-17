package httpMock

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestMockery(t *testing.T) {
	handler := Mockery(func() {
		// DO NOTHING
	})
	assert.NotNil(t, handler, "handler is nil")
}

func TestEndpoint(t *testing.T) {
	handler := Mockery(func() {
		Endpoint("/foo/bar", func() {
			// DO NOTHING
		})
	})
	assert.NotNil(t, handler, "handler is nil")
	if assert.IsType(t, &http.ServeMux{}, handler, "mockery is not an http.ServeMux") {
		serveMux, _ := handler.(*http.ServeMux)
		testReq, err := http.NewRequest("GET", "http://localhost/foo/bar", nil)
		assert.NoError(t, err)
		pathHandler, pattern := serveMux.Handler(testReq)
		assert.NotEmpty(t, pattern, "pattern should not be empty: %s", pattern)
		assert.NotNil(t, pathHandler, "path handler should be defined")
		assert.IsType(t, &mock{}, pathHandler, "path handler is not a mock")
	}
}

func TestMethod(t *testing.T) {
	handler := Mockery(func() {
		Endpoint("/foo/bar", func() {
			Method("GET", func() {
				Header("FOO", "BAR")
				RespondWithFile(500, "error.json")
			})
		})
	})
	assert.NotNil(t, handler, "handler is nil")
	if assert.IsType(t, &http.ServeMux{}, handler, "mockery is not an http.ServeMux") {
		serveMux, _ := handler.(*http.ServeMux)
		testReq, err := http.NewRequest("GET", "http://localhost/foo/bar", nil)
		assert.NoError(t, err)
		pathHandler, pattern := serveMux.Handler(testReq)
		assert.NotEmpty(t, pattern, "pattern should not be empty: %s", pattern)
		assert.NotNil(t, pathHandler, "path handler should be defined")
		if assert.IsType(t, &mock{}, pathHandler, "path handler is not a mock") {
			pathMock, _ := pathHandler.(*mock)
			getMock, ok := pathMock.methods["GET"]
			assert.True(t, ok, "No GET method found")
			fooValue, ok := getMock.headers["FOO"]
			assert.True(t, ok, "No FOO header")
			assert.Equal(t, "BAR", fooValue)
			assert.Equal(t, "error.json", getMock.responseFileName)
			assert.Equal(t, 500, getMock.statusCode)
		}
	}
}
