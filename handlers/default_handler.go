package handlers

import (
	"fmt"
	logging "github.com/ipfs/go-log"
	"html"
	"net/http"
)

var log = logging.Logger("rest-api-handlers")


type DefaultHandler struct {
}

func (hh *DefaultHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err  := fmt.Fprintf(w, "%s is not implemented", html.EscapeString(r.RequestURI)); err != nil {
		log.Error("called unimplemented endpoint")
	}
}