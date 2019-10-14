package v1

import (
	"fmt"
	"html"
	"net/http"

	"github.com/filecoin-project/go-http-api/types"
)

// HelloHandler is a handler for the hello endpoint.
// It is intended to test connection to the API
type HelloHandler struct {
}

// ServeHTTP handles an HTTP request.
func (hh *HelloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := fmt.Fprintf(w, "%s, world!", html.EscapeString(r.RequestURI)); err != nil {
		log.Error(err)
		if _, err = fmt.Fprint(w, types.MarshalErrors([]string{err.Error()})); err != nil {
			log.Error(err)
		}
	}
}
