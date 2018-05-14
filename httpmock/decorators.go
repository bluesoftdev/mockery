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
		//w.Header().Add("Content-Length", fmt.Sprintf("%d", data.Len()))
		w.WriteHeader(status)
		w.Write(data.Bytes())
	}))
}

func RespondWithFile(status int, fileName string) {
	DecorateHandler(NoopHandler, http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		//stat, err := os.Stat(fileName)
		//if err != nil {
		//	log.Printf("ERROR while serving up a file: %+v", err)
		//	w.WriteHeader(500)
		//	return
		//}
		//w.Header().Add("Content-Length", fmt.Sprintf("%d", stat.Size()))
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

func RespondWithString(status int, body string) {
	DecorateHandler(NoopHandler, http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		//w.Header().Add("Content-Length", fmt.Sprintf("%d", len(body)))
		w.WriteHeader(status)
		w.Write([]byte(body))
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