package httpmock

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestLogRequest(t *testing.T) {
	var counter countingHandler
	currentMockHandler = &counter
	LogRequest()
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
	assert.Equal(t, "Request:\nGET /foo HTTP/0.0\r\nHost: localhost\r\n\r\n",logMessage[20:])
}
