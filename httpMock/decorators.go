package httpMock

import (
	"os"
	"fmt"
	"net/http"
	"log"
	"io"
	"encoding/json"
	"bytes"
)

func Header(name, value string) {
	DecorateHandler(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.Header().Add(name, value)
	}), NoopHandler)
}

func Trailer(name, value string) {
	DecorateHandler(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.Header().Add("Trailer", name)
	}), http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.Header().Add(name, value)
	}))
}

func RespondWithJson(status int, jsonBody interface{} ) {
	var data bytes.Buffer
	json.NewEncoder(&data).Encode(jsonBody)
	DecorateHandler(NoopHandler, http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Add("Content-Length", fmt.Sprintf("%d", data.Len()))
		w.WriteHeader(status)
		w.Write(data.Bytes())
	}))
}

func RespondWithFile(status int, fileName string) {
	DecorateHandler(NoopHandler, http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		stat, err := os.Stat(fileName)
		if err != nil {
			log.Printf("ERROR while serving up a file: %+v", err)
			w.WriteHeader(500)
			return
		}
		w.Header().Add("Content-Length", fmt.Sprintf("%d", stat.Size()))
		file, err := os.Open(fileName)
		if err != nil {
			log.Printf("ERROR while serving up a file: %+v", err)
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(status)
		_, err = io.Copy(w, file)
		if err != nil {
			log.Printf("ERROR while serving up a file: %+v", err)
		}
	}))
}

func RespondWithString(status int, body string) {
	DecorateHandler(NoopHandler, http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.Header().Add("Content-Length", fmt.Sprintf("%d", len(body)))
		w.WriteHeader(status)
		w.Write([]byte(body))
	}))
}
func Respond(status int) {
	DecorateHandler(NoopHandler, http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.WriteHeader(status)
	}))
}