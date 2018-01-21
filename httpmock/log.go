package httpmock

import (
	"fmt"
	"net/http"
	"net/http/httputil"
)

func LogRequest() {
	outerMockHandler := currentMockHandler
	currentMockHandler = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		bytes, err := httputil.DumpRequest(r, true)
		if err == nil {
			fmt.Print("Request:")
			fmt.Println(string(bytes))
		}
		outerMockHandler.ServeHTTP(rw,r)
	})
}