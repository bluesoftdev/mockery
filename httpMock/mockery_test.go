package httpMock

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
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
			_, ok := pathMock.methods["GET"]
			assert.True(t, ok, "No GET method found")
		}
	}
}

func TestServeHTTP(t *testing.T) {
	handler := Mockery(func() {
		Endpoint("/foo/bar", func() {
			Method("GET", func() {
				Header("FOO", "BAR")
				RespondWithFile(500, "error.json")
			})
		})
		Endpoint("/foo/bar/", func() {
			Method("GET", func() {
				Header("FOO", "BAR")
				RespondWithFile(200, "ok.json")
			})
		})
	})

	mockWriter := httptest.NewRecorder()
	mockRequest := httptest.NewRequest("GET", "/foo/bar", nil)

	handler.ServeHTTP(mockWriter, mockRequest)

	assert.Equal(t, 500, mockWriter.Code)
	assert.Equal(t, "{\"error\": \"This is an error\"}", mockWriter.Body.String())

	mockWriter = httptest.NewRecorder()
	mockRequest = httptest.NewRequest("GET", "/foo/bar/snafu", nil)

	handler.ServeHTTP(mockWriter, mockRequest)

	assert.Equal(t, 200, mockWriter.Code)
	assert.Equal(t, "{\"ok\": \"everything is ok!\"}", mockWriter.Body.String())
	assert.Equal(t, "BAR", mockWriter.Header().Get("FOO"))
}

func BenchmarkServeHTTP(b *testing.B) {
	handler := Mockery(func() {
		Endpoint("/foo/bar", func() {
			Method("GET", func() {
				Header("FOO", "BAR")
				RespondWithFile(500, "error.json")
			})
		})
		Endpoint("/foo/bar/", func() {
			Method("GET", func() {
				Header("FOO", "BAR")
				RespondWithFile(200, "ok.json")
			})
		})
	})

	for i := 0; i < b.N; i++ {
		mockWriter := httptest.NewRecorder()
		mockRequest := httptest.NewRequest("GET", "/foo/bar/snafu", nil)
		handler.ServeHTTP(mockWriter, mockRequest)
	}
}
