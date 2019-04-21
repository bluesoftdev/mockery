package httpmock

import (
	"github.com/bluesoftdev/go-http-matchers/predicate"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestEndpointForConditionNoMatch(t *testing.T) {
	mockery := Mockery(func() {
		EndpointForCondition(predicate.False(), func() {
			Respond(200)
		})
	})

	// build the request
	testURL, _ := url.ParseRequestURI("http://localhost/foo")
	request := &http.Request{
		Method: "GET",
		Header: http.Header{},
		URL:    testURL,
	}

	// Execute the request
	mockWriter := httptest.NewRecorder()
	mockery.ServeHTTP(mockWriter, request)
	response := mockWriter.Result()

	// Test the result.
	assert.Equal(t, 404, response.StatusCode)
}

func TestEndpointForConditionMatch(t *testing.T) {
	mockery := Mockery(func() {
		EndpointForCondition(predicate.True(), func() {
			Respond(200)
		})
	})

	// build the request
	testURL, _ := url.ParseRequestURI("http://localhost/foo")
	request := &http.Request{
		Method: "GET",
		Header: http.Header{},
		URL:    testURL,
	}

	// Execute the request
	mockWriter := httptest.NewRecorder()
	mockery.ServeHTTP(mockWriter, request)
	response := mockWriter.Result()

	// Test the result.
	assert.Equal(t, 200, response.StatusCode)
}

func TestEndpointPatternNotFound(t *testing.T) {
	mockery := Mockery(func() {
		EndpointPattern("/foo/.+", func() {
			Respond(200)
		})
	})

	// build the request
	testURL, _ := url.ParseRequestURI("http://localhost/foo")
	request := &http.Request{
		Method: "GET",
		Header: http.Header{},
		URL:    testURL,
	}

	// Execute the request
	mockWriter := httptest.NewRecorder()
	mockery.ServeHTTP(mockWriter, request)
	response := mockWriter.Result()

	// Test the result.
	assert.Equal(t, 404, response.StatusCode)
}

func TestEndpointPatternMatch(t *testing.T) {
	mockery := Mockery(func() {
		EndpointPattern("/fo{2}", func() {
			Respond(200)
		})
	})

	// build the request
	testURL, _ := url.ParseRequestURI("http://localhost/foo")
	request := &http.Request{
		Method: "GET",
		Header: http.Header{},
		URL:    testURL,
	}

	// Execute the request
	mockWriter := httptest.NewRecorder()
	mockery.ServeHTTP(mockWriter, request)
	response := mockWriter.Result()

	// Test the result.
	assert.Equal(t, 200, response.StatusCode)
}

func TestEndpointMatch(t *testing.T) {
	mockery := Mockery(func() {
		EndpointPattern("/foo", func() {
			Respond(200)
		})
	})

	// build the request
	testURL, _ := url.ParseRequestURI("http://localhost/foo")
	request := &http.Request{
		Method: "GET",
		Header: http.Header{},
		URL:    testURL,
	}

	// Execute the request
	mockWriter := httptest.NewRecorder()
	mockery.ServeHTTP(mockWriter, request)
	response := mockWriter.Result()

	// Test the result.
	assert.Equal(t, 200, response.StatusCode)
}

func TestEndpointNoMatch(t *testing.T) {
	mockery := Mockery(func() {
		EndpointPattern("/bar", func() {
			Respond(200)
		})
	})

	// build the request
	testURL, _ := url.ParseRequestURI("http://localhost/foo")
	request := &http.Request{
		Method: "GET",
		Header: http.Header{},
		URL:    testURL,
	}

	// Execute the request
	mockWriter := httptest.NewRecorder()
	mockery.ServeHTTP(mockWriter, request)
	response := mockWriter.Result()

	// Test the result.
	assert.Equal(t, 404, response.StatusCode)
}

func TestEndpointForConditionWithPriority(t *testing.T) {
	mockery := Mockery(func() {
		EndpointForConditionWithPriority(2, predicate.True(), func() {
			Respond(201)
		})
		EndpointForConditionWithPriority(1, predicate.True(), func() {
			Respond(200)
		})
	})

	// build the request
	testURL, _ := url.ParseRequestURI("http://localhost/foo")
	request := &http.Request{
		Method: "GET",
		Header: http.Header{},
		URL:    testURL,
	}

	// Execute the request
	mockWriter := httptest.NewRecorder()
	mockery.ServeHTTP(mockWriter, request)
	response := mockWriter.Result()

	// Test the result.
	assert.Equal(t, 200, response.StatusCode)
}
