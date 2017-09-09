package httpMock

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"container/list"
	"errors"
)

interface CanHandle {
	CanHandle(r *http.Request) bool
}

var handlerStack = list.New()

func pushHandler(handler http.Handler) {
	handlerStack.PushBack(handler);
}

func peekHandler() http.Handler {
	if handlerStack.Len() == 0 {
		panic("ther is no current handler!!!")
	}
	e := handlerStack.Back()
	return e.Value.(http.Handler)
}

func popHandler() http.Handler {
	if handlerStack.Len() == 0 {
		panic("ther is no current handler!!!")
	}
	e := handlerStack.Back()
	handlerStack.Remove(e)
	return e.Value.(http.Handler)
}

func unwindHandlersUpTo(_mock http.Handler) {
	for peekHandler() != _mock {
		popHandler()
	}
}



type mock struct {
	url     string
	methods map[string]http.Handler
}

var currentMockery *http.ServeMux = nil

func Mockery(configFunc func()) http.Handler {
	mockery := http.NewServeMux()
	currentMockery = mockery
	pushHandler(mockery)
	defer popHandler()
	defer func() { currentMockery = nil }()
	configFunc()
	unwindHandlersUpTo(mockery);
	return mockery
}

type mock struct {
	url     string

	methods map[string]http.Handler
}

func Endpoint(url string, configureFunc func()) {
	_mock := &mock{url: url, methods: make(map[string]http.Handler)}
	currentMock = _mock
	pushHandler(_mock)
	defer popHandler()
	defer func() { currentMock = nil }()
	configureFunc()

	handler := peekHandler();
	unwindHandlersUpTo(_mock);
	currentMockery.Handle(url, handler)
}

func (m* mock) ServeHTTP(w http.ResponseWriter, request *http.Request) {

}

type mockMethod struct {
	method           string
	delegate http.Handler
}

func respondWith404(w http.ResponseWriter, request *http.Request) {
	w.WriteHeader(404)
}

func Method(method string, configFunc func()) {
	_mockMethod := &mockMethod{method: method, delegate: respondWith404 }
	pushHandler(_mockMethod)
	configFunc()
	_mockMethod.delegate = popHandler()
}

func (mm* mockMethod) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	mm.delegate.ServeHTTP(w,request)
}

func (mm* mockMethod) CanHandler(request *http.Request) {
	return request.Method == mm.method
}

type headerHandler struct {
	name, value string
	delegate http.Handler
}

func (hh *headerHandler) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	w.Header().Add(hh.name,hh.value)
	hh.ServeHTTP(w,request)
}

func Header(name, value string) {
	pushHandler(&headerHandler{name,value,popHandler()})
}

type respondWithFile struct {
	status int
	fileName string
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
	w.WriteHeader(rwf.status)
	_, err = io.Copy(w, file)
	if err != nil {
		log.Printf("ERROR while serving up a file: %+v", err)
	}
}

func RespondWithFile(status int, fileName string) {
	pushHandler(&respondWithFile{status,fileName})
}

type respondWithStatus struct {
	status int
}

func (rws* respondWithStatus) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	w.WriteHeader(rws.status)
}

func Respond(status int) {
	pushHandler(&respondWithStatus{status})
}


