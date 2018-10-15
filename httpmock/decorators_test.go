package httpmock

import (
	"testing"
	"io"
	"os"
	"net/http"
	"net/url"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
)

type countingHandler int
func (ch *countingHandler) ServeHTTP(w http.ResponseWriter,r *http.Request) {
	*ch += 1
}

func TestRespondWithReader(t *testing.T) {
	var counter countingHandler = 0
	currentMockHandler = &counter
	RespondWithReader(200, func() io.Reader {
		file, err := os.Open("./testdata/response.xml")
		if err != nil {
			t.Error("opening response file", err)
			return nil
		}
		return file
	})
	testUrl, _ :=  url.ParseRequestURI("http://localhost/foo")
	request := &http.Request{
		Method: "GET",
		Header: http.Header{},
		URL:    testUrl,
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