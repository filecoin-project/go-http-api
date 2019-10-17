package v1

import (
	"fmt"
	"html"
	"net/http"
)

// DefaultHandler is a fallback handler for a supported API endpoint that was
// not provided a callback when the API was created.  Note this is distinct from
// 404 Not Found which is the response to a request for a non-existent, unsupported
// endpoint.
type DefaultHandler struct{}

// ServeHTTP handles an HTTP request.
func (hh *DefaultHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ep := html.EscapeString(r.RequestURI)
	if _, err := fmt.Fprintf(w, "%s is not implemented", ep); err != nil {
		log.Errorf("Called unimplemented endpoint %s", ep)
	}
}
