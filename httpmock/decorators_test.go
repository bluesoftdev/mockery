package httpmock

import (
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
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

func TestHeader(t *testing.T) {
	currentMockHandler = NoopHandler

}
