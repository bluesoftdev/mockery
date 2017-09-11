package httpMock

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var currentMockery *http.ServeMux = nil

func Mockery(configFunc func()) http.Handler {
	mockery := http.NewServeMux()
	currentMockery = mockery
	defer func() { currentMockery = nil }()
	configFunc()
	return mockery
}

type mock struct {
	url     string
	methods map[string]http.Handler
}

var currentMock *mock = nil

func Endpoint(url string, configureFunc func()) {
	_mock := &mock{url: url, methods: make(map[string]http.Handler)}
	currentMock = _mock
	defer func() { currentMock = nil }()
	configureFunc()
	currentMockery.Handle(url, _mock)
}

func (m *mock) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	handler, ok := m.methods[request.Method]
	if ok {
		handler.ServeHTTP(w, request)
	} else {
		w.WriteHeader(404)
	}
}

type mockMethod struct {
	method           string
}

var (
	currentMockMethod        *mockMethod = nil
	currentMockMethodHandler http.Handler
)

func (mm *mockMethod) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	// Nothing to do.
	// The decorators should have taken care of everything by this point.
}

func Method(method string, configFunc func()) {
	_mockMethod := &mockMethod{method: method}
	currentMockMethod = _mockMethod
	currentMockMethodHandler = _mockMethod
	defer func() { currentMockMethod = nil; currentMockMethodHandler = nil }()
	configFunc()
	currentMock.methods[method] = currentMockMethodHandler
}

func Header(name, value string) {
	DecorateHandler(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.Header().Add(name, value)
	}), NoopHandler)
}

func Trailer(name, value string) {
	DecorateHandler(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.Header().Add("Trailer", name)
	}),http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.Header().Add(name, value)
	}))
}

type respondWithFile struct {
	statusCode int
	fileName   string
}

func (rwf *respondWithFile) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	stat, err := os.Stat(rwf.fileName)
	if err != nil {
		log.Printf("ERROR while serving up a file: %+v", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Add("Content-Length", fmt.Sprintf("%d", stat.Size()))
	file, err := os.Open(rwf.fileName)
	if err != nil {
		log.Printf("ERROR while serving up a file: %+v", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(rwf.statusCode)
	_, err = io.Copy(w, file)
	if err != nil {
		log.Printf("ERROR while serving up a file: %+v", err)
	}
}

func RespondWithFile(status int, fileName string) {
	DecorateHandler(NoopHandler,&respondWithFile{status, fileName})
}

type respondWithStatus struct {
	statusCode int
}

func (rws *respondWithStatus) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	w.WriteHeader(rws.statusCode)
}

func Respond(status int) {
	DecorateHandler(NoopHandler,&respondWithStatus{status})
}

func CurrentHandler() http.Handler {
	return currentMockMethodHandler
}

var NoopHandler http.HandlerFunc = func(w http.ResponseWriter, request *http.Request) {
}

func DecorateHandler(preHandler, postHandler http.Handler) {
	delegate := CurrentHandler()
	currentMockMethodHandler = http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		preHandler.ServeHTTP(w,request)
		delegate.ServeHTTP(w,request)
		postHandler.ServeHTTP(w, request)
	})
}
