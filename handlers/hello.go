package handlers

import (
	"fmt"
	"html"
	"net/http"
)

type HelloHandler struct {
}

func (hh *HelloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err  := fmt.Fprintf(w, "%s, world!", html.EscapeString(r.RequestURI)); err != nil {
		log.Error(err)
	}
}
