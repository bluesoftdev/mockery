package httpmock

import (
	"log"
	"net/http"
	"net/http/httputil"
)

// LogRequest will cause the request information to be logged to the console.
func LogRequest() {
	DecorateHandlerBefore(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		bytes, err := httputil.DumpRequest(r, true)
		if err == nil {
			log.Printf("Request:\n%s", string(bytes))
		}
	}))
}
