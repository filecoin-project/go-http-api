package v1

import (
	"github.com/carbonfive/go-filecoin-rest-api/handlers"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/carbonfive/go-filecoin-rest-api/types"
)

// ActorHandler is a handler for the actors/{actorId} endpoint
type ActorHandler struct {
	Callback func(actorId string) (*types.Actor, error)
}

// ServeHTTP handles an HTTP request.
func (a *ActorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	actorId := chi.URLParam(r, "actorId")
	actor, err := a.Callback(actorId)
	handlers.Respond(w, actor, err)
}
