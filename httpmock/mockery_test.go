package httpmock_test

import (
	. "github.com/bluesoftdev/mockery/httpmock"

	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"regexp"
	"testing"
)

func TestServeHTTP(t *testing.T) {
	handler := Mockery(func() {
		Header("SNAFU", "BAZ")
		pathPattern := regexp.MustCompile("/foo/bar/snafu.*")
		EndpointForCondition(
			And(PathMatches(pathPattern), MethodIs("GET")),
			func() {
				Header("FOO", "SNAFU")
				RespondWithFile(200, "ok.json")
			})
		Endpoint("/foo/bar", func() {
			Method("GET", func() {
				Header("FOO", "BAR")
				RespondWithFile(500, "error.json")
			})
		})
		Endpoint("/foo/bar/", func() {
			Method("GET", func() {
				Header("FOO", "BAZ")
				RespondWithFile(200, "ok.json")
			})
		})
	})

	mockWriter := httptest.NewRecorder()
	mockRequest := httptest.NewRequest("GET", "/foo/bar", nil)

	handler.ServeHTTP(mockWriter, mockRequest)

	assert.Equal(t, 500, mockWriter.Code)
	assert.Equal(t, "{\"error\": \"This is an error\"}", mockWriter.Body.String())
	assert.Equal(t, "BAZ", mockWriter.Header().Get("SNAFU"))

	mockWriter = httptest.NewRecorder()
	mockRequest = httptest.NewRequest("GET", "/foo/bar/snafu", nil)

	handler.ServeHTTP(mockWriter, mockRequest)

	assert.Equal(t, 200, mockWriter.Code)
	assert.Equal(t, "{\"ok\": \"everything is ok!\"}", mockWriter.Body.String())
	assert.Equal(t, "SNAFU", mockWriter.Header().Get("FOO"))

	mockWriter = httptest.NewRecorder()
	mockRequest = httptest.NewRequest("GET", "/foo/bar/fubar", nil)

	handler.ServeHTTP(mockWriter, mockRequest)

	assert.Equal(t, 200, mockWriter.Code)
	assert.Equal(t, "{\"ok\": \"everything is ok!\"}", mockWriter.Body.String())
	assert.Equal(t, "BAZ", mockWriter.Header().Get("FOO"))
}

func BenchmarkServeHTTP(b *testing.B) {
	handler := Mockery(func() {
		Endpoint("/foo/bar", func() {
			Method("GET", func() {
				Header("FOO", "BAR")
				RespondWithFile(500, "error.json")
			})
		})
		EndpointPattern("/foo/bar/snafu", func() {

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
