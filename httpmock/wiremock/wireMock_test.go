package wiremock_test

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http/httptest"
	"testing"
	. "github.com/bluesoftdev/mockery/httpmock"
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

	testRequest := httptest.NewRequest("GET", "http://localhost/testfilemapping?foo=bar", nil)
	responseWriter := httptest.NewRecorder()
	mockery.ServeHTTP(responseWriter, testRequest)
	response := responseWriter.Result()
	assert.Equal(t, 200, response.StatusCode)
	assert.Equal(t, "text/plain", response.Header.Get("Content-Type"))
	body, err := ioutil.ReadAll(response.Body)
	assert.NoError(t, err)
	assert.Equal(t, []byte("Testing...\n"), body)
}

func TestWireMockEndpointsFileMapping2(t *testing.T) {
	mockery := Mockery(func() {
		WireMockEndpoints("./wiremock")
	})

	testRequest := httptest.NewRequest("POST", "http://localhost/testfilemapping?foo=bar", nil)
	responseWriter := httptest.NewRecorder()
	mockery.ServeHTTP(responseWriter, testRequest)
	response := responseWriter.Result()
	assert.Equal(t, 200, response.StatusCode)
	assert.Equal(t, "text/plain", response.Header.Get("Content-Type"))
	body, err := ioutil.ReadAll(response.Body)
	assert.NoError(t, err)
	assert.Equal(t, []byte("Testing POST...\n"), body)
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
	js := f.([]interface{})
	js0, ok := js[0].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "value", js0["key"])
	js1, ok := js[1].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "value2", js1["key"])
}


func TestWireMockEndpointsPriorities(t *testing.T) {
	mockery := Mockery(func() {
		WireMockEndpoints("./wiremock")
	})

	testRequest := httptest.NewRequest("GET", "http://localhost/testpriority/12345", nil)
	responseWriter := httptest.NewRecorder()
	mockery.ServeHTTP(responseWriter, testRequest)
	response := responseWriter.Result()

	assert.Equal(t, 200, response.StatusCode)

	data, err := ioutil.ReadAll(response.Body)
	assert.NoError(t, err)
	var f interface{}
	err = json.Unmarshal(data, &f)
	assert.NoError(t, err)
	js, ok := f.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t,"priority100", js["key"])

	testRequest = httptest.NewRequest("GET", "http://localhost/testpriority/23451", nil)
	responseWriter = httptest.NewRecorder()
	mockery.ServeHTTP(responseWriter, testRequest)
	response = responseWriter.Result()

	assert.Equal(t, 200, response.StatusCode)

	data, err = ioutil.ReadAll(response.Body)
	assert.NoError(t, err)
	err = json.Unmarshal(data, &f)
	assert.NoError(t, err)
	js, ok = f.(map[string]interface{})
	assert.True(t,ok)
	assert.Equal(t,"priority101", js["key"])
}

func TestWireMockEndpointsHeaderMatching(t *testing.T) {
	mockery := Mockery(func() {
		WireMockEndpoints("./wiremock")
	})

	testRequest := httptest.NewRequest("GET", "http://localhost/testheadermapping", nil)
	testRequest.Header.Add("Accept","application/xml;encoding=utf-8")
	responseWriter := httptest.NewRecorder()
	mockery.ServeHTTP(responseWriter, testRequest)
	response := responseWriter.Result()

	assert.Equal(t, 404, response.StatusCode)

	testRequest = httptest.NewRequest("GET", "http://localhost/testheadermapping", nil)
	testRequest.Header.Add("Accept","application/json;encoding=utf-8")
	responseWriter = httptest.NewRecorder()
	mockery.ServeHTTP(responseWriter, testRequest)
	response = responseWriter.Result()

	assert.Equal(t, 200, response.StatusCode)
}