package v1

import (
	"fmt"
	"github.com/carbonfive/go-filecoin-rest-api/handlers/api_errors"
	"html"
	"net/http"
)

type HelloHandler struct {
}

func (hh *HelloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := fmt.Fprintf(w, "%s, world!", html.EscapeString(r.RequestURI)); err != nil {
		fmt.Fprint(w, api_errors.MarshalErrors([]string{err.Error()})) // nolint: errcheck
		log.Error(err)
	}
}
