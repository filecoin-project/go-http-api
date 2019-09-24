package server

import (
	"log"
	"net/http"
)

type HTTPServer struct {
	hello http.Handler
}

func NewHTTPServer(hello http.Handler) *HTTPServer {
	return &HTTPServer{hello: hello}
}

func (s *HTTPServer) Run() {
	http.Handle("/hello", s.hello)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
