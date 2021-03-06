package v1

import (
	"net/http"

	"github.com/go-chi/chi"

	"github.com/filecoin-project/go-http-api/handlers"
	"github.com/filecoin-project/go-http-api/types"
)

// ActorHandler is a handler for the actors/{actorId} endpoint
type ActorHandler struct {
	Callback func(actorId string) (*types.Actor, error)
}

// ServeHTTP handles an HTTP request.
func (a *ActorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	actorID := chi.URLParam(r, "actorID")
	actor, err := a.Callback(actorID)
	handlers.Respond(w, actor, err)
}
