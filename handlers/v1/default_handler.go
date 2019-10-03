package v1

import (
	"fmt"
	"html"
	"net/http"

	logging "github.com/ipfs/go-log"
)

var log = logging.Logger("rest-api-handlers")

// DefaultHandler is a fallback handler for a supported API endpoint that was
// not provided a callback when the API was created.  Note this is distinct from
// 404 Not Found which is the response to a request for a non-existent, unsupported
// endpoint.
type DefaultHandler struct {
}

func (hh *DefaultHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := fmt.Fprintf(w, "%s is not implemented", html.EscapeString(r.RequestURI)); err != nil {
		log.Error("called unimplemented endpoint")
	}
}
