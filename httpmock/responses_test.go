package httpmock

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

const testString = "The quick brown fox jumped over the lazy dogs."

var responseTests = []struct {
	Name         string
	WriteFn      func()
	ExpectedBody string
	Status       int
}{
	{
		"String",
		func() {
			WriteStatusAndBody(200, testString)
		},
		testString,
		200,
	},
	{
		"[]byte",
		func() {
			WriteStatusAndBody(200, []byte(testString))
		},
		testString,
		200,
	},
	{
		"Reader",
		func() {
			WriteStatusAndBody(200, bytes.NewBufferString(testString))
		},
		testString,
		200,
	},
	{
		"ReadCloser",
		func() {
			WriteStatusAndBody(200, ioutil.NopCloser(bytes.NewBufferString(testString)))
		},
		testString,
		200,
	},
	{
		"ReaderFunc",
		func() {
			WriteStatusAndBody(200, func() io.Reader {
				return bytes.NewBufferString(testString)
			})
		},
		testString,
		200,
	},
	{
		"ReadCloserFunc",
		func() {
			WriteStatusAndBody(200, func() io.ReadCloser {
				return ioutil.NopCloser(bytes.NewBufferString(testString))
			})
		},
		testString,
		200,
	},
	{
		"interface{}",
		func() {
			WriteStatusAndBody(200, &testStruct{"joe", 28})
		},
		"{\"Name\":\"joe\",\"Age\":28}",
		200,
	},
	{
		"Respond",
		func() {
			Respond(200)
		},
		"",
		200,
	},
	{
		"Created",
		func() {
			Created()
		},
		"",
		201,
	},
	{
		"Bad Request",
		func() {
			RespondWithBadRequest(testString)
		},
		testString,
		400,
	},
	{
		"Internal Server Error",
		func() {
			RespondWithInternalServerError(testString)
		},
		testString,
		500,
	},
	{
		"Not Found",
		func() {
			NotFound()
		},
		"",
		404,
	},
}

func TestWriteStatusAndBody(t *testing.T) {
	for _, tst := range responseTests {
		t.Run(tst.Name, func(t *testing.T) {
			mockery := Mockery(func() {
				Endpoint("/foo", func() {
					Method("GET", func() {
						tst.WriteFn()
					})
				})
			})
			testURL, _ := url.ParseRequestURI("http://localhost/foo")
			request := &http.Request{
				Method: "GET",
				Header: http.Header{},
				URL:    testURL,
			}
			mockWriter := httptest.NewRecorder()
			mockery.ServeHTTP(mockWriter, request)
			result := mockWriter.Result()
			assert.Equal(t, tst.Status, result.StatusCode)
			body, err := ioutil.ReadAll(result.Body)
			if assert.NoError(t, err) {
				assert.Equal(t, tst.ExpectedBody, string(body))
			}
		})
	}
}
