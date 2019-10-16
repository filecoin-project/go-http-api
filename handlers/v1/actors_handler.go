package v1

import (
	"net/http"

	"github.com/filecoin-project/go-http-api/handlers"
	"github.com/filecoin-project/go-http-api/types"
)

// ActorsHandler is the handler for the actors endpoint
type ActorsHandler struct {
	Callback func() ([]*types.Actor, error)
}

// ServeHTTP handles an HTTP request.
func (a *ActorsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	result, err := a.Callback()
	handlers.Respond(w, result, err)
}
