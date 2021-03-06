package v1

import (
	"fmt"
	"net/http"

	"github.com/filecoin-project/go-http-api/handlers"
)

// HelloHandler is a handler for the hello endpoint.
// It is intended to test connection to the API
type HelloHandler struct {
}

// ServeHTTP handles an HTTP request.
func (hh *HelloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := fmt.Fprint(w, "HELLO"); err != nil {
		handlers.Respond(w, "hello", err)
	}
}
