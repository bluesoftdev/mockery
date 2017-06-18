package httpMock

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type mockMethod struct {
	method           string
	statusCode       int
	responseFileName string
	headers          map[string]string
}

type mock struct {
	url     string
	methods map[string]http.Handler
}

var currentMockery *http.ServeMux = nil

func Mockery(configFunc func()) http.Handler {
	mockery := http.NewServeMux()
	currentMockery = mockery
	defer func() { currentMockery = nil }()
	configFunc()
	return mockery
}

var currentMock *mock = nil

func Endpoint(url string, configureFunc func()) {
	_mock := &mock{url: url, methods: make(map[string]http.Handler)}
	currentMock = _mock
	defer func() { currentMock = nil }()
	configureFunc()
	currentMockery.Handle(url, _mock)
}

var (
	currentMockMethod        *mockMethod = nil
	currentMockMethodHandler http.Handler
)

func Method(method string, configFunc func()) {
	_mockMethod := &mockMethod{method: method, statusCode: 200, responseFileName: "", headers: make(map[string]string)}
	currentMockMethod = _mockMethod
	currentMockMethodHandler = _mockMethod
	defer func() { currentMockMethod = nil; currentMockMethodHandler = nil }()
	configFunc()
	currentMock.methods[method] = currentMockMethodHandler
}

func Header(name, value string) {
	currentMockMethod.headers[name] = value
}

func RespondWithFile(status int, fileName string) {
	currentMockMethod.statusCode = status
	currentMockMethod.responseFileName = fileName
}

func (m *mock) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	handler, ok := m.methods[request.Method]
	if ok {
		handler.ServeHTTP(w, request)
	} else {
		w.WriteHeader(404)
	}
}

func (mm *mockMethod) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	for hn, hv := range mm.headers {
		w.Header().Add(hn, hv)
	}
	stat, err := os.Stat(mm.responseFileName)
	if err != nil {
		log.Printf("ERROR while serving up a file: %+v", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Add("Content-Length", fmt.Sprintf("%d", stat.Size()))
	file, err := os.Open(mm.responseFileName)
	if err != nil {
		log.Printf("ERROR while serving up a file: %+v", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(mm.statusCode)
	_, err = io.Copy(w, file)
	if err != nil {
		log.Printf("ERROR while serving up a file: %+v", err)
	}
}
