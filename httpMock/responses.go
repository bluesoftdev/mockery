package httpMock

import (
	"net/http"
	"encoding/json"
)

var NotFound http.HandlerFunc = func(w http.ResponseWriter,request *http.Request) {
	w.WriteHeader(404)
}

func WriteStatusAndBody(w http.ResponseWriter, status int, body interface{}) {
	bytes, err := json.Marshal(body)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(`{ "error": "`+err.Error()+`"`))
		return
	}
	w.WriteHeader(400)
	w.Write(bytes)
}

func BadRequest(body interface{}) http.HandlerFunc  {
	return func(w http.ResponseWriter,request *http.Request) {
		WriteStatusAndBody(w,400,body)
	}
}

func InternalServerError(body interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter,request *http.Request) {
		WriteStatusAndBody(w,500,body)
	}
}
