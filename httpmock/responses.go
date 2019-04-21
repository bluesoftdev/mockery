package httpmock

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

// WriteStatusAndBody writes the given status with the given body.  It is expected that any headers needed by the
// response have been added as this will being the sending of the response.
func WriteStatusAndBody(status int, body interface{}) {
	var bodyProvider func() io.ReadCloser
	var contentType string
	switch bdy := body.(type) {
	case []byte:
		bodyProvider = func() io.ReadCloser { return ioutil.NopCloser(bytes.NewBuffer(bdy)) }
	case string:
		bodyProvider = func() io.ReadCloser { return ioutil.NopCloser(bytes.NewBufferString(bdy)) }
	case io.Reader:
		bodyProvider = func() io.ReadCloser { return ioutil.NopCloser(bdy) }
	case func() io.Reader:
		bodyProvider = func() io.ReadCloser { return ioutil.NopCloser(bdy()) }
	case func() io.ReadCloser:
		bodyProvider = bdy
	default:
		bdyBytes, err := json.Marshal(bdy)
		if err != nil {
			panic("unable to marshal Body to json!")
		}
		bodyProvider = func() io.ReadCloser { return ioutil.NopCloser(bytes.NewBuffer(bdyBytes)) }
		Header("Content-Type", "application/json")
	}
	DecorateHandler(NoopHandler, http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		if contentType != "" {
			w.Header().Set("Content-Type", "application/json")
		}
		w.WriteHeader(status)
		io.Copy(w, bodyProvider())
	}))
}

// Created is a shortcut for returning Created (201) http response with no body.
func Created() {
	Respond(201)
}

// RespondWithBadRequest is a shortcut for returning Bad Request (400) with the given body.
func RespondWithBadRequest(body interface{}) {
	WriteStatusAndBody(400, body)
}

// RespondWithInternalServerError is a shortcut for returning Internal Server Error (500) with the given body.
func RespondWithInternalServerError(body interface{}) {
	WriteStatusAndBody(500, body)
}

// NotFound is a shortcut for returning Not Found (404) with no body.
func NotFound() {
	Respond(404)
}
