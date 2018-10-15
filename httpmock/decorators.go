package httpmock

import (
	"os"
	//"fmt"
	"net/http"
	"log"
	"io"
	"encoding/json"
	"bytes"
	"runtime"
	"fmt"
)

// Adds a header to the response, may be called at any time.
func Header(name, value string) {
	DecorateHandler(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.Header().Add(name, value)
	}), NoopHandler)
}


// Adds a trailer to the response, must be called after the response body has been specified.
func Trailer(name, value string) {
	DecorateHandler(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.Header().Add("Trailer", name)
	}), http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.Header().Add(name, value)
	}))
}

// Adds a response code and body to the response.  The jsonBody parameter is JSON encoded using the json encoder in the
// encoding/json package.
func RespondWithJson(status int, jsonBody interface{} ) {
	var data bytes.Buffer
	json.NewEncoder(&data).Encode(jsonBody)
	DecorateHandler(NoopHandler, http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		w.Write(data.Bytes())
	}))
}

// Responds with the status code given and the content of the file specified.
func RespondWithFile(status int, fileName string) {
	DecorateHandler(NoopHandler, http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
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
		err = file.Close()
		if err !=  nil {
			log.Printf("ERROR while serving up a file: %+v", err)
		}
	}))
}

// Responds with the status code given and the body
func RespondWithString(status int, body string) {
	DecorateHandler(NoopHandler, http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.WriteHeader(status)
		w.Write([]byte(body))
	}))
}

func RespondWithReader(status int, bodyProducer func() io.Reader) {
	DecorateHandler(NoopHandler, http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.WriteHeader(status)
		io.Copy(w,bodyProducer())
	}))
}

func Respond(status int) {
	DecorateHandler(NoopHandler, http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.WriteHeader(status)
	}))
}


func LogLocation(comment string) {
	frames := make([]uintptr,1)
	runtime.Callers(2, frames)
	fun := runtime.FuncForPC(frames[0]-1)
	var fileLocation string
	if fun == nil {
		fileLocation = "Unknown"
	} else {
		file, line := fun.FileLine(frames[0])
		fileLocation = fmt.Sprintf("%s:%d(%s)", file, line, fun.Name())
	}
	DecorateHandler(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		log.Printf("Endpoint Defined at %s: %s",fileLocation,comment)
	}),NoopHandler)
}