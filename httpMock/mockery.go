package httpMock

import (
	"net/http"
	"io"
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
	methods map[string]*mockMethod
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
	_mock := &mock{url: url, methods: make(map[string]*mockMethod) }
	currentMock = _mock
	defer func() { currentMock = nil }()
	configureFunc()
	currentMockery.Handle(url,_mock)
}

var currentMockMethod *mockMethod = nil

func Method(method string, configFunc func()) {
	_mockMethod := &mockMethod{method: method, statusCode: 200, responseFileName: "", headers: make(map[string]string)}
	currentMockMethod = _mockMethod
	defer func() { currentMockMethod = nil }()
	configFunc()
	currentMock.methods[method] = _mockMethod
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
	file, err := os.Open(mm.responseFileName)
	if err != nil {
		// TODO: Log the error
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(mm.statusCode)
	_, err = io.Copy(w, file)
	if err != nil {
		// TODO: Log the error but we've already written the headers so we will just hope for the best.
	}
}

// handler, err := Mockery(func() {
//   Endpoint("/foo/bar/",func() {
//     Method("GET", func() {
//       Header("foo","bar")
//       RespondWithFile(200,"foo.json")
//     })
//   })
// })
// ...
// http.Server(8080, handler)...