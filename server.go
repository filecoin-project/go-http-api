package server

import (
	"context"
	"fmt"
	"github.com/carbonfive/go-filecoin-rest-api/handlers"
	"log"
	"net/http"
	"reflect"
)

type V1Callbacks struct {
	NodeID func() (string, error)                                  `path:"node_id"`
	Block  func(cid string, msgs bool, rcpts bool) (string, error) `path:"block"`
}

type HTTPAPI struct {
	ctx   context.Context
	srv   *http.Server
	smux  *http.ServeMux
	hello http.Handler
	v1cb  *V1Callbacks
}

// NewHTTPServer creates and returns a *HTTPAPI using the provided *V1Callbacks
func NewHTTPServer(ctx context.Context, cb1 *V1Callbacks, port int) *HTTPAPI {
	smux := http.NewServeMux()
	s := &http.Server{
		Addr:           ":" + fmt.Sprintf("%d", port),
		Handler: smux,
	}

	return &HTTPAPI{
		ctx:   ctx,
		srv:   s,
		smux: smux,
		hello: &handlers.HelloHandler{},
		v1cb:  cb1,
	}
}

func (s *HTTPAPI) Run() {
	// TODO make this a documented connection check
	s.AddHandler("/hello", s.hello)

	cb1t := reflect.TypeOf(*s.v1cb)
	cb1v := reflect.ValueOf(*s.v1cb)

	for i := 0; i < cb1t.NumField(); i++ {
		field := cb1t.Field(i)

		// get the path associated with the callback
		if path, ok := field.Tag.Lookup("path"); ok {

			// check if the callback is present
			// if so, provide it to the correct handler and otherwise
			// use the default handler.
			if cb1v.Field(i).IsNil() {
				s.AddHandler("/" + path, &handlers.DefaultHandler{})
			} else {
				switch path {
				case "node_id":
					s.AddHandler("/" + path, &handlers.NodeID{Callback: s.v1cb.NodeID})
				}
			}
		}
	}

	log.Fatal(s.srv.ListenAndServe())
}

func (s *HTTPAPI) Shutdown() {
	if err := s.srv.Shutdown(s.ctx); err != nil {
		log.Fatal("Server Shutdown failed: ", err)
	}
}

func (s *HTTPAPI) AddHandler(path string, hdlr http.Handler) {
	s.smux.Handle(path, hdlr)
}
