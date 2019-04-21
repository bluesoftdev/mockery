package httpmock

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"regexp"
	"testing"
)

type countingHandler int

func (ch *countingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	*ch++
}

func TestRespondWithReader(t *testing.T) {
	var counter countingHandler
	currentMockHandler = &counter
	RespondWithReader(200, func() io.Reader {
		file, err := os.Open("./testdata/response.xml")
		if err != nil {
			t.Error("opening response file", err)
			return nil
		}
		return file
	})
	testURL, _ := url.ParseRequestURI("http://localhost/foo")
	request := &http.Request{
		Method: "GET",
		Header: http.Header{},
		URL:    testURL,
	}
	mockWriter := httptest.NewRecorder()
	currentMockHandler.ServeHTTP(mockWriter, request)
	result := mockWriter.Result()
	assert.Equal(t, 200, result.StatusCode)
	assert.Equal(t, countingHandler(1), counter)
}

type testStruct struct {
	Name string
	Age  int
}

func TestRespondWithJson(t *testing.T) {
	var counter countingHandler
	currentMockHandler = &counter
	RespondWithJson(200, &testStruct{"joe", 28})
	testURL, _ := url.ParseRequestURI("http://localhost/foo")
	request := &http.Request{
		Method: "GET",
		Header: http.Header{},
		URL:    testURL,
	}
	mockWriter := httptest.NewRecorder()
	currentMockHandler.ServeHTTP(mockWriter, request)
	result := mockWriter.Result()
	assert.Equal(t, 200, result.StatusCode)
	assert.Equal(t, countingHandler(1), counter)
	assert.Equal(t, "application/json", result.Header.Get("Content-Type"))
	bodyBytes, err := ioutil.ReadAll(result.Body)
	if assert.NoError(t, err) {
		assert.Equal(t, "{\"Name\":\"joe\",\"Age\":28}", string(bodyBytes))
	}
}

func TestRespondWithString(t *testing.T) {
	var counter countingHandler
	currentMockHandler = &counter
	RespondWithString(200, "The quick brown fox jumped over the lazy dogs.")
	testURL, _ := url.ParseRequestURI("http://localhost/foo")
	request := &http.Request{
		Method: "GET",
		Header: http.Header{},
		URL:    testURL,
	}
	mockWriter := httptest.NewRecorder()
	currentMockHandler.ServeHTTP(mockWriter, request)
	result := mockWriter.Result()
	assert.Equal(t, 200, result.StatusCode)
	assert.Equal(t, countingHandler(1), counter)
	bodyBytes, err := ioutil.ReadAll(result.Body)
	if assert.NoError(t, err) {
		assert.Equal(t, "The quick brown fox jumped over the lazy dogs.", string(bodyBytes))
	}
}

func TestRespondWithFile(t *testing.T) {
	var counter countingHandler
	currentMockHandler = &counter
	RespondWithFile(200, "./testData/ok.json")
	testURL, _ := url.ParseRequestURI("http://localhost/foo")
	request := &http.Request{
		Method: "GET",
		Header: http.Header{},
		URL:    testURL,
	}
	mockWriter := httptest.NewRecorder()
	currentMockHandler.ServeHTTP(mockWriter, request)
	result := mockWriter.Result()
	assert.Equal(t, 200, result.StatusCode)
	assert.Equal(t, countingHandler(1), counter)
	bodyBytes, err := ioutil.ReadAll(result.Body)
	if assert.NoError(t, err) {
		assert.Equal(t, "{\"ok\": \"everything is ok!\"}", string(bodyBytes))
	}
}

func TestRespondWithFileNotFound(t *testing.T) {
	var counter countingHandler
	currentMockHandler = &counter
	RespondWithFile(200, "testData/notok.json")
	testURL, _ := url.ParseRequestURI("http://localhost/foo")
	request := &http.Request{
		Method: "GET",
		Header: http.Header{},
		URL:    testURL,
	}
	mockWriter := httptest.NewRecorder()
	currentMockHandler.ServeHTTP(mockWriter, request)
	result := mockWriter.Result()
	assert.Equal(t, 500, result.StatusCode)
	assert.Equal(t, countingHandler(1), counter)
	bodyBytes, err := ioutil.ReadAll(result.Body)
	if assert.NoError(t, err) {
		assert.Equal(t, "", string(bodyBytes))
	}
}

func TestHeader(t *testing.T) {
	currentMockHandler = NoopHandler
	Header("X-Test", "Bar")
	RespondWithString(200, "The quick brown fox jumped over the lazy dogs.")
	testURL, _ := url.ParseRequestURI("http://localhost/foo")
	request := &http.Request{
		Method: "GET",
		Header: http.Header{},
		URL:    testURL,
	}
	mockWriter := httptest.NewRecorder()
	currentMockHandler.ServeHTTP(mockWriter, request)
	result := mockWriter.Result()
	assert.Equal(t, "Bar", result.Header.Get("X-Test"))
}

func TestTrailer(t *testing.T) {
	currentMockHandler = NoopHandler
	RespondWithString(200, "The quick brown fox jumped over the lazy dogs.")
	Trailer("X-Test", "Bar")
	testURL, _ := url.ParseRequestURI("http://localhost/foo")
	request := &http.Request{
		Method: "GET",
		Header: http.Header{},
		URL:    testURL,
	}
	mockWriter := httptest.NewRecorder()
	currentMockHandler.ServeHTTP(mockWriter, request)
	result := mockWriter.Result()
	assert.Equal(t, "X-Test", result.Header.Get("Trailer"))
	assert.Equal(t, "Bar", result.Trailer.Get("X-Test"))
}

func TestLogLocation(t *testing.T) {
	var counter countingHandler
	currentMockHandler = &counter
	LogLocation("This is a log message")
	RespondWithString(200, "The quick brown fox jumped over the lazy dogs.")
	testURL, _ := url.ParseRequestURI("http://localhost/foo")
	request := &http.Request{
		Method: "GET",
		Header: http.Header{},
		URL:    testURL,
	}
	mockWriter := httptest.NewRecorder()
	var str bytes.Buffer
	log.SetOutput(&str)
	currentMockHandler.ServeHTTP(mockWriter, request)
	logMessage := str.String()
	assert.Regexp(t, regexp.MustCompile("[0-9/: ]{20}Endpoint Defined at [/a-zA-Z0-9_\\-]+/decorators_test.go\\:[0-9]+"+
		"\\(github\\.com/bluesoftdev/mockery/httpmock.TestLogLocation\\)\\: This is a log message\n"), logMessage)
}
