package server

import (
	"context"
	"fmt"
	"github.com/carbonfive/go-filecoin-rest-api/handlers/v1"
	"github.com/carbonfive/go-filecoin-rest-api/types"
	"github.com/go-chi/chi"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"reflect"
	"time"
)

const DefaultPort = ":8080"

// V1Callbacks is a struct for callbacks configurable for the given API endpoint,
// shown by the 'path' tag
type V1Callbacks struct {
	GetActorByID func(string) (*types.Actor, error)
	GetActors    func() ([]byte, error)
	GetBlockByID func(string) ([]byte, error)
	GetNode      func() ([]byte, error)
}

// HTTPAPI is a struct containing all the things needed to serve the Filecoin HTTP API
type HTTPAPI struct {
	ctx    context.Context          // provided context
	srv    *http.Server             // HTTP server
	chimux *chi.Mux                 // Muxer
	hello  http.Handler             // hello handler
	v1h    *map[string]http.Handler // v1 handlers
}

// NewHTTPAPI creates and returns a *HTTPAPI using the provided context, *V1Callbacks,
// and desired port. If port <= 0, port 8080 will be used.
func NewHTTPAPI(ctx context.Context, cb1 *V1Callbacks, port int) *HTTPAPI {
	chimux := chi.NewRouter()

	gmux := mux.NewRouter()
	gmux.PathPrefix("/api/filecoin/v1")

	lport := DefaultPort
	if port > 0 {
		lport = fmt.Sprintf(":%d", port)
	}

	s := &http.Server{
		Addr:         lport,
		Handler:      chimux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	return &HTTPAPI{
		ctx:    ctx,
		srv:    s,
		chimux: chimux,
		hello:  &v1.HelloHandler{},
		v1h:    SetupV1Handlers(cb1),
	}
}

// Run sets up the route handlers using the provided callbacks, and starts
// the HTTP API server.
func (s *HTTPAPI) Run() {
	s.Route()

	go func() {
		if err := s.srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()
}

// Route sets up all the routes for the API
func (s *HTTPAPI) Route() {
	hdls := *s.v1h

	s.chimux.Route("/api/filecoin/v1", func(r chi.Router) {
		// TODO make this a documented connection check
		r.Handle("/hello", s.hello)

		r.Get("/control/node", hdls["GetNode"].ServeHTTP)
		r.Route("/actors", func(r chi.Router) {
			r.Get("/", hdls["GetActors"].ServeHTTP)
			r.Get("/{actorId}", hdls["GetActorByID"].ServeHTTP)
		})
	})
}

// Shutdown shuts down the http.Server.
func (s *HTTPAPI) Shutdown() error {
	return s.srv.Shutdown(s.ctx)
}

// Router returns the chimux
func (s *HTTPAPI) Router() chi.Router {
	return s.chimux
}

// SetupV1Handlers takes a V1Callback struct and creates a handler,
// setting the callback for each
func SetupV1Handlers(cb *V1Callbacks) *map[string]http.Handler {
	cb1t := reflect.TypeOf(*cb)
	cb1v := reflect.ValueOf(*cb)

	numCallbacks := cb1t.NumField()
	handlers := make(map[string]http.Handler, numCallbacks)

	for i := 0; i < numCallbacks; i++ {
		fieldName := cb1t.Field(i).Name

		fieldValue := cb1v.Field(i)
		if fieldValue.IsNil() {
			handlers[fieldName] = &v1.DefaultHandler{}
		} else {
			switch fieldName {
			case "GetActors":
				handlers[fieldName] = &v1.ActorsHandler{Callback: cb.GetActors}
			case "GetActorByID":
				handlers[fieldName] = &v1.ActorHandler{Callback: cb.GetActorByID}
			case "GetNode":
				handlers[fieldName] = &v1.NodeHandler{Callback: cb.GetNode}
			default:
				handlers[fieldName] = &v1.DefaultHandler{}
			}

		}
	}
	return &handlers
}
