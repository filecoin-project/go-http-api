package handlers

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"rsc.io/quote"
)

type HelloHandler struct {
}

func (hh *HelloHandler) Hello() string {
	return quote.Hello()
}

func (hh *HelloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err  := fmt.Fprintf(w, "%s, world!", html.EscapeString(r.RequestURI)); err != nil {
		log.Fatal(err)
	}
}



