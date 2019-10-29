package v1

import (
	"math/big"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/filecoin-project/go-http-api/handlers"
)

// ActorNonceHandler is the handler for the GET /actors/{actorID}/nonce endpoint.
type ActorNonceHandler struct {
	Callback func(string) (*big.Int, error)
}

// ServeHTTP handles an HTTP request.
func (anh *ActorNonceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	actorID := chi.URLParam(r, "actorId")

	actor, err := anh.Callback(actorID)
	handlers.Respond(w, actor, err)
}
