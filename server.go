package server

import (
	"context"
	"fmt"
	"github.com/carbonfive/go-filecoin-rest-api/handlers"
	"net/http"
	"reflect"
)

const DefaultPort = ":8080"

// V1Callbacks is a struct for callbacks configurable for the given API endpoint,
// shown by the 'path' tag
type V1Callbacks struct {
	NodeID func() (string, error)                                  `path:"node_id"`
	Block  func(cid string, msgs bool, rcpts bool) (string, error) `path:"block"`
}

// HTTPAPI is a struct containing all the things needed to serve the Filecoin HTTP API
type HTTPAPI struct {
	ctx   context.Context
	srv   *http.Server
	smux  *http.ServeMux
	hello http.Handler
	v1cb  *V1Callbacks
}

// NewHTTPServer creates and returns a *HTTPAPI using the provided context, *V1Callbacks,
// and desired port. If port <= 0, port 8080 will be used.
func NewHTTPServer(ctx context.Context, cb1 *V1Callbacks, port int) *HTTPAPI {
	smux := http.NewServeMux()

	lport := DefaultPort
	if port > 0 {
		lport = fmt.Sprintf(":%d", port)
	}

	s := &http.Server{
		Addr:    lport,
		Handler: smux,
	}

	return &HTTPAPI{
		ctx:   ctx,
		srv:   s,
		smux:  smux,
		hello: &handlers.HelloHandler{},
		v1cb:  cb1,
	}
}

// Run sets up the route handlers using the provided callbacks, and starts
// the HTTP API server.
func (s *HTTPAPI) Run() error {
	// TODO make this a documented connection check
	s.AddHandler("/hello", s.hello)

	// This allows us to avoid a lot of duplicate code, especially for dealing with callback
	// functions that have different signatures.
	cb1t := reflect.TypeOf(*s.v1cb)
	cb1v := reflect.ValueOf(*s.v1cb)

	for i := 0; i < cb1t.NumField(); i++ {
		field := cb1t.Field(i)

		// get the path associated with the callback
		if path, ok := field.Tag.Lookup("path"); ok {

			if cb1v.Field(i).IsNil() {
				s.AddHandler("/"+path, &handlers.DefaultHandler{})
			} else {
				switch path {
				case "node_id":
					s.AddHandler("/"+path, &handlers.NodeID{Callback: s.v1cb.NodeID})
				}
			}
		}
	}

	return s.srv.ListenAndServe()
}

// Shutdown shuts down the http.Server.
func (s *HTTPAPI) Shutdown() error {
	return s.srv.Shutdown(s.ctx)
}

// AddHandler adds the handler to the http.Server's ServeMux
func (s *HTTPAPI) AddHandler(path string, hdlr http.Handler) {
	s.smux.Handle(path, hdlr)
}
