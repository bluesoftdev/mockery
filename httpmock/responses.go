package httpmock

import (
	"encoding/json"
	"net/http"
)

// Writes the given status with the given body.  It is expected that any headers needed by the response have been added
// as this will being the sending of the response.
func WriteStatusAndBody(status int, body interface{}) {
	var bytes []byte
	var contentType string
	if bodyStr, ok := body.(string); ok {
		bytes = ([]byte)(bodyStr)
	} else {
		var err error
		bytes, err = json.Marshal(body)
		if err != nil {
			panic("unable to marshal Body to json!")
		}
		contentType = "application/json"
	}
	DecorateHandler(NoopHandler, http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		if contentType != "" {
			w.Header().Set("Content-Type", "application/json")
		}
		w.WriteHeader(status)
		w.Write(bytes)
	}))
}

// Shortcut or returning Created (201) http response with no body.
func Created() {
	Respond(201)
}

// Shortcut for returning Bad Request (400) with the given body.
func RespondWithBadRequest(body interface{}) {
	WriteStatusAndBody(400, body)
}

// Shortcut for returning Internal Server Error (500) with the given body.
func RespondWithInternalServerError(body interface{}) {
	WriteStatusAndBody(500, body)
}

// Shortcut for returning Not Found (404) with no body.
func NotFound() {
	Respond(404)
}
