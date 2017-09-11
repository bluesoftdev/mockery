package httpMock

import (
	"encoding/json"
	"net/http"
)

func WriteStatusAndBody(status int, body interface{}) {
	var bytes []byte
	var contentType string
	if bodyStr, ok := body.(string); ok {
		bytes = ([]byte)(bodyStr)
	} else {
		var err error
		bytes, err = json.Marshal(body)
		if err != nil {
			panic("unable to marshal body to json!")
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

func RespondWithBadRequest(body interface{}) {
	WriteStatusAndBody(400, body)
}

func RespondWithInternalServerError(body interface{}) {
	WriteStatusAndBody(500, body)
}

func NotFound() {
	Respond(404)
}
