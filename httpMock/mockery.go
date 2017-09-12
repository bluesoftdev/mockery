package httpMock

import (
	"net/http"
	"sort"
)

type mockeryHandler struct {
	priority  int
	predicate RequestPredicate
	handler   http.Handler
}

type ByPriority []*mockeryHandler

func (a ByPriority) Len() int           { return len(a) }
func (a ByPriority) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByPriority) Less(i, j int) bool { return a[i].priority < a[j].priority }

type mockery struct {
	mux      *http.ServeMux
	handlers ByPriority
}

func (m *mockery) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	for _, h := range m.handlers {
		if h.predicate.Accept(request) {
			h.handler.ServeHTTP(w, request)
			return
		}
	}
	w.WriteHeader(404)
}

func (m *mockery) Handle(path string, handler http.Handler) {
	if m.mux == nil {
		m.mux = http.NewServeMux()
		m.HandleForCondition(DEFAULT_PRIORITY, PredicateFunc(func(r interface{}) bool {
			_, p := m.mux.Handler(r.(*http.Request))
			return p != ""
		}), m.mux)
	}
	m.mux.Handle(path, handler)
}

func (m *mockery) HandleForCondition(priority int, predicate RequestPredicate, handler http.Handler) {
	m.handlers = append(m.handlers, &mockeryHandler{priority, predicate, handler})
}

var (
	currentMockery     *mockery = nil
	currentMockHandler http.Handler
)

func Mockery(configFunc func()) http.Handler {
	currentMockery = &mockery{handlers: make(ByPriority, 0, 10)}
	currentMockHandler = NoopHandler
	defer func() { currentMockery = nil }()
	configFunc()
	sort.Stable(currentMockery.handlers)
	return currentMockery
}

func CurrentHandler() http.Handler {
	return currentMockHandler
}

var NoopHandler http.HandlerFunc = func(w http.ResponseWriter, request *http.Request) {
}

func DecorateHandler(preHandler, postHandler http.Handler) {
	delegate := CurrentHandler()
	currentMockHandler = http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		preHandler.ServeHTTP(w, request)
		delegate.ServeHTTP(w, request)
		postHandler.ServeHTTP(w, request)
	})
}
