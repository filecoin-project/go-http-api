package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	logging "github.com/ipfs/go-log"

	"github.com/go-chi/chi"

	v1 "github.com/filecoin-project/go-http-api/handlers/v1"
)

// DefaultPort is the default port for the REST HTTP API
const DefaultPort = ":8080"

var log = logging.Logger("rest-api-server")

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

// NewHTTPAPI creates and returns a *HTTPAPI using the provided context, *Callbacks,
// and desired port. If port <= 0, port 8080 will be used.
func NewHTTPAPI(ctx context.Context, cb1 *v1.Callbacks, config Config) *HTTPAPI {
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
		v1h:    cb1.BuildHandlers(),
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
			r.Route("/messages", func(r chi.Router) {
				r.Get("/{messageCid}/wait", handlers["WaitForMessage"].ServeHTTP)
				r.Post("/", handlers["CreateMessage"].ServeHTTP)
			})
			r.Post("/signed-messages", handlers["SendSignedMessage"].ServeHTTP)
		})
		r.Route("/actors", func(r chi.Router) {
			r.Get("/", handlers["GetActors"].ServeHTTP)
			r.Route("/{actorId}", func(r chi.Router) {
				r.Get("/", handlers["GetActorByID"].ServeHTTP)
				r.Get("/nonce", handlers["GetActorNonce"].ServeHTTP)
			})
		})
	})
}

// Shutdown shuts down the http.Server.
func (s *HTTPAPI) Shutdown() error {
	return s.srv.Shutdown(s.ctx)
}

// Router returns the chimux router.  Currently used for testing (see server_test)
func (s *HTTPAPI) Router() chi.Router {
	return s.chimux
}

// Config returns a copy of the config
func (s *HTTPAPI) Config() Config {
	return s.config
}
