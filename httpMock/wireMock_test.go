package httpMock

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http/httptest"
	"testing"
)

func TestWireMockEndpoints(t *testing.T) {
	mockery := Mockery(func() {
		WireMockEndpoints("./wiremock")
	})

	testRequest := httptest.NewRequest("GET", "http://localhost/testmapping", nil)
	responseWriter := httptest.NewRecorder()
	mockery.ServeHTTP(responseWriter, testRequest)
	response := responseWriter.Result()
	assert.Equal(t, 200, response.StatusCode)
	assert.Equal(t, "text/plain", response.Header.Get("Content-Type"))
	body, err := ioutil.ReadAll(response.Body)
	assert.NoError(t, err)
	assert.Equal(t, []byte("default test mapping"), body)
}

func TestWireMockEndpointsFileMapping(t *testing.T) {
	mockery := Mockery(func() {
		WireMockEndpoints("./wiremock")
	})

	testRequest := httptest.NewRequest("GET", "http://localhost/testfilemapping", nil)
	responseWriter := httptest.NewRecorder()
	mockery.ServeHTTP(responseWriter, testRequest)
	response := responseWriter.Result()
	assert.Equal(t, 200, response.StatusCode)
	assert.Equal(t, "text/plain", response.Header.Get("Content-Type"))
	body, err := ioutil.ReadAll(response.Body)
	assert.NoError(t, err)
	assert.Equal(t, []byte("Testing...\n"), body)
}

func TestWireMockEndpointsJsonMapping(t *testing.T) {
	mockery := Mockery(func() {
		WireMockEndpoints("./wiremock")
	})

	testRequest := httptest.NewRequest("GET", "http://localhost/testjsonmapping", nil)
	responseWriter := httptest.NewRecorder()
	mockery.ServeHTTP(responseWriter, testRequest)
	response := responseWriter.Result()
	assert.Equal(t, 200, response.StatusCode)
	assert.Equal(t, "application/json", response.Header.Get("Content-Type"))
	data, err := ioutil.ReadAll(response.Body)
	assert.NoError(t, err)
	var f interface{}
	err = json.Unmarshal(data, &f)
	assert.NoError(t, err)
	js := f.(map[string]interface{})
	assert.Equal(t, "value", js["key"])
}
