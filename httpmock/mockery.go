package httpmock

import (
	"net/http"
	"sort"
)

type mockeryHandler struct {
	priority  int
	predicate Predicate
	handler   http.Handler
}

type byPriority []*mockeryHandler

func (a byPriority) Len() int           { return len(a) }
func (a byPriority) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byPriority) Less(i, j int) bool { return a[i].priority < a[j].priority }

type mockery struct {
	mux      *http.ServeMux
	handlers byPriority
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

func (m *mockery) HandleForCondition(priority int, predicate Predicate, handler http.Handler) {
	m.handlers = append(m.handlers, &mockeryHandler{priority, predicate, handler})
}

var (
	currentMockery     *mockery = nil
	currentMockHandler http.Handler
)

// Mockery contains the top level dispatcher.  This method establishes the root handler and the configFunc is called to
// create handlers for the various mocks.  Once the config method returns some clean up actions will occur and the
// mock handler will be returned.
func Mockery(configFunc func()) http.Handler {
	currentMockery = &mockery{handlers: make(byPriority, 0, 10)}
	currentMockHandler = NoopHandler
	defer func() { currentMockery = nil }()
	configFunc()
	sort.Stable(currentMockery.handlers)
	return currentMockery
}

// CurrentHandler returns the current handler that should be decorated with any additional behaviors.
func CurrentHandler() http.Handler {
	return currentMockHandler
}

// NoopHandler is a handler that does nothing.
var NoopHandler http.HandlerFunc = func(w http.ResponseWriter, request *http.Request) {
}

// DecorateHandler is used by DSL methods to inject pre & post actions to the current handler.  For instance, the
// Header(string,string) function adds a preHandler that adds a Header to the ResponseWriter.  To use this function
// to create new DSL Methods, follow this pattern:
//
//    func Header(name, value string) {
//      return Decoratehandler(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
//				w.Header.Add(name,value)
//    	}), NoopHandler)
//    }
//
func DecorateHandler(preHandler, postHandler http.Handler) {
	delegate := CurrentHandler()
	currentMockHandler = http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		preHandler.ServeHTTP(w, request)
		delegate.ServeHTTP(w, request)
		postHandler.ServeHTTP(w, request)
	})
}
