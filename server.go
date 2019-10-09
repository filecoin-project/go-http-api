package server

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"time"

	logging "github.com/ipfs/go-log"

	"github.com/go-chi/chi"

	v1 "github.com/carbonfive/go-filecoin-rest-api/handlers/v1"
	"github.com/carbonfive/go-filecoin-rest-api/types"
)

const DefaultPort = ":8080"

var log = logging.Logger("rest-api-server")

// V1Callbacks is a struct for callbacks configurable for the given API endpoint,
// shown by the 'path' tag
// To add a new endpoint:
//   * Write a new handler to use a new V1Callback, with tests
//   * Add a new callback name/signature to V1Callbacks
//   * Add a case to SetupV1Handlers that uses the callback
type V1Callbacks struct {
	GetActorByID   func(string) (*types.Actor, error)
	GetActors      func() ([]*types.Actor, error)
	GetBlockByID   func(string) (*types.Block, error)
	CreateMessage func(*types.Message) (*types.Message, error)
	GetMessageByID func(string) (*types.Message, error)
	GetNode        func() (*types.Node, error)
}

// HTTPAPI is a struct containing all the things needed to serve the Filecoin HTTP API
type HTTPAPI struct {
	ctx    context.Context          // provided context
	srv    *http.Server             // HTTP server
	chimux *chi.Mux                 // Muxer
	hello  http.Handler             // hello handler
	v1h    *map[string]http.Handler // v1 handlers
	config Config                   // API configuration
}

// Config is implementation-specific configuration for the API
type Config struct {
	Port        int
	TLSCertPath string
	TLSKeyPath  string
}

// NewHTTPAPI creates and returns a *HTTPAPI using the provided context, *V1Callbacks,
// and desired port. If port <= 0, port 8080 will be used.
func NewHTTPAPI(ctx context.Context, cb1 *V1Callbacks, config Config) *HTTPAPI {
	chimux := chi.NewRouter()

	lport := DefaultPort
	if config.Port > 0 {
		lport = fmt.Sprintf(":%d", config.Port)
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
		config: config,
	}
}

// Run sets up the route handlers using the provided callbacks, and starts
// the HTTP API server.
func (s *HTTPAPI) Run() *HTTPAPI {
	s.Route()

	go func() {
		if s.config.TLSCertPath != "" && s.config.TLSKeyPath != "" {
			if err := s.srv.ListenAndServeTLS(s.config.TLSCertPath, s.config.TLSKeyPath); err != nil {
				log.Error(err)
			}
		} else {
			if err := s.srv.ListenAndServe(); err != nil {
				log.Error(err)
			}
		}
	}()
	return s
}

// Route sets up all the routes for the API
func (s *HTTPAPI) Route() {
	handlers := *s.v1h

	s.chimux.Route("/api/filecoin/v1", func(r chi.Router) {
		// TODO make this a documented connection check
		r.Handle("/hello", s.hello)

		r.Get("/control/node", handlers["GetNode"].ServeHTTP)
		r.Route("/chain", func(r chi.Router) {
			r.Get("/blocks/{blockId}", handlers["GetBlockByID"].ServeHTTP)
			r.Get("/executed-messages/{executedMessageId}", handlers["GetMessageByID"].ServeHTTP)
			r.Post("/messages", handlers["CreateMessage"].ServeHTTP)
		})
		r.Route("/actors", func(r chi.Router) {
			r.Get("/", handlers["GetActors"].ServeHTTP)
			r.Get("/{actorId}", handlers["GetActorByID"].ServeHTTP)
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

// Config returns a copy of the config
func (s *HTTPAPI) Config() Config {
	return s.config
}

// SetupV1Handlers takes a V1Callback struct and iterates over all
// functions, creating a handler with a callback for each supported endpoint.
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
			case "GetBlockByID":
				handlers[fieldName] = &v1.BlockHandler{Callback: cb.GetBlockByID}
			case "GetActors":
				handlers[fieldName] = &v1.ActorsHandler{Callback: cb.GetActors}
			case "GetActorByID":
				handlers[fieldName] = &v1.ActorHandler{Callback: cb.GetActorByID}
			case "CreateMessage":
				handlers[fieldName] = &v1.CreateMessageHandler{Callback: cb.CreateMessage }
			case "GetMessageByID":
				handlers[fieldName] = &v1.MessageHandler{Callback: cb.GetMessageByID}
			case "GetNode":
				handlers[fieldName] = &v1.NodeHandler{Callback: cb.GetNode}
			default:
				log.Errorf("skipping unknown handler: %s.", fieldName)
			}
		}
	}
	return &handlers
}
